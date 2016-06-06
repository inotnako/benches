# benches

some bench on golang

# Results

```
$ go test -bench=. -benchtime=5s -count=3

  BenchmarkGrpc-4 	   30000	    265268 ns/op // like 6000 req/sec
  BenchmarkGrpc-4 	   30000	    253664 ns/op
  BenchmarkGrpc-4 	   30000	    233668 ns/op
  BenchmarkHttp-4 	   50000	    192015 ns/op  // like 10000 req/sec
  BenchmarkHttp-4 	   50000	    193716 ns/op
  BenchmarkHttp-4 	   50000	    192249 ns/op
  BenchmarkHttp2-4	   30000	    206757 ns/op // like 6000 req/sec
  BenchmarkHttp2-4	   30000	    206644 ns/op
  BenchmarkHttp2-4	   30000	    210661 ns/op
  ok  	github.com/antonikonovalov/benches/grpc/client	90.587s

$ psql/benches=# select count(*) from test_messages;
  count
---------
 420909
(1 row)


```



