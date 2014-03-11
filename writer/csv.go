package writer

import (
    "encoding/csv"
    "github.com/lisael/logparser_ng/parser"
    "os"
)


type SVFormater struct{
    writer      *csv.Writer
    fieldnames  []string
	buffer		int
}

func NewSVFormater(filename string, separator rune, fieldnames []string, bufferSize int) (s *SVFormater){
    var w *os.File
    var err error
    if filename == ""{
        w = os.Stdout
    } else {
        w, err = os.Create(filename)
        if err != nil { panic(err) }
    }
    s = new(SVFormater)
    s.writer = csv.NewWriter(w)
    s.writer.Comma = separator
    s.writer.Write(fieldnames)
    s.fieldnames = fieldnames
	s.buffer = bufferSize
    return
}


func (s *SVFormater)Pipe(input chan *parser.ParsingContext) (stop chan bool){
    stop = make(chan bool)
    lines := [][]string{}
    go func(){
        for pctx := range input{
            line := []string{}
            for _, n := range s.fieldnames{
                 _ = pctx.Tokens
                 line = append(line, pctx.Tokens[n])
            }
            lines = append(lines, line)
            if len(lines) == s.buffer{
                s.writer.WriteAll(lines)
                lines = [][]string{}
			}
            select{
            case parser.PctxPool <- pctx:
            default:
                pctx = nil
            }
		}
        s.writer.WriteAll(lines)
        stop <- true
        close(stop)
	}()
	return
}
			
