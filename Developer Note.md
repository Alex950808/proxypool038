# Developer Notes

## 处理订阅源：getter
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

接口Getter与err状态组成一个creator。
为了方便外部程序辨认creator类型，在init()中初始化一个map，key为sourceType字符串，value为creator。

程序运行时，package app由配置文件读取到source.yaml，由sourceType map到对应的creator类型，同时使用sourceOption(通常是url)初始化一个creator。

所有Getter最后存于package app的Getters中。

## Database
