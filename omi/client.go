package omi

import (
	"fmt"
	"net/http"
	//"net/http/cookiejar"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"github.com/imdario/mergo"
	"io/ioutil"
	"log"
	"net/http/cookiejar"
	"net/url"
	"os"
)

type ClientConfig struct {
	Username string
	Password string
	BaseUrl  string `json:"base_url"`
}

type Client struct {
	SecureModifyToken string
	HttpClient        *http.Client
	Uri               *url.URL
	Jar               *cookiejar.Jar
	ClientConfig      ClientConfig
	Log               *log.Logger
}

var (
	EventListPath = "/opr-web/rest/9.10/event_list"
	PingPath      = "/opr-web/rest"
)

func NewClient(configContents []byte) *Client {
	client := new(Client)
	var config ClientConfig
	if err := json.Unmarshal(configContents, &config); err != nil {
		return client
	}
	client.ClientConfig = config
	client.Jar, _ = cookiejar.New(nil)
	client.HttpClient = &http.Client{
		Jar: client.Jar,
	}
	client.Uri, _ = url.Parse(client.ClientConfig.BaseUrl)
	client.Log = log.New(os.Stdout, "GOMI: ", 0)
	return client
}

func (client *Client) Ping() bool {
	path := fmt.Sprintf("%s%s", client.ClientConfig.BaseUrl, PingPath)
	req, err := http.NewRequest("GET", path, nil)
	req.SetBasicAuth(client.ClientConfig.Username, client.ClientConfig.Password)
	resp, err := client.HttpClient.Do(req)
	for _, cook := range client.Jar.Cookies(client.Uri) {
		if cook.Name == "secureModifyToken" {
			client.SecureModifyToken = cook.Value
		}
	}
	if err != nil {
		return false
	}
	if client.SecureModifyToken != "" && resp.StatusCode == 200 {
		return true
	}
	return false
}

func (client *Client) Do(method, path string, body []byte) (*http.Response, error) {
	fullPath := fmt.Sprintf("%s%s", client.ClientConfig.BaseUrl, path)
	req, err := http.NewRequest(method, fullPath, bytes.NewBuffer([]byte(body)))
	req.SetBasicAuth(client.ClientConfig.Username, client.ClientConfig.Password)
	req.Header.Add("Accept", "application/xml")
	if (method == "POST") || (method == "PUT") {
		req.Header.Set("Content-Type", "application/xml")
	}
	if client.SecureModifyToken != "" {
		req.Header.Add("X-Secure-Modify-Token", client.SecureModifyToken)
	}

	resp, err := client.HttpClient.Do(req)

	if err != nil {
		client.Log.Printf("Error during HTTP Call to OMI: %v", err)
		return nil, err
	}
	if (resp.StatusCode == 403) || (resp.StatusCode == 401) {
		if client.Ping() {
			resp, err = client.Do(method, path, body)
		} else {
			err = fmt.Errorf("Could not Ping after 401/403", err)
			return nil, err
		}

	}
	if err != nil {
		client.Log.Printf("Error during HTTP Call to OMI: %v", err)
		return nil, err
	}
	return resp, err
}

func (client *Client) Get(path string) (*http.Response, error) {
	return client.Do("GET", path, nil)
}

func (client *Client) Post(path string, body []byte) (*http.Response, error) {
	return client.Do("POST", path, body)
}
func (client *Client) Put(path string, body []byte) (*http.Response, error) {
	return client.Do("PUT", path, body)
}

func (client *Client) GetEventList() (*EventList, error) {
	var el EventList
	resp, err := client.Get(EventListPath)
	data, err := ioutil.ReadAll(resp.Body)

	err = xml.Unmarshal(data, &el)
	if err != nil {
		return nil, err
	}
	el.Client = client
	for _, event := range el.Event {
		event.Client = client
	}
	return &el, err
}

func (client *Client) QueryEventList(query string) (*EventList, error) {
	var el EventList
	encodedQuery := &url.URL{Path: query}
	queryPath := fmt.Sprintf("%v?query=%v", EventListPath, encodedQuery.String())
	resp, err := client.Get(queryPath)
	data, err := ioutil.ReadAll(resp.Body)
	err = xml.Unmarshal(data, &el)
	if err != nil {
		return nil, err
	}
	el.Client = client
	for _, event := range el.Event {
		event.Client = client
	}
	return &el, err
}

func (client *Client) NewEvent() (*Event, error) {
	var event Event
	if err := mergo.Merge(&event, DefaultEvent); err != nil {
		return nil, fmt.Errorf("Could not get default event")
	}
	return &event, nil
}

func (client *Client) CreateEvent(event *Event) error {
	//event.Client = client
	result, err := xml.Marshal(event)
	if err != nil {
		return fmt.Errorf("Could not marshal Event %v", err)
	}
	resp, err := client.Post(EventListPath, result)
	if err != nil {
		return err
	}
	respByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = xml.Unmarshal(respByte, &event)
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) UpdateEvent(event *Event) error {
	var updateEvent = event
	event.Xmlns = EventXmlns
	updateEvent.clean()
	if updateEvent.Id == "" {
		return fmt.Errorf("Cannot UpdateEvent without Id")
	}
	result, err := xml.Marshal(updateEvent)
	if err != nil {
		return fmt.Errorf("Could not marshal Event %v", err)
	}
	eventPath := fmt.Sprintf("%v/%v", EventListPath, updateEvent.Id)
	resp, err := client.Put(eventPath, result)
	if err != nil {
		return err
	}
	respByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {

		return err
	}
	err = xml.Unmarshal(respByte, &event)
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) AddAnnotation(event *Event, annotation *Annotation) error {
	annotation.Xmlns = EventXmlns
	if event.Id == "" {
		return fmt.Errorf("Cannot AddAnnotation without Id")
	}
	annotationXml, err := xml.Marshal(annotation)
	if err != nil {
		return fmt.Errorf("Could not marshal Event %v", err)
	}
	requestPath := fmt.Sprintf("%v/%v/annotation_list", EventListPath, event.Id)
	resp, err := client.Post(requestPath, annotationXml)
	if err != nil {
		return err
	}
	respByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = xml.Unmarshal(respByte, &annotation)
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) AddCustomAttribute(event *Event, ca *CustomAttribute) error {
	ca.Xmlns = EventXmlns
	if event.Id == "" {
		return fmt.Errorf("Cannot AddCustomAttribute without Id")
	}
	caXml, err := xml.Marshal(ca)
	if err != nil {
		return fmt.Errorf("Could not marshal Event %v", err)
	}
	requestPath := fmt.Sprintf("%v/%v/custom_attribute_list", EventListPath, event.Id)
	resp, err := client.Post(requestPath, caXml)
	if err != nil {
		return err
	}
	respByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = xml.Unmarshal(respByte, &ca)
	if err != nil {
		return err
	}
	return nil
}
