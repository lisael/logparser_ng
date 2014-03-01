package config

import (
    "fmt"
    plib "logparser_ng/parser"
    "errors"
)

func MakeFormater(data string){
//TODO
}

func MakeParser(data string) (*plib.Parser, error){
    mark := 0
    cs, p, pe, eof := 0, 0, len(data), len(data)
    _ = eof
    parser := plib.NewParser()
    tokenName := ""
    factoryName := ""
    argList := []string{}
    errorSampleStart := 0
    _ = errorSampleStart
    var err error
    // only for debug
    argName := ""

    sendText := func(){
        if p-mark > 0 {
            /*fmt.Printf("TEXT: `%s`\n", data[mark:p])*/
            parser.AddTextParser(string(data[mark:p]))
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
            //print("TOKEN: ")
            //println(tokenName)
        }
        action facname {
            factoryName = data[mark:p]
            /*print("FACTORY: ")*/
            /*println(factoryName)*/
        }
        action arg {
            if p - 1 - mark > 0 {
                argName = data[mark:p-1]
                if len(argName) > 1{
                    if argName[0] == '"' || argName[0] == '\''{
                        argName = argName[1:len(argName)-1]
                    }
                }
                argList = append(argList, argName)
                //print(" ARG: ")
                //println(argName)
            }
        }
        action error {
            if p > 10 {
                errorSampleStart = p-10
            }
            return nil, errors.New(fmt.Sprintf("error in config at char %d, after `%s`", p, data[errorSampleStart:p]))
        }

        action simple_token {
            tokenName = data[mark:p]
            /*print("TOKEN: ")*/
            /*println(tokenName)*/
            if tokenName[0] == '_'{
                // special tokens are actually unnamed tokens
                // eg |_ignore| => |:_ignore()|
                factoryName = tokenName
                tokenName = ""
            } else {
                factoryName="non_blank"
            }
        }

        action emit {
            //fmt.Printf("emit token `%s`\n", tokenName)
            err = parser.MakeSubparser(tokenName, factoryName, argList)
            argList = []string{}
            if err != nil {
                return nil, err
            }
        }

        DIGIT = [0-9];
        LETTER = [a-zA-Z_];
        IDENTIFIER = LETTER (DIGIT|LETTER)* ;
        TOKEN_SEP = "|";
        FAC_SEP = ":";
        WHITESPACE = ( ' ' | '\t' ) +;
        ARG_SEP = WHITESPACE* "," WHITESPACE*;
        TOKEN_NAME = IDENTIFIER?;
        FAC_NAME = IDENTIFIER;
        TEXT = [^|] + ;
        SIMPLE_TOKEN =  TOKEN_SEP
                        TOKEN_NAME >start_token :>>
                        TOKEN_SEP @simple_token %mark;
        SIMPLE_ARG = (DIGIT|LETTER)+; 
        STRING_ARG = ('"' [^"]* '"') | ("'" [^']* "'"); #'//restore go syntax highlighting
        ARG = (SIMPLE_ARG | STRING_ARG);
        ARGLIST = ( ARG("" :>>  ARG_SEP %arg %mark ARG)*)?;
        FACTORY = FAC_NAME :>> "(" @facname %mark ARGLIST :>> ")" %arg;
        TOKEN = (SIMPLE_TOKEN
                |( TOKEN_SEP
                   TOKEN_NAME >start_token :>>
                   FAC_SEP @tokname %mark
                   FACTORY
                   TOKEN_SEP %mark)) %emit $err(error);

        main := ( TOKEN | TEXT ) * $err(error);

        write init;
        write exec;
    }%%
    sendText()
    err = parser.Finalize()
    return parser, err
}

// |ip:ipv4()| - - [|date:catch_all()|] "|ignore| |url:nonblank()| |ignore| |response_code:nonblank()| |ignore| |ignore| |user_agent:quoted(\"\\\"\",true)|

