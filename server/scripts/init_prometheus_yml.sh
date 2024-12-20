inityml(){
# init prometheus yml  
cat>$2<<EOF
global:
  scrape_interval: 15s
  evaluation_interval: 15s
rule_files:
  - /etc/prometheus/rules.yaml
scrape_configs:
  - job_name: node_exporter
    http_sd_configs:
    - url: http://$1/plugin/prometheus/api/target
      refresh_interval: 60s
EOF
cat>/etc/prometheus/rules.yaml<<EOF
groups:
- name: node
  rules:
EOF
echo "初始化"
# restart
systemctl restart prometheus;

}

main(){
	inityml $1 $2
}
main $1 $2