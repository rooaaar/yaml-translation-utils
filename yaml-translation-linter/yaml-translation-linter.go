package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"strings"

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

	if err := checkKeys(ref, tra, ""); err != nil {
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

func checkKeys(ref interface{}, tra interface{}, parentPath string) error {
	if _, ok := ref.(string); ok {
		if _, ok := tra.(string); ok {
			return nil
		}
	}

	refMap, ok := ref.(map[interface{}]interface{})

	if ok {
		if traMap, ok := tra.(map[interface{}]interface{}); ok {
			var errs []string
			for key, val := range refMap {
				path := parentPath + "." + key.(string)
				if _, ok := traMap[key]; !ok {
					errs = append(errs, "cannot find '"+path+"' in translations")
				} else if err := checkKeys(val, traMap[key], path); err != nil {
					errs = append(errs, err.Error())
				}
			}

			for key := range traMap {
				if _, ok := refMap[key]; !ok {
					path := parentPath + "." + key.(string)
					errs = append(errs, "found redunant translation '"+path+"'")
				}
			}

			if len(errs) > 0 {
				return errors.New(strings.Join(errs, "\n"))
			}

			return nil
		}

		return errors.New("cannot convert translation to map: " + parentPath)
	}

	return errors.New("cannot convert reference to map: " + parentPath)
}
