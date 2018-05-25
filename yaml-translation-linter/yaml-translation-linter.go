package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

func main() {
	refData := readFile("en.yml")
	traData := readFile("fa.yml")

	var ref interface{}
	var tra interface{}

	if err := yaml.Unmarshal(refData, &ref); err != nil {
		log.Fatal(err)
	}

	if err := yaml.Unmarshal(traData, &tra); err != nil {
		log.Fatal(err)
	}

	if err := checkKeys(ref, tra, ""); err != nil {
		log.Fatal(err)
	}

	// for k, v := range ref {
	// 	fmt.Printf("key[%s] value[%s]\n", k, v)
	// }

	// refLines := strings.Split(ref, "\n")
	// traLines := strings.Split(tra, "\n")

	// if len(refLines) != len(traLines) {
	// 	log.Fatal("translation lines (" + strconv.Itoa(len(traLines)) + ") are not equal to reference lines (" + strconv.Itoa(len(refLines)) + ").")
	// }

	// fmt.Println(tra)
	fmt.Println("ok")
}

func readFile(file string) []byte {
	data, err := ioutil.ReadFile(file)

	if err != nil {
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
			for key, val := range refMap {
				path := parentPath + "." + key.(string)
				if _, ok := traMap[key]; !ok {
					return errors.New("cannot find '" + path + "' in translations")
				}
				if err := checkKeys(val, traMap[key], path); err != nil {
					return err
				}
			}

			for key := range traMap {
				if _, ok := refMap[key]; !ok {
					path := parentPath + "." + key.(string)
					return errors.New("find redunant translation '" + path + "'")
				}
			}

			return nil
		}

		return errors.New("cannot convert translation to map")
	}

	return errors.New("cannot convert reference to map")
}
