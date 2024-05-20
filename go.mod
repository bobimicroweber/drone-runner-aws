module github.com/drone-runners/drone-runner-aws

go 1.21

toolchain go1.21.5

replace github.com/docker/docker => github.com/docker/engine v17.12.0-ce-rc1.0.20200309214505-aa6a9891b09c+incompatible

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.5.0
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.2.2
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute v1.0.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork v1.1.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources v1.1.1
	github.com/Masterminds/squirrel v1.5.4
	github.com/aws/aws-sdk-go v1.44.248
	github.com/buildkite/yaml v2.1.0+incompatible
	github.com/cenkalti/backoff/v4 v4.2.1
	github.com/dchest/uniuri v1.2.0
	github.com/drone/drone-go v1.7.1
	github.com/drone/envsubst v1.0.3
	github.com/drone/runner-go v1.12.0
	github.com/drone/signal v1.0.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-chi/chi/v5 v5.0.8
	github.com/google/go-cmp v0.6.0
	github.com/google/uuid v1.3.0
	github.com/google/wire v0.5.0
	github.com/harness/lite-engine v0.5.69
	github.com/hashicorp/nomad/api v0.0.0-20230421025320-b4e6a70fe69b
	github.com/jmoiron/sqlx v1.3.5
	github.com/joho/godotenv v1.5.1
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/lib/pq v1.10.8
	github.com/maragudk/migrate v0.4.3
	github.com/mattn/go-isatty v0.0.18
	github.com/mattn/go-sqlite3 v1.14.16
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.19.0
	github.com/rs/zerolog v1.29.1
	github.com/sirupsen/logrus v1.9.0
	github.com/stretchr/testify v1.9.0
	github.com/syndtr/goleveldb v1.0.0
	github.com/wings-software/dlite v1.0.0-rc.10
	golang.org/x/exp v0.0.0-20230420155640-133eef4313cb
	golang.org/x/oauth2 v0.16.0
	golang.org/x/sync v0.3.0
	google.golang.org/api v0.126.0
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	gopkg.in/yaml.v2 v2.4.0
)

require (
	cloud.google.com/go/compute v1.20.1 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.3.0 // indirect
	github.com/AzureAD/microsoft-authentication-library-for-go v1.0.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/corpix/uarand v0.2.0 // indirect
	github.com/golang-jwt/jwt/v4 v4.5.0 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/s2a-go v0.1.4 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.2.3 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/harness/godotenv/v3 v3.0.0 // indirect
	github.com/hashicorp/cronexpr v1.1.1 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hetznercloud/hcloud-go/v2 v2.8.0 // indirect
	github.com/icrowley/fake v0.0.0-20221112152111-d7b7e2276db2 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/pkg/browser v0.0.0-20210911075715-681adbf594b8 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.48.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/rogpeppe/go-internal v1.11.0 // indirect
	github.com/shoenig/test v0.6.4 // indirect
	golang.org/x/time v0.3.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230530153820-e85fd2cbaebc // indirect
	gopkg.in/square/go-jose.v2 v2.6.0 // indirect
	k8s.io/utils v0.0.0-20230406110748-d93618cff8a2 // indirect
)

require (
	github.com/99designs/basicauth-go v0.0.0-20230316000542-bf6f9cbbf0f8 // indirect
	github.com/99designs/httpsignatures-go v0.0.0-20170731043157-88528bf4ca7e // indirect
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751 // indirect
	github.com/alecthomas/units v0.0.0-20211218093645-b94a6e3cc137 // indirect
	github.com/bmatcuk/doublestar v1.3.4 // indirect
	github.com/coreos/go-semver v0.3.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/digitalocean/godo v1.98.0
	github.com/docker/distribution v2.8.2+incompatible // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/gofrs/uuid v4.4.0+incompatible // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/googleapis/gax-go/v2 v2.11.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/natessilva/dag v0.0.0-20180124060714-7194b8dcc5c4 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.opencensus.io v0.24.0 // indirect
	golang.org/x/crypto v0.22.0 // indirect
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20230530153820-e85fd2cbaebc // indirect
	google.golang.org/grpc v1.55.0 // indirect
	google.golang.org/protobuf v1.32.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
