package main

import (
	"bytes"
	"fmt"
	"learninglabs/sp-6node/internal/model"
	"log"
	"text/template"

	"os"

	"gopkg.in/yaml.v3"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: %s <vars.yml>", os.Args[0])
	}

	varsFile := os.Args[1]

	// Load YAML
	data, err := os.ReadFile(varsFile)
	if err != nil {
		log.Fatal(err)
	}

	var cfg model.DeviceConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatal(err)
	}

	// Parse ALL templates
	tmpl, err := template.ParseGlob("templates/*.xml.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err = tmpl.ParseGlob("templates/protocols/*.xml.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	// Render starting from "configuration"
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "configuration", cfg); err != nil {
		log.Fatal(err)
	}

	// Print rendered XML
	fmt.Println(buf.String())
}
