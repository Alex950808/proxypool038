# Note 

## 处理订阅源：getter类
有关订阅源的package位于pkg/getter。

订阅源的类型为接口Getter，实现Getter至少需要实现Get()和Get2chan()。 
- Get() 返回一个ProxyList
- Get2chan() 还没研究

已实现的Getter（以sourceType命名）
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

所有字段名依据clash的配置文件标准设计。比如
```
type ShadowsocksR struct {
	Base          // 节点基本信息
	Password      string `yaml:"password" json:"password"`
	Cipher        string `yaml:"cipher" json:"cipher"`
	Protocol      string `yaml:"protocol" json:"protocol"`
	ProtocolParam string `yaml:"protocol-param,omitempty" json:"protocol_param,omitempty"`
	Obfs          string `yaml:"obfs" json:"obfs"`
	ObfsParam     string `yaml:"obfs-param,omitempty" json:"obfs_param,omitempty"`
	Group         string `yaml:"group,omitempty" json:"group,omitempty"`
}
```

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
所有节点存储到cache中。

cache中的key设计有：
- allproxies: 所有节点（包括不可用节点）
- proxies: 可用节点
- clashproxies: clash支持的节点。第一次运行时是把proxies复制过来的。

问题是对于失效的节点也存储，运行时间久了无用的cache会非常多。可以考虑删除对失效节点的存放。

如果配置文件填写了database url会存储到database中。Database连接参数自行修改源码。  
目前database数据只在抓取时使用，不用于处理Web请求。

### 使用数据库
安装postgresql，建立相应user和database。

```
	dsn := "user=proxypool password=proxypool dbname=proxypool port=5432 sslmode=disable TimeZone=Asia/Shanghai"
```

程序运行时建立会proxies表。
每次运行时读出节点，爬虫完成后再存储可用的节点进去。

直接暴力删除整表再全部写入，因为可用的节点注定不会多，不用担心性能问题。

## Web界面

静态的assets文件模板由zip压缩后存为字符串的形式，如

```
var _assetsHtmlSurgeHtml="[]byte("\x1f\x8b\x...")"
```

以上字节解压后是一个go的HTML模板。解压时，由gzip的reader写入byte.Buffer，再转换为Bytes写入相应文件。

因此想修改html文件请写好后自行压缩并替换字节。

不知道原作者为何一定要把模板设计为html generator，毕竟感觉直接使用html file也没什么问题。我猜可能是为了增加不可读性吧，免得有人用了这个项目不标明出处。那其实也不必这样做，直接在最后注入原项目地址就成了。

## 本地测试
需要注意：
- 修改了config的domain
- 修改了source，注释掉较慢的源

增加了对config-local文件的解析。
bindata中增加了clash-config-local.yaml字段，也增加了模板。