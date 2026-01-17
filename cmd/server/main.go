package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"reflect"

	log "github.com/golang/glog"
	"google.golang.org/grpc"

	pb "github.com/openconfig/gnmi/proto/gnmi"
	gnmi "github.com/zhang-xingyun/gnmi_test/gnxi/gnmi"
	"github.com/zhang-xingyun/gnmi_test/gnxi/gnmi/modeldata"
	gpb "github.com/zhang-xingyun/gnmi_test/gnxi/gnmi/modeldata/gostruct"
)

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "错误: 必须指定一个监听地址作为参数")
		fmt.Fprintln(os.Stderr, "用法: go run cmd/server/main.go :5000")
		os.Exit(1)
	}

	address := args[0]

	// 创建 gRPC 服务器
	g := grpc.NewServer()

	// 创建模型实例
	model := gnmi.NewModel(
		modeldata.ModelData,
		reflect.TypeOf((*gpb.Device)(nil)),
		gpb.SchemaTree["Device"],
		gpb.Unmarshal,
		gpb.ΛEnumTypes,
	)

	// 创建空配置
	var config []byte

	// 创建 gNMI 服务器实例，不设置回调函数
	s, err := gnmi.NewServer(model, config, nil)
	if err != nil {
		log.Exitf("Failed to create server: %v", err)
	}

	// 注册 gNMI 服务
	pb.RegisterGNMIServer(g, s)

	// 启动监听
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Exitf("Failed to listen on address %s: %v", address, err)
	}
	defer listener.Close()

	fmt.Printf("gNMI server listening on %s\n", address)
	if err := g.Serve(listener); err != nil {
		log.Exitf("Failed to serve gNMI server: %v", err)
	}
}
