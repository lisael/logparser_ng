package formater

import (
    "time"
    "github.com/lisael/logparser_ng/parser"
)

type dateParser func(string) time.Time

type dateFormater func(time.Time) string

func parseStringDateFactory(pattern) dateParser{
    p := func(s string) (result time.Time){
        return time.Parse(pattern, s)
    }
    return p
}

var parsers map[string]dateParser = map[string]datePaser{}
var formaters map[string]dateFormater = map[string]dateFormater{}

type DateFormater struct{
    fieldName   string
    parser      dateParser
    formater_   dateFormater
}

func NewDateFormater(fieldname string, inputFormat string, outputFormat) (d *DateFormater){
    d = new(DateFormater)
    d.fieldName = fieldname
    var f dateFormater
    var fok bool
    // try to find a registerd date parser
    f, fok = formaters[outputFormat]
    if fok == false{
        // it may be a literal date layout
        f = func(output time.Time) string, err{
            s := output.Format(outputFormat)
            return s
        }
    }
    d.formater_ = f
    var p dateParser
    var pok bool
    // try to find a registerd date parser
    p, pok = parsers[inputFormat]
    if pok == false{
        // it may be a literal date layout
        p = func(input string) time.Time, err{
            t, err := time.Parse(inputFormat, input)
            return t, err
        }
    }
    d.parser = p
    return d
}

func (self *DateFormater) Process(parser.Resu)
