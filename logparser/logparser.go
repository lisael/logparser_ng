package main

import (
	"github.com/codegangsta/cli"
	"os"
    lpconf "logparser_ng/config"
)

const VERSION string = "0.1-devel"

func parseLog(input_file string, output_file string, config string) {
	println(input_file, output_file, config)
    p, err := lpconf.MakeParser(config)
    if err != nil {
        panic(err)
    }
    _ = p
}

func main() {
	app := cli.NewApp()
	app.Name = "Logparser-ng"
	app.Usage = "Parse a log and return a TSV of interesting patterns"
	app.Version = VERSION
	app.Flags = []cli.Flag{
		cli.StringFlag{"output, o", "", "Output file"},
		cli.StringFlag{"config, c", "", "Imput format"},
	}
	app.Action = func(c *cli.Context) {
	    input_file := ""
		if len(c.Args()) > 0 {
			input_file = c.Args()[0]
		}
        output_file := c.String("output")
        config := c.String("config")
		parseLog(input_file, output_file, config)
	}
	app.Run(os.Args)
}
