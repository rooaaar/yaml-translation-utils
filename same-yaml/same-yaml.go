package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/sijad/yaml-translation-utils/yamlutils"
	"gopkg.in/yaml.v2"
)

var refPath string
var traPath string

func init() {
	flag.StringVar(&refPath, "ref", "", "reference file path")
	flag.StringVar(&traPath, "tra", "", "translation file path")
}

func main() {
	flag.Parse()

	if refPath == "" {
		log.Fatal("reference file not defined.")
	}

	if traPath == "" {
		log.Fatal("translation file not defined.")
	}

	ref := readYamlFile(refPath)
	tra := readYamlFile(traPath)

	if err := yamlutils.IdenticalKeys(ref, tra, ""); err != nil {
		log.Fatal(err)
	}
}

func readYamlFile(file string) interface{} {
	bytes, err := ioutil.ReadFile(file)

	if err != nil {
		log.Fatal(err)
	}

	var data interface{}

	if err := yaml.Unmarshal(bytes, &data); err != nil {
		log.Fatal(err)
	}

	return data
}
