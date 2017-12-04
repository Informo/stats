package main

import (
	"flag"
	"fmt"
	"strconv"

	"informo-statisticator/entrypoints"
)

var (
	homeserver  = flag.String("homeserver", "127.0.0.1", "The homeserver's FQDN")
	port        = flag.Int("port", 8448, "The homeserver's port")
	accessToken = flag.String("access-token", "", "The Matrix access token to use to connect to the network")
	noTLS       = flag.Bool("no-tls", false, "If set to true, traffic will be sent with no TLS (plain HTTP)")
)

func main() {
	flag.Parse()

	if *accessToken == "" {
		panic("No access token provided")
	}

	homeserverURL := "http"
	if !*noTLS {
		homeserverURL = homeserverURL + "s"
	}
	homeserverURL = homeserverURL + "://" + *homeserver + ":" + strconv.Itoa(*port)

	ep, err := entrypoints.GetEntryPoints(homeserverURL, *accessToken)
	if err != nil {
		panic(err)
	}

	count := entrypoints.CountEntryNodes(ep)

	fmt.Printf("Counted %d entrypoints on %d nodes\n", len(ep), count)
}
