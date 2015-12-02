# beego 的 http server...
在工作中有用到golang，后来遇到了beego 重构了一下我的应用。感觉棒棒的~ 应用强壮了不少。所以我打算以最新的stable v1.5.0 来剖析下beego的源代码，因为知其然知其所以然.我们才能更好的使用beego ，同时提高我们的golang能力 。加我的群吧 一起学习 群号:511634754


##* beego 的源代码里，为了让大家顺畅阅读更好理解。我会删掉那些本节我没分析到的代码。

####这节我们先从 beego 的 http server 说起。我们会抽丝剥茧，让一个最简单的的 "beego" 跑起来。最基础最核心其实就两样东西.  http.Server  和  Handler   他们的关系~就类似 麦当劳 和 麦当劳里的服务员妹子。 先有了麦当劳..然后 你给钱妹子。。妹子给你冰淇淋。。也就  请求和输出..

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
		//  HttpAddr  和  HttpPort 这两个变量是在 config.go 里定义的 大写字母开头哦~ 是全局变量
		//  记录了http server 要绑定的地址和 端口号,方便在其他模块进行调用

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

##### config.go文件
```javascript
var (
	BeeApp                 *App
```

同时在 init 函数里 里进行初始化

```javascript
// import的时候其实是执行了该包里面的init函数。。应该懂吧
//执行了 NewApp() 函数~ 那 NewApp 函数执行了啥呢。。我们一层一层的脱了她的衣服....接着看
func init() {

	BeeApp = NewApp()
```

NewApp 函数定义在 app.go 文件里。 看代码...

##### app.go文件
```javascript
type App struct {
	Server   *http.Server  // 这个就是 麦当劳..
	Handlers *ControllerRegistor  //这个就是麦当劳里面的服务员妹子。具体代码请看后面
}

func NewApp() *App {

	cr := NewControllerRegister() //new 一个 Handlers，方便用来处理 http服务 的输入和输出

	//初始化了 App结构。http.Server  这个就是golang自带的http server了~ 太好理解了。这样
	//app.Server 就是 等于  http.Server了~ 如果我们不用 beego 的时候。 写个http server 是不是直接调用 http.Server.ListenAndServe()就很容易实现一个 类似nginx 的基础http 服务器 。 那可能有的人说  http.ListenAndServe() 这样就可以啦。嗯 不过看下源代码就知道  http.ListenAndServe 其实 也是调用更底层的  http.Server.ListenAndServe 。 beego 为了灵活性所以用更底层的 http.Server.ListenAndServe

	app := &App{Handlers: cr,Server: &http.Server{}}
	return app
}

func (app *App) Run() {

	endRunning := make(chan bool, 1)

	go func() {

		// 组装好绑定地址和端口
		addr := fmt.Sprintf("%s:%d", HttpAddr, HttpPort)

		//前面我们说了 app.Server 其实就是 http.Server。 那么我们看下官网手册 http.Server 这个结构里有啥
		//type Server struct {
   	 	//Addr           string        // TCP address to listen on, ":http" if empty
    	//Handler        Handler       // handler to invoke, http.DefaultServeMux if nil
    	//ReadTimeout    time.Duration // maximum duration before timing out read of the request
    	//WriteTimeout   time.Duration // maximum duration before timing out write of the response
    	//MaxHeaderBytes int           // maximum size of request headers, DefaultMaxHeaderBytes if 0
    	//TLSConfig      *tls.Config   // optional TLS config, used by ListenAndServeTLS

 		// 我们这里可以只关注。两个变量  Addr 和  Handler
 		// Addr 就是我们要绑定的地址和端口
 		// Handler 就是我们的处理器, GET  POST  PUT 等请求就是需要他接收和输出.. 这么理解吧  http.Server 这个结构就像是麦当劳.. 而 Handler 就是服务员小妹妹，她负责收钱 和给你冰淇淋...

 		//确定要绑定ip和端口
		app.Server.Addr = addr
		//确定这个http容器里负责处理 输入和输出的方法.. 就是那个 麦当劳服务员小妹妹,你给她钱。她给你...
		app.Server.Handler = app.Handlers

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

	//channel 默认是阻塞。 利用这点。阻塞宿主程序。 否则~~ 宿主都退出了 而 go func 里面的程序~自然也就不存在了
	<-endRunning
}
```


上面是创建了 http 服务容器。接下来就是 接待输入和输出的自定义方法。beego是怎么设计的呢。只有一个目的就是实现ServeHTTP 这个方法。只有实现了这个方法，那么才符合 app.Server.Handler = app.Handlers。我们看下官网手册 http 这节对于Handler的定义。

```javascript
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}
```

其实就是一个接口。里面只有一个方法.. 你懂了吗？ interface 你可以理解为Handler发动机设计图. 无论是歼20还是歼16，只要按照图纸做出来的发动机，当然里面的细节可以根据具体的战斗机需要进行调整(实现了 ServeHTTP这个方法)..我们都可以说 这款战机的发动机是 Handler发动机 -_-! 讲了好多废话..



##### router.go 
```javascript
//beego 从命名中可以get到~~  这个是控制器注册器
type ControllerRegistor struct {

}


//模仿 new  go 例牌
func NewControllerRegister() *ControllerRegistor {
	return &ControllerRegistor{}
}

//就是这个东西。http容器里的输入输出我们如何把玩？就看你怎么实现 ServeHTTP 这个方法了。
func (p *ControllerRegistor) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

	//在这里。你可以自由的发挥了。比如  写你自己的 路由算法 等等等...

	//这段是我自己加的。为了能输出清晰点。。。
	rw.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header
    rw.WriteHeader(http.StatusOK)
    io.WriteString(rw,"加我的群 golang 一起学习 群号:511634754")
}
```

###  记得跑一下代码。在整体理解下哦。完整可运行代码，请看 nixuehan.go   我会一步一步带着大家实现一个 "beego 框架" 而且只要一个文件 



