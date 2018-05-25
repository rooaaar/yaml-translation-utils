package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

const transifexAPI = "https://www.transifex.com/api/2/project/"

var apiKey string
var projectSlug string
var lang string
var outDir string

func init() {
	flag.StringVar(&apiKey, "api", "", "Transifex API key")
	flag.StringVar(&projectSlug, "project-slug", "", "Transifex API key")
	flag.StringVar(&lang, "lang", "", "Language code")
	flag.StringVar(&outDir, "out", "", "output directory to save resouces")
}

func main() {
	flag.Parse()

	if apiKey == "" {
		log.Fatal("API key is not defined.")
	}

	if projectSlug == "" {
		log.Fatal("project slug is not defined.")
	}

	if lang == "" {
		log.Fatal("Language code is not defined.")
	}

	if outDir == "" {
		log.Fatal("output dir is not specified.")
	}

	stat, err := os.Stat(outDir)
	if err != nil {
		log.Fatal(err)
	}

	if !stat.IsDir() {
		log.Fatal("output dir is not a directory.")
	}

	if fullpath, err := filepath.Abs(filepath.Dir(outDir)); err != nil {
		log.Fatal(err)
	} else {
		outDir = fullpath
	}

	resources := getProjectResouces()

	for _, res := range resources {
		writeTranslationToFile(res)
		fmt.Println(res)
	}
}

func transifexAPIRequest(path string) *http.Response {
	client := &http.Client{}
	url := transifexAPI + projectSlug + path
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatal(err)
	}

	req.SetBasicAuth("api", apiKey)

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 || resp.ContentLength == 0 {
		log.Fatal("[status code " + strconv.Itoa(resp.StatusCode) + "] err in calling transifex API: " + url)
	}

	return resp
}

type resource struct {
	Slug string `json:"slug"`
	Name string `json:"name"`
}

func writeTranslationToFile(res resource) {
	resp := transifexAPIRequest("/resource/" + res.Slug + "/translation/" + lang + "?mode=translator&file")

	out, err := os.Create(outDir + "/" + filepath.Base(res.Name))
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := out.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	io.Copy(out, resp.Body)
}

func getProjectResouces() (v []resource) {
	resp := transifexAPIRequest("/resources/")

	dec := json.NewDecoder(resp.Body)
	if err := dec.Decode(&v); err != nil {
		log.Fatal(err)
	}

	return
}
