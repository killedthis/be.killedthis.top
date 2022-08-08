package main

import "github.com/segmentio/go-route53"
import "github.com/mitchellh/goamz/aws"
import "encoding/json"
import "fmt"
import "os"

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	auth, err := aws.EnvAuth()
	check(err)

	fmt.Println(os.Args[1])

	dns := route53.New(auth, aws.EUCentral)

	res, err := dns.Zone("Z07647923LTTIH8HQ17W6").Add("CNAME", os.Args[1]+".killedthis.top", "killedthis.top")
	check(err)

	b, err := json.MarshalIndent(res, "", "  ")
	check(err)

	os.Stdout.Write(b)
}
