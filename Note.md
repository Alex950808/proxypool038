# Note

## 处理订阅源：getter类
有关订阅源的package位于pkg/getter。

订阅源的类型为接口Getter，实现Getter至少需要实现Get()和Get2chan()。 
- Get() 返回一个ProxyList
- Get2chan() 还没研究

多态：已实现的Getter（以sourceType命名）
- subscribe（该实现比接口Getter多了个url）
- tgchannel
- web_fanqiangdang
- web_fuzz
- web_fuzz_sub

接口Getter与err状态组成一个creator，方便错误处理。
为了方便外部程序辨认creator类型，在init()中初始化一个map，key为sourceType字符串，value为creator。

程序运行时，package app由配置文件读取到source.yaml，由sourceType map到对应的creator类型，同时使用sourceOption(通常是url)初始化一个creator。

所有Getter最后存于package app的Getters中。

## proxy类
节点的接口为interface proxy，由struct Base实现其基类，Vmess等实现多态。

Proxylist是proxy数组加上一系列批量处理proxy的方法。



## 抓取
task.go的Crawl.go实现抓取。

1. 并发抓取订阅源，加载历史节点
2. 节点去重，去除Clash不支持的类型，重命名
3. 存储所有节点（包括不可用节点）到database和cache
4. 检测IP可用性  
  尽管已经对IP的有效性测试过，但并不保证节点在客户端上可用，因为可能其他参数有误。
5. 存储可用的节点到cache

## 存储
所有节点存储到cache中，key为all_proxies。问题是对于失效的节点也存储，运行时间久了无用的cache会非常多。可以考虑删除对失效节点的存放。

如果配置文件填写了database url会存储到database中。

## Web界面
还没看
