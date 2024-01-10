package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/mitchellh/goamz/aws"
	"github.com/segmentio/go-route53"
)

func check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

type DNSRECORD []struct {
	Name    string   `json:"Name"`
	Type    string   `json:"Type"`
	Records []string `json:"Records"`
}

var dnsrecord DNSRECORD

func main() {
	auth, err := aws.EnvAuth()
	check(err)

	log.Printf("Checking for: %s\n", os.Args[1]+".killedthis.top")

	dns := route53.New(auth, aws.EUCentral)
	check(err)

	resr, err := dns.Zone("Z07647923LTTIH8HQ17W6").RecordsByName(os.Args[1] + ".killedthis.top")
	check(err)

	re, err := json.MarshalIndent(resr, "", "  ")
	check(err)

	err = json.Unmarshal([]byte(re), &dnsrecord)
	if err != nil {
		log.Panic(err)
	}

	if len(dnsrecord) < 1 {
		fmt.Printf("Creating: %s\n", os.Args[1]+".killedthis.top")
		resc, err := dns.Zone("Z07647923LTTIH8HQ17W6").Add("CNAME", os.Args[1]+".killedthis.top", "killedthis.top")
		check(err)

		cr, err := json.MarshalIndent(resc, "", "  ")
		check(err)

		os.Stdout.Write(cr)
	} else {
		log.Printf("Found: %s\n", dnsrecord[0].Name)
	}
}
