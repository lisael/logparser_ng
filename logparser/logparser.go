package main

import (
	"github.com/codegangsta/cli"
	"os"
    conf "logparser_ng/config"
    "logparser_ng/formater"
    _ "logparser_ng/parser"
    "logparser_ng/utils"
)

const VERSION string = "0.1-devel"


func parseLog(input_file string, output_file string, pconfig string, fconfig string) {
    // set the reader
    reader := utils.NewFileReader(input_file)

    // set the parser
    // just to test, I don't know yet how shell escaping works...
    // TODO: test passing config in command line
    //p, err := conf.MakeParser(pconfig)
    p, err := conf.MakeParser("|ip:ipv4()| - - [|date:until(\"]\",false)|] |_ignore| |url| |_ignore| |http_code| |_ignore| |_ignore| \"|user_agent:until('\"',false)|\"")
    if err != nil { panic(err) }

    // set the writter
    var formater_ *formater.SVFormater
    switch fconfig{
    case "CSV":
        formater_ = formater.NewSVFormater(output_file, rune(','), p.FieldNames())
    case "TSV":
        formater_ = formater.NewSVFormater(output_file, rune('\t'), p.FieldNames())
    }

    // launch the pipeline
    stop := formater_.Pipe(p.Pipe(reader.ReadLines()))
    <- stop
    return
}

func main() {
	app := cli.NewApp()
	app.Name = "Logparser-ng"
	app.Usage = "Parse a log and return a TSV of interesting patterns"
	app.Version = VERSION
	app.Flags = []cli.Flag{
		cli.StringFlag{"output, o", "", "Output file"},
		cli.StringFlag{"parser-config, c", "", "Input format"},
		cli.StringFlag{"formater-config, f", "TSV", "Output format"},
	}
	app.Action = func(c *cli.Context) {
		input_file := ""
		if len(c.Args()) > 0 {
			input_file = c.Args()[0]
		}
        output_file := c.String("output")
        pconfig := c.String("parser-config")
        fconfig := c.String("formater-config")
		parseLog(input_file, output_file, pconfig, fconfig)
	}
	app.Run(os.Args)
}
