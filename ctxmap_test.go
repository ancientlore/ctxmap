package ctxmap

import (
	"context"
	"net/http"
	"testing"
)

/*
Based on work from: https://github.com/gorilla/context

	Copyright 2012 The Gorilla Authors. All rights reserved.
	Use of this source code is governed by a BSD-style
	license that can be found in the LICENSE file.
*/

type sessionKey string

func TestContext(t *testing.T) {
	assertEqual := func(val interface{}, exp interface{}) {
		if val != exp {
			t.Errorf("Expected %v, but got %v.", exp, val)
		}
	}

	r, _ := http.NewRequest("GET", "http://localhost:8080/", nil)
	emptyR, _ := http.NewRequest("GET", "http://localhost:8080/", nil)

	// Get()
	assertEqual(Get(r), nil)

	// Set()
	ctx := context.Background()
	Set(r, ctx)
	assertEqual(Get(r), ctx)

	ctx2 := context.WithValue(ctx, sessionKey("session"), make(map[string]interface{}))
	Set(r, ctx2)
	assertEqual(Get(r), ctx2)

	//GetOk
	value, ok := GetOk(r)
	assertEqual(value, ctx2)
	assertEqual(ok, true)

	value, ok = GetOk(emptyR)
	assertEqual(value, nil)
	assertEqual(ok, false)

	Set(emptyR, nil)
	value, ok = GetOk(emptyR)
	assertEqual(value, nil)
	assertEqual(ok, true)

	// Delete()
	Delete(emptyR)
	assertEqual(Get(emptyR), nil)

	Delete(r)
	assertEqual(Get(r), nil)
}

func parallelReader(r *http.Request, iterations int, wait, done chan struct{}) {
	<-wait
	for i := 0; i < iterations; i++ {
		ctx := Get(r)
		// to be fair in our tests and comparing against gorilla/context, let's fetch a value
		if ctx != nil {
			_ = ctx.Value("Foo")
		}
	}
	done <- struct{}{}

}

func parallelWriter(value context.Context, r *http.Request, iterations int, wait, done chan struct{}) {
	<-wait
	for i := 0; i < iterations; i++ {
		Set(r, value)
	}
	done <- struct{}{}

}

func benchmarkMutex(b *testing.B, numReaders, numWriters, iterations int) {

	b.StopTimer()
	reqs := make([]*http.Request, numWriters)
	for i := 0; i < numWriters; i++ {
		reqs[i], _ = http.NewRequest("GET", "http://localhost:8080/", nil)
	}
	done := make(chan struct{})
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		wait := make(chan struct{})

		for i := 0; i < numReaders; i++ {
			go parallelReader(reqs[i%numWriters], iterations, wait, done)
		}

		for i := 0; i < numWriters; i++ {
			// to be fair in our tests and comparing against gorilla/context, let's set  value
			ctx := context.WithValue(context.Background(), sessionKey("Foo"), "1234")
			go parallelWriter(ctx, reqs[i], iterations, wait, done)
		}

		close(wait)

		for i := 0; i < numReaders+numWriters; i++ {
			<-done
		}
	}
}

// 1 reader, 1 writer, 32 iteratons each = 64 operations
func BenchmarkMutexSameReadWrite1(b *testing.B) {
	benchmarkMutex(b, 1, 1, 32)
}

// 2 readers, 2 writers, 32 iterations each = 128 operations
func BenchmarkMutexSameReadWrite2(b *testing.B) {
	benchmarkMutex(b, 2, 2, 32)
}

// 4 readers, 4 writers, 32 iterations each = 256 operations
func BenchmarkMutexSameReadWrite4(b *testing.B) {
	benchmarkMutex(b, 4, 4, 32)
}

// 2 readers, 8 writers, 32 iterations each = 320 operations
func BenchmarkMutex1(b *testing.B) {
	benchmarkMutex(b, 2, 8, 32)
}

// 16 readers, 4 writers, 64 iterations each = 1280 operations
func BenchmarkMutex2(b *testing.B) {
	benchmarkMutex(b, 16, 4, 64)
}

// 1 reader, 2 writers, 128 iterations each = 384 operations
func BenchmarkMutex3(b *testing.B) {
	benchmarkMutex(b, 1, 2, 128)
}

// 128 readers, 32 writers, 256 iterations each = 40960 operations
func BenchmarkMutex4(b *testing.B) {
	benchmarkMutex(b, 128, 32, 256)
}

// 1024 readers, 2048 writers, 64 iterations each = 196608 operations
func BenchmarkMutex5(b *testing.B) {
	benchmarkMutex(b, 1024, 2048, 64)
}

// 2048 readers, 1024 writers, 512 iterations each = 1572864 operations
func BenchmarkMutex6(b *testing.B) {
	benchmarkMutex(b, 2048, 1024, 512)
}
