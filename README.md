
gosp - Go Server Pages
--------

Gosp is a strongly typed templating language. Put the bug finding abilities
and speed of the golang compiler to work for you. It eschews reflection
and the  el (expression language) approach taken by the go stdlib in favor
of leaning on go's type inference and compile time bug finding.


Features
==========

* typed page parameters (no request.setParam thank goodness)
* Compile time type checking on your templates!
* Output is streamable to any io.Writer.
* No hoops to jump through for registering functions, no new syntax for calling them.
* Deployment is as easy as deploying any other go code/package.
* Performance is largely determined by what the go runtime can provide.
* Tested


Disadvantages
=============

* Ugly syntax
* Doesn't integrate with go build and friends (makefile's required or gospc *gosp)
* Go stdlib has text/template which is more elegant and builtin

Example
=======

    cat > template/tmpl.gosp
    @(import "runtime")
    hello, my name is shane.
    This file was compiled using <%=runtime.Compiler%>
    <%for i:=0; i<10; i++ {%>
    Yo <%=i%>
    <%}%>
    done
    Control-D

    gospc  template/*gosp
    

Docs
====

[godoc](http://godoc.org/github.com/shanemhansen/gosp)

A plain text file is a valid gosp program. Go code can be embedded
using the well-known <% %> syntax. Go values can be printed (through the fmt package)
using the well-known <%= %> syntax. All lines at the beginning of the file starting with '@'
are processed for import statements ( @(import "package") ) and parameters to pass to
the generated go function ( @(param1 string, param2 int) ).

The generated go code is relatively clean. Something like:

    func(param1 string, param2 int) { return func(out io.Writer) { fmt.Fprintf(param1) } }

The gospc compiler takes '.gosp' files given on the command line and generated .go files
next to them. The package name is derived from the foldername the gosp file is in. The files
are intended to live next to each other. The go filename is CamelCased to generate the template. The idea
is that related templates will live in a package.

Calling the generated template is simple and elegant.

    template.Test("param1", 2)(os.Stdout)

Composing templates (for pages with layouts) is similarly easy

    template.Master(template.Slave("shane"))(os.Stdout)

The first call binds the parameters to the function returning a func(io.Writer) which
is a little easier to pass around.


FAQ
===

Q. Are gosp files a substitute for php or asp files?
A. Not really, It's just a template. Executing that template in response
to a request is a web framework's job.

Q. What's the performance?
A. We built gosp pages with an eye towards performance, and we suspect they'd
perform just as well as hand written output code, but we haven't submitted our code to any benchmarks.
The only abstraction we force on you is that you have to write your data to an interface (io.Writer).

Q. What's the tooling/feedback cycle?
A. We currently run the gospc compiler on save and use emacs's wonderful flymake mode to find problems in the generated go files.
Eventually I'd like to use source maps like cgo does, but that would require parsing the go code we're outputting to determine
appropriate places to emit line number declarations, right now we don't parse the go code.

Q. How does this handle things a "big" templating framework might need, like themes, asset pipelines, etc.
A. The plan is to build all that fuctionality into libraries, orthonogonal but designed to be pluggable with gosp.

Q. Why such a sucky name?
A. Naming things is hard.

Q. Can I call gospc from a library?
A. Sure, check out the gosp.Compiler function

Thanks
=====

The import and declaration syntax was stolen from the scala play
framework. http://www.playframework.com/documentation/2.0/ScalaTemplates
