package main

import (
	"log"

	yaml "gopkg.in/yaml.v2"
)

func unmarshal(input []byte) interface{} {
	var main interface{}

	mainFile := make(map[interface{}]interface{})
	err := yaml.Unmarshal(input, &mainFile)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	main = mainFile
	return main
}
