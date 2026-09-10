module github.com/elastic/spigot

go 1.26.7

require (
	github.com/andrewkroh/sys v0.0.0-20151128191922-287798fe3e43
	github.com/aws/aws-sdk-go v1.55.8
	github.com/aws/aws-sdk-go-v2/config v1.33.4
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.23.5
	github.com/aws/aws-sdk-go-v2/service/s3 v1.113.0
	github.com/elastic/beats/v7 v7.0.0-alpha2.0.20260907165736-d1f6783e7314
	github.com/elastic/go-ucfg v0.9.2
	github.com/google/uuid v1.6.0
	github.com/stretchr/testify v1.12.1
	go.uber.org/multierr v1.11.0
)

require (
	github.com/aws/aws-sdk-go-v2 v1.47.0 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.7.20 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.20.4 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.20.0 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.5.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.8.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.5.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.19 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.11.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.14.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.20.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/signin v1.10.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.38.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.43.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.50.0 // indirect
	github.com/aws/smithy-go v1.28.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/elastic/elastic-agent-libs v0.46.3 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	go.elastic.co/ecszap v1.0.3 // indirect
	go.uber.org/zap v1.28.0 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.41.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

// Because we need to import Beats/v7 for the windows/stream hook there is a problem
// and because of that we need to do the following replace directive, see https://github.com/elastic/beats/issues/21188
replace github.com/Shopify/sarama => github.com/elastic/sarama v1.24.1-elastic
