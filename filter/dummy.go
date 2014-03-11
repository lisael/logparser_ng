package filter

import (
	"logparser_ng/parser"
	"runtime"
	"time"
)


type dumyWork struct{
	data *parser.ParsingContext
	result chan *parser.ParsingContext
}

type DumyFilter struct{
	workers		int // max number of workers
	running		int // running workers
	works		chan *dumyWork // work queue
	results		chan chan *parser.ParsingContext // pending results
	stop		chan bool
}

//////////////////// plumbing
func NewDumyFilter(buffSize int) *DumyFilter{
    self := new(DumyFilter)
	// this task is CPU bound
	concurency := runtime.NumCPU()
	self.workers = concurency
	self.works = make(chan *dumyWork, buffSize) 
	self.results = make(chan chan *parser.ParsingContext, buffSize) 
	self.stop = make(chan bool)
	for i:=0; i<self.workers; i++ {
		go self.worker()
		self.running ++
	}
    return self
}

func (self *DumyFilter) worker(){
	for{
		select{
		case <- self.stop:
			return
		case w := <-self.works:
			rm, _ := self.Process(w.data)
			w.result <- rm
		}
	}
}

func (self *DumyFilter) Pipe(input chan *parser.ParsingContext) (output chan *parser.ParsingContext){
    output = make(chan *parser.ParsingContext, 1000)
    go func(){
        for in := range input{
            r := make(chan *parser.ParsingContext, 1)
            self.results <- r
			w := new(dumyWork)
			w.data = in
			w.result = r
			self.works <- w
        }
		close(self.results)
    }()
    go func(){
        for res := range self.results{
            rm := <- res
			if rm != nil {
				output <- rm
            }
        }
        close(output)
		for self.running > 0{
			self.stop <- true
			self.running --
		}
        close(self.works)
    }()
    return
}

////////////// features


func (self DumyFilter) Process(data *parser.ParsingContext) (*parser.ParsingContext, error){
	ts := time.Now()
	for i:=0; i<10000; i++{}
	tf := time.Now().Sub(ts).Nanoseconds()
	if tf % 2 == 0{
		return data, nil
	} else {
		return nil, nil
	}
}
