package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/dimitertodorov/gomi/omi"
	"github.com/prometheus/common/version"
	"io/ioutil"
	"os"
)

func main() {

	var (
		showVersion = flag.Bool("version", false, "Print version information.")
		configFile  = flag.String("config.file", "gomi.json", "Gomi configuration file name.")
	)
	flag.Parse()

	if *showVersion {
		fmt.Fprintln(os.Stdout, version.Print("gomi"))
		os.Exit(0)
	}

	//t := template.New("Person template")

	configContents, err := ioutil.ReadFile(*configFile)
	if err != nil {
		panic(fmt.Sprintf("Could not read config file %v", err))
	}
	client := omi.NewClient(configContents)
	if eventList, err := client.GetEventList(); err != nil {
		panic(fmt.Errorf("Cannot Get Events %v", err))
	} else {
		//t, err := t.Parse(tmpl)
		//if err != nil {
		//	fmt.Printf("Err %v", err)
		//	return
		//}
		//
		//err = t.Execute(os.Stdout, eventList)
		//if err != nil {
		//	fmt.Printf("Err %v", err)
		//	return
		//}
		for _, event := range eventList.Event {
			st, _ := json.MarshalIndent(event, "", "\t")
			fmt.Printf("Got Event %v\n", string(st[:]))
		}
	}

}
