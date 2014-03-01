package config 

import (
    "testing"
    "fmt"
)

type makeParserTest struct{
    config  string
    ok      bool
    tokenNr int
}

func TestMakeParser(t *testing.T) {
    var ok string

    tests := []makeParserTest{
        makeParserTest{ "hello |my_token| world", true, 3 },
        makeParserTest{ "hello |_ignore| world", true, 3 },
        makeParserTest{ "hello |url:non_blank(coucou, true)| world", true, 3 },
        makeParserTest{ "hello |url:non_blank('coucou')| world", true, 3 },
        makeParserTest{ "hello |url:non_blank(\"coucou hello\", 'hello coucou')| world",
                        true, 3 },
        makeParserTest{ "hello |:non_blank()| world", true, 3 },
        makeParserTest{ "hello |url:any()| world", true, 3 },
        makeParserTest{ "hello |bad token| world", false, 0 },
        makeParserTest{ "hello |_bad| world", false, 3 },
        makeParserTest{ "hello |my_token| world |unclosed", false, 0 },
    }
    for _, testCase := range tests {
        p, err := MakeParser(testCase.config)
        if (err == nil) != testCase.ok{
            if testCase.ok { ok = fmt.Sprintf("pass: \n    %s", err)} else { ok = "not pass"}
            t.Errorf("config `%s` should %s", testCase.config, ok)
        }
        if p != nil {
            nr := p.SubparsersCount()
            if nr != testCase.tokenNr {
                t.Errorf("config `%s` created %d tokens. should be %d", testCase.config, nr, testCase.tokenNr)
            }
        }
    }
}
