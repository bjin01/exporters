package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/bjin01/autoapi/getyaml"
	"github.com/bjin01/go-xmlrpc"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	namespace = "suma" // For Prometheus metrics.
)

var (
	jobLableNames = []string{"type"}
)

type metricInfo struct {
	Desc *prometheus.Desc
	Type prometheus.ValueType
}

type metrics map[int]metricInfo

func newJobsMetric(metricName string, docString string, t prometheus.ValueType, constLabels prometheus.Labels) metricInfo {
	return metricInfo{
		Desc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "jobs", metricName),
			docString,
			jobLableNames,
			constLabels,
		),
		Type: t,
	}
}

func newSystemsMetric(metricName string, docString string, t prometheus.ValueType, constLabels prometheus.Labels) metricInfo {
	return metricInfo{
		Desc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "systems", metricName),
			docString,
			jobLableNames,
			constLabels,
		),
		Type: t,
	}
}

var (
	sumaUp         = prometheus.NewDesc(prometheus.BuildFQName(namespace, "", "up"), "Was the last scrape of SUSE Manager successful.", nil, nil)
	systemsMetrics = metrics{
		2: newSystemsMetric("physical_systems", "Number of physical bare metal systems in SUSE Manager.", prometheus.GaugeValue, nil),
		3: newSystemsMetric("virtual_systems", "Number of virtual systems in SUSE Manager.", prometheus.GaugeValue, nil),
		4: newSystemsMetric("active_systems", "Number of active online systems in SUSE Manager.", prometheus.GaugeValue, nil),
		5: newSystemsMetric("offline_systems", "Number of inactive systems in SUSE Manager.", prometheus.GaugeValue, nil),
		6: newSystemsMetric("outdated_systems", "Number of out of date systems in SUSE Manager.", prometheus.GaugeValue, nil),
	}
	jobMetrics = metrics{
		2: newJobsMetric("pending_jobs", "Current number of active pending jobs in SUSE Manager.", prometheus.GaugeValue, nil),
		3: newJobsMetric("completed_jobs", "Current number of completed jobs in SUSE Manager.", prometheus.GaugeValue, nil),
		4: newJobsMetric("failed_jobs", "Current number of failed jobs in SUSE Manager.", prometheus.GaugeValue, nil),
		5: newJobsMetric("archived_jobs", "Current number of archived jobs in SUSE Manager.", prometheus.CounterValue, nil),
	}
)

type Exporter struct {
	suma_server_url     string
	username            string
	password            string
	mutex               sync.RWMutex
	up                  prometheus.Gauge
	totalScrapes        prometheus.Counter
	suma_jobMetrics     map[int]metricInfo
	suma_systemsMetrics map[int]metricInfo
}

func NewExporter(suma_server_url string, username string, password string, jobmetrics map[int]metricInfo, systemsMetrics map[int]metricInfo) *Exporter {
	return &Exporter{
		suma_server_url: suma_server_url,
		username:        username,
		password:        password,
		up: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "up",
			Help:      "Was the last scrape of SUSE Manager successful.",
		}),
		totalScrapes: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "exporter_scrapes_total",
			Help:      "Current total SUMA scrapes.",
		}),
		suma_jobMetrics:     jobmetrics,
		suma_systemsMetrics: systemsMetrics,
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range jobMetrics {
		ch <- m.Desc

	}
	for _, m := range systemsMetrics {
		ch <- m.Desc

	}
	ch <- e.totalScrapes.Desc()
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mutex.Lock() // To protect metrics from concurrent collects.
	defer e.mutex.Unlock()

	up := e.scrape(ch)
	ch <- prometheus.MustNewConstMetric(sumaUp, prometheus.GaugeValue, up)
	ch <- e.totalScrapes

}

func (e *Exporter) scrape(ch chan<- prometheus.Metric) (up float64) {
	e.totalScrapes.Inc()

	for _, metric := range e.suma_jobMetrics {
		value, labelValue := e.query_suma(metric.Desc.String())
		ch <- prometheus.MustNewConstMetric(metric.Desc, metric.Type, value, labelValue...)
	}
	for _, metric := range e.suma_systemsMetrics {
		value, labelValue := e.query_suma(metric.Desc.String())
		ch <- prometheus.MustNewConstMetric(metric.Desc, metric.Type, value, labelValue...)
	}
	return 1
}

