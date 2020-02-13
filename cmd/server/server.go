package main

import (
	"log"

	"github.com/InVisionApp/case-study/api"
)

func main() {
	log.Fatal(api.New("").Run())
}
