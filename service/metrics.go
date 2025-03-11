package service

import "github.com/prometheus/client_golang/prometheus"

const (
	serviceName = "orderd"
	component   = "service"

	totalOrdersCreatedMetric = "total_orders"
	totalFetchedOrdersMetric = "total_fetched_orders"

	successResultLabel = "success"
	failedResultLabel  = "failed"
)

var (
	ordersCreatedMetrics = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: serviceName,
		Subsystem: component,
		Name:      totalOrdersCreatedMetric,
	}, []string{"result"})

	ordersFetchedMetrics = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: serviceName,
		Subsystem: component,
		Name:      totalFetchedOrdersMetric,
	}, []string{"result", "order_id"})
)
