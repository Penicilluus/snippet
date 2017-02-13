配置https(使用Let’s Encrypt)

跨站脚本攻击(Cross-Site Scripting，XSS)

`概念定义`：通过存在安全漏洞的web网站注册用户的浏览器内运行非法的HTML标签或JavaScript进行的一种攻击。利用预先设置的陷阱触发的被动攻击

- 主要危害利用虚假输入表单骗取用户个人信息
- 利用脚本窃取用户的cookie值，在被害者不知情下，帮助攻击者发送恶意请求
- 显示伪造的文章和图片

SQL注入攻击

`概念定义`：是指针对web应用使用的数据库，通过运行非法的SQL而产生的攻击。

跨站点请求伪造

`跨站点请求伪造(Cross-Site Request Forgeries,CSRF)`：该**攻击**可以在受害者毫不知情的情况下以受害者名义伪造请求发送给受**攻击**站点，从而在并未授权的情况下执行在权限保护之下的操作，有很大的危害性

　　1.登录受信任网站A，并在本地生成Cookie。

　　2.在不登出A的情况下，访问危险网站B。

​	3.危险网站B向网站A发送请求，用户携带本地生成的Cookie，网站A不能区别是用户本身访问还是，危险网站B访问。

针对CSRF浏览器不允许跨站请求，即不同协议、域名、端口即为跨站。

一起来聊聊数据加密技术

https://gold.xitu.io/post/584e3d1c0ce463005c625aff



参考链接

[https://letsencrypt.org/](https://letsencrypt.org/)

[https://github.com/Neilpang/acme.sh](https://github.com/Neilpang/acme.sh)