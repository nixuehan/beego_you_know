package main

import (
	"strconv"
	"strings"
	"net/http"
	"net"
	"io"
	"fmt"
)

var (
	HttpAddr string
	HttpPort int
	BeeApp   *App
)

type App struct {
	Handlers *ControllerRegistor
	Server   *http.Server
}

func NewApp() *App {
	cr := NewControllerRegister()
	app := &App{Handlers: cr,Server: &http.Server{}}
	return app
}

func (app *App) Run() {

	endRunning := make(chan bool, 1)

	addr := fmt.Sprintf("%s:%d", HttpAddr, HttpPort)
	app.Server.Addr = addr
	app.Server.Handler = app.Handlers

	go func() {
		ln, err := net.Listen("tcp4",app.Server.Addr)
		if err != nil {
			endRunning <- true
			return
		}
		err = app.Server.Serve(ln)
		if err != nil {
			endRunning <- true
			return
		}
	}()
	<-endRunning
}


type ControllerRegistor struct {

}

func NewControllerRegister() *ControllerRegistor {
	return &ControllerRegistor{}
}

func (p *ControllerRegistor) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	//在这里。你可以自由的发挥了。比如  写你自己的 路由算法 等等等...

	//这段是我自己加的。为了能输出清晰点。。。
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
    rw.WriteHeader(http.StatusOK)
    io.WriteString(rw,"加我的群 golang 一起学习 群号:511634754")
}


func init() {
	BeeApp = NewApp()
}

func Run(params ...string) {
	if len(params) > 0 && params[0] != "" {
		strs := strings.Split(params[0], ":")
		if len(strs) > 0 && strs[0] != "" {
			HttpAddr = strs[0]
		}
		if len(strs) > 1 && strs[1] != "" {
			HttpPort, _ = strconv.Atoi(strs[1])
		}
	}
	
	BeeApp.Run()
}

func main() {
    Run("192.168.202.5:9413")
}