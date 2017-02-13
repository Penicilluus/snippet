###### gitlab安装

` Install and configure the necessary dependencies`：

`````shell
sudo apt-get install curl openssh-server ca-certificates postfix
`````

`下载安装包`：

[gitlab-ce_8.14.0-ce.0_amd64.deb](https://packages.gitlab.com/gitlab/gitlab-ce/packages/ubuntu/precise/gitlab-ce_8.14.0-ce.0_amd64.deb)

```shell
sudo dpkg -i gitlab_7.4.2-omnibus-1_amd64.deb
sudo gitlab-ctl reconfigure

gitlab: GitLab should be reachable at http://precise64
gitlab: Otherwise configure GitLab for your system by editing /etc/gitlab/gitlab.rb file
gitlab: And running reconfigure again.

sudo mkdir -p /etc/gitlab
sudo touch /etc/gitlab/gitlab.rb
sudo chmod 600 /etc/gitlab/gitlab.rb
sudo gedit /etc/gitlab/gitlab.rb

external_url "192.168.1.10:8800"
修改 In your gitlab.rb you will need to update unicorn['port'] = 8080
如果8080已经被占用
Username: root
Password: 5iveL!fe

sudo gitlab-ctl restart # 重启gitlab
sudo gitlab-ctl tail # 查看所有日志
```

参考链接

- [(记录)gitlab安装配置| Sky's Blog](http://skyao.github.io/2015/02/16/git-gitlab-setup/)
- [gitlab 中文文档](https://doc.gitlab.cc/)