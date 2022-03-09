package main

import (
	"flag"
	"github.com/cihub/seelog"
	"github.com/yincongcyincong/go12306/action"
	"github.com/yincongcyincong/go12306/utils"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
)

var (
	runType   = flag.String("run_type", "command", "web：网页模式")
	wxrobot   = flag.String("wxrobot", "", "企业微信机器人通知")
	mustDevice = flag.String("must_device", "0", "强制生成设备信息")
)

func main() {
	flag.Parse()
	Init()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	select {
	case <-sigs:
		seelog.Info("用户登出")
		utils.SaveConf()
		action.LoginOut()
	}
}
