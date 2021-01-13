package basic

import (
	"github.com/sidazhang123/f10-go/basic/config"
)

var pluginFuncs []func()

func Init(opts ...config.Option) {
	config.Init(opts...)
	for _, f := range pluginFuncs {
		f()
	}

}
func Register(f func()) {
	pluginFuncs = append(pluginFuncs, f)
}
