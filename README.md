
gosp - Go Server Pages
--------

It's the "templating" language you love to hate. Embedded go in the
spirit of JSP/ASP/PHP/perl's Text::Template. Go Server Pages compile
ahead of time to regular go code, control structures, and print statements.
No reflection needed. Is this project "serious"? Not really but I might
as well share.


Advantages
==========

* Output is streamable (Eventually to any io.Writer).
* No hoops to jump through for registering functions, no new syntax for calling them.
* Deployment is as easy as deploying any other go code/package.
* Performance should be pretty good. Generated assembly is as tight as the go compiler can make it.

Disadvantages
=============

* Ugly syntax
* Doesn't integrate with go build and friends
* Go stdlib has text/template which is more elegant and builtin

Example
=======

    go run compiler.go < test.tmpl > out.go
    go run out.go


Bugs
====

* Output is currently hard coded to stdout. If anyone actually wants to use this I'll fix it.

