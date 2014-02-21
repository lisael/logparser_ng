package parser

import (
    "errors"
    "fmt"
    /*"regexp"*/
)

type Subparser func([]rune, *Parser) ([]rune, error)

type factory func([]string) (Subparser, error)

type ResultMap map[string][]rune

var factories map[string]factory = map[string]factory{
	"_ignore": NonBlankFactory,
	"non_blank": NonBlankFactory,
	"ipv4": IPV4Factory,
    "any": AnyFactory,
}

func RegisterFactory(fact factory, name string){
    factories[name] = fact
}

type Parser struct{
    subparsers  []Subparser
    data        []rune
    eof         int
    idx         int
	resultNames	[]string
}

var defereds []*DeferredFactoryDef = []*DeferredFactoryDef{}
var texts []string = []string{}

func (p *Parser) MakeSubparser(tokenName string, factory_name string, args []string) error {
    var currentDefered *DeferredFactoryDef
    // TODO: raise an error for the number of args.
	// try to add the token name to result names
	fmt.Printf("making `%s` subparser for token `%s`", factory_name, tokenName )
	if tokenName != ""{
		for _, name := range p.resultNames{
			if name == tokenName{
				return errors.New(fmt.Sprintf("Token name `%s` allready in use" ))
			}
		}
	}
	p.resultNames = append(p.resultNames, tokenName)
    // find the factory and instanciate a subparser
	factory := factories[factory_name]
    if factory == nil {
        return errors.New(fmt.Sprintf("Factory `%s` doesn't exit", factory_name))
    }
    println("calling factory")
	subp, err := factory(args)
	if err != nil{
        switch err.(type){
        case *DeferredFactoryDef:
            currentDefered = err.(*DeferredFactoryDef)
            currentDefered.Args = args
            defereds = append(defereds, currentDefered)
            p.subparsers = append(p.subparsers, subp)
            return nil
        default:
            return err
        }
	}
    p.subparsers = append(p.subparsers, subp)
    return nil
}

func (p *Parser) AddTextParser(txt string){
	// we just have to move forward without emitting anything
	fmt.Printf("making text subparser for text `%s`", txt )
	p.subparsers = append(p.subparsers, SkipFactory(len(txt)))
	p.resultNames = append(p.resultNames, "")
	fmt.Println("  Done")
}

func (p *Parser) SubparsersNumber() int{
    return len(p.subparsers)
}

func (p *Parser) ParseOnce(data []rune) (ResultMap, error){
	results := ResultMap{}
	p.eof = len(data)
	p.idx = 0
	var nextSubparser Subparser
	for idx, name := range p.resultNames{
		nextSubparser = p.subparsers[idx]
		result, err := nextSubparser(data[p.idx:], p)
		if err != nil {
			return nil, err
		}
		// do not save unnamed chunks
		if name != ""{
			results[name] = result
		}
	}
	return results, nil
}

func NewParser() *Parser{
    p := new(Parser)
    p.subparsers = []Subparser{}
	p.resultNames = []string{}
    return p
}

func SkipFactory(length int) Subparser {
	return func(data []rune, p *Parser)([]rune, error){
		p.idx += length
		return nil, nil
	}
}