func (e *Exporter) query_suma(metric_desc string) (value float64, labels []string) {
	x := 0
	physicals_checked := false

	client := xmlrpc.NewClient(e.suma_server_url)

	f, err := client.Call("auth.login", e.username, e.password)

	if err != nil {
		log.Fatal("Couldn't login to suse manager host.")
	}

	if strings.Contains(metric_desc, "failed_jobs") {
		a := e.get_suma_values(client, f.String(), "schedule.listFailedActions")
		int_val, _ := a.(int)
		labelNames := []string{"failed_jobs"}
		return float64(int_val), labelNames
	}

	if strings.Contains(metric_desc, "pending_jobs") {
		a := e.get_suma_values(client, f.String(), "schedule.listInProgressActions")
		int_val, _ := a.(int)
		labelNames := []string{"pending_jobs"}
		return float64(int_val), labelNames
	}

	if strings.Contains(metric_desc, "completed_jobs") {
		a := e.get_suma_values(client, f.String(), "schedule.listCompletedActions")
		int_val, _ := a.(int)
		labelNames := []string{"completed_jobs"}
		return float64(int_val), labelNames
	}

	if strings.Contains(metric_desc, "archived_jobs") {
		a := e.get_suma_values(client, f.String(), "schedule.listArchivedActions")
		int_val, _ := a.(int)
		labelNames := []string{"archived_jobs"}
		return float64(int_val), labelNames
	}

	if strings.Contains(metric_desc, "physical_systems") {
		a := e.get_suma_values(client, f.String(), "system.listPhysicalSystems")
		x = a.(int)
		physicals_checked = true
		int_val, _ := a.(int)
		labelNames := []string{"physical_systems"}
		return float64(int_val), labelNames
	}

	if strings.Contains(metric_desc, "virtual_systems") {
		a := e.get_suma_values(client, f.String(), "system.listSystems")
		// x is the number of physical systems, a is total number of systems.
		int_val := 0
		// Need to do be sure that physical systems number is already known, if not we call listPhysicalSystems
		if physicals_checked == true {
			b := a.(int) - x
			int_val = b
		} else {
			a1 := e.get_suma_values(client, f.String(), "system.listPhysicalSystems")
			x = a1.(int)
			b := a.(int) - x
			int_val = b
		}
		labelNames := []string{"virtual_systems"}
		return float64(int_val), labelNames
	}

	if strings.Contains(metric_desc, "active_systems") {
		a := e.get_suma_values(client, f.String(), "system.listActiveSystems")
		int_val, _ := a.(int)
		labelNames := []string{"active_systems"}
		return float64(int_val), labelNames
	}

	if strings.Contains(metric_desc, "offline_systems") {
		a := e.get_suma_values(client, f.String(), "system.listInactiveSystems")
		int_val, _ := a.(int)
		labelNames := []string{"offline_systems"}
		return float64(int_val), labelNames
	}

	if strings.Contains(metric_desc, "outdated_systems") {
		a := e.get_suma_values(client, f.String(), "system.listOutOfDateSystems")
		int_val, _ := a.(int)
		labelNames := []string{"outdated_systems"}
		return float64(int_val), labelNames
	}

	client.Call("auth.logout", f.String())
	labelNames := []string{"something went wrong"}
	return 0.00, labelNames
}

func (e *Exporter) get_suma_values(client xmlrpc.Client, sessionkey string, api_method string) interface{} {

	method := api_method
	log.Printf("Calling: %v\n", method)
	u, err := client.Call(method, sessionkey)
	if err != nil {
		log.Fatal("Couldn't get values: " + method)
	}

	a := getVal(u)
	return a

}

func getVal(v xmlrpc.Value) interface{} {
	for _, x := range v.Values() {
		for _, y := range x.Members() {
			fmt.Printf("%v: ", y.Name())
			getvalue3(y.Value())
		}
	}
	return len(v.Values())
}

func getvalue3(v xmlrpc.Value) interface{} {
	z := v.Kind()
	y := v
	var return_val interface{}

	switch f := z; f {
	case 1:
		//GetMembers(y.Members(), searchfield)
	case 2:
		fmt.Printf("\t%v\n", y.Bytes())
	case 3:
		fmt.Printf("\t%v\n", y.Bool())
	case 4:
		fmt.Printf("\t%s\n", y.Time())
		//return y.Time
	case 5:
		fmt.Printf("%v\n", y.Double())
	case 6: //this is a int type
		fmt.Printf("\t integer %v\n", y.Int())
		//return y.Int()

	case 7: //this is a string type
		fmt.Printf("\t%v\n", y.String())
		//return y.String()

	case 8: //this is a member type

		//return_val = GetVal(y, searchfield, 0)
	default:

		//return_val = GetVal(y, searchfield, 0)
	}
	return return_val

}

func main() {
	listenAddress := ":9102"

	metricsPath := "/metrics"

	cfgPath, err := getyaml.ParseFlags()
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := getyaml.NewConfig(cfgPath)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Scraping %v as %v. exporter on port: %v", cfg.Server.ApiUrl, cfg.Server.Username, listenAddress)

	exporter := NewExporter(cfg.Server.ApiUrl, cfg.Server.Username, cfg.Server.Password, jobMetrics, systemsMetrics)
	prometheus.MustRegister(exporter)

	http.Handle(metricsPath, promhttp.Handler())
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}
