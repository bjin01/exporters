# Prometheus exporter - SUSE Manager / Uyuni

This is my first prometheus exporter that scrapes SUSE Manager to get some jobs and system information visualized in Grafana.

* Number of Jobs in SUSE Manager (pending, failed, completed, archived)
* Number of systems (active, inactive, outdated etc.)
* Number of Systems by OS Version
* Top 10 Systems with highest system scores - means systems with number of patches outdated

![suma dashboard](https://github.com/bjin01/exporters/blob/main/sample-dashboard.png "suma dashboard")
## Usage:
### For admins:
Download the compiled binary to your local SUSE Manager host.
The binary has been compiled and there is no need to install go library on your host. Just download it and use it.

```
cd /usr/local/bin
wget https://raw.githubusercontent.com/bjin01/exporters/main/binary/suma-exporter
chmod +x /usr/local/bin/suma-exporter
```
Create a configuration file and store it at a place where only root user has access as the file will contain the login password to suse manager.
In my systemd service unit file I used a configuration file located here: ```/etc/suma-exporters/config.yml```

Make sure you create a yaml config file that consist your SUSE Manager port, host and login information as below:
The port number should be higher than 9100 in order to not conflict with other prometheus ports. 

```
server:
  apiurl: http://your-host/rpc/api
  username: admin
  password: 12345678
  port: 9102
```
Add the suma-exporter as a systemd service on SUSE Manager host:
work as root user or with sudo
```
cd /etc/systemd/system
wget https://raw.githubusercontent.com/bjin01/exporters/main/systemd/suma-exporter.service
systemctl daemon-reload
systemctl start suma-exporter.service
```
Test the running suma-exporter:
```curl -v http://your-host:9102/metrics```
You should get the metrics results displayed as shown further below.

Last but not least do not forget to add the suma-exporter into /etc/prometheus/prometheus.yaml scraping job and restart prometheus.service on prometheus host.
```
# Scrape configurations
scrape_configs:
  # --------------------
  # Monitor bjsuma.bo2go.home
  # --------------------
  - job_name: 'mgr-server'
    scrape_interval: 360s
    static_configs:
      - targets:
        - bjsuma.bo2go.home:9102 # bo suma exporter 
```
__Cautious:__ Do not set too short scraping interval too low which will cause performance issues on SUSE Manager as the exporter has to make several xmlrpc api calls with each scraping. For __large systems with over more than 200 managed systems__ I would recommend to set scraping interval less than 360 seconds. 

The next logical step would be to add the new metrics into grafana dashboard.

Feel free to import the [grafana-dashboard-panel.json](https://github.com/bjin01/exporters/blob/main/grafana-dashboard-panel.json) to your grafana.

### For developers you have to runn below command to start a test and further coding.
Make sure you installed go1.13 or higher on your host. 

```
bjsuma:~ # export PATH=$PATH:/usr/local/go/bin
bjsuma:~ # go version
go version go1.14.3 linux/amd64

```

```
cd src/github.com/bjin01/exporters
go run *.go -config ./config.yml
```

Or for users executing the binary:
```./exporter -config ./config.yml```

Testing:
```curl http://your-host:9102/metrics```

With the curl output you should get a list of metrics that has go and suma metrics I expose. Check if you see the suma metrics and values.
```
# HELP suma_exporter_scrapes_total Current total SUMA scrapes.
# TYPE suma_exporter_scrapes_total counter
suma_exporter_scrapes_total 1
# HELP suma_scores_system_currency system currency of the top10 nodes
# TYPE suma_scores_system_currency gauge
suma_scores_system_currency{critical="0",hostname="bjlx15",important="8",total_scores="288"} 288
suma_scores_system_currency{critical="0",hostname="my15sp1test.bo2go.home",important="5",total_scores="152"} 152
suma_scores_system_currency{critical="0",hostname="testrhel72.bo2go.home",important="10",total_scores="591"} 591
suma_scores_system_currency{critical="1",hostname="pxetest.bo2go.home",important="9",total_scores="252"} 252
suma_scores_system_currency{critical="1",hostname="testrhel02.bo2go.home",important="7",total_scores="418"} 418
suma_scores_system_currency{critical="2",hostname="caasp05.bo2go.home",important="17",total_scores="631"} 631
suma_scores_system_currency{critical="2",hostname="pampam.bo2go.home",important="3",total_scores="131"} 131
suma_scores_system_currency{critical="2",hostname="smt1.bo2go.home",important="3",total_scores="156"} 156
suma_scores_system_currency{critical="6",hostname="azure-sap-test.bo2go.home",important="45",total_scores="1694"} 1694
suma_scores_system_currency{critical="6",hostname="tomcat2.bo2go.home",important="35",total_scores="1251"} 1251
# HELP suma_jobs_archived_jobs Current number of archived jobs in SUSE Manager.
# TYPE suma_jobs_archived_jobs counter
suma_jobs_archived_jobs{type="archived_jobs"} 80
# HELP suma_jobs_base_product Number of each base product in SUSE Manager
# TYPE suma_jobs_base_product gauge
suma_jobs_base_product{type=""} 1
suma_jobs_base_product{type="CentOS 8 x86_64"} 1
suma_jobs_base_product{type="SLES 12 SP4 x86_64"} 2
suma_jobs_base_product{type="SLES 15 SP1 x86_64"} 5
suma_jobs_base_product{type="SLES with RES 7 x86_64"} 2
suma_jobs_base_product{type="SLES4SAP 15 SP1 x86_64"} 5
suma_jobs_base_product{type="SLES4SAP 15 SP2 x86_64"} 1
suma_jobs_base_product{type="Ubuntu 20.04"} 1
# HELP suma_jobs_completed_jobs Current number of completed jobs in SUSE Manager.
# TYPE suma_jobs_completed_jobs gauge
suma_jobs_completed_jobs{type="completed_jobs"} 174
# HELP suma_jobs_failed_jobs Current number of failed jobs in SUSE Manager.
# TYPE suma_jobs_failed_jobs gauge
suma_jobs_failed_jobs{type="failed_jobs"} 8
# HELP suma_jobs_pending_jobs Current number of active pending jobs in SUSE Manager.
# TYPE suma_jobs_pending_jobs gauge
suma_jobs_pending_jobs{type="pending_jobs"} 2
# HELP suma_systems_active_systems Number of active online systems in SUSE Manager.
# TYPE suma_systems_active_systems gauge
suma_systems_active_systems{type="active_systems"} 6
# HELP suma_systems_offline_systems Number of inactive systems in SUSE Manager.
# TYPE suma_systems_offline_systems gauge
suma_systems_offline_systems{type="offline_systems"} 12
# HELP suma_systems_outdated_systems Number of out of date systems in SUSE Manager.
# TYPE suma_systems_outdated_systems gauge
suma_systems_outdated_systems{type="outdated_systems"} 16
# HELP suma_systems_physical_systems Number of physical bare metal systems in SUSE Manager.
# TYPE suma_systems_physical_systems gauge
suma_systems_physical_systems{type="physical_systems"} 1
# HELP suma_systems_virtual_systems Number of virtual systems in SUSE Manager.
# TYPE suma_systems_virtual_systems gauge
suma_systems_virtual_systems{type="virtual_systems"} 17
# HELP suma_up Was the last scrape of SUSE Manager successful.
# TYPE suma_up gauge
suma_up 1
```

Feedbacks are highly appreciated.


