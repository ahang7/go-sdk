package cli

// FlagInterface 命令行读取配置
type FlagInterface interface {
	// Flags 添加命令行
	Flags() (fs FlagSet)
	// Validate 验证
	Validate() []error
}
