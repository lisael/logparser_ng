package main

import (
    pconf "logparser_ng/config"
    "time"
    "testing"
)

func TestParser(t *testing.T){
    p, err := pconf.MakeParser("|ip:ipv4()| - - [|date:until(\"]\",false)|] |_ignore| |url| |_ignore| |http_code| |_ignore| |_ignore| \"|user_agent:until('\"',false)|\"")
    if err != nil{
        t.Errorf(err.Error())
        return
    }
	s := "42.42.42.42 - - [03/Jan/2014:06:25:33 +0100] \"GET http://www.example.com/stuff.html HTTP/1.1\" 302 26 \"-\" \"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)\""
    r := []rune(s)
    result, _ := p.ParseOnce(r)
    ip := "42.42.42.42"
    if string(result["ip"]) != ip {
        t.Errorf("Found ip=`%s`. Expected `%s`", string(result["ip"]), ip)
    }
    date := "03/Jan/2014:06:25:33 +0100"
    if string(result["date"]) != date {
        t.Errorf("Found date=`%s`. Expected `%s`", string(result["date"]), date)
    }
    url := "http://www.example.com/stuff.html"
    if string(result["url"]) != url {
        t.Errorf("Found url=`%s`. Expected `%s`", string(result["url"]), url)
    }
    code := "302"
    if string(result["http_code"]) != code {
        t.Errorf("Found http_code=`%s`. Expected `%s`", string(result["http_code"]), code)
    }
    UA := "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"
    if string(result["user_agent"]) != UA {
        t.Errorf("Found user_agent=`%s`. Expected `%s`", string(result["user_agent"]), UA)
    }
}

func TestParserPerfs(t *testing.T){
    p, err := pconf.MakeParser("|ip:ipv4()| - - [|date:until(\"]\",false)|] |_ignore| |url| |_ignore| |http_code| |_ignore| |_ignore| \"|user_agent:until('\"',false)|\"")
    if err != nil{
        t.Errorf(err.Error())
        return
    }
	s := "42.42.42.42 - - [03/Jan/2014:06:25:33 +0100] \"GET http://www.example.com/stuff.html HTTP/1.1\" 302 26 \"-\" \"Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)\""
    r := []rune(s)
    ts1 := time.Now()
	const rounds = 200000
	for i := 0; i < rounds; i++ {
        _, _ = p.ParseOnce(r)
	}
    ts2 := time.Now()
    per_round := ts2.Sub(ts1).Nanoseconds() / rounds
    // on my 5 yo 1 core 32bits 1,5GHz laptop it doesn't exceed 3600ns
    // Aerrmm... Please donate.
    if per_round > 5000 {
        t.Errorf("Parsing is too slow (%dns/line)", per_round)
    }
}
