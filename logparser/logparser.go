package main

import (
	"os"
    "runtime"
    "github.com/codegangsta/cli"
    conf "github.com/lisael/logparser_ng/config"
    "github.com/lisael/logparser_ng/writer"
    "github.com/lisael/logparser_ng/parser"
    "github.com/lisael/logparser_ng/utils"
    "time"
    "fmt"
	//"logparser_ng/filter"
)

const VERSION string = "0.1-devel"

func stat() error{
    //memstats := cmd.Flag.Lookup("memstats").Value.String()
    memstats := "mem.dat"
    if memstats != "" {
        //interval := cmd.Flag.Lookup("meminterval").Value.Get().(time.Duration)
        interval := time.Duration(10000000)

        fileMemStats, err := os.Create(memstats)
        if err != nil {
            return err
        }

        fileMemStats.WriteString("# Time\tHeapSys\tHeapAlloc\tHeapIdle\tHeapReleased\n")

        go func() {
            var stats runtime.MemStats

            start := time.Now().UnixNano()

            for {
                runtime.ReadMemStats(&stats)
                if fileMemStats != nil {
                    fileMemStats.WriteString(fmt.Sprintf("%d\t%d\t%d\t%d\t%d\n",
                        (time.Now().UnixNano()-start)/1000000, stats.HeapSys, stats.HeapAlloc, stats.HeapIdle, stats.HeapReleased))
                    time.Sleep(interval)
                } else {
                    break
                }
            }
        }()
    }
    return nil
}


func parseLog(input_file string, output_file string, pconfig string, fconfig string) {
    // set the reader
    reader := utils.NewFileReader(input_file)
    //j := new(utils.Janitor)

    // error handling
    errchan := make(chan *parser.ParsingContext)
    stopErrors := make(chan bool)
    go func(){
        for {
            select{
            case pctx := <- errchan:
                panic(pctx.Error)
            case <- stopErrors:
                return
            }
        }
    }()
    // set the parser
    // just to test, I don't know yet how shell escaping works...
    // TODO: test passing config in command line
    //parser_, err := conf.MakeParser(pconfig)
    parser_, err := conf.MakeParser("|ip:ipv4()| - - [|date:until(\"]\",false)|] |_ignore| |url| |_ignore| |http_code| |_ignore| |_ignore| \"|user_agent:until('\"',false)|\"", 1000, errchan)
    if err != nil { panic(err) }

    // set the writter
    var writer_ *writer.SVFormater
    switch fconfig{
    case "CSV":
        writer_ = writer.NewSVFormater(output_file, rune(','), parser_.FieldNames(), 5000)
    case "TSV":
        writer_ = writer.NewSVFormater(output_file, rune('\t'), parser_.FieldNames(), 5000)
    }
	//df := filter.NewDumyFilter()

    // launch the pipeline
    stop := writer_.Pipe(parser_.Pipe(reader.ReadLines()))
    <- stop
    stopErrors <- true
    return
}

func main() {
    runtime.GOMAXPROCS(runtime.NumCPU() * 12)
    _ = stat()
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
