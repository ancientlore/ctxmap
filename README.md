ctxmap
=======

[![Go Reference](https://pkg.go.dev/badge/github.com/ancientlore/ctxmap.svg)](https://pkg.go.dev/github.com/ancientlore/ctxmap)

Package ctxmap implements a registry for global context.Context for use in web applications.

Based on work from github.com/gorilla/context, this package simplifies the storage by mapping
a pointer to an http.Request to a context.Context. This allows applications to use Google's
standard context mechanism to pass state around their web applications, while sticking to
the standard http.HandlerFunc implementation for their middleware implementations.

As a result of the simplification, the runtime overhead of the package is reduced by 35 to 50
percent in my tests. However, it would be common to store a map of values or a pointer to
a structure in the Context object, and my testing does not account for time taken beyond
calling Context.Value().

| Benchmark                    | Readers | Writers | Iterations | Map Ops | context   | ctxmap    |
|:-----------------------------|--------:|--------:|-----------:|--------:|----------:|----------:|
| BenchmarkMutexSameReadWrite1 | 1       | 1       | 32         | 64      | 208.91 ns | 102.63 ns |
| BenchmarkMutexSameReadWrite2 | 2       | 2       | 32         | 128     | 211.27 ns | 103.97 ns |
| BenchmarkMutexSameReadWrite4 | 4       | 4       | 32         | 256     | 216.21 ns | 101.54 ns |
| BenchmarkMutex1              | 2       | 8       | 32         | 320     | 252.26 ns | 92.62 ns  |
| BenchmarkMutex2              | 16      | 4       | 64         | 1280    | 166.85 ns | 108.31 ns |
| BenchmarkMutex3              | 1       | 2       | 128        | 384     | 221.67 ns | 78.47 ns  |
| BenchmarkMutex4              | 128     | 32      | 256        | 40960   | 179.70 ns | 107.31 ns |
| BenchmarkMutex5              | 1024    | 2048    | 64         | 196608  | 233.41 ns | 90.10 ns  |
| BenchmarkMutex6              | 2048    | 1024    | 512        | 1572864 | 183.25 ns | 92.33 ns  |
