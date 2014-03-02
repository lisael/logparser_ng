package parser

import (
    "errors"
    "fmt"
    "sync"
)

type ParsingContext struct{
    eof     int
    idx     int
    data    []rune
}

type Subparser func(*ParsingContext) ([]rune, error)

type factory func([]string) (Subparser, error)

type ResultMap map[string][]rune

var factories map[string]factory = map[string]factory{
	"_ignore": NonBlankFactory,
	"non_blank": NonBlankFactory,
	"ipv4": IPV4Factory,
    "any": AnyFactory,
    "until": UntilFactory,
}

func RegisterFactory(fact factory, name string){
    factories[name] = fact
}

type ParserBuildContext struct{
    Defereds []*DeferredFactoryDef
    Texts []*string
}

func NewParserBuildContext() *ParserBuildContext{
    p := new(ParserBuildContext)
    p.Defereds = []*DeferredFactoryDef{}
    p.Texts = []*string{}
    return p
}

type Parser struct{
    subparsers  []Subparser
	resultNames	[]string
    ctx         *ParserBuildContext
}


func (p *Parser) FieldNames() []string{
    names := []string{}
    for _, name := range p.resultNames{
        if name != ""{
            names = append(names, name)
        }
    }
    return names
}

func (p *Parser) MakeSubparser(tokenName string, factory_name string, args []string) error {
    var currentDefered *DeferredFactoryDef
    // TODO: raise an error for the number of args.
	// try to add the token name to result names
    //fmt.Printf("making `%s` subparser for token `%s`", factory_name, tokenName )
	if tokenName != ""{
		for _, name := range p.resultNames{
			if name == tokenName{
				return errors.New(fmt.Sprintf("Token name `%s` already in use" ))
			}
		}
	}
	p.resultNames = append(p.resultNames, tokenName)
    // find the factory and instanciate a subparser
	factory := factories[factory_name]
    if factory == nil {
        return errors.New(fmt.Sprintf("Factory `%s` doesn't exit", factory_name))
    }
	subp, err := factory(args)
	if err != nil{
        switch err.(type){
        case *DeferredFactoryDef:
            currentDefered = err.(*DeferredFactoryDef)
            currentDefered.Args = args
            p.ctx.Defereds = append(p.ctx.Defereds, currentDefered)
            p.subparsers = append(p.subparsers, subp)
            p.ctx.Texts = append(p.ctx.Texts, nil)
            return nil
        default:
            return err
        }
	}
    p.subparsers = append(p.subparsers, subp)
    p.ctx.Texts = append(p.ctx.Texts, nil)
    p.ctx.Defereds = append(p.ctx.Defereds, nil)
    return nil
}

func (p *Parser) AddTextParser(txt string){
	// we just have to move forward without emitting anything
	//fmt.Printf("making text subparser for text `%s`", txt )
	p.subparsers = append(p.subparsers, SkipFactory(len(txt)))
	p.resultNames = append(p.resultNames, "")
    p.ctx.Texts = append(p.ctx.Texts, &txt)
    p.ctx.Defereds = append(p.ctx.Defereds, nil)
	//fmt.Println("  Done")
}

func (p *Parser) Finalize() error{
    // call defered factories
    var before, after *string = nil, nil
    for _, s := range p.ctx.Texts {
        nn := "nil"
        if s  == nil {s = &nn}
    }
    for idx, def := range p.ctx.Defereds{
        if def == nil { before = p.ctx.Texts[idx]; continue}
        if len(p.ctx.Texts) >= idx+1 {
            after = p.ctx.Texts[idx+1]
        }
        if before == nil {
            def.Args = append(def.Args, "")
        } else { def.Args = append(def.Args, *before) }
        if after == nil {
            def.Args = append(def.Args, "")
        } else { def.Args = append(def.Args, *after) }
        subp, err := def.Fact(def.Args)
        if err != nil {return err}
        p.subparsers[idx] = subp
    }
    return nil
}

func (p *Parser) SubparsersCount() int{
    return len(p.subparsers)
}

var pctxPool chan *ParsingContext = make(chan *ParsingContext, 1000000)
var pctxLock *sync.Mutex= new(sync.Mutex)

func (p *Parser)getParsingContext() *ParsingContext{
    var pctx *ParsingContext
    pctxLock.Lock()
    if len(pctxPool) != 0 {
        pctx = <- pctxPool
    }
    pctxLock.Unlock()
    if pctx == nil {
        pctx = new(ParsingContext)
    }
    return pctx
}

func (p *Parser) ParseOnce(data []rune) (*ResultMap, error){
	results := ResultMap{}
    pctx := p.getParsingContext()
	pctx.eof = len(data)
	pctx.idx = 0
    pctx.data = data
	var nextSubparser Subparser
	for idx, name := range p.resultNames{
		nextSubparser = p.subparsers[idx]
		result, err := nextSubparser(pctx)
		if err != nil {
			return nil, err
		}
		// do not save unnamed chunks
		if name != ""{
			results[name] = result
		}
	}
    pctxPool <- pctx
	return &results, nil
}

func (p *Parser) Pipe(input chan *string) (output chan *ResultMap){
    output = make(chan *ResultMap)
    buffer := make(chan chan *ResultMap, 100000)
    go func(){
        for line := range input{
            r := make(chan *ResultMap, 1)
            buffer <- r
            go func(data []rune, rc chan *ResultMap){
                result, _ := p.ParseOnce(data)
                rc <- result
            }([]rune(*line), r)
            line = nil
        }
        close(buffer)
        input = nil
    }()
    go func(){
        for res := range buffer{
            rm := <- res
            res = nil
            output <- rm
        }
        close(output)
        buffer=nil
    }()
    return
}


func NewParser() *Parser{
    p := new(Parser)
    p.subparsers = []Subparser{}
	p.resultNames = []string{}
    p.ctx = NewParserBuildContext()
    return p
}

func SkipFactory(length int) Subparser {
	return func(pctx *ParsingContext)([]rune, error){
		pctx.idx += length
		return nil, nil
	}
}

