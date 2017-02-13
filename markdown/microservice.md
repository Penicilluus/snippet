###### Before adopting microservice

- Which languages and technologies should I adopt?
- Where and how do I deploy my microservices?
- How do I perform service-discovery in this environment?
- How do I manage my data?
- How do I design my application to handle failure?* (Yes! It will fail!) *
- How do I address authentication, monitoring and tracing?


[Go:微服务架构的基石](http://ppt.geekbang.org/slide/show/615)

- 负载均衡：seesaw、caddy
- 服务网关：tyk、fabio、vulcand
- 进程间通信：RESTful、RPC、自定义
  - REST框架：beego、gin、Iris、micro、go-kit、goa
  - RPC框架：grpc、thrift、hprose
  - 自定义：协议，编解码
- 服务发现：etcd、consul、serf
- 调度系统：kubernetes、swarm、mesos
- 消息队列：NSQ、Nats
- APM（应用性能监控）：appdash、Cloudinsight、opentracing
- 配置管理：etcd、consul、mgmt
- 日志分析：Beats、Heka
- 服务监控：open-falcon、prometheus
- CI/CD：Drone
- 熔断器：gateway、Hystrix-go

###### 参考

- [kubernetes笔记](https://www.zybuluo.com/dujun/note/58625)

