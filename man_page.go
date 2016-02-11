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

type ManPageOutput struct {
	ManualName  string
	Date        string
	CommandName string
	Section     int
	OutputPath  string
}

func (o ManPageOutput) Output(providers []TerraformProvider) error {
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
				if arg.Default != "" {
					attributes = append(attributes, fmt.Sprintf("Default: %s", arg.Default))
				}

				output.WriteString(fmt.Sprintf("%s\n", strings.Join(attributes, ",")))
				output.WriteString(".br\n")
				output.WriteString(fmt.Sprintf("%s\n", arg.Description))
			}

			output.WriteString("\n\n.SH ATTRIBUTE REFERENCE\n")
			output.WriteString("The following attributes are exported:\n")

			for _, attr := range r.Attributes {
				output.WriteString(".TP\n")
				output.WriteString(fmt.Sprintf(`.BR %s\ (\fI%s\fR)%s`, attr.Name, attr.Type, "\n"))
				output.WriteString(".br\n")
				output.WriteString(fmt.Sprintf("%s\n", attr.Description))
			}

			if len(r.Examples) == 0 {
				continue
			}

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

				output.WriteString(fmt.Sprintf("%d: %s:\n", i+1, ex.Description))
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

			filename := filepath.Join(o.OutputPath, fmt.Sprintf("%s.%d", r.Name, o.Section))
			err := ioutil.WriteFile(filename, []byte(output.String()), 0644)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	return nil
}
