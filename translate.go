package main

import (
	"flag"
	"fmt"
	"github.com/bregydoc/gtranslate"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var (
	from = flag.String("from", "en", "from")
	to   = flag.String("to", "tr", "to")
	f    = flag.String("f", "", "yml file or folder path")
)

type custom map[interface{}]interface{}

func logIfNotOk(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {

	flag.Parse()

	if *from == "" || *to == "" || *f == "" {
		log.Fatalln("from and to arguments can't be empty")
	}

	fileInfo, err := os.Stat(*f)
	logIfNotOk(err)

	if fileInfo.Mode().IsDir() {
		files, err := ioutil.ReadDir(*f)
		logIfNotOk(err)

		for _, file := range files {
			translateFile(*f + file.Name())
		}
	} else {
		translateFile(*f)
	}

}

func translateFile(fileName string) {
	if !strings.HasSuffix(fileName, "yml") || strings.HasSuffix(fileName, "yaml") {
		return
	}

	var (
		body []byte
		err  error
	)

	file, err := os.Open(fileName) // For read access.
	logIfNotOk(err)
	defer file.Close()

	body, err = ioutil.ReadAll(file)
	logIfNotOk(err)

	var yamlFile custom
	err = yaml.Unmarshal(body, &yamlFile)
	logIfNotOk(err)

	transformed := translate(yamlFile)

	body, err = yaml.Marshal(&transformed)
	logIfNotOk(err)
	println(string(body))

	newFileName := file.Name() + "-" + *to
	err = ioutil.WriteFile(newFileName, body, 0644)
	logIfNotOk(err)
	log.Println(file.Name()+" is done.")
}

func translate(yml interface{}) (transformed interface{}) {

	switch field := yml.(type) {
	case custom:
		transformed = make(custom)
		for key := range field {
			transformed.(custom)[key] = translate((field)[key])
		}
	case []interface{}:
		transformed = make([]interface{}, 0)
		for key := range field {
			transformed = append(transformed.([]interface{}), translate((field)[key]))
		}
	case string:
		if strings.Contains(field, "%") || field == "" {
			return field
		}

		translation, err := internalTranslate(field)
		logIfNotOk(err)
		return translation
	default:
		println(fmt.Sprintf("%T", field))
		return field
	}

	return transformed
}

func internalTranslate(text string) (string, error) {
	return gtranslate.TranslateWithParams(
		text,
		gtranslate.TranslationParams{
			From: *from,
			To:   *to,
		},
	)
}
