# 分布式缓存 TinyGroupCache

## 使用分布式缓存系统的目的：

​	提高了整个项目的并发量、响应速度。

## 使用了分布式缓存系统之后要面对的新问题：

- 缓存穿透（缓存和数据库中都没有的数据）

- 缓存击穿（热点数据存储到期，大量并发请求直接访问数据库）

- 缓存雪崩（大量的缓存在同一时间集体失效，大量的查询直接透传到数据库层面）
- 缓存一致性

​	详细问题描述和可行解决方案：

​		[分布式缓存系统必须要解决的四大问题-阿里云开发者社区 (aliyun.com)](https://developer.aliyun.com/article/1009128)

​		[分布式缓存面临的常见问题及其解决方案_分布式缓存解决方案-CSDN博客](https://blog.csdn.net/c15158032319/article/details/117848048)

## TinyGroupCache原型：

​	 **groupcache**（对其功能进行了裁剪）

## 预期功能：

​	资源控制、淘汰策略、并发、分布式节点通信

## 支持特性：

- 单机缓存和基于 HTTP、gRPC的分布式缓存
- 最近最少访问(Least Recently Used, LRU) 缓存策略
- 使用 Go 锁机制防止缓存击穿
- 使用一致性哈希选择节点，实现负载均衡
- 使用 protobuf 优化节点间二进制通信