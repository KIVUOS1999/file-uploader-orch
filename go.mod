module github.com/KIVUOS1999/file-uploader-orch

go 1.23.2

require (
	github.com/KIVUOS1999/easyApi v0.0.0-20241117070720-954caed24eaa
	github.com/KIVUOS1999/easyLogs v1.0.0
	github.com/KIVUOS1999/file-uploader-db v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.6.0
	google.golang.org/api v0.211.0
)

replace github.com/KIVUOS1999/easyApi => ../../easyApi

replace github.com/KIVUOS1999/file-uploader-db => ../file-uploader-db

require (
	cloud.google.com/go/auth v0.12.1 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.6 // indirect
	cloud.google.com/go/compute/metadata v0.5.2 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/google/s2a-go v0.1.8 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.4 // indirect
	github.com/googleapis/gax-go/v2 v2.14.0 // indirect
	github.com/gorilla/mux v1.8.1 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.54.0 // indirect
	go.opentelemetry.io/otel v1.29.0 // indirect
	go.opentelemetry.io/otel/metric v1.29.0 // indirect
	go.opentelemetry.io/otel/trace v1.29.0 // indirect
	golang.org/x/crypto v0.30.0 // indirect
	golang.org/x/net v0.32.0 // indirect
	golang.org/x/oauth2 v0.24.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241206012308-a4fef0638583 // indirect
	google.golang.org/grpc v1.67.1 // indirect
	google.golang.org/protobuf v1.35.2 // indirect
)
