============
Logparser-NG
============

Go rewrite of Botify versatile LogParser. Improve usability and parallelism
with a clean codebase and easy configuration

Features
========

TODO

Usage
=====

TODO

Example
+++++++

``$bin/logparser -c "|IP:ipv4()| |_ignore| \"|tok2:until('\"',false)| tail"``

Build
=====

Lacks a Makefile at the moment, so:

  - make a clean go env (/bin and /src)

  - clone the project into /src

  - ``$ go get``

  - ``$ ragel -Z -G2 -o src/logparser_ng/config/config.go src/logparser_ng/config/config.rl``

  - ``$ ragel -Z -G2 -o src/logparser_ng/parser/subparsers.go src/logparser_ng/parser/subparsers.rl

  - ``$ go install logparser_ng/logparser``

  - ``$ bin/logparser -h``

It's Open Source
================

Fork It
+++++++

``$ git clone https://github.com/lisael/logparser-ng``

Hack It
+++++++

A parser is basically a list of subparser functions (or closure) called
sequentially while parsing a string. It stores a global vars such as current
read char index in the input.  Once the parsing is completed, the Parser returns
a map[string][]rune which is passed to an output formater

A parser can be built using a config string looking like

``|ip:ipv4()| some stuff |_ignore| [|date:until(],false)|]``

Each token is in the form ``|name:factory(args)|`` (except ``|_ignore|`` which
is a special token and is equivalent to ``|:_ignore()|``)

``name`` is the name of the match in the result map. ``factory`` is the name of
a subparser factory, which has been previously registered using
parser.RergisterFactory(name, factory).

The factory takes a []string containing the passed args. It must return a
parser.Subparser function, or an error. If the error is a
parser.DeferredFactoryDef, the factory will be called only when all the config
is parsed, preceding and following text are append to its args (read
AnyFactory() code).
