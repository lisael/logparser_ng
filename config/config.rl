package config

import (
    "fmt"
    plib "logparser_ng/parser"
    "errors"
)


func MakeParser(data string) (*plib.Parser, error){
    println("building parser...")
    mark := 0
    cs, p, pe, eof := 0, 0, len(data), len(data)
    _ = eof
    parser := plib.NewParser()
    tokenName := ""
    factoryName := ""
    argList := []string{}
    errorSampleStart := 0
    _ = errorSampleStart
    // only for debug
    argName := ""

    sendText := func(){
        if p-mark > 0 {
            fmt.Printf("TEXT: `%s`\n", data[mark:p])
        }
    }

    %%{
        machine configparser;
        write data;
        
        action start_token {
                p--
                sendText()
                p++
                mark = p
        }
        action mark { mark = p }
        action tokname {
            tokenName = data[mark:p]
            print("TOKEN: ")
            println(tokenName)
        }
        action facname {
            factoryName = data[mark:p]
            print("FACTORY: ")
            println(factoryName)
        }
        action arg {
            argName = data[mark:p]
            argList = append(argList, argName)
            print("ARG: ")
            println(argName)
        }
        action error {
            if p > 10 {
                errorSampleStart = p-10
            }
            return nil, errors.New(fmt.Sprintf("error in config at char %d, after `%s`", p, data[errorSampleStart:p]))
        }

        action simple_token {
            tokenName = data[mark:p]
            if tokenName == "ignore" || tokenName == "ignore_rest" {
                print("TOKEN: ")
                println(tokenName)
            } else {
                return nil, errors.New(fmt.Sprintf("`%s` at column %d is not a known simple token.\nAllowed simple tokens are `ignore` and `ignore_rest`",mark, data[mark:p]))
            }
        }

        DIGIT = [0-9];
        LETTER = [a-zA-Z_];
        IDENTIFIER = LETTER (DIGIT|LETTER)* ;
        TOKEN_SEP = "|";
        FAC_SEP = ":";
        WHITESPACE = ( ' ' | '\t' ) +;
        ARG_SEP = WHITESPACE* "," WHITESPACE* %mark;
        TOKEN_NAME = IDENTIFIER;
        FAC_NAME = IDENTIFIER;
        TEXT = [^|] + ;
        SIMPLE_TOKEN =  TOKEN_SEP
                        TOKEN_NAME >start_token :>>
                        TOKEN_SEP @simple_token %mark;
        ARG = (DIGIT|LETTER)+ %arg; 
        ARGLIST = ( ARG ( ARG_SEP ARG)*)*;
        FACTORY = FAC_NAME :>> "(" @facname %mark ARGLIST ")";
        TOKEN = SIMPLE_TOKEN
                |( TOKEN_SEP
                   TOKEN_NAME >start_token :>>
                   FAC_SEP @tokname %mark
                   FACTORY
                   TOKEN_SEP %mark) $err(error);

        main := ( TOKEN | TEXT ) * $err(error);

        write init;
        write exec;
    }%%
    sendText()
    return parser, nil
}

// |ip:ipv4()| - - [|date:catch_all()|] "|ignore| |url:nonblank()| |ignore| |response_code:nonblank()| |ignore| |ignore| |user_agent:quoted(\"\\\"\",true)|

