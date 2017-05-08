package main

import (
	"flag"
	"fmt"
	"github.com/dimitertodorov/gomi/omi"
	"github.com/prometheus/common/version"
	"io/ioutil"
	"os"
	"strings"
	log "github.com/Sirupsen/logrus"
)

var (
	hostState2severity = map[string]string{
		"UP":          "normal",
		"DOWN":        "major",
		"UNREACHABLE": "major",
	}

	serviceState2severity = map[string]string{
		"OK":       "normal",
		"WARNING":  "warning",
		"UNKNOWN":  "minor",
		"CRITICAL": "critical",
	}
)

func init() {

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

//Forward events from Nagios to OMI
func main() {
	var (
		showVersion = flag.Bool("version", false, "Print version information.")
		configFile  = flag.String("config.file", "gomi.json", "Gomi configuration file name.")
		//Try these domains for OMI resolver.
		dnsDomains = flag.String("dns-domains", "cihs.ad.gov.on.ca,service.cihs.gov.on.ca,mgmt.cihs.gov.on.ca", "DNS Domains to try for CI Resolution")

		notificationType = flag.String("notification-type", "host", "host or service")
		alertType        = flag.String("type", "PROBLEM", "")

		contact      = flag.String("contact", "", "")
		contactEmail = flag.String("contactemail", "", "")

		author   = flag.String("author", "nagios", "")
		comments = flag.String("comments", "", "")

		escalated = flag.Int("escalated", 0, "")

		hostname        = flag.String("host", "", "")
		hostAddress     = flag.String("hostaddress", "", "")
		hostAlias       = flag.String("hostalias", "", "")
		hostDisplayName = flag.String("hostdisplayname", "", "")

		hostState     = flag.String("hoststate", "", "")
		hostStateId   = flag.Int("hoststateid", 0, "")
		hostStateType = flag.String("hoststatetype", "", "")

		lastHostState   = flag.String("lasthoststate", "", "")
		lastHostStateId = flag.Int("lasthoststateid", 0, "")

		hostEventId       = flag.String("hosteventid", "", "")
		hostProblemId     = flag.Int("hostproblemid", 0, "")
		lastHostProblemId = flag.Int("lasthostproblemid", 0, "")

		hostOutput     = flag.String("hostoutput", "", "")
		longHostOutput = flag.String("longhostoutput", "", "")

		dateTime = flag.String("datetime", "", "")

		service            = flag.String("service", "", "")
		serviceState       = flag.String("servicestate", "", "")
		serviceStateId     = flag.String("servicestateid", "", "")
		lastServiceState   = flag.String("lastservicestate", "", "")
		lastServiceStateId = flag.String("lastservicestateid", "", "")

		currentAttempt = flag.Int("currentattempt", 0, "")
		maxAttempts    = flag.Int("maxattempts", 0, "")

		createFlag = flag.Bool("create", false, "Create or Not Create")
	)
	flag.Parse()

	configContents, err := ioutil.ReadFile(*configFile)
	if err != nil {
		panic(fmt.Sprintf("Could not read config file %v", err))
	}
	client := omi.NewClient(configContents)

	if *showVersion {
		fmt.Fprintln(os.Stdout, version.Print("nagomi"))
		os.Exit(0)
	}

	*longHostOutput = strings.Replace(*longHostOutput, `\n`, "\n", -1)

	var hints []string
	nameParts := strings.Split(*hostname, ".")
	shortName := nameParts[0]
	hints = append(hints, fmt.Sprintf("@@%v", shortName))
	domains := strings.Split(*dnsDomains, ",")
	for _, domain := range domains {
		hints = append(hints, fmt.Sprintf("@@%v.%v", shortName, domain))
	}
	joinedHints := strings.Join(hints, "|")
	ciHints := omi.RelatedCiHints{Hint: []string{joinedHints}}

	newEvent := omi.Event{
		Xmlns:           omi.EventXmlns,
		Title:           fmt.Sprintf("NAGIOS: %v", *hostOutput),
		Description:     fmt.Sprintf("NAGIOS: %v", *longHostOutput),
		Key:             fmt.Sprintf("NAGIOS:%v:%v:%v", *hostname, *alertType, *hostProblemId),
		CloseKeyPattern: fmt.Sprintf("NAGIOS:%v:%v:%v", *hostname, *alertType, *lastHostProblemId),
		State:           "open",
		Severity:        hostState2severity[*hostState],
		RelatedCiHints:  &ciHints,
		Application:     "Nagios",
		Object:          "Nagios Object",
	}

	if *createFlag {
		if err := client.CreateEvent(&newEvent); err != nil {
			panic(fmt.Sprintf("Error: %s", err))
		}
		log.Info(fmt.Sprintf("Created Event %v", newEvent.Id))
	} else {
		log.Info("Not Creating Event")
	}

	log.Debug(fmt.Sprintf("%v", *notificationType))
	log.Debug(fmt.Sprintf("%v", *alertType))
	log.Debug(fmt.Sprintf("%v", *dnsDomains))
	log.Debug(fmt.Sprintf("%v", *contact))
	log.Debug(fmt.Sprintf("%v", *contactEmail))
	log.Debug(fmt.Sprintf("%v", *escalated))
	log.Debug(fmt.Sprintf("%v", *hostname))
	log.Debug(fmt.Sprintf("%v", *hostAddress))
	log.Debug(fmt.Sprintf("%v", *hostAlias))
	log.Debug(fmt.Sprintf("%v", *hostDisplayName))
	log.Debug(fmt.Sprintf("%v", *configFile))
	log.Debug(fmt.Sprintf("%v", *hostState))
	log.Debug(fmt.Sprintf("%v", *hostStateId))
	log.Debug(fmt.Sprintf("%v", hostState2severity[*hostState]))
	log.Debug(fmt.Sprintf("%v", *hostStateType))
	log.Debug(fmt.Sprintf("%v", *lastHostState))
	log.Debug(fmt.Sprintf("%v", *lastHostStateId))
	log.Debug(fmt.Sprintf("%v", *hostEventId))
	log.Debug(fmt.Sprintf("%v", *hostProblemId))
	log.Debug(fmt.Sprintf("%v", *lastHostProblemId))

	log.Debug(fmt.Sprintf("%v", *hostOutput))
	log.Debug(fmt.Sprintf("%v", *longHostOutput))
	log.Debug(fmt.Sprintf("%v", *dateTime))
	log.Debug(fmt.Sprintf("%v", *service))
	log.Debug(fmt.Sprintf("%v", *serviceState))

	log.Debug(fmt.Sprintf("%v", *serviceStateId))
	log.Debug(fmt.Sprintf("%v", *lastServiceState))
	log.Debug(fmt.Sprintf("%v", serviceState2severity[*lastServiceState]))
	log.Debug(fmt.Sprintf("%v", *lastServiceStateId))
	log.Debug(fmt.Sprintf("%v", *service))
	log.Debug(fmt.Sprintf("%v", *serviceState))
	log.Debug(fmt.Sprintf("%v", *currentAttempt))
	log.Debug(fmt.Sprintf("%v", *maxAttempts))

	log.Debug(fmt.Sprintf("%v", *author))
	log.Debug(fmt.Sprintf("%v", *comments))

}
