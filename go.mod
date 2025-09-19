module github.com/zhs007/slotsgamecore7

go 1.24.0

require (
	devt.de/krotik/common v1.5.1
	github.com/bytedance/sonic v1.12.10
	github.com/fatih/color v1.18.0
	github.com/golang/protobuf v1.5.4
	github.com/google/cel-go v0.22.1
	github.com/jarcoal/httpmock v1.3.1
	github.com/stretchr/testify v1.10.0
	github.com/valyala/fasthttp v1.58.0
	github.com/valyala/fastrand v1.1.0
	github.com/xuri/excelize/v2 v2.9.1
	github.com/zhs007/goutils v0.2.2
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.58.0
	go.opentelemetry.io/otel v1.34.0
	go.opentelemetry.io/otel/exporters/jaeger v1.17.0
	go.opentelemetry.io/otel/sdk v1.34.0
	go.uber.org/zap v1.27.0
	golang.org/x/net v0.40.0
	gonum.org/v1/gonum v0.15.1
	google.golang.org/grpc v1.71.3
	google.golang.org/protobuf v1.36.9
	gopkg.in/yaml.v2 v2.4.0
)

// (replace removed) Previously we used an explicit replace to map the old
// v2 path to the v3 module. The replace is removed to rely on normal module
// resolution. If build or tidy shows conflicts, we can restore or use an
// alternate targeted replace.

require (
	cel.dev/expr v0.19.3 // indirect
	github.com/andybalholm/brotli v1.1.1 // indirect
	github.com/antlr4-go/antlr/v4 v4.13.1 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/bytedance/sonic/loader v0.3.0 // indirect
	github.com/cloudwego/base64x v0.1.6 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/klauspost/cpuid/v2 v2.2.11 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/richardlehane/mscfb v1.0.4 // indirect
	github.com/richardlehane/msoleps v1.0.4 // indirect
	github.com/stoewer/go-strcase v1.2.1 // indirect
	github.com/tiendc/go-deepcopy v1.6.1 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/xuri/efp v0.0.1 // indirect
	github.com/xuri/nfp v0.0.1 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/otel/metric v1.34.0 // indirect
	go.opentelemetry.io/otel/trace v1.34.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/arch v0.0.0-20210923205945-b76863e36670 // indirect
	golang.org/x/crypto v0.38.0 // indirect
	golang.org/x/exp v0.0.0-20250911091902-df9299821621 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250908214217-97024824d090 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250908214217-97024824d090 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
