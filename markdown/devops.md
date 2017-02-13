###### supervisor

实用 Python编写的client/server服务，是Linux/Unix系统下的一个进程管理工具

supervisor安装完成后会生成三个执行程序：supervisortd、supervisorctl、echo_supervisord_conf，分别是supervisor的守护进程服务（用于接收进程管理命令）、客户端（用于和守护进程通信，发送管理进程的指令）、生成初始配置文件程序。

```sh
# 启动supervisor
supervisord -c /etc/supervisord.conf
supervisortd -c ~/config/supervisor/apollo.supervisord.conf start all

supervisorctl -c ~/config/supervisor/apollo.supervisord.conf status
supervisorctl -c ~/config/supervisor/apollo.supervisord.conf restart all
# 查看supervisor服务是否启动
ps -ef | grep supervisor | grep -v grep
```

[Supervisor](https://github.com/Supervisor)  [Initscripts](https://github.com/Supervisor/initscripts)





参考链接

[Linux工具快速教程](http://linuxtools-rst.readthedocs.io/zh_CN/latest/index.html)

[Kubernetes 集群安装指南](https://blog.eood.cn/kubernetes)