// Entry point, only wiring, a Go app only has one main.go

package main

import (
	"log"

	"github.com/maxdeekay/nook/internal/server"
)

// stdlib (standard library), default libs
// := is a short for inferring type from the value
// fmt.Println writes to stdout, no prefix, followed by a newline. General-purpose "print this value" function. Space-separated args
// log.Println writes to stderr with an automatic timestamp prefix, intended for diagnostic/logging output
// Print (no newline)
// Printf (formatted with %s, %d etc)

func main() {
	srv := server.New(server.Config{Port: 8080})

	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
