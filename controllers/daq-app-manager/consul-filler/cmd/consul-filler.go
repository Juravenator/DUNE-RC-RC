package main

import (
	"bufio"
	"os"
	"text/template"
	"github.com/Masterminds/sprig"
	"strings"
	"consul-filler.ccm.dunescience.org/internal"
)

func main() {
	input := ""
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		escapedLine := ""

		for line != "" {
			tmplStartI := strings.Index(line, "{{")
			if tmplStartI == -1 {
				escapedLine += line
				break
			}
			tmplEndI := strings.Index(line[tmplStartI:], "}}")
			if tmplEndI == -1 {
				escapedLine += line
				break
			}
			tmplEndI += tmplStartI
			escapedLine += line[:tmplStartI]
			escapedLine += strings.ReplaceAll(line[tmplStartI:tmplEndI+2], "\\\"", "\"")
			line = line[tmplEndI+2:]
		}

		input += escapedLine + "\n"
	}

	// fmt.Fprintf(os.Stderr, input)

	inputTemplate, err := template.New("").Funcs(sprig.TxtFuncMap()).Funcs(funcs).Parse(input)
	if err != nil {
		panic(err)
	}

	err = inputTemplate.Execute(os.Stdout, templateData{
		RunNumber: "5",
	})
	if err != nil {
		panic(err)
	}
}

type templateData struct {
	RunNumber string
}

var funcs = template.FuncMap{
	"firstServiceAddr": internal.GetFirstServiceAddress,
	// "makeStringArray":    util.TemplateStringArray,
	// "capitalize":         util.Capitalize,
	// "noSpaces":           util.NoSpaces,
	// "splitCommand":       splitCommand,
	// "lastCommand":        lastCommand,
	// "penultimateCommand": penultimateCommand,
}