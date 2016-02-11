package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"unicode"
)

type MarkdownOutput struct {
	Layout       string
	SectionName  string
	SidebarClass string
	OutputPath   string
}

const frontMatterFormat = `---
layout: "%s"
page_title: "%s: %s"
sidebar_current: "%s"
description: |-
	%s 
---`

func (o MarkdownOutput) Output(providers []TerraformProvider) error {
	for _, p := range providers {
		for _, r := range p.Resources {
			var output bytes.Buffer

			frontMatter := fmt.Sprintf(frontMatterFormat,
				o.Layout,
				o.SectionName,
				r.Name,
				o.SidebarClass,
				r.Description)
			output.WriteString(fmt.Sprintf("%s\n\n", frontMatter))

			output.WriteString(fmt.Sprintf("# %s\n\n", strings.Replace(r.Name, "_", "\\_", -1)))
			output.WriteString(fmt.Sprintf("%s\n\n", r.Description))

			if len(r.Examples) > 0 {
				output.WriteString("## Example Usages\n\n")
				for _, ex := range r.Examples {
					var lines []string
					scanner := bufio.NewScanner(strings.NewReader(ex.Text))
					for scanner.Scan() {
						lines = append(lines, scanner.Text())
					}

					unindentedFirstLine := strings.TrimLeftFunc(lines[1], unicode.IsSpace)
					unindentLength := len(lines[1]) - len(unindentedFirstLine)

					output.WriteString(fmt.Sprintf("### %s\n", ex.Description))
					output.WriteString("\n```\n")
					output.WriteString(fmt.Sprintf("%s\n", unindentedFirstLine))
					for _, v := range lines[2:] {
						if len(v) > unindentLength {
							unindented := v[unindentLength:]
							output.WriteString(fmt.Sprintf("%s\n", unindented))
						} else {
							output.WriteString(fmt.Sprintf("%s\n", v))
						}
					}
					output.WriteString("```\n")
				}

				output.WriteString("\n## Argument Reference\n\nThe following arguments are supported:\n\n")

				for _, arg := range r.Arguments {
					attributes := []string{arg.Type}
					if arg.Required {
						attributes = append(attributes, "Required")
					}
					if arg.ForceNew {
						attributes = append(attributes, "Change forces new resource")
					}

					output.WriteString(fmt.Sprintf("* `%s` - (%s)", arg.Name, strings.Join(attributes, ",")))
					if arg.Default != "" {
						output.WriteString(fmt.Sprintf("  Default: `%s`", arg.Default))
					}
					output.WriteString(fmt.Sprintf("\n    %s\n\n", arg.Description))
				}

				output.WriteString("## Attribute Reference\n\nThe following attributes are exported:\n\n")

				for _, attr := range r.Attributes {
					output.WriteString(fmt.Sprintf("* `%s` - (%s) - %s \n", attr.Name, attr.Type, attr.Description))
				}
			}

			filename := filepath.Join(o.OutputPath, fmt.Sprintf("%s.md", r.Name))
			err := ioutil.WriteFile(filename, []byte(output.String()), 0644)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	return nil
}
