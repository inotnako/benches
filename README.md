# benches

some bench on golang



# Results

 Mac OS; 2,8 GHz Intel Core i5; 8 GB 1600 МGHz DDR3; Fusion Drive

```
$ go get -v -u ./...
$ make build
// set ulimit -n 10000 for each console
$  ulimit -n 10000
// run services
$ bin/srvgrpc , bin/srvpq, bin/srvpq2

$ go test -bench=. -benchtime=5s -cpu=4

  BenchmarkGrpc-4 	   30000	    265268 ns/op // like 6000 req/sec
  BenchmarkHttp-4 	   50000	    192015 ns/op  // like 10000 req/sec
  BenchmarkHttp2-4	   30000	    206757 ns/op // like 6000 req/sec

$ psql/benches=# select count(*) from test_messages;
  count
---------
 420909
(1 row)


```



