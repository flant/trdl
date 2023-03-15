module github.com/werf/trdl/client

go 1.20

require (
	bou.ke/monkey v1.0.2
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2
	github.com/gookit/color v1.5.2
	github.com/inconshreveable/go-update v0.0.0-20160112193335-8152e7eb6ccf
	github.com/rodaine/table v1.1.0
	github.com/spaolacci/murmur3 v1.1.0
	github.com/spf13/cobra v1.6.1
	github.com/spf13/pflag v1.0.5
	github.com/theupdateframework/go-tuf v0.5.2
	github.com/werf/lockgate v0.0.0-20211004100849-f85d5325b201
	github.com/werf/logboek v0.5.4
	gopkg.in/yaml.v3 v3.0.1
	mvdan.cc/xurls v1.1.0
)

require (
	github.com/avelino/slugify v0.0.0-20180501145920-855f152bd774 // indirect
	github.com/gofrs/flock v0.7.1 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/mvdan/xurls v1.1.0 // indirect
	github.com/secure-systems-lab/go-securesystemslib v0.4.0 // indirect
	github.com/syndtr/goleveldb v1.0.1-0.20220721030215-126854af5e6d // indirect
	github.com/xo/terminfo v0.0.0-20210125001918-ca9a967f8778 // indirect
	golang.org/x/crypto v0.0.0-20211117183948-ae814b36b871 // indirect
	golang.org/x/net v0.0.0-20220607020251-c690dde0001d // indirect
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	golang.org/x/text v0.3.7 // indirect
)

replace github.com/theupdateframework/go-tuf => github.com/werf/3p-go-tuf v0.0.0-20230315082915-5fc159235553
