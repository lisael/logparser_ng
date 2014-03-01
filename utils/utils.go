package utils

import (
    "bufio"
    "os"
    "io"
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
    output = make(chan *string)
    go func (){
        for {
            line, err := r.reader.ReadString('\n')
            if err != nil {
                if err != io.EOF {
                    panic(err)
                }
                break
            }
            output <- &line
        }
        close(output)
    }()
    return
}

