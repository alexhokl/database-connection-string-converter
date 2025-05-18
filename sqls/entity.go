package sqls

type Connection struct {
	Alias   string            `yaml:"alias"`
	Driver  string            `yaml:"driver"`
	Proto   string            `yaml:"proto"`
	User    string            `yaml:"user"`
	Passwd  string            `yaml:"passwd"`
	Host    string            `yaml:"host"`
	Port    int               `yaml:"port"`
	DBName  string            `yaml:"dbName"`
	Options map[string]string `yaml:"options,omitempty"`
}

type Config struct {
	LowercaseKeywords bool         `yaml:"lowercaseKeywords"`
	Connections       []Connection `yaml:"connections"`
}
