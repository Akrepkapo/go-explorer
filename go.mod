module github.com/IBAX-io/go-explorer

go 1.16

require (
	github.com/IBAX-io/go-ibax v0.0.0
	github.com/centrifugal/gocent v2.1.0+incompatible
	github.com/spf13/cobra v1.0.0
	github.com/vmihailenco/msgpack/v5 v5.3.4
	gopkg.in/yaml.v2 v2.3.0
	gorm.io/driver/postgres v1.1.0
	gorm.io/gorm v1.21.13
)

replace github.com/IBAX-io/go-ibax => ../go-ibax
