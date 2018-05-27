package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type yamlErrorCode string

const (
	notSorted                    yamlErrorCode = "should be sorted alphabetically"
	notLower                                   = "key should be in lowercase"
	notAlphaNumericDashUnderline               = "key should only contain letters, numbers, dash and underline"
)

type yamlError struct {
	Err  yamlErrorCode
	Path []string
}

var filePath string
var alphaBeticLevel int

func init() {
	flag.StringVar(&filePath, "file", "", "yaml file path")
	flag.IntVar(&alphaBeticLevel, "level", 0, "minimal level to check alphabetically sorted")
}

func main() {
	flag.Parse()

	if filePath == "" {
		log.Fatal("yaml file not defined.")
	}

	logger := log.New(os.Stderr, "", 0)
	errs := lint(readYamlFile(filePath))

	if errs != nil {
		for _, err := range errs {
			path := strings.Join(err.Path, ".")
			logger.Printf("%s: %s", path, err.Err)
		}
		log.Fatalln("linter faild.")
	}
}

func lint(node yaml.MapSlice) []yamlError {
	return getNodeErrors(node, nil, 0)
}

func getNodeErrors(node yaml.MapSlice, path []string, level int) (errs []yamlError) {
	var biggestKey string
	var lastKey string

	for _, kv := range node {

		key := kv.Key.(string)

		if level >= alphaBeticLevel {
			if key > biggestKey {
				biggestKey = key
			}

			if key < lastKey || key < biggestKey {
				errs = append(errs, newError(notSorted, append(path, key)))
			}

			lastKey = key
		}

		for _, err := range getKeyErrors(key) {
			errs = append(errs, newError(err, append(path, key)))
		}

		// check inner keys if any.
		if nextNode, ok := kv.Value.(yaml.MapSlice); ok {
			if nextErrs := getNodeErrors(nextNode, append(path, key), level+1); nextErrs != nil {
				errs = append(errs, nextErrs...)
			}
		}
	}

	return
}

func getKeyErrors(key string) (errs []yamlErrorCode) {
	if strings.ToLower(key) != key {
		errs = append(errs, notLower)
	}

	if !isAlphaUnderline(key) {
		errs = append(errs, notAlphaNumericDashUnderline)
	}

	return
}

func isAlphaUnderline(s string) bool {
	for _, r := range s {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') && r != '_' && r != '-' {
			return false
		}
	}
	return true
}

func newError(err yamlErrorCode, path []string) yamlError {
	return yamlError{
		err,
		path,
	}
}

func readYamlFile(file string) yaml.MapSlice {
	bytes, err := ioutil.ReadFile(file)

	if err != nil {
		log.Fatal(err)
	}

	var data yaml.MapSlice

	if err := yaml.Unmarshal(bytes, &data); err != nil {
		log.Fatal(err)
	}

	return data
}
