============
Logparser-NG
============

Go rewrite of Botify versatile LogParser. Improve usability and paralelism
with a clean codebase

Features
========

TODO

Usage
=====

TODO

Example
+++++++

``$bin/logparser -c "coucou |ignore||tok2:plop(1,2)| hello"``

Build
=====

Lacks a Makefile at the moment, so:

  - make a clean go env (/bin and /src)

  - clone the project in /src

  - ``$go get``

  - ``$ragel -Z -G2 -o src/logparser_ng/config/config.go src/logparser_ng/config/config.rl``

  - ``$go install logparser_ng/logparser``

  - ``bin/logparser -h``

