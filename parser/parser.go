package parser

import (
    "errors"
    "fmt"
    /*"regexp"*/
)

type Subparser func([]rune) (*Parser, error)

type factory func([]string) Subparser

var factories map[string]factory

func RegisterFactory(fact factory, name string){
    factories[name] = fact
}

type Parser struct{
    subparsers  []Subparser
    data        []rune
    eof         int
    idx         int
}

func (p *Parser) MakeFunction(factory_name string, args []string) error {
    factory, ok := factories[factory_name]
    if ok == false {
        return errors.New(fmt.Sprintf("Factory `%s` doesn't exit", factory_name))
    }
    p.subparsers = append(p.subparsers, factory(args))
    return nil
}

/*func (p *Parser) AddTextParser(txt string){*/
	/*if txt[0] == '^'{*/
        /*re, err := regexp.Compile(t)*/
        /*return func*/
	/*} else {*/
		/*p.pattern  = p.pattern + regexp.QuoteMeta(txt)*/
	/*}*/
/*}*/

/*func (p *Parser) Parse(data )*/

func NewParser() *Parser{
    p := new(Parser)
    p.subparsers = []Subparser{}
    return p
}

