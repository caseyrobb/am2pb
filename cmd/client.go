package main

import (
	"log"
	"os"

	am2pb "github.com/caseyrobb/am2pb"
)

func main() {
	bearerToken := os.Getenv("BEARER_TOKEN")
	if bearerToken == "" {
		log.Fatal("BEARER_TOKEN environment variable not set")
	}

	am2pb.StartServer(bearerToken)
}
