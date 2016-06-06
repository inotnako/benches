
run_srvpq:
	bin/srvpq

build:
	go build -o bin/srvpq http/serverpq/main.go
	go build -o bin/srvpq2 http2/serverpq/main.go
	go build -o bin/srvgrpc grpc/main.go


bench:
	echo "POST http://0.0.0.0:4567/create" | vegeta attack -duration=5s -rate=4500 -workers=100 | tee results.bin | vegeta report -reporter=plot > plot.html
	cat results.bin | vegeta report -reporter='hist[0,2ms,4ms,6ms,30ms,100ms]'
	cat results.bin | vegeta report


bench_http2:
	echo "POST https://0.0.0.0:4568/create" | vegeta attack  -insecure -http2 -duration=5s -rate=4500 -workers=100 | tee results_http2.bin | vegeta report -reporter=plot > plot_http2.html
	cat results_http2.bin | vegeta report -reporter='hist[0,2ms,4ms,6ms,30ms,100ms]'
	cat results_http2.bin | vegeta report

proto:
	protoc --go_out=plugins=grpc:. grpc/proto/*.proto


bench_grpc:
