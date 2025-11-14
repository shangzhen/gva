package config

type Elasticsearch struct {
	Addresses   []string `mapstructure:"addresses" json:"addresses" yaml:"addresses"`                   // Elasticsearch 集群地址
	Username    string   `mapstructure:"username" json:"username" yaml:"username"`                     // 用户名
	Password    string   `mapstructure:"password" json:"password" yaml:"password"`                     // 密码
	Index       string   `mapstructure:"index" json:"index" yaml:"index"`                             // 默认索引
	MaxRetries  int      `mapstructure:"max-retries" json:"max-retries" yaml:"max-retries"`           // 最大重试次数
	Timeout     int      `mapstructure:"timeout" json:"timeout" yaml:"timeout"`                       // 请求超时时间 (秒)
	Compression string   `mapstructure:"compression" json:"compression" yaml:"compression"`           // gzip压缩
	CertFile    string   `mapstructure:"cert-file" json:"cert-file" yaml:"cert-file"`                 // 证书文件路径
}
