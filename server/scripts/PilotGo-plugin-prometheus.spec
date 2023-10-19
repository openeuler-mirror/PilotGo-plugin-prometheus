%define         debug_package %{nil}

Name:           PilotGo-plugin-prometheus
Version:        1.0.1
Release:        2
Summary:        PilotGo prometheus plugin provides cluster monitor and alert.
License:        MulanPSL-2.0
URL:            https://gitee.com/openeuler/PilotGo-plugin-prometheus
Source0:        https://gitee.com/openeuler/PilotGo-plugin-prometheus/%{name}-%{version}.tar.gz
# tar -xvf Source0
# cd %{name}-%{version}/web
# run 'yarn install and yarn build' in it
# tar -czvf %{name}-web.tar.gz ../web/dist/
Source1:        PilotGo-plugin-prometheus-web.tar.gz
BuildRequires:  systemd
BuildRequires:  golang
Requires:       prometheus2
Requires:       golang-github-prometheus-node_exporter
Provides:       pilotgo-plugin-prometheus = %{version}-%{release}

%description
PilotGo prometheus plugin provides cluster monitor and alert.


%prep
%autosetup -p1 -n %{name}-%{version}
tar -xzvf %{SOURCE1}

%build
cd server
GO111MODULE=on go build -mod=vendor -o PilotGo-plugin-prometheus main.go

%install
mkdir -p %{buildroot}/opt/PilotGo/plugin/prometheus/{server/{scripts,log},web/dist}
install -D -m 0755 server/PilotGo-plugin-prometheus %{buildroot}/opt/PilotGo/plugin/prometheus/server
install -D -m 0644 server/scripts/init_prometheus_yml.sh %{buildroot}/opt/PilotGo/plugin/prometheus/server/scripts/init_prometheus_yml.sh
install -D -m 0644 server/scripts/prometheus.yml %{buildroot}/opt/PilotGo/plugin/prometheus/server/scripts/prometheus.yml
install -D -m 0644 server/config.yml.templete %{buildroot}/opt/PilotGo/plugin/prometheus/server/config.yml
install -D -m 0644 server/scripts/PilotGo-plugin-prometheus.service %{buildroot}%{_unitdir}/PilotGo-plugin-prometheus.service
cp -rf web/dist %{buildroot}/opt/PilotGo/plugin/prometheus/web

%post
%systemd_post PilotGo-plugin-prometheus.service

%preun
%systemd_preun PilotGo-plugin-prometheus.service

%postun
%systemd_postun PilotGo-plugin-prometheus.service

%files
%dir /opt/PilotGo
%dir /opt/PilotGo/plugin
%dir /opt/PilotGo/plugin/prometheus
%dir /opt/PilotGo/plugin/prometheus/server
%dir /opt/PilotGo/plugin/prometheus/server/log
%dir /opt/PilotGo/plugin/prometheus/server/scripts
/opt/PilotGo/plugin/prometheus/server/PilotGo-plugin-prometheus
/opt/PilotGo/plugin/prometheus/server/scripts/init_prometheus_yml.sh
/opt/PilotGo/plugin/prometheus/server/scripts/prometheus.yml
/opt/PilotGo/plugin/prometheus/server/config.yml
%{_unitdir}/PilotGo-plugin-prometheus.service
/opt/PilotGo/plugin/prometheus/web/dist


%changelog
* Tue Sep 26 2023 jianxinyu <jiangxinyu@kylinos.cn> - 1.0.1-2
- Compatible with golang version 1.17

* Fri Sep 01 2023 jianxinyu <jiangxinyu@kylinos.cn> - 1.0.1-1
- Package init

