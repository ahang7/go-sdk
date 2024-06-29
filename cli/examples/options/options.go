package options

import (
	"github.com/ahang7/go-sdk/cli"
	"github.com/ahang7/go-sdk/cli/examples/exoptions"
)

type Options struct {
	MySQLOption *exoptions.MySQLOptions `json:"mysql" mapstructure:"mysql"`
}

func (o *Options) Flags() (fs cli.FlagSet) {
	o.MySQLOption.AddFlags(fs.Flags("mysql"))

	return
}

func (o *Options) Validate() []error {
	var errs []error
	errs = append(errs, o.MySQLOption.Validate())

	return errs
}

func NewAppOption() *Options {
	o := &Options{
		MySQLOption: exoptions.NewMySQLOptions(),
	}

	return o
}

var _ cli.FlagInterface = (*Options)(nil)
