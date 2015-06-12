ctxmap
=======

[![Build Status](https://travis-ci.org/ancientlore/ctxmap.svg?branch=master)](https://travis-ci.org/ancientlore/ctxmap)
[![Coverage Status](https://coveralls.io/repos/ancientlore/ctxmap/badge.svg)](https://coveralls.io/r/ancientlore/ctxmap)
[![GoDoc](https://godoc.org/github.com/ancientlore/ctxmap?status.png)](https://godoc.org/github.com/ancientlore/ctxmap)
[![status](https://sourcegraph.com/api/repos/github.com/ancientlore/ctxmap/.badges/status.png)](https://sourcegraph.com/github.com/ancientlore/ctxmap)
[gocover](http://gocover.io/github.com/ancientlore/ctxmap)

Package ctxmap implements a registry for global context.Context for use in web applications.

Based on work from github.com/gorilla/context, this package simplifies the storage by mapping
a pointer to an http.Request to a context.Context. This allows applications to use Google's
standard context mechanism to pass state around their web applications, while sticking to
the standard http.HandlerFunc implementation for their middleware implementations.

As a result of the simplification, the runtime overhead of the package is reduced by 30 to 40
percent in my tests. However, it would be common to store a map of values or a pointer to
a structure in the Context object, and my testing does not account for time taken beyond
calling Context.Value().
