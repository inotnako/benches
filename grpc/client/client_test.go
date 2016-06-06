package client

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"testing"

	"crypto/tls"
	pb "github.com/antonikonovalov/benches/grpc/proto"
	"net"
	"net/http"

	"io"
	"io/ioutil"
	"time"
)

func BenchmarkGrpc(b *testing.B) {
	conn, err := grpc.Dial(`0.0.0.0:4569`, grpc.WithInsecure())
	if err != nil {
		b.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewCreatorClient(conn)

	req := &pb.MessageRequest{Msg: `hello grpc!`}
	ctx := context.Background()
	//b.SetParallelism(50)
	b.RunParallel(func(pbt *testing.PB) {
		for pbt.Next() {
			_, err = c.Create(ctx, req)
			if err != nil {
				b.Errorf(`error grpc: %s`, err)
			}
		}
	})

	//
	//for i := 0; i < b.N; i++ {
	//	_, err = c.Create(ctx, req)
	//	if err != nil {
	//		b.Logf(`error grpc: %s`, err)
	//	}
	//}
}

const (
	// DefaultRedirects is the default number of times an Attacker follows
	// redirects.
	DefaultRedirects = 10
	// DefaultTimeout is the default amount of time an Attacker waits for a request
	// before it times out.
	DefaultTimeout = 30 * time.Second
	// DefaultConnections is the default amount of max open idle connections per
	// target host.
	DefaultConnections = 10000
	// DefaultWorkers is the default initial number of workers used to carry an attack.
	DefaultWorkers = 10
	// NoFollow is the value when redirects are not followed but marked successful
	NoFollow = -1
)

var (
	// DefaultLocalAddr is the default local IP address an Attacker uses.
	DefaultLocalAddr = net.IPAddr{IP: net.IPv4zero}
	// DefaultTLSConfig is the default tls.Config an Attacker uses.
	DefaultTLSConfig = &tls.Config{InsecureSkipVerify: true}
)

func BenchmarkHttp(b *testing.B) {
	dialer := &net.Dialer{
		LocalAddr: &net.TCPAddr{IP: DefaultLocalAddr.IP, Zone: DefaultLocalAddr.Zone},
		KeepAlive: 30 * time.Second,
		Timeout:   DefaultTimeout,
	}
	client := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial:  dialer.Dial,
			ResponseHeaderTimeout: DefaultTimeout,
			TLSClientConfig:       DefaultTLSConfig,
			TLSHandshakeTimeout:   10 * time.Second,
			MaxIdleConnsPerHost:   DefaultConnections,
		},
	}

	//b.SetParallelism(100)
	b.RunParallel(func(pbt *testing.PB) {
		for pbt.Next() {
			req, _ := http.NewRequest(`POST`, `http://0.0.0.0:4567/create`, nil)
			resp, err := client.Do(req)
			if err != nil {
				b.Errorf(`error http: %s`, err)
			} else {
				io.Copy(ioutil.Discard, resp.Body)
				resp.Body.Close()
			}
		}
	})
}

func BenchmarkHttp2(b *testing.B) {
	dialer := &net.Dialer{
		LocalAddr: &net.TCPAddr{IP: DefaultLocalAddr.IP, Zone: DefaultLocalAddr.Zone},
		KeepAlive: 30 * time.Second,
		Timeout:   DefaultTimeout,
	}
	client := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial:  dialer.Dial,
			ResponseHeaderTimeout: DefaultTimeout,
			TLSClientConfig:       DefaultTLSConfig,
			TLSHandshakeTimeout:   10 * time.Second,
			MaxIdleConnsPerHost:   DefaultConnections,
		},
	}

	//b.SetParallelism(100)
	b.RunParallel(func(pbt *testing.PB) {
		for pbt.Next() {
			req, _ := http.NewRequest(`POST`, `https://0.0.0.0:4568/create`, nil)

			resp, err := client.Do(req)
			if err != nil {
				b.Errorf(`error http2: %s`, err)
			} else {
				io.Copy(ioutil.Discard, resp.Body)
				resp.Body.Close()
			}
		}
	})
}
