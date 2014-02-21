package parser

import (
    "errors"
    "fmt"
)

// DeferredFactoryDef are factories that are called
// after the whole config parsing. see AnyFactory code.
type DeferredFactoryDef struct{
    Fact    factory
    Args    []string
}
    
// well, it's not really an error...
func (d DeferredFactoryDef)Error() string{return ""}

////////////// non blank
func NonBlankParser(data []rune, pr *Parser) ([]rune, error){
    startIdx := pr.idx
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
    return data[startIdx: pr.idx], nil
}

// fake factory, requires no args
func NonBlankFactory([]string) (Subparser, error){
    return NonBlankParser, nil
}

///////////// IP
func IPV4Parser(data []rune, pr *Parser) ([]rune, error){
    p := pr.idx
    cs :=0
    pe, eof := pr.eof, pr.eof
    ts, te, act := 0, 0, 0
    _, _, _ = ts, te, act
    
    %%{
        machine ipv4;
        write data;

        action ip_error{
            return nil, errors.New(fmt.Sprintf("Error while parsing IP at column %d", p))
        }
        action ret{
            goto ret_ip
        
        }
        CHUNK = ( ("25"[0-5]) | ("2" [0-4][0-9]) | ("1"?[0-9][0-9]) | ([0-9]) );
        
        IP = CHUNK "." CHUNK "." CHUNK "." CHUNK;
        main := |* IP $err(ip_error) => ret;
        #any* => ret;
                *|;

        write init;
        write exec;
    }%%

ret_ip:
    // next char must not be a digit
    if eof-p > 1 {
        if data[p+1] > 47 && data[p+1] < 58{
            return nil, errors.New(fmt.Sprintf("Error while parsing IP at column %d", p))
        }
    }
    result := data[pr.idx:p+1]
    pr.idx = p+1
    return result, nil
}

// fake factory, requires no args
func IPV4Factory(args []string) (Subparser, error){
    return IPV4Parser, nil
}

/////////////// date

func AnyFactory(args []string) (Subparser, error){
    if len(args) == 0 {
        d := new(DeferredFactoryDef)
        d.Fact = AnyFactory
        return nil, d
    }
    println(len(args))
    println(args[0])
    next := []rune(args[0])
    _ = next
    return nil, nil
}





