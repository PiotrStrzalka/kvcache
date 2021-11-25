package main

import (
	"io/ioutil"
	"log"

	"sigs.k8s.io/yaml"
)

var input = "openapi.yaml"

func main() {
	y, err := ioutil.ReadFile(input)
	if err != nil {
		log.Fatalf("Cannot open %s\n, error occured %w", input, err)
	}

	j, err := yaml.YAMLToJSON(y)
	if err != nil {
		log.Fatalf("Problem during conversion %w", err)
	}

	ioutil.WriteFile("swagger.json", j, 644)
}
