package parser

import (
    "errors"
    "fmt"
)


////////////// ignore
func IgnoreParser(data []rune, pr *Parser) ([]rune, error){
start:
    if pr.idx == pr.eof{ goto break_ }
    switch data[pr.idx] {
    case ' ', '\t':
        goto break_
    default:
        pr.idx ++
        goto start
    }
break_:
    return data, nil
}

// fake factory, requires no args
func IgnoreFactory([]string) (Subparser, error){
    return IgnoreParser, nil
}

///////////// IP
func IPV4Parser(data []rune, pr *Parser) ([]rune, error){
    p := pr.idx
    cs :=0
    pe, eof := pr.eof, pr.eof
    
    %%{
        machine ipv4;
        write data;

        action ip_error{
            return nil, errors.New(fmt.Sprintf("Error while parsing IP at column %d", p))
        }

        CHUNK = [1-2]?[0-9]?[0-9];
        
        main := CHUNK "." CHUNK "." CHUNK "." CHUNK $err(ip_error);

        write init;
        write exec;
    }%%
    result := data[pr.idx:p]
    pr.idx = p
    return result, nil
}

// fake factory, requires no args
func IPV4Factory([]string) (Subparser, error){
    return IPV4Parser, nil
}
