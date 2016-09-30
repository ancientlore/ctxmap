/*
Package ctxmap implements a registry for global context.Context for use in web applications.

Based on work from github.com/gorilla/context, this package simplifies the storage by mapping
a pointer to an http.Request to a context.Context. This allows applications to use Google's
standard context mechanism to pass state around their web applications, while sticking to
the standard http.HandlerFunc implementation for their middleware implementations.

As a result of the simplification, the runtime overhead of the package is reduced by 30 to 40
percent in my tests. However, it would be common to store a map of values or a pointer to
a structure in the Context object, and my testing does not account for time taken beyond
calling Context.Value().
*/
package ctxmap

import (
	"context"
	"net/http"
	"sync"
)

/*
Based on work from: https://github.com/gorilla/context

	Copyright 2012 The Gorilla Authors. All rights reserved.
	Use of this source code is governed by a BSD-style
	license that can be found in the LICENSE file.

The code was simplified to only store a context.Context,
which allows for better thread performance when looking
up values.
*/

var (
	mutex sync.RWMutex
	data  = make(map[*http.Request]context.Context)
)

// Set stores a context for a given request.
func Set(r *http.Request, ctx context.Context) {
	mutex.Lock()
	data[r] = ctx
	mutex.Unlock()
}

// Get returns a context for given request.
func Get(r *http.Request) context.Context {
	mutex.RLock()
	ctx := data[r]
	mutex.RUnlock()
	return ctx
}

// GetOk returns stored context and presence state like multi-value return of map access.
func GetOk(r *http.Request) (context.Context, bool) {
	mutex.RLock()
	ctx, ok := data[r]
	mutex.RUnlock()
	return ctx, ok
}

// Delete removes the context for a given request.
//
// This is usually called by a handler wrapper to clean up request
// variables at the end of a request lifetime. See ClearHandler().
func Delete(r *http.Request) {
	mutex.Lock()
	delete(data, r)
	mutex.Unlock()
}

// ClearHandler wraps an http.Handler and clears request values at the end
// of a request lifetime.
func ClearHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer Delete(r)
		h.ServeHTTP(w, r)
	})
}
