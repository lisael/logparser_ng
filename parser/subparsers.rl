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
func NonBlankParser(pctx *ParsingContext) (string, error){
    startIdx := pctx.idx
start:
    if pctx.idx == pctx.eof{ goto break_ }
    switch pctx.Data[pctx.idx] {
    case ' ', '\t':
        goto break_
    default:
        pctx.idx ++
        goto start
    }
break_:
    return pctx.Data[startIdx: pctx.idx], nil
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
    next_txt := args[1]
    next_eof := len(next_txt)
	return func(pctx *ParsingContext)(string, error){
        idx := pctx.idx
        end := pctx.idx
        next_idx := 0
        var next_char uint8
start:
        if idx == pctx.eof{ goto break_ }
        next_char = next_txt[next_idx]
        switch pctx.Data[idx] {
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
        ret := pctx.Data[pctx.idx: end]
        pctx.idx = end
        return ret, nil
    }, nil

    return nil, nil
}

////////////// until
func UntilFactory(args []string) (Subparser, error){
    next_txt := args[0]
    next_eof := len(next_txt)
    include := args[1] == "true"
    // almost same code as `any`, but to improve perfs avoiding tests
    // we copy/paste it 
    if !include {
        return func(pctx *ParsingContext)(string, error){
            start := pctx.idx
            end := pctx.idx
            next_idx := 0
            var next_char uint8
start:
            // may be an error here...
            // TODO
            if pctx.idx == pctx.eof{ goto break_ }
            next_char = next_txt[next_idx]
            switch pctx.Data[pctx.idx] {
            case next_char:
                end = pctx.idx - next_idx
                next_idx++
                if next_idx == next_eof{ goto break_ }
                pctx.idx ++
                goto start
            default:
                pctx.idx ++
                end = pctx.idx
                next_idx = 0
                goto start
            }
break_:
            ret := pctx.Data[start : end]
            return ret, nil
        }, nil
    } else {
        return func(pctx *ParsingContext)(string, error){
            start := pctx.idx
            next_idx := 0
            var next_char uint8
start:
            // TODO error...
            if pctx.idx == pctx.eof{ goto break_ }
            next_char = next_txt[next_idx]
            switch pctx.Data[pctx.idx] {
            case next_char:
                next_idx ++
                pctx.idx ++
                if next_idx == next_eof{ goto break_ }
                goto start
            default:
                pctx.idx ++
                next_idx = 0
                goto start
            }
break_:
            ret := pctx.Data[start : pctx.idx]
            return ret, nil
        }, nil
    }
}



///////////// IP
func IPV4Parser(pctx *ParsingContext) (string, error){
    data := pctx.Data
    p := pctx.idx
    cs :=0
    pe, eof := pctx.eof, pctx.eof
    ts, te, act := 0, 0, 0
    _, _, _ = ts, te, act
    
    %%{
        machine ipv4;
        write data;

        action ip_error{
            return "", errors.New(fmt.Sprintf("Error while parsing IP at column %d", p))
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
            return "", errors.New(fmt.Sprintf("Error while parsing IP at column %d", p))
        }
    }
    result := data[pctx.idx:p+1]
    pctx.idx = p+1
    return result, nil
}

// fake factory, requires no args
func IPV4Factory(args []string) (Subparser, error){
    return IPV4Parser, nil
}





