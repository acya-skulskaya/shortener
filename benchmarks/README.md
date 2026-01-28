Закомментировала тесты с удалением и получением единичного URL, т.к. запуск командой `$ go test -bench . -benchtime=100s` бесконечно висит

Если запускать в аргументом `-benchtime=10000x`, а не `-benchtime=100s`, то тесты проходят, хотя и долго

```shell
$ go test -bench . -benchmem -benchtime=10000x
goos: darwin
goarch: arm64
pkg: github.com/acya-skulskaya/shortener/benchmarks
cpu: Apple M2 Pro
BenchmarkDeleteUserURLs-10         10000               108.3 ns/op             0 B/op          0 allocs/op
BenchmarkGet-10                    10000               125.4 ns/op             0 B/op          0 allocs/op
BenchmarkGetUserURLs-10            10000           1993318 ns/op               0 B/op          0 allocs/op
BenchmarkStoreBatch-10             10000           9612069 ns/op            1536 B/op         29 allocs/op
BenchmarkStore-10                  10000           2585892 ns/op             256 B/op          7 allocs/op
PASS
ok      github.com/acya-skulskaya/shortener/benchmarks  218.993s

```