package main

type ResourceArgument struct {
	Name        string `xml:"name,attr"`
	Type        string `xml:"type,attr"`
	Required    bool   `xml:"required,attr"`
	ForceNew    bool   `xml:"forcenew,attr"`
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
	Description string `xml:"description,attr"`

	Arguments  []ResourceArgument  `xml:"arguments>argument"`
	Attributes []ResourceAttribute `xml:"attributes>attribute"`
	Examples   []ResourceExample   `xml:"examples>example"`
}

type ResourceExample struct {
	Description string `xml:"description,attr"`
	Text        string `xml:",chardata"`
}

type TerraformProvider struct {
	Name      string              `xml:"name,attr"`
	Resources []TerraformResource `xml:"resources>resource"`
}
