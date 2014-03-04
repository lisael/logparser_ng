package utils

import (
    "bufio"
    "os"
	"io"
    "runtime"
)

type FileReader struct{
    reader  *bufio.Reader
}

func NewFileReader(filename string) (r *FileReader){
    r = new(FileReader)
    var f *os.File
    var err error
    if filename == ""{
        f = os.Stdin
    } else {
        f, err = os.Open(filename)
        if err != nil { panic(err) }
    }
    r.reader = bufio.NewReader(f)
    return
}

func (r *FileReader)ReadLines() (output chan *string){
    output = make(chan *string, 1000)
    /*s := "42.42.42.42 - - [03/Jan/2014:06:25:33 +0100] \"GET http://www.example.com/stuff.html HTTP/1.1\" 302 26 \"-\" \"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)\""*/
	
    go func (){
		/*for i:=0; i<1000000; i++{*/
		for {
			line, err := r.reader.ReadString('\n')
			if err != nil {
				if err != io.EOF {
					panic(err)
				}
				break
			}
			output <- &line
            /*output <- &s*/
        }
        close(output)
    }()
    return
}

// not used at the moment, the memory leak came from a delayed fs flush
type Janitor struct{}

func (j *Janitor)Pipe(input chan *string) (output chan *string){
    output = make(chan *string)
    go func (){
        i := 0
        for line := range input{
            i++
            if i == 200000 {
                //println("GC")
                runtime.GC()
                //println(runtime.NumGoroutine())
                i=0
            }
            output <- line
        }
        close(output)
        input = nil
        runtime.GC()
    }()
    return
}

