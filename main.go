package main

import (
	"encoding/xml"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type OutputFormatter interface {
	Output([]TerraformProvider) error
}

var outputPath string

func init() {
	const (
		defaultOutputPath = "output/"
		usageOutputPath   = "base directory for output"
	)
	flag.StringVar(&outputPath, "output-path", defaultOutputPath, usageOutputPath)
	flag.StringVar(&outputPath, "o", defaultOutputPath, usageOutputPath+" (shorthand)")
}

func main() {
	flag.Parse()
	var inputs []TerraformProvider

	for _, f := range flag.Args() {
		var v TerraformProvider

		xmlInput, err := ioutil.ReadFile(f)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		err = xml.Unmarshal(xmlInput, &v)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
		inputs = append(inputs, v)
	}

	markdownOutput := filepath.Join(outputPath, "markdown/")
	manOutput := filepath.Join(outputPath, "man/")
	err := os.MkdirAll(markdownOutput, 0744)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = os.MkdirAll(manOutput, 0744)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	outputs := []OutputFormatter{
		&ManPageOutput{
			ManualName:  "Terraform Manual",
			Date:        "Feb 11, 2016",
			CommandName: "TERRAFORM",
			Section:     1,
			OutputPath:  manOutput,
		},
		&MarkdownOutput{
			Layout:       "triton",
			SectionName:  "Triton",
			SidebarClass: "docs-triton-firewall",
			OutputPath:   markdownOutput,
		},
	}
	for _, output := range outputs {
		output.Output(inputs)
	}
}
