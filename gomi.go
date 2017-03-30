package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"net/http/cookiejar"
	"net/http"
	"net/url"
	"github.com/dimitertodorov/gomi/omi"
	"encoding/xml"
	"github.com/prometheus/common/version"
)

func main() {

	var (
		showVersion = flag.Bool("version", false, "Print version information.")
		//configFile = flag.String("config.file", "alertmanager.yml", "Alertmanager configuration file name.")
		//dataDir    = flag.String("storage.path", "data/", "Base path for data storage.")
		//omiUrl     = flag.String("omi.url", "", "mesh peer ID (default: MAC address)")
		//omiUsername   = flag.String("omi.username", "admin", "mesh peer nickname")
		//omiPassword   = flag.String("omi.password", "admin", "password to join the peer network (empty password disables encryption)")
	)
	flag.Parse()

	if *showVersion {
		fmt.Fprintln(os.Stdout, version.Print("gomi"))
		os.Exit(0)
	}

	cookieJar, _ := cookiejar.New(nil)

	client := &http.Client{
		Jar: cookieJar,
	}

	req, err := http.NewRequest("GET", "https://dev-omi.tools.cihs.gov.on.ca/opr-web/rest/9.10/event_list", nil)
	req.SetBasicAuth("admin","admin")

	resp, err := client.Do(req)

	_, err = ioutil.ReadAll(resp.Body)

	// error handle
	if err != nil {
		fmt.Printf("error = %s \n", err);
	}

	// Print response
	omiPath, err := url.Parse("https://dev-omi.tools.cihs.gov.on.ca")
	cooks := client.Jar.Cookies(omiPath)

	for _, v := range cooks {
		fmt.Printf("COOKIE: %s = %s\n\n", v.Name, v.Value);
	}

	resp, err = client.Get("https://dev-omi.tools.cihs.gov.on.ca/opr-web/rest/9.10/event_list")


	data, err := ioutil.ReadAll(resp.Body)

	fmt.Printf("%s ", data)

	var el omi.EventList
	if err := xml.Unmarshal(data, &el); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(el)




}
