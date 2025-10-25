// Package usecase is a nice package
package usecase

import "github.com/prometheus/client_golang/prometheus"

var OrdersCreatedCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "orders_created_total_from_orders",
		Help: "Total number of created orders from order-service",
	},
	[]string{"status", "region"},
)

func init() {
	prometheus.MustRegister(OrdersCreatedCounter)
}
