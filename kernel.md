# beego 的 http server...
在工作中有用到golang，后来遇到了beego 重构了一下我的应用。感觉棒棒的~ 应用强壮了不少。所以我打算以最新的stable v1.5.0 来剖析下beego的源代码，因为知其然知其所以然.我们才能更好的使用beego ，同时提高我们的golang能力 。


##* beego 的源代码里，为了让大家顺畅阅读更好理解。我会删掉那些本节我没分析到的代码。

###这节我们先从 beego 的 http server 说起。我们会抽丝剥茧，让一个最简单的的 "beego" 跑起来。

这段是官网贴的启动beego的代码。 那我们就从 Run 这个函数开始

```javascript
package main

import "github.com/astaxie/beego"

func main() {
    beego.Run()
}
```

##### beego.go文件

```javascript

// Run 函数 有一个可变参数  ...string 相当于 []string{}  这个不懂..那就别继续往下看了~先复习golang基础
// 此参数是告知 http server 要绑定的 host 和 port 。 例如: 127.0.0.1:9413

func Run(params ...string) {
	
	if len(params) > 0 && params[0] != "" {

		// -_- 看看人家是怎么把 地址和端口拆到两个变量里去的...
		//  HttpAddr  和  HttpPort 这两个变量是在 config.go 里定义的 大写字母开头哦~ 是全局变量, 记录了http server 要绑定的地址和 端口号 

		strs := strings.Split(params[0], ":")
		if len(strs) > 0 && strs[0] != "" {
			HttpAddr = strs[0]
		}
		if len(strs) > 1 && strs[1] != "" {
			HttpPort, _ = strconv.Atoi(strs[1])
		}
	}

	//要绑定的地址和端口都得到了之后。正式启动。 调用BeeApp 里面的方法 Run 
	BeeApp.Run()
}
```

BeeApp 这个结构很重要。很多底层的都封装在这个结构里。我们看下 它的 Run 是啥东西。 BeeApp 是一个指向 App结构的变量。 在 config.go 文件里定义

```javascript
var (
	BeeApp                 *App
```

同时在 init 里进行初始化

```javascript

// import的时候其实是执行了该包里面的init函数。。应该懂吧
//执行了 NewApp() 函数~ 那 NewApp 函数执行了啥呢。。接着看
func init() {

	BeeApp = NewApp()
```

NewApp 函数定义在 app.go 文件里。 看代码...
```javascript

type App struct {
	
	Server   *http.Server
}



func NewApp() *App {

	//初始化了 App结构。http.Server  这个就是golang自带的http server了~ 太好理解了。这样
	//app.Server 就是 等于  http.Server了~ 如果我们不用 beego 的时候。 写个http server 是不是直接调用 http.Server.ListenAndServe()就很容易实现一个 类似nginx 的基础http 服务器 。 那可能有的人说  http.ListenAndServe() 这样就可以啦。嗯 不过看下源代码就知道  http.ListenAndServe 其实 也是调用更底层的  http.Server.ListenAndServe 。 beego 为了灵活性所以用更底层的 http.Server.ListenAndServe

	app := &App{Server: &http.Server{}}
	return app
}

func (app *App) Run() {

	//启动完成chan.. 方便宿主知道 http server 启动完成。。然后他想改嘛就干嘛...
	endRunning := make(chan bool, 1)

	go func() {
		// app.Server 前面说了~其实就是 http.Server 。所以 app.Server.Addr ..就是 app.Server结构里的一个变量.存地址的。app.Server.Addr里面存的就是前面我们Run时候传进来的地址 。前面有段代码 app.Server.Addr = HttpAddr  

		//下面就是标准的启动一个 golang http server 的流程了...
		ln, err := net.Listen("tcp4", app.Server.Addr)
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
}
```







