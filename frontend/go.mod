module github.com/odigos-io/odigos/frontend

go 1.23.0

toolchain go1.24.1

require (
	github.com/99designs/gqlgen v0.17.70
	github.com/gin-contrib/cors v1.7.4
	github.com/gin-gonic/gin v1.10.0
	github.com/glebarez/sqlite v1.11.0
	github.com/go-logr/logr v1.4.2
	github.com/odigos-io/odigos/api v0.0.0
	github.com/odigos-io/odigos/common v1.0.63
	github.com/odigos-io/odigos/destinations v0.0.0-20240223090638-df3328a088bc
	github.com/odigos-io/odigos/k8sutils v1.0.170
	github.com/stretchr/testify v1.10.0
	github.com/vektah/gqlparser/v2 v2.5.23
	go.opentelemetry.io/collector/component v1.28.1
	go.opentelemetry.io/collector/component/componenttest v0.122.1
	go.opentelemetry.io/collector/confmap v1.28.1
	go.opentelemetry.io/collector/confmap/xconfmap v0.122.1
	go.opentelemetry.io/collector/exporter v0.122.1
	go.opentelemetry.io/collector/exporter/exportertest v0.122.1
	go.opentelemetry.io/collector/exporter/otlpexporter v0.122.1
	go.opentelemetry.io/collector/exporter/otlphttpexporter v0.122.1
	go.opentelemetry.io/collector/pdata v1.28.1
	go.opentelemetry.io/collector/receiver/otlpreceiver v0.122.1
	go.opentelemetry.io/collector/receiver/receivertest v0.122.1
	go.opentelemetry.io/otel v1.35.0
	golang.org/x/sync v0.12.0
	k8s.io/api v0.32.3
	k8s.io/apimachinery v0.32.3
	k8s.io/client-go v0.32.3
	sigs.k8s.io/yaml v1.4.0
)

require (
	github.com/cenkalti/backoff/v5 v5.0.2 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/glebarez/go-sqlite v1.21.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-sqlite3 v1.14.24 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	go.opentelemetry.io/collector/receiver/receiverhelper v0.122.1 // indirect
	modernc.org/libc v1.22.5 // indirect
	modernc.org/mathutil v1.5.0 // indirect
	modernc.org/memory v1.5.0 // indirect
	modernc.org/sqlite v1.23.1 // indirect
)

require (
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/golang/snappy v1.0.0 // indirect
	github.com/hashicorp/go-version v1.7.0 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/knadh/koanf/maps v0.1.1 // indirect
	github.com/knadh/koanf/providers/confmap v0.1.0 // indirect
	github.com/knadh/koanf/v2 v2.1.2 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/mostynb/go-grpc-compression v1.2.3 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rs/cors v1.11.1 // indirect
	go.opentelemetry.io/collector v0.122.1 // indirect
	go.opentelemetry.io/collector/config/configauth v0.122.1 // indirect
	go.opentelemetry.io/collector/config/configcompression v1.28.1 // indirect
	go.opentelemetry.io/collector/config/configgrpc v0.122.1 // indirect
	go.opentelemetry.io/collector/config/confighttp v0.122.1 // indirect
	go.opentelemetry.io/collector/config/confignet v1.28.1 // indirect
	go.opentelemetry.io/collector/config/configopaque v1.28.1 // indirect
	go.opentelemetry.io/collector/config/configretry v1.28.1 // indirect
	go.opentelemetry.io/collector/config/configtls v1.28.1 // indirect
	go.opentelemetry.io/collector/consumer v1.28.1
	go.opentelemetry.io/collector/extension v1.28.1 // indirect
	go.opentelemetry.io/collector/featuregate v1.28.1 // indirect
	go.opentelemetry.io/collector/receiver v1.28.1 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.60.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.60.0 // indirect
	go.opentelemetry.io/otel/metric v1.35.0 // indirect
	go.opentelemetry.io/otel/sdk v1.35.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.35.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250218202821-56aae31c358a // indirect
	google.golang.org/grpc v1.71.0 // indirect
)

require (
	github.com/agnivade/levenshtein v1.2.1 // indirect
	github.com/bytedance/sonic v1.12.6 // indirect
	github.com/bytedance/sonic/loader v0.2.1 // indirect
	github.com/cloudwego/base64x v0.1.4 // indirect
	github.com/cloudwego/iasm v0.2.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/emicklei/go-restful/v3 v3.11.0 // indirect
	github.com/evanphx/json-patch/v5 v5.9.11 // indirect
	github.com/fxamacker/cbor/v2 v2.7.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.7 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-openapi/jsonpointer v0.21.0 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-openapi/swag v0.23.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.23.0 // indirect
	github.com/goccy/go-json v0.10.4 // indirect
	github.com/goccy/go-yaml v1.16.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/gnostic-models v0.6.8 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.9 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/pierrec/lz4/v4 v4.1.22 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/sosodev/duration v1.3.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/collector/client v1.28.1 // indirect
	go.opentelemetry.io/collector/component/componentstatus v0.122.1 // indirect
	go.opentelemetry.io/collector/consumer/consumererror v0.122.1 // indirect
	go.opentelemetry.io/collector/consumer/consumererror/xconsumererror v0.122.1 // indirect
	go.opentelemetry.io/collector/consumer/consumertest v0.122.1 // indirect
	go.opentelemetry.io/collector/consumer/xconsumer v0.122.1 // indirect
	go.opentelemetry.io/collector/exporter/exporterhelper/xexporterhelper v0.122.1 // indirect
	go.opentelemetry.io/collector/exporter/xexporter v0.122.1 // indirect
	go.opentelemetry.io/collector/extension/extensionauth v0.122.1 // indirect
	go.opentelemetry.io/collector/extension/xextension v0.122.1 // indirect
	go.opentelemetry.io/collector/internal/sharedcomponent v0.122.1 // indirect
	go.opentelemetry.io/collector/internal/telemetry v0.122.1 // indirect
	go.opentelemetry.io/collector/pdata/pprofile v0.122.1 // indirect
	go.opentelemetry.io/collector/pipeline v0.122.1 // indirect
	go.opentelemetry.io/collector/pipeline/xpipeline v0.122.1 // indirect
	go.opentelemetry.io/collector/receiver/xreceiver v0.122.1 // indirect
	go.opentelemetry.io/otel/trace v1.35.0 // indirect
	golang.org/x/arch v0.13.0 // indirect
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/net v0.37.0 // indirect
	golang.org/x/oauth2 v0.25.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/term v0.30.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	golang.org/x/time v0.7.0 // indirect
	google.golang.org/protobuf v1.36.5 // indirect
	gopkg.in/evanphx/json-patch.v4 v4.12.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/driver/sqlite v1.5.7
	gorm.io/gorm v1.25.12
	k8s.io/klog/v2 v2.130.1 // indirect
	k8s.io/kube-openapi v0.0.0-20241105132330-32ad38e42d3f // indirect
	k8s.io/utils v0.0.0-20241104100929-3ea5e8cea738 // indirect
	sigs.k8s.io/controller-runtime v0.20.3 // indirect
	sigs.k8s.io/json v0.0.0-20241010143419-9aa6b5e7a4b3 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.6.0 // indirect
)

replace (
	github.com/odigos-io/odigos/api => ../api
	github.com/odigos-io/odigos/common => ../common
	github.com/odigos-io/odigos/destinations => ../destinations
	github.com/odigos-io/odigos/k8sutils => ../k8sutils
)
