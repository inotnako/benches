# Result

```
go test -bench=. -benchtime=1s ./kv

PASS
BenchmarkKVCtx_Set-4           	 3000000	       419 ns/op
BenchmarkKVCtx_Get-4           	 3000000	       451 ns/op
BenchmarkKVChe_Set-4           	 5000000	       321 ns/op
BenchmarkKVChe_Get-4           	 5000000	       398 ns/op
BenchmarkKVChm_Set-4           	 5000000	       306 ns/op
BenchmarkKVChm_Get-4           	 5000000	       283 ns/op
BenchmarkKV_Set-4              	20000000	       116 ns/op
BenchmarkKV_Get-4              	10000000	       191 ns/op
BenchmarkKV_Channel_Set-4      	 2000000	       669 ns/op
BenchmarkKV_Channel_Get-4      	 2000000	       874 ns/op
BenchmarkKVR_Get-4             	20000000	        75.1 ns/op
BenchmarkMap_Successively_Set-4	 5000000	       438 ns/op
ok  	github.com/antonikonovalov/benches/kv	498.973s

```