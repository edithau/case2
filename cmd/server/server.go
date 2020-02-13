package main

import (
	"log"

	"github.com/InVisionApp/interview-test/api"
)

func main() {
	log.Fatal(api.New("").Run())
}
