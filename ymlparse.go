package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"gopkg.in/yaml.v2"
)

var command = "fish.and.chips"

func fileExists(filename string) bool {
	prefix := "./"
	prefixBuffer := make([]byte, len(filename)+2)
	prefixBufferPointer := copy(prefixBuffer, prefix)
	copy(prefixBuffer[prefixBufferPointer:], filename)

	prefixed := string(prefixBuffer)

	_, err := os.Stat(prefixed)
	// fmt.Println(prefixed, "?", err == nil)
	return err == nil
}

func main() {

	// parse arguments
	defaultInputFilename, defaultOverrideFilename := "docker-compose.yml", "docker-compose.override.yml"
	var inputFilename, overrideFilename string
	flag.StringVar(&inputFilename, "i", defaultInputFilename, "the input file to read from")
	flag.StringVar(&overrideFilename, "override", defaultOverrideFilename, "the override file to read from; will attempt to find the '.override.yml' of the input file if omitted")

	flag.Parse()
	overflow := flag.Args()
	if len(overflow) == 0 {
		log.Fatalf("error: must supply arguments to drill down through the yml file")
	}
	bits := overflow

	// handle automatically generating the '.override.yml'
	if overrideFilename == defaultOverrideFilename && inputFilename != defaultInputFilename {
		if !fileExists(overrideFilename) {
			replaceYmlSuffix := regexp.MustCompile("\\.yml$")
			newOverrideFilename := replaceYmlSuffix.ReplaceAllString(inputFilename, ".override.yml")
			if fileExists(newOverrideFilename) {
				overrideFilename = newOverrideFilename
			}
		}
	}

	// Load files in
	data, err := ioutil.ReadFile(inputFilename)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	var dataoverride []byte
	if _, existsErr := os.Stat(overrideFilename); existsErr == nil {
		dataoverride, err = ioutil.ReadFile(overrideFilename)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
	}

	// Unmarshal
	var main interface{}
	mainFile := make(map[interface{}]interface{})
	err = yaml.Unmarshal([]byte(data), &mainFile)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	main = mainFile

	var override interface{}
	overrideFile := make(map[interface{}]interface{})
	err = yaml.Unmarshal([]byte(dataoverride), &overrideFile)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	override = overrideFile

	// Drill down with the command string
	for _, element := range bits {
		if main != nil {
			if val, ok := main.(map[interface{}]interface{})[element]; ok {
				main = val
			} else {
				main = nil
			}
		}

		if override != nil {
			if val, ok := override.(map[interface{}]interface{})[element]; ok {
				override = val
			} else {
				override = nil
			}
		}
	}

	// Output the result, or die with exit 1
	if override != nil {
		fmt.Print(override)
	} else if main != nil {
		fmt.Print(main)
	} else {
		os.Exit(1)
	}

}
