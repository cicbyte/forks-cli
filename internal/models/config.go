package models

type AppConfig struct {
	Version   string `yaml:"version"`              // 版本号，用于升级时判断
	Server    string `yaml:"server,omitempty"`     // Forks 服务端地址
	Token     string `yaml:"token,omitempty"`      // API Token
	BackupDir string `yaml:"backup_dir,omitempty"` // 备份目录
	Log struct {
		Level      string `yaml:"level"`
		MaxSize    int    `yaml:"maxSize"`
		MaxBackups int    `yaml:"maxBackups"`
		MaxAge     int    `yaml:"maxAge"`
		Compress   bool   `yaml:"compress"`
	} `yaml:"log"`
}
