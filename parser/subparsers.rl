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

/////////////// Any
func AnyFactory(args []string) (Subparser, error){
    if len(args) == 0 {
        d := new(DeferredFactoryDef)
        d.Fact = AnyFactory
        return nil, d
    }
    next_txt := []rune(args[1])
    next_eof := len(next_txt)
	return func(data []rune, p *Parser)([]rune, error){
        idx := p.idx
        end := p.idx
        next_idx := 0
        var next_char rune
start:
        if idx == p.eof{ goto break_ }
        next_char = next_txt[next_idx]
        switch data[idx] {
        case next_char:
            end = idx - next_idx
            next_idx++
            if next_idx == next_eof{ goto break_ }
            idx ++
            goto start
        default:
            idx ++
            end = idx
            next_idx = 0
            goto start
        }
break_:
        ret := data[p.idx: end]
        p.idx = end
        return ret, nil
    }, nil

    return nil, nil
}

////////////// until
func UntilFactory(args []string) (Subparser, error){
    next_txt := []rune(args[0])
    next_eof := len(next_txt)
    include := args[1] == "true"
    // almost same code as `any`, but to improve perfs avoiding tests
    // we copy/paste it 
    if !include {
        return func(data []rune, p *Parser)([]rune, error){
            start := p.idx
            end := p.idx
            next_idx := 0
            var next_char rune
start:
            // may be an error here...
            if p.idx == p.eof{ goto break_ }
            next_char = next_txt[next_idx]
            switch data[p.idx] {
            case next_char:
                end = p.idx - next_idx
                next_idx++
                if next_idx == next_eof{ goto break_ }
                p.idx ++
                goto start
            default:
                p.idx ++
                end = p.idx
                next_idx = 0
                goto start
            }
break_:
            ret := data[start : end]
            return ret, nil
        }, nil
    } else {
        return func(data []rune, p *Parser)([]rune, error){
            start := p.idx
            next_idx := 0
            var next_char rune
start:
            // TODO error...
            if p.idx == p.eof{ goto break_ }
            next_char = next_txt[next_idx]
            switch data[p.idx] {
            case next_char:
                next_idx ++
                p.idx ++
                if next_idx == next_eof{ goto break_ }
                goto start
            default:
                p.idx ++
                next_idx = 0
                goto start
            }
break_:
            ret := data[start : p.idx]
            return ret, nil
        }, nil
    }
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





