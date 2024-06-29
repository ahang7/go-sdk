package cli

import (
	"github.com/ahang7/go-sdk/log/zlog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
)

const (
	configFlagName = "config"
	configFileType = "yaml"
)

var configFile string
var configIn string

func init() {
	pflag.StringVarP(&configFile, configFlagName, "c", configFile, "set the configuration file, the default configuration file type is yaml")
}

func addConfigFile(prefixFlag string, configName string, fs *pflag.FlagSet) {
	fs.AddFlag(pflag.Lookup(configFlagName))

	viper.AutomaticEnv()
	viper.SetEnvPrefix(strings.Replace(strings.ToUpper(prefixFlag), "-", "_", -1))
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	cobra.OnInitialize(func() {
		zlog.Info("init viper")
		if configFile != "" {
			viper.SetConfigFile(configFile)
		} else {
			if configIn != "" {
				viper.AddConfigPath(configIn)
			} else {
				// 默认为当前包下的config包
				defaultIn := getRootDir()
				viper.AddConfigPath(defaultIn)
			}
			viper.SetConfigFile(configName)
			viper.SetConfigType(configFileType)
		}

		if err := viper.ReadInConfig(); err != nil {
			zlog.Fatal("viper read config failed", zlog.Errors(err))
		}
	})
}

func getRootDir() string {
	pwd, err := os.Getwd()
	if err != nil {
		zlog.Fatal("get Root Dir failed", zlog.Errors(err))
	}
	var infer func(dir string) string
	infer = func(dir string) string {
		modFile := filepath.Join(dir, "go.mod")
		if exist(modFile) {
			return dir
		}

		parent := filepath.Dir(dir)
		return infer(parent)
	}

	return infer(pwd)
}

func exist(dir string) bool {
	_, err := os.Stat(dir)
	return err == nil || os.IsExist(err)
}

// SetConfigIn 设置配置文件路径
func SetConfigIn(in string) {
	configIn = in
}
