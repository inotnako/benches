# Result

```

go test -bench=. -benchtime=1s ./kv
BenchmarkKV_Set-4        1000000              2310 ns/op
BenchmarkKVR_Get-4      20000000                69.3 ns/op
BenchmarkKV_Get-4       10000000               150 ns/op

```