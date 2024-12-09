# pm-ssl-management

## 简介

pm-ssl-management 是一个基于 go 的 ssl 证书监控工具，用于监控阿里云免费 ssl 证书的有效期，并在证书即将过期时自动续签证书。

前端代码：[https://github.com/passerma/pm-ssl-management-web](https://github.com/passerma/pm-ssl-management-web)

## 使用

1. 新建 **app.yml** 文件，参考目录下 **app.yml**
2. 启动项目: `./pm-ssl-management`
3. 默认启动在 **4013** 端口，可以通过配置 yml 文件来自定义端口

## 编译

```sh
sh build.sh
chmod +x ./pm-ssl-management
./pm-ssl-management
```
