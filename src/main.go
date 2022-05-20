package main

import (
	"fmt"
	"www/src/router"
	"log"
    "github.com/PaloAltoNetworks/pango"
	"github.com/go-ping/ping"
)

func main() {
	var err error

    c := &pango.Firewall{Client: pango.Client{
        Hostname: "extrapass.cdg.co.th",
        Username: "003459",
        Password: "cdgs$8133Tongx",
        Logging: pango.LogAction | pango.LogOp,
    }}
    if err = c.Initialize(); err != nil {
        log.Printf("Failed to initialize client: %s", err)
        return
    }
    log.Printf("Initialize ok")
	p, err := ping.NewPinger("10.255.1.150")
	if err != nil {
		log.Printf("Cannot connect to 10.254.1.177: %s", err)
	} else{
		log.Printf("Ping complete... %v",p)
	}
	


	fmt.Println("Wellcome to Golang")
	fmt.Println("Start to Programming golang")
	rt := router.Init()

	rt.Logger.Fatal(rt.Start(":8080"))
}
