package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	order "github.com/lakhansamani/ecom-grpc-apis/order/v1"
	user "github.com/lakhansamani/ecom-grpc-apis/user/v1"

	"github.com/lakhansamani/ecom-grpc-orderd/db"
	"github.com/lakhansamani/ecom-grpc-orderd/service"
)

const (
	serviceName = "orderd"
	metricsPort = ":9092"
)

// Register Metrics
var (
	grpcMetrics       = grpcprom.NewServerMetrics()
	grpcClientMetrics = grpcprom.NewClientMetrics()
)

func main() {
	// Read .env file as environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, using environment variables")
	}

	// DB URL
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is required")
	}
	// Initialize database
	dbProvider := db.New(dbURL)

	// Get User Service URL
	userServiceURL := os.Getenv("USER_SERVICE_URL")
	if userServiceURL == "" {
		log.Fatal("USER_SERVICE_URL is required")
	}
	// Jaeger URL
	jaegerURL := os.Getenv("JAEGER_URL")
	if jaegerURL == "" {
		log.Fatal("JAEGER_URL is required")
	}

	ctx := context.Background()
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure(), otlptracegrpc.WithEndpoint(jaegerURL))
	if err != nil {
		log.Fatalf("failed to create OTLP trace exporter: %v", err)
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)

	defer func() {
		tracerProvider.Shutdown(context.Background())
	}()

	otel.SetTracerProvider(tracerProvider)

	// Create OpenTelemetry gRPC stats handler
	serverHandler := otelgrpc.NewServerHandler(
		otelgrpc.WithTracerProvider(tracerProvider),
	)
	prometheus.MustRegister(grpcMetrics)
	// Create a new gRPC server
	server := grpc.NewServer(
		grpc.StatsHandler(serverHandler),
		grpc.ChainUnaryInterceptor(
			grpcMetrics.UnaryServerInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			grpcMetrics.StreamServerInterceptor(),
		),
	)

	openTelemetryClientHandler := otelgrpc.NewClientHandler(
		otelgrpc.WithTracerProvider(tracerProvider),
	)
	// Create UserServiceClient using grpc
	grpcConn, err := grpc.NewClient(
		userServiceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(openTelemetryClientHandler),
		// Add gRPC Client Interceptors for Prometheus Metrics
		grpc.WithUnaryInterceptor(grpcClientMetrics.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(grpcClientMetrics.StreamClientInterceptor()),
	)
	if err != nil {
		log.Fatalf("Failed to dial UserService: %v", err)
	}
	defer grpcConn.Close()

	userServiceClient := user.NewUserServiceClient(grpcConn)

	// Register OrderService with gRPC
	orderService := service.New(
		service.Config{},
		service.Dependencies{
			DBProvider:    dbProvider,
			UserService:   userServiceClient,
			TraceProvider: tracerProvider,
		})
	order.RegisterOrderServiceServer(server, orderService)

	// Start gRPC server
	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	// Start Prometheus HTTP server
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println("Prometheus metrics server running on port", metricsPort)
		log.Fatal(http.ListenAndServe(metricsPort, nil))
	}()
	reflection.Register(server)
	log.Println("gRPC Server is running on port 50052...")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
