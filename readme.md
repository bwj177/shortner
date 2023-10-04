
Jay短链系统
项目介绍：    概述：该项目基于go-zero框架编写，主要用于将繁琐的长链转换成短链    
场景：营销短信、app push    
解决痛点：短信内容超长，需要拆分发送造成浪费；长链生成二维码模糊；许多平台有发送消息长度限制 


技术栈：go-zero,redis,mysql,nginx,(kafka,prometheus,grafana)


```基于MySQL主键实现了高可用的发号器组件（可分布式分片扩展），以62进制为基础转链

在转链前进行特殊词过滤（前缀树）和防止循环转链、重复转链的校验处理 

查看链接服务使用singleflight防止缓存击穿 ，采用布谷鸟过滤器防止缓存穿透 对短链设置有效期，

参考Redis过期删除策略对短链采取惰性删除+定期删除进行过期短链删除
扩展：基于prometheus+grafama对系统进行监控
```