package main

import (
	"caozhipan/nsq-prometheus-exporter/controllers"
	"flag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"net/http"
	"time"
)

var (
	nsqLookupdAddress = flag.String("nsq.lookupd.address", "127.0.0.1:4161", "nsqllookupd address list with comma")
)

func main() {
	flag.Parse()

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for {
			controllers.SyncNodeList(*nsqLookupdAddress)
			<-ticker.C
		}
	}()

	prometheus.MustRegister(controllers.Collector)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9527", nil))

}
