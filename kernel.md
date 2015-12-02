# beego 的 http server...
在工作中有用到golang，后来遇到了beego 重构了一下我的应用。感觉棒棒的~ 应用强壮了不少。所以我打算剖析下beego的源代码...因为知其然知其所以然.我们才能更好的使用beego ，同时提高我们的golang能力 。

```golang
package main

import "github.com/astaxie/beego"

func main() {
    beego.Run()
}
```