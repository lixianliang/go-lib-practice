package main

import (
	"fmt"
	"net"
	"net/netip"

	"github.com/prometheus/client_golang/prometheus"
)

var dialTakesTime = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: "kodo",
		Subsystem: "outbound_addr",
		Name:      "dial_duration_ms",
		Buckets:   []float64{1, 5, 10, 30, 50, 100, 300, 500, 1000, 3000, 5000, 10000, 30000}, // 单位 毫秒
		Help:      "outbound add dial duration",
	},
	[]string{"srcisp", "ip", "dstisp"},
)

func init() {
	prometheus.MustRegister(dialTakesTime)
}

func Pick() net.Addr {
	return nil
}

func main() {
	addr := netip.MustParseAddr("10.0.0.1")
	if addr.IsPrivate() {
		fmt.Println("Private IP")
	}

	addr = netip.MustParseAddr("127.0.0.1")
	if addr.IsPrivate() {
		fmt.Println("Private IP")
	}

	a := Pick()
	fmt.Printf("%v\n", a)
	fmt.Printf("%s\n", a)
	// fmt.Println(a.String()) panic

	// label 必须填满。
	dialTakesTime.With(prometheus.Labels{"srcisp": "ss", "ip": "dd", "dstisp": ""}).Observe(float64(10))

}
