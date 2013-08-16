
gosp - Go Server Pages
--------

It's the templating language you love to hate. A simple templating
langage in the style of PHP/JSP what compiles straight to go which is then
linked into your executable. No runtime reflection required. If you'd like
to catch more of your templating mistakes at compile time rather than runtime,
gosp might be for you!


Advantages
==========

* Compile time type checking on your templates!
* Output is streamable to any io.Writer.
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

    cat > template/test.gosp
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

Work in progress. A plain text file is a valid gosp program. Go code can be embedded
using the well-known <% %> syntax. Go values can be printed (through the fmt package)
using the well-known <%= %> syntax. All lines at the beginning of the file starting with '@'
are processed for import statements ( @(import "package") ) and parameters to pass to
the generated go function ( @(param1 string, param2 int) ).

The gospc compiler takes '.gosp' files given on the command line and generated .go files
next to them. The package name currently defaults to 'template'. You should use the -package
flag to set it correctly. The go filename is CamelCased to generate the template. The idea
is that related templates will live in a package.

Calling the generated template is relatively simple.

    template.Test("param1", 2)(os.Stdout)

The first call binds the parameters to the function returning a func(io.Writer) which
is a little easier to pass around.

Thanks
=====

The import and declaration syntax was stolen from the scala play
framework. http://www.playframework.com/documentation/2.0/ScalaTemplates

