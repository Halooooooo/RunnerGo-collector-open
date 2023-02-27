package main

import (
	"flag"
	"github.com/Runner-Go-Team/RunnerGo-collector-open/internal"
	"github.com/Runner-Go-Team/RunnerGo-collector-open/internal/pkg/conf"
	log2 "github.com/Runner-Go-Team/RunnerGo-collector-open/internal/pkg/log"
	"github.com/Runner-Go-Team/RunnerGo-collector-open/internal/pkg/server"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var mode int
var configFile string

func main() {

	flag.IntVar(&mode, "m", 0, "读取环境变量还是读取配置文件")
	flag.StringVar(&configFile, "c", "./dev.yaml", "配置文件")
	if !flag.Parsed() {
		flag.Parse()
	}
	internal.InitProjects(mode, configFile)

	runtime.GOMAXPROCS(runtime.NumCPU())

	collectorService := &http.Server{
		Addr: conf.Conf.Http.Host,
	}
	go server.Execute(conf.Conf.Kafka.Host)

	go func() {
		if err := collectorService.ListenAndServe(); err != nil {
			log2.Logger.Error("collector:", err)
			return
		}
	}()

	/// 接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log2.Logger.Info("注销成功")

}
