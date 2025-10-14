package zap_logger

type Config struct {
	Logger  *Logger  `yaml:"logger"`
	File    *File    `yaml:"file"`
	Console *Console `yaml:"console"`
	Loki    *Loki    `yaml:"loki"`
}

type Logger struct {
	Encoder string `yaml:"encoder"`
	Level   string `yaml:"level"`
}

type File struct {
	Enable     bool   `yaml:"enable"`
	Encoder    string `yaml:"encoder"`
	Level      string `yaml:"level"`
	Filename   string `yaml:"name"`
	MaxSize    int    `yaml:"max_size"`
	MaxAge     int    `yaml:"max_age"`
	MaxBackups int    `yaml:"max_backups"`
	LocalTime  bool   `yaml:"local_time"`
	Compress   bool   `yaml:"compress"`
}

type Console struct {
	Enable  bool   `yaml:"enable"`
	Encoder string `yaml:"encoder"`
	Level   string `yaml:"level"`
	Color   bool   `yaml:"color"`
}

type Loki struct {
	Enable      bool   `yaml:"enable"`
	Encoder     string `yaml:"encoder"`
	Level       string `yaml:"level"`
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	Source      string `yaml:"source"`
	Service     string `yaml:"service"`
	Job         string `yaml:"job"`
	Environment string `yaml:"environment"`
}
