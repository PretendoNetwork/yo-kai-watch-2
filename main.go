package main

import (
	"sync"

	"github.com/PretendoNetwork/yo-kai-watch-2/nex"
)

var wg sync.WaitGroup

func main() {
	wg.Add(2)

	// TODO - Add gRPC server
	go nex.StartAuthenticationServer()
	go nex.StartSecureServer()

	wg.Wait()
}
