package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	hemgw16a "local.packages/HEMGW16A"
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
	todatBuyWattHour = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "today_buy_watthour",
		Help:      "Today buy WattHout (wh)",
	})
)

// Update update
func Update(gw hemgw16a.GW) {
	for {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("Runtime Error:", err)
			}
		}()
		todatBuyWattHour.Set(gw.GetTodayBuyWatt())
		nowWattHour.Set(float64(gw.GetNowWatt()))
	}
}

func main() {
	log.Println("HEMGW16A exporter started.")

	// argv
	port := flag.String("port", "0.0.0.0:9100", "Port number to listen on")
	url := flag.String("url", "http://192.168.3.121", "HEMGW16A URL")
	username := flag.String("username", "root", "HEMGW16A user name")
	password := flag.String("password", "pass", "HEMGW16A user password")
	flag.Parse()

	var gw hemgw16a.GW
	gw.Init(*url, *username, *password)

	go Update(gw)

	prometheus.MustRegister(nowWattHour)
	prometheus.MustRegister(todatBuyWattHour)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*port, nil))

}
