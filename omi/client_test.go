package omi

import (
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"io/ioutil"
	"net/http"
	"regexp"
	"testing"
	"time"
)

var (
	//TURN On Real TEST Here
	mockOnly = false

	configFile = "../testdata/gomi.json"

	secureModifyCookie = &http.Cookie{
		Name:   "secureModifyToken",
		Value:  "5bcaf756371567ba91a9d8195770e5739afda223",
		Path:   "/",
		Domain: "dev-omi.tools.cihs.gov.on.ca",
	}
	secureModifyTokenMatchFunc = func(req *http.Request, ereq *gock.Request) (bool, error) {
		header := req.Header.Get("X-Secure-Modify-Token")
		if header == "" {
			return false, nil
		} else {
			return true, nil
		}
	}
	missingSecureModifyTokenMatchFunc = func(req *http.Request, ereq *gock.Request) (bool, error) {
		header := req.Header.Get("X-Secure-Modify-Token")
		if header == "" {
			return true, nil
		} else {
			return false, nil
		}
	}

	//UUID
	runUuid = uuid.NewV4()

	//RelateCi
	//RelatedCiHints
	ciHints = RelatedCiHints{Hint: []string{"@@CTSbiGdcEMadm34"}}

	newEvent = Event{
		Xmlns:           EventXmlns,
		Title:           fmt.Sprintf("GOMI DEMO EVENT %v", runUuid.String()),
		Key:             fmt.Sprintf("GOLANG OMI EVENT:%s:major", runUuid.String()),
		State:           "open",
		Severity:        "major",
		CloseKeyPattern: fmt.Sprintf("GOLANG OMI EVENT:%s:<*>normal", runUuid.String()),
		Application:     "TORCA_GOMI",
		Object:          "test_object",
		RelatedCiHints:  &ciHints,
	}

	//Matchers
	secureModifyTokenMatcher        = gock.NewBasicMatcher()
	missingSecureModifyTokenMatcher = gock.NewBasicMatcher()

	xmlMap = make(map[string][]byte)

	configContents []byte
)

func init() {
	secureModifyTokenMatcher.Add(secureModifyTokenMatchFunc)
	missingSecureModifyTokenMatcher.Add(missingSecureModifyTokenMatchFunc)
	files, _ := ioutil.ReadDir("../testdata")
	r, _ := regexp.Compile("(.*)\\.(xml)")
	for _, f := range files {
		if r.MatchString(f.Name()) {
			fileXml, _ := ioutil.ReadFile(fmt.Sprintf("../testdata/%v", f.Name()))
			xmlMap[r.FindStringSubmatch(f.Name())[1]] = fileXml
		}
	}
	//Load Sample config
	configContents, _ = ioutil.ReadFile(configFile)

}

func TestClientConfig(t *testing.T) {
	var config ClientConfig
	byt := []byte(`{
  "username": "admin",
  "password": "admin",
  "base_url": "https://dev-omi.tools.cihs.gov.on.ca"
}`)
	if err := json.Unmarshal(byt, &config); err != nil {
		t.Fatalf("%v", err)
	}
	assert.Equal(t, "admin", config.Username, "admin")
}

func TestNewClient(t *testing.T) {
	client := NewClient(configContents)
	assert.Equal(t, "admin", client.ClientConfig.Username, "Client should initialize properly")
}

func TestGet(t *testing.T) {
	client := NewClient(configContents)
	if mockOnly {
		defer gock.Off()
		mockSuccess(client, t, "GET", EventListPath, "event_list", 200)
	}

	if events, err := client.GetEventList(); err != nil {
		t.Fatalf("Error getting events %v", err)
	} else {
		assert.NotEqual(t, "", events.Event[0].Id, "get event list")
	}

}

func TestGetUnauth(t *testing.T) {
	client := NewClient(configContents)

	if mockOnly {
		defer gock.Off() // Flush pending mocks after test execution
		mock403(client, t, EventListPath, "GET")
		mockPing(client, t)
		mockSuccess(client, t, "GET", EventListPath, "event_list", 200)
	}

	_, err := client.GetEventList()

	if err != nil {
		t.Fatalf("Error getting events %v", err)
	}
}

func TestNewEvent(t *testing.T) {
	configContents, err := ioutil.ReadFile(configFile)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	client := NewClient(configContents)
	event, err := client.NewEvent()
	assert.Equal(t, event.Severity, "")
}

