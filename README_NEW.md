Clash客户端支持：
- Clash for Windows（需要Clash Core1.3以上）
- ClashX（需要Clash Core1.3以上）
- 不支持ClashXR与ClashR等非原生Clash Core客户端。

TODO:
 
-[ ] 遇到无节点的provider，Clash会报错。
    还是减少一下不常用的国家吧。可以用动态生成config解决，~~但懒得写~~。
    提供了clash-config-local模板，可以自行替换域名。
    之后考虑添加缺省字段DIRECT解决这个问题。
    

## New

2020-10-09
1. 增加本地http运行的支持  
    > clash的本地配置文件位于clash/localconfig
