package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"unicode"

	"golang.org/x/net/html"
)

type ManPageOutput struct {
	ManualName  string
	Date        string
	CommandName string
	Section     int
	OutputPath  string
}

func (o ManPageOutput) Output(providers []TerraformProvider) error {

	var description string
	var references []string
	var err error

	for _, p := range providers {
		for _, r := range p.Resources {
			var output bytes.Buffer

			output.WriteString(fmt.Sprintf(`.TH %s %d "%s" "%s" "%s"`,
				o.CommandName,
				o.Section,
				o.Date,
				o.ManualName,
				r.Name))
			output.WriteString("\n")
			output.WriteString(fmt.Sprintf(".SH NAME\n%s\\ \\- %s\n\n", r.Name, r.Description))

			output.WriteString(".SH SYNOPSIS\n")
			output.WriteString(fmt.Sprintf(`resource "\fB%s\fR" "\fIname\fR" {%s`, r.Name, "\n"))
			output.WriteString(`    \fIoptions\fR`)
			output.WriteString("\n")
			output.WriteString(`.br`)
			output.WriteString("\n}\n\n")

			output.WriteString(".SH ARGUMENT REFERENCE\n")
			output.WriteString("The following arguments are supported:\n")

			for _, arg := range r.Arguments {
				output.WriteString(".TP\n")
				output.WriteString(fmt.Sprintf(`.BR %s\ (\fI%s\fR)%s`, arg.Name, arg.Type, "\n"))

				attributes := make([]string, 0)
				if arg.Required {
					attributes = append(attributes, "Required")
				}
				if arg.ForceNew {
					attributes = append(attributes, "Change forces new resource")
				}
				if arg.Computed {
					attributes = append(attributes, "Computed")
				}
				if arg.Default != "" {
					attributes = append(attributes, fmt.Sprintf("Default: %s", arg.Default))
				}

				output.WriteString(fmt.Sprintf("%s\n", strings.Join(attributes, ",")))
				output.WriteString(".br\n")
				description, references, err = htmlToTroff(arg.Description, references)
				if err != nil {
					return err
				}
				output.WriteString(fmt.Sprintf("%s\n", description))
			}

			if len(r.Attributes) > 0 {
				output.WriteString("\n\n.SH ATTRIBUTE REFERENCE\n")
				output.WriteString("The following attributes are exported:\n")

				for _, attr := range r.Attributes {
					output.WriteString(".TP\n")
					output.WriteString(fmt.Sprintf(`.BR %s\ (\fI%s\fR)%s`, attr.Name, attr.Type, "\n"))
					output.WriteString(".br\n")
					description, references, err = htmlToTroff(attr.Description, references)
					if err != nil {
						return err
					}
					output.WriteString(fmt.Sprintf("%s\n", description))
				}
			}

			if len(r.Examples) != 0 {
				output.WriteString("\n\n.SH EXAMPLES\n")
				output.WriteString(fmt.Sprintf("The following examples demonstrate the use of the %s resource.\n\n.br\n", r.Name))

				for i, ex := range r.Examples {
					var lines []string
					scanner := bufio.NewScanner(strings.NewReader(ex.Text))
					for scanner.Scan() {
						lines = append(lines, scanner.Text())
					}

					unindentedFirstLine := strings.TrimLeftFunc(lines[1], unicode.IsSpace)
					unindentLength := len(lines[1]) - len(unindentedFirstLine)

					description, references, err = htmlToTroff(ex.Description, references)
					if err != nil {
						return err
					}
					output.WriteString(fmt.Sprintf("%d: %s:\n", i+1, description))
					output.WriteString(".PP\n.nf\n.RS\n")
					output.WriteString(fmt.Sprintf("%s\n", unindentedFirstLine))
					for _, v := range lines[2:] {
						if len(v) > unindentLength {
							unindented := v[unindentLength:]
							output.WriteString(fmt.Sprintf("%s\n", unindented))
						} else {
							output.WriteString(fmt.Sprintf("%s\n", v))
						}
					}
					output.WriteString("\n.RE\n.fi\n.PP\n")
				}
			}

			if len(references) > 0 {
				output.WriteString("\n\n.SH REFERENCES\n")
				for index, url := range references {
					output.WriteString(fmt.Sprintf(`.BR %d\ -\ \fI%s\fR\n`, index+1, url))
				}
			}

			filename := filepath.Join(o.OutputPath, fmt.Sprintf("%s.%d", r.Name, o.Section))
			err := ioutil.WriteFile(filename, []byte(output.String()), 0644)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	return nil
}

func htmlToTroff(htmlIn string, references []string) (string, []string, error) {
	var output bytes.Buffer

	tokenizer := html.NewTokenizer(strings.NewReader(htmlIn))

ParseLoop:
	for {
		tt := tokenizer.Next()

		switch {
		case tt == html.ErrorToken:
			if tokenizer.Err() != io.EOF {
				return "", references, tokenizer.Err()
			}
			break ParseLoop
		case tt == html.TextToken:
			t := tokenizer.Token()
			output.WriteString(t.Data)
		case tt == html.StartTagToken:
			t := tokenizer.Token()
			switch t.Data {
			case "a":
				for _, a := range t.Attr {
					if a.Key == "href" {
						references = append(references, a.Val)
						break
					}
				}
				output.WriteString("")
			case "i":
				output.WriteString("_")
			case "b":
				output.WriteString("_")
			}
		case tt == html.EndTagToken:
			t := tokenizer.Token()
			switch t.Data {
			case "i":
				output.WriteString("_")
			case "a":
				output.WriteString(fmt.Sprintf("[%d]", len(references)))
			case "b":
				output.WriteString("_")
			}
		}
	}

	return output.String(), references, nil
}
