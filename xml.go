package main

import (
	"strings"
	"unicode"
)

type ResourceArgument struct {
	Name        string `xml:"name,attr"`
	Type        string `xml:"type,attr"`
	Required    bool   `xml:"required,attr"`
	ForceNew    bool   `xml:"forcenew,attr"`
	Computed    bool   `xml:"computed,attr"`
	Default     string `xml:"default,attr"`
	Description string `xml:",chardata"`
}

type ResourceAttribute struct {
	Name        string `xml:"name,attr"`
	Type        string `xml:"type,attr"`
	Description string `xml:",chardata"`
}

type TerraformResource struct {
	Name        string `xml:"name,attr"`
	Description string `xml:"description"`

	Arguments  []*ResourceArgument  `xml:"arguments>argument"`
	Attributes []*ResourceAttribute `xml:"attributes>attribute"`
	Examples   []*ResourceExample   `xml:"examples>example"`
}

type ResourceExample struct {
	Description string `xml:"description,attr"`
	Text        string `xml:",chardata"`
}

type TerraformProvider struct {
	Name      string               `xml:"name,attr"`
	Resources []*TerraformResource `xml:"resources>resource"`
}

func (p *TerraformProvider) Tidy() {
	for _, resource := range p.Resources {
		resource.Description = strings.TrimFunc(resource.Description, unicode.IsSpace)
	}
}
