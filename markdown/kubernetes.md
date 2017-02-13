Kubernetes是Google开源的容器集群管理系统。它构建于docker技术之上，为容器化的应用提供资源调度、部署运行、服务发现、扩容缩容等整一套功能，本质上可看作是基于容器技术的mini-PaaS平台

![kubernetes](/Users/tamchen/Documents/project/material/kubernetes.jpg)![kuberctl](/Users/tamchen/Documents/project/material/kuberctl.jpg)

达到的目标

- 所有的业务服务都运行在docker上，使用kubernetes进行编排
- 容器计算资源由平台统一管理，扩容弹性升级
- 与CI&CD完全集成（build once run anyway）和jenkins集成自动化部署
- 分布式架构，保证扩展性；
- 逻辑集中式的控制平面 + 物理分布式的运行平面；
- 一套资源调度系统，管理哪个容器该分配到哪个节点上；
- 一套对容器内服务进行抽象和 HA 的系统。



post pods ： user - kubecfg - apiserver - etcd - scheduler - kuberctl - docker 

post controllers ：user - kubecfg - apiserver - etcd - controllermanager - kuberctl - docker

post services ： user - kubecfg - apiserver - etcd - controllermanager  - proxy - socket - IPtables

服务注册发现、负载均衡、日志收集、监控报警

- apiserver 是整个系统的对外接口，提供 RESTful 方式供客户端和其它组件调用；
- scheduler 负责对资源进行调度，分配某个 pod 到某个节点上；
- controller-manager 负责管理控制器，包括 endpoint-controller（刷新服务和 pod 的关联信息）和 replication-controller（维护某个 pod 的复制为配置的数值）。
- kubelet 是工作节点执行操作的 agent，负责具体的容器生命周期管理，根据从数据库中获取的信息来管理容器，并上报 pod 运行状态等；
- proxy 为 pod 上的服务提供访问的代理。
- etcd 是所有状态的存储数据库；
- `gcr.io/google_containers/pause:0.8.0` 是 Kubernetes 启动后自动 pull 下来的测试镜像。

通过lab机制实现负载均衡

Deployment/RC的定义其实是Pod创建模板(Template)+Pod副本数量的声明(replicas)：

###### pod用例

- 内容管理系统，文件和数据载入器，本地缓存管理等；
- 日志和检查点备份，压缩，轮换，快照系统等；
- 数据变化监视，日志末端数据读取，日志和监控适配器，事件打印等；
- 代理，桥接和适配器；
- 控制器，管理器，配置编辑器和更新器。

###### flannel网络

![packet-01](/Users/tamchen/Documents/project/material/packet-01.png)

https://github.com/coreos/flannel

###### 服务配置中 port/node port/target port

- port表示service暴露在cluster ip上的端口
- nodePort 是提供给集群外部客户访问service的入口
- targetPort是pod上的端口，从port和nodePort上到来的数据最终经过kube-proxy流入到后端pod的targetPort上进入容器。

service是如何通过port和nodePort向内外提供服务的呢？通过 Kube-Proxy 在本地node上建立的iptables规则

###### 参考

[Docker从入门到实战](http://www.kancloud.cn/maozhenggang/docker-api/94298)