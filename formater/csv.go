package formater

import (
    "io"
    "encoding/csv"
    "logparser_ng/parser"
    "os"
    "sync"
)


type SVFormater struct{
    writer      *csv.Writer
    fieldnames  []string
}

func NewSVFormater(filename string, separator rune, fieldnames []string) (s *SVFormater){
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
    return
}


func (s *SVFormater)Pipe(input chan *parser.ResultMap) (stop chan bool){
    stop = make(chan bool)
    buffer := make(chan chan bool, 100000)
    wl := new(sync.Mutex)
    lines := [][]string{}
    go func(){
        for res := range input{
            r := make(chan bool, 1)
            buffer <- r
            go func(rmp *parser.ResultMap, rc chan bool){
                rm := *rmp
                line := []string{}
                for _, n := range s.fieldnames{
                    line = append(line, string(rm[n]))
                }
                wl.Lock()
                lines = append(lines, line)
                rmp = nil
                if len(lines) == 15000{
                    s.writer.WriteAll(lines)
                    lines = [][]string{}
                }
                wl.Unlock()
                rc <- true
            }(res, r)
        }
        close(buffer)
        input = nil
    }()
    go func(){
        for res := range buffer{<-res }
        wl.Lock()
        println(len(lines))
        s.writer.WriteAll(lines)
        wl.Unlock()
        s.writer.Flush()
        stop <- true
        close(stop)
    }()
    return


    //go func(){
        //for rm := range input{
            //buffer <- rm
        //}
        //close(buffer)
    //}()
    //go func(){
        //for rmp := range buffer{
            //rm := *rmp
            //line := []string{}
            //for _, n := range s.fieldnames{
                //line = append(line, string(rm[n]))
            //}
            //s.writer.Write(line)
        //}
        ////stop <- true
    //}()
    //return stop
}


func SVFormater_old(sep rune, w io.Writer, fieldnames []string) chan parser.ResultMap{
    cw := csv.NewWriter(w)
    input := make(chan parser.ResultMap)
    _ = cw.Write(fieldnames)
    count := 0
    bufferize := func(line []string, start chan bool, stop chan bool){
        <- start
        count ++
        cw.Write(line)
        stop <- true
    }
    start_chan := make(chan bool, 1)
    start_chan <- true
    for res := range input{
        line := []string{}
        for _, n := range fieldnames{
            line = append(line, string(res[n]))
        }
        end_chan := make(chan bool, 1)
        go bufferize(line, start_chan, end_chan)
        start_chan = end_chan
    }
    return input
}
