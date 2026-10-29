English | [简体中文](./README.md)

# PilotGo-plugin-prometheus

#### Introduction
The PilotGo prometheus plugin provides cluster monitoring dashboards and alerting.

#### Software Architecture
System environment: openEuler 23.09  
Development languages: Golang > 1.18, Vue3, NodeJs > 14


#### Development / Contribution Steps
Install golang-prometheus-node_exporter, prometheus > 2.28 and mysql in advance  
```bash
$ yum install -y golang-prometheus-node_exporter  
$ yum install -y prometheus > 2.28  
$ yum install -y mysql
```

1.  Clone the repository  
```bash
$ git clone https://atomgit.com/openeuler/PilotGo-plugin-prometheus.git
```
2.  Modify the configuration file  
```bash
$ mv ./server/config.yml.templete ./server/config.yml
$ vim ./server/config.yml
```
plugin_prometheus: url is the IP on which the plugin service starts  
prometheus_server: address where prometheus is installed  
http_server: address on which the plugin service starts  
pilotgo_server: address of the PilotGo server service  
mysql: host is the IP of the host where the DB is installed; user is the DB login user; password is the DB access password

3.  Build the frontend  
```bash
$ cd ./web
$ yarn install  
$ yarn build
```
4.  Start the plugin
```bash
$ cd ./server
$ go run main.go
```
Open http://ip:8090/plugin/prometheus in a browser

#### Supplementary Links
1.  [PilotGo User Manual](https://atomgit.com/openeuler/PilotGo/tree/master/docs/zh/usage_instructions.md)
2.  PilotGo [Code Repository](https://atomgit.com/openeuler/PilotGo)
3.  PilotGo [Package Repository](https://atomgit.com/src-openeuler/PilotGo)
