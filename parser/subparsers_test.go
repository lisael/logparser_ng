package parser

import (
    "time"
    "testing"
    "regexp"
    "fmt"
)

///////////// utils
var p *Parser = NewParser()

type testCase struct{
    data    []rune
    result  []rune
    end_idx int
    ok      bool
}

func PerformSubparserTests(sp Subparser, tests []testCase) []string{
    terrors := []string{}
    var ok string
    for cid, case_ := range tests {
        p.idx = 0
        p.eof = len(case_.data)
        result, err := sp(case_.data, p)
        if (err == nil) != case_.ok{
            if case_.ok { ok = "pass"} else { ok = "not pass"}
            println(err)
            errors := fmt.Sprintf("case_ %d: `%s` should %s", cid, string(case_.data), ok)
            terrors = append(terrors, errors, )
        }
        if result != nil {
            if string(result) != string(case_.result) {
                terrors = append( terrors, fmt.Sprintf("case %d: got `%s`. expected `%s`", cid, string(result), string(case_.result)))
            }
            if case_.end_idx < 0 { case_.end_idx = len(result)}
            if p.idx != case_.end_idx{
                terrors = append(terrors, fmt.Sprintf("case %d: bad parser index. got `%d`. expected `%d`", cid, p.idx, case_.end_idx))
            }
        }
    }
    return terrors
}

func BenchParserPerfs(sp Subparser, datas string, rx string) int64{
	const rounds = 1000
    data := []rune(datas)
    p.eof=len(datas)
    ts1 := time.Now()
	for i := 0; i < rounds; i++ {
        p.idx=0
		sp(data, p)
	}
	tp := time.Now().Sub(ts1).Nanoseconds()

    re , _ := regexp.Compile(rx)
    ts1 = time.Now()
	for i := 0; i < rounds; i++ {
        re.MatchString(datas)
	}
    tre := time.Now().Sub(ts1).Nanoseconds()
    return tre / tp
}


/////////////// tests
func TestIPV4Parser(t *testing.T) {
    tests := []testCase{
        testCase{ []rune("123.123.123.123"),
                  []rune("123.123.123.123"),
                  -1,
                  true},
        testCase{ []rune("123.123.123.123 .."),
                  []rune("123.123.123.123"),
                  -1,
                  true},
        testCase{ []rune("0.0.0.0.."),
                  []rune("0.0.0.0"),
                  -1,
                  true},
        testCase{ []rune("255.255.0.255"),
                  []rune("255.255.0.255"),
                  -1,
                  true},
        testCase{ []rune("423.123.123.123"),
                  []rune{},
                  -1,
                  false},
        testCase{ []rune("123.123.123.423"),
                  []rune{},
                  -1,
                  false},
    }
    res := PerformSubparserTests(IPV4Parser, tests)
    for _, r := range res{
        t.Errorf(r)
    }
}

func TestNonBlankParser(t *testing.T) {
    tests := []testCase{
        testCase{ []rune("hello world"),
                  []rune("hello"),
                  -1,
                  true},
        testCase{ []rune(" hello world"),
                  []rune(""),
                  -1,
                  true},
        testCase{ []rune("hello"),
                  []rune("hello"),
                  -1,
                  true},
    }
    res := PerformSubparserTests(NonBlankParser, tests)
    for _, r := range res{
        t.Errorf(r)
    }
}

func TestAny(t *testing.T){
    any, _ := AnyFactory([]string{"", " world"})
    tests := []testCase{
        testCase{ []rune("hello my world, bitch"),
                  []rune("hello my"),
                  -1,
                  true},
        testCase{ []rune(" hello world"),
                  []rune(" hello"),
                  -1,
                  true},
        testCase{ []rune("hello"),
                  []rune("hello"),
                  -1,
                  true},
    }
    res := PerformSubparserTests(any, tests)
    for _, r := range res{
        t.Errorf(r)
    }

}

func TestUntil(t *testing.T){
    until, _ := UntilFactory([]string{" world", "false"})
    tests := []testCase{
        testCase{ []rune("hello my world, bitch"),
                  []rune("hello my"),
                  13,
                  true},
        testCase{ []rune(" hello world"),
                  []rune(" hello"),
                  11,
                  true},
        testCase{ []rune("hello"),
                  []rune("hello"),
                  -1,
                  true},
    }
    res := PerformSubparserTests(until, tests)
    for _, r := range res{
        t.Errorf(r)
    }
    until, _ = UntilFactory([]string{" world", "true"})
    tests = []testCase{
        testCase{ []rune("hello my world, bitch"),
                  []rune("hello my world"),
                  -1,
                  true},
        testCase{ []rune(" hello world"),
                  []rune(" hello world"),
                  -1,
                  true},
        testCase{ []rune("hello"),
                  []rune("hello"),
                  -1,
                  true},
    }
    res = PerformSubparserTests(until, tests)
    for _, r := range res{
        t.Errorf(r)
    }
}



// Quick and dirty benchs... how do the subparsers compare to regexps ?

func TestAnyParserPerfs(t *testing.T){
    any_re := "(?P<any>.*) trailing stuff"
    any, _ := AnyFactory([]string{"", " trailing stuff"})
    ratio := BenchParserPerfs(any, "some stuff trailing stuff", any_re)
    // I benched it as ~70x faster than regexp
    if ratio < 30{
        t.Errorf("Parser too slow. only %d× faster than regexp.", ratio)
    }
}

func TestNonBlankParserPerfs(t *testing.T){
    nonblank_re := "[^ \\t]+"
    ratio := BenchParserPerfs(NonBlankParser, "qwerty", nonblank_re)
    // I benched it as 35x faster than regexp
    if ratio < 15{
        t.Errorf("Parser too slow. only %d× faster than regexp.", ratio)
    }
}

func TestIPV4ParserPerfs(t *testing.T){
    ipv4_re := "(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)"
    // this one is slower
    //ipv4_re = "(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)"
    ratio := BenchParserPerfs(IPV4Parser, "123.123.123.123", ipv4_re)
    // I benched it as ~300x faster than regexp
    if ratio < 130{
        t.Errorf("Parser too slow. only %d× faster than regexp.", ratio)
    }
}

