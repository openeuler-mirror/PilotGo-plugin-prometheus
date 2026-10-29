[English](./README.en.md) | 简体中文

# PilotGo-plugin-prometheus

#### 介绍
PilotGo prometheus 插件提供了集群监控看板和告警功能。

#### 软件架构
系统环境：openEuler 23.09  
开发语言：Golang > 1.18、Vue3、NodeJs > 14


#### 贡献/开发步骤
请提前安装 golang-prometheus-node_exporter、prometheus > 2.28、mysql 
```bash
$ yum install -y golang-prometheus-node_exporter  
$ yum install -y prometheus > 2.28  
$ yum install -y mysql
```

1.  拉取仓库代码  
```bash
$ git clone https://atomgit.com/openeuler/PilotGo-plugin-prometheus.git
```
2.  修改配置文件  
```bash
$ mv ./server/config.yml.templete ./server/config.yml
$ vim ./server/config.yml
```
plugin_prometheus: url地址ip为插件服务启动ip  
prometheus_server: 安装prometheus的地址  
http_server: 插件服务启动地址   
pilotgo_server:  PilotGo server服务地址  
mysql: host为安装DB的主机ip；user为DB的登录用户；password为DB访问密码

3.  前端编译  
```bash
$ cd ./web
$ yarn install 
$ yarn build
```
4.  启动插件
```bash
$ cd ./server
$ go run main.go
```
网页查看 http://ip:8090/plugin/prometheus

#### 补充链接
1.  [PilotGo使用手册](https://atomgit.com/openeuler/PilotGo/tree/master/docs/zh/usage_instructions.md)
2.  PilotGo[代码仓](https://atomgit.com/openeuler/PilotGo)
3.  PilotGo[软件包仓](https://atomgit.com/src-openeuler/PilotGo)
