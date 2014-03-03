============
Logparser-NG
============

Go rewrite of Botify versatile LogParser. Improve usability and parallelism
with a clean codebase and easy configuration

Features
========

Provides lib to build custom log parsers and a CLI that is flexible enough
for most common cases

TODO: list all features

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

  - ``$ ragel -Z -G2 -o src/logparser_ng/config/config.go src/logparser_ng/config/config.rl``
  
  - ``$ ragel -Z -G2 -o src/logparser_ng/parser/subparsers.go src/logparser_ng/parser/subparsers.rl``
  
  - ``$ go get``

  - ``$ go install logparser_ng/logparser``

  - ``$ bin/logparser -h``

It's Open Source
================

Fork It
+++++++

``$ git clone https://github.com/lisael/logparser-ng.git``

Hack It
+++++++

Pipeline
--------

Each part of a logparser is a so-called service. A service exposes a 'Process'
method and a method (usually called 'Pipe') that recieves from an input chan,
sends data to its Process method (maybe asynchronously) and returns an output
chan filled with the results.

Logpaser has just to configure a bunch ouf services (reader, parser, filter, 
formater and writer) and pipe them together like so (in logparser/main.go)

.. code-block:: go

  stopChan := writer.Pipe(formater.Pipe(filter.Pipe(parser.Pipe(reader.Read(inputPath)))))
  <- stopChan


The first lines of result are written before the last lines of the input are
read. Hence, parsing 8M lines (1,4GB) takes 2 minutes and less of 200MB of RAM
on a dual core proc.

The pipe line (well, and go itself) allows heavy concurrency while keeping
readable and testable code.

Parser
------

A parser is basically a list of subparser functions (or closure) called
sequentially while parsing a string. Once the parsing is completed, the
Parser returns a map[string][]rune which is passed to an output formater

A parser for lines like::

  42.42.42.42 some stuff 412 [Sun Mar 2 11:37:05 CET 2014]

can be built using a config string looking like::

  |ip:ipv4()| some stuff |_ignore| [|date:until(],false)|]
  
it will return::

  map[string]rune[]{
    "ip": "42.42.42.42",
    "date": "Sun Mar 2 11:37:05 CET 2014"
  }

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

Break it
++++++++

I bet you don't need any advice.

Make it awesome !
+++++++++++++++++

TODO
====

- enable communication from (error handling) and to (processing management:
  Pause/Resume...) pipeline elements
  
- add formaters

- add filters

- count parsing errors and stop when a given ratio is exceeded

- make TextParsers (text outside of tokens) strict about the actual text
  (it only skips the length of the text without checking)

- doc

- update tests (subparsers test don't build any more)
