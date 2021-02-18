# Prometheus exporter - SUSE Manager / Uyuni

This is my first prometheus exporter that scrapes SUSE Manager to get some jobs and system information.

### Usage:
Make sure you installed go1.13 or higher on your host. 

```
bjsuma:~ # export PATH=$PATH:/usr/local/go/bin
bjsuma:~ # go version
go version go1.14.3 linux/amd64

```
I have not compiled my exporter.go as binary yet so you have to runn below command to start a test.

```go run exporter.go -config ./config.yml```

Feedbacks are highly appreciated.


