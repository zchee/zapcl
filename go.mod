module github.com/zchee/zapcl

go 1.20

replace (
	go.opentelemetry.io/otel => go.opentelemetry.io/otel v1.13.0
	go.opentelemetry.io/otel/trace => go.opentelemetry.io/otel/trace v1.13.0
)

require (
	cloud.google.com/go/compute/metadata v0.2.3
	github.com/goccy/go-json v0.10.0
	github.com/google/go-cmp v0.5.9
	go.opentelemetry.io/otel/trace v1.13.0
	go.uber.org/zap v1.24.0
	golang.org/x/sys v0.6.1-0.20230304190818-494aa493ccb0
	google.golang.org/genproto v0.0.0-20230222225845-10f96fb3dbec
	google.golang.org/protobuf v1.28.1
)

require (
	cloud.google.com/go/compute v1.18.0 // indirect
	cloud.google.com/go/logging v1.7.0 // indirect
	cloud.google.com/go/longrunning v0.4.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	go.opentelemetry.io/otel v1.13.0 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.9.1-0.20230215063618-4504ef7e0048 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/grpc v1.53.0 // indirect
)
