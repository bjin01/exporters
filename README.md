# Prometheus exporter - SUSE Manager / Uyuni

This is my first prometheus exporter that scrapes SUSE Manager to get some jobs and system information.

![suma dashboard](https://github.com/adam-p/markdown-here/raw/master/src/common/images/icon48.png "suma dashboard")
### Usage:
Make sure you installed go1.13 or higher on your host. 

```
bjsuma:~ # export PATH=$PATH:/usr/local/go/bin
bjsuma:~ # go version
go version go1.14.3 linux/amd64

```

Make sure you create a yaml config file that consist your SUSE Manager host and login information as below:
```
server:
  apiurl: http://your-host/rpc/api
  username: admin
  password: 12345678
```

For developers you have to runn below command to start a test and further code.

```go run exporter.go -config ./config.yml```

Or for users executing the binary:
```./exporter -config ./config.yml```

Testing:
```curl http://your-host:9102/metrics```

With the curl output you should get a list of metrics that has go and suma metrics I expose. Check if you see the suma metrics and values.
```
# HELP suma_exporter_scrapes_total Current total SUMA scrapes.
# TYPE suma_exporter_scrapes_total counter
suma_exporter_scrapes_total 143
# HELP suma_jobs_archived_jobs Current number of archived jobs in SUSE Manager.
# TYPE suma_jobs_archived_jobs counter
suma_jobs_archived_jobs{type="archived_jobs"} 80
# HELP suma_jobs_completed_jobs Current number of completed jobs in SUSE Manager.
# TYPE suma_jobs_completed_jobs gauge
suma_jobs_completed_jobs{type="completed_jobs"} 168
# HELP suma_jobs_failed_jobs Current number of failed jobs in SUSE Manager.
# TYPE suma_jobs_failed_jobs gauge
suma_jobs_failed_jobs{type="failed_jobs"} 8
# HELP suma_jobs_pending_jobs Current number of active pending jobs in SUSE Manager.
# TYPE suma_jobs_pending_jobs gauge
suma_jobs_pending_jobs{type="pending_jobs"} 2
# HELP suma_systems_active_systems Number of active online systems in SUSE Manager.
# TYPE suma_systems_active_systems gauge
suma_systems_active_systems{type="active_systems"} 4
# HELP suma_systems_offline_systems Number of inactive systems in SUSE Manager.
# TYPE suma_systems_offline_systems gauge
suma_systems_offline_systems{type="offline_systems"} 14
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

If data is correct then you can go to your prometheus server and add this target into your scrap job.

This is my prometheus.yaml config I added:
```
# Scrape configurations
scrape_configs:
  # --------------------
  # Monitor bjsuma.bo2go.home
  # --------------------
  - job_name: 'mgr-server'
    static_configs:
      - targets:
        - bjsuma.bo2go.home:9102 # bo suma exporter 

```
Restart prometheus service and check the metrics and values using prometheus expression browser
```http://localhost:9090/graph```

If you can see some results then feel free to continue with grafana dashboard.
Feel free to import the [grafana-dashboard-panel.json](https://github.com/bjin01/exporters/blob/main/grafana-dashboard-panel.json)

Feedbacks are highly appreciated.


