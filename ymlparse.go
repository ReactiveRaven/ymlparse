package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

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
	if overrideFilename == defaultOverrideFilename && inputFilename != defaultInputFilename && !fileExists(overrideFilename) {
		replaceYmlSuffix := regexp.MustCompile("\\.yml$")
		newOverrideFilename := replaceYmlSuffix.ReplaceAllString(inputFilename, ".override.yml")
		overrideFilename = newOverrideFilename
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
	main := unmarshal(data)
	override := unmarshal(dataoverride)

	// Drill down with the command string
	for _, element := range bits {
		main = drilldownObject(main, element)
		override = drilldownObject(override, element)
	}

	// Output the result, or die with exit 1
	if override != nil {
		marshalled, err := yaml.Marshal(override)
		if err != nil {
			log.Fatalf("error: failed to marshal response: %s", err)
		}
		fmt.Print(string(marshalled))
	} else if main != nil {
		marshalled, err := yaml.Marshal(main)
		if err != nil {
			log.Fatalf("error: failed to marshal response: %s", err)
		}
		fmt.Print(string(marshalled))
	} else {
		log.Fatalf(
			"error: could not find path '%s' in either '%s' or '%s'",
			strings.Join(bits, "','"),
			inputFilename,
			overrideFilename)
	}

}
