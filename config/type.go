package config

type mySQL struct {
	Addr     string
	Database string
	Username string
	Password string
	Charset  string
}

type redis struct {
	Addr     string
	Password string
	DB       int
}

type oss struct {
	Endpoint        string
	AccessKeyID     string `mapstructure:"accessKey-id"`
	AccessKeySecret string `mapstructure:"accessKey-secret"`
	BucketName      string
	MainDirectory   string `mapstructure:"main-directory"`
}

// SMTP 配置
type smtp struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
	FromName string `mapstructure:"from_name"`
}

// Verify（验证码）配置
type verify struct {
	CodeLength int `mapstructure:"code_length"`
	TTLSeconds int `mapstructure:"ttl_seconds"`
}

type server struct {
	Addr string
	Port int
}

type config struct {
	MySQL  mySQL
	Redis  redis
	OSS    oss
	Smtp   smtp   `mapstructure:"smtp"`
	Verify verify `mapstructure:"verify"`
	Server server
}
