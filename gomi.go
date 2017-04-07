package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"github.com/dimitertodorov/gomi/omi"
	"github.com/prometheus/common/version"
)

func main() {

	var (
		showVersion = flag.Bool("version", false, "Print version information.")
		configFile = flag.String("config.file", "gomi.json", "Gomi configuration file name.")
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


	configContents, err := ioutil.ReadFile(*configFile)
	if err != nil {
		panic(fmt.Sprintf("Could not read config file %v", err))
	}
	client := omi.NewClient(configContents)
	if eventList, err := client.GetEventList(); err != nil {
		panic(fmt.Errorf("Cannot Get Events %v", err))
	}else{
		for i, event := range eventList.Event {
			fmt.Printf("[%v] Got Event %v - %v\n", i, event.Title, event.Id)
		}
	}






}
