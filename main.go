// go run -tags etcd main.go
//  完全没搞懂为什么要这样启动？？？

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	metrics "github.com/rcrowley/go-metrics"
	"github.com/smallnest/rpcx/protocol"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/serverplugin"
)

type Args struct {
	A int
	B int
}

type Reply struct {
	C int
}

type Arith int

type Arith2 int

func (t *Arith2) Mul(ctx context.Context, args *Args, reply *Reply) error {
	reply.C = args.A * args.B * 1
	return nil
}

func (t *Arith) Mul(ctx context.Context, args *Args, reply *Reply) error {
	reply.C = args.A * args.B
	return nil
}

func (t *Arith) Add(ctx context.Context, args *Args, reply *Reply) error {
	reply.C = args.A + args.B
	return nil
}

var (
	addr     = flag.String("addr", "localhost:8972", "server address")
	etcdAddr = flag.String("etcdAddr", "localhost:2379", "etcd address")
	basePath = flag.String("base", "/rpcx_test", "prefix path")
)

func main() {

	// s1 := server.NewServer()
	// s1.RegisterName("Arith", new(Arith), "")
	// go s1.Serve("tcp", "127.0.0.1:8972")

	s2 := server.NewServer()
	addRegistryPlugin(s2)
	addYouJinPlugin(s2) /// test
	s2.RegisterName("Arith", new(Arith2), "")
	s2.Serve("tcp", *addr)
	defer s2.Close()
}

func addRegistryPlugin(s *server.Server) {
	r := &serverplugin.EtcdRegisterPlugin{
		ServiceAddress: "tcp@" + *addr,
		EtcdServers:    []string{*etcdAddr},
		BasePath:       *basePath,
		Metrics:        metrics.NewRegistry(),
		// UpdateInterval: time.Minute,
		UpdateInterval: time.Second,
	}
	err := r.Start()
	if err != nil {
		log.Fatal(err)
	}
	s.Plugins.Add(r)
}

// test
func addYouJinPlugin(s *server.Server) {
	r := &YouJinPlugin{}
	s.Plugins.Add(r)
}

type YouJinPlugin struct {
}

func (youjin *YouJinPlugin) PostReadRequest(ctx context.Context, r *protocol.Message, e error) error {
	fmt.Println("---  接收到消息  ---")
	// if rand.Seed(time.Now().UnixNano()); rand.Intn(10)%2 == 0 {
	// 	fmt.Println("===  返回错误!  ===")
	// 	return errors.New("不行不行啊~~~")
	// }
	return nil
}