func TestCreateEvent(t *testing.T) {
	var event = newEvent
	configContents, err := ioutil.ReadFile(configFile)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	client := NewClient(configContents)

	if mockOnly {
		defer gock.Off()
		mock403(client, t, EventListPath, "POST")
		mockPing(client, t)
		mockSuccess(client, t, "POST", EventListPath, "event", 202)
	}

	if err := client.CreateEvent(&event); err != nil {
		t.Fatalf("Error: %s", err)
	}

	assert.NotEqual(t, "", event.Id)
}

func TestUpdate(t *testing.T) {
	client := NewClient(configContents)
	var event Event
	if mockOnly {
		defer gock.Off()

		mockSuccess(client, t, "GET", EventListPath, "event_list", 200)
		mock403(client, t, EventListPath, "PUT")
		mockPing(client, t)
		mockSuccess(client, t, "PUT", EventListPath, "event", 201)
	} else {
		time.Sleep(time.Second * 6)
	}

	query := fmt.Sprintf("key LIKE \"$25%v$25\"", runUuid.String())
	eventList, err := client.QueryEventList(query)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}

	event = eventList.Event[0]

	event.Description = "TESTDESC"

	if err := client.UpdateEvent(&event); err != nil {
		t.Fatalf("Error: %s", err)
	}

	assert.Equal(t, "TESTDESC", event.Description)
	assert.NotEqual(t, "", event.Id)
}

func TestCreateAddAnnotation(t *testing.T) {
	client := NewClient(configContents)
	var event Event
	if mockOnly {
		defer gock.Off()

		mockSuccess(client, t, "GET", EventListPath, "event_list", 200)
		mock403(client, t, EventListPath, "POST")
		mockPing(client, t)
		mockSuccess(client, t, "POST", fmt.Sprintf("%v/(.*)/annotation_list", EventListPath), "annotation", 201)
	}
	query := fmt.Sprintf("key LIKE \"$25%v$25\"", runUuid.String())
	eventList, err := client.QueryEventList(query)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	event = eventList.Event[0]

	annotation := DefaultAnnotation
	annotation.Text = "From GOLAN Text"
	annotation.Author = "Dimiter"

	if err := client.AddAnnotation(&event, &annotation); err != nil {
		t.Fatalf("Error: %s", err)
	}

	assert.NotEqual(t, "", annotation.Id)
}

func TestCreateAddCustomAttribute(t *testing.T) {
	var event Event
	client := NewClient(configContents)
	if mockOnly {
		defer gock.Off()
		mockSuccess(client, t, "GET", EventListPath, "event_list", 200)
		mock403(client, t, EventListPath, "POST")
		mockPing(client, t)
		mockSuccess(client, t, "POST", fmt.Sprintf("%v/(.*)/custom_attribute_list", EventListPath), "custom_attribute", 201)
	}

	query := fmt.Sprintf("key LIKE \"$25%v$25\"", runUuid.String())
	eventList, err := client.QueryEventList(query)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	event = eventList.Event[0]

	ca := new(CustomAttribute)
	ca.Name = "GOLANG_RUN"
	ca.Value = runUuid.String()

	if err := client.AddCustomAttribute(&event, ca); err != nil {
		t.Fatalf("Error: %s", err)
	}

	assert.Equal(t, "GOLANG_RUN", ca.Name)
}

func mockSuccess(client *Client, t *testing.T, method, path, source string, code int) error {
	if method == "GET" {
		gock.New(client.ClientConfig.BaseUrl).
			Get(path).
			Reply(code).
			BodyString(string(xmlMap[source][:]))
	} else if method == "POST" {
		gock.New(client.ClientConfig.BaseUrl).
			SetMatcher(secureModifyTokenMatcher).
			Post(path).
			Reply(code).
			BodyString(string(xmlMap[source][:]))
	} else if method == "PUT" {
		gock.New(client.ClientConfig.BaseUrl).
			SetMatcher(secureModifyTokenMatcher).
			Put(path).
			Reply(code).
			BodyString(string(xmlMap[source][:]))
	}
	return nil
}

func mockPing(client *Client, t *testing.T) {
	gock.New(client.ClientConfig.BaseUrl).
		Get("/opr-web/rest").
		Reply(200).
		BodyString("200").
		SetHeader("Set-Cookie", secureModifyCookie.String())
}

func mock403(client *Client, t *testing.T, path, method string) {
	if method == "GET" {
		gock.New(client.ClientConfig.BaseUrl).
			Get(path).
			Reply(403).
			BodyString("UNAUTHORIZED")
	} else if method == "POST" {
		gock.New(client.ClientConfig.BaseUrl).
			Post(path).
			Reply(403).
			BodyString("UNAUTHORIZED")
	} else if method == "PUT" {
		gock.New(client.ClientConfig.BaseUrl).
			Put(path).
			Reply(403).
			BodyString("UNAUTHORIZED")
	}

}
