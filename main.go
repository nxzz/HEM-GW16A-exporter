package main

import (
	"flag"
	"log"
	"net/http"

	hemgw16a "./HEMGW16A"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics
const (
	namespace = "HWMGW16A"
)

var (
	nowWattHour = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "now_watthour",
		Help:      "Now WattHout (wh)",
	})
	todatBuyWattHour = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "today_buy_watthour",
		Help:      "Today buy WattHout (wh)",
	})
)

// HWMGW16ACollector コレクタ
type HWMGW16ACollector struct {
	gw hemgw16a.GW
}

// Collect 取得と投げるとこ
func (c HWMGW16ACollector) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(
		todatBuyWattHour.Desc(),
		prometheus.CounterValue,
		float64(c.gw.GetTodayBuyWatt()),
	)
	ch <- prometheus.MustNewConstMetric(
		nowWattHour.Desc(),
		prometheus.GaugeValue,
		float64(c.gw.GetNowWatt()),
	)
}

// Describe お掃除
func (c HWMGW16ACollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- todatBuyWattHour.Desc()
	ch <- nowWattHour.Desc()
}

func main() {
	log.Println("HEMGW16A exporter started.")

	// argv
	port := flag.String("port", "0.0.0.0:9100", "Port number to listen on")
	url := flag.String("url", "http://192.168.3.121", "HEMGW16A URL")
	username := flag.String("username", "root", "HEMGW16A user name")
	password := flag.String("password", "pass", "HEMGW16A user password")
	flag.Parse()

	var collector HWMGW16ACollector
	collector.gw.Init(*url, *username, *password)

	prometheus.MustRegister(collector)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*port, nil))

}
