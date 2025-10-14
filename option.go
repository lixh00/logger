package zap_logger

type Encoder string

const (
	JsonEncoder    Encoder = "json"
	ConsoleEncoder Encoder = "console"
)

type Option func(conf *Config)

func WithEncoder(encoder Encoder) Option {
	return func(conf *Config) {
		conf.Logger.Encoder = string(encoder)
	}
}

func WithLevel(level string) Option {
	return func(conf *Config) {
		conf.Logger.Level = level
	}
}

func WithEnableFile(enable bool) Option {
	return func(conf *Config) {
		conf.File.Enable = enable
	}
}

func WithFilename(filename string) Option {
	return func(conf *Config) {
		conf.File.Filename = filename
	}
}

func WithFileMaxSize(maxSize int) Option {
	return func(conf *Config) {
		conf.File.MaxSize = maxSize
	}
}

func WithFileMaxAge(maxAge int) Option {
	return func(conf *Config) {
		conf.File.MaxAge = maxAge
	}
}

func WithFileMaxBackups(maxBackups int) Option {
	return func(conf *Config) {
		conf.File.MaxBackups = maxBackups
	}
}

func WithFileLocaltime(localtime bool) Option {
	return func(conf *Config) {
		conf.File.LocalTime = localtime
	}
}

func WithFileCompress(compress bool) Option {
	return func(conf *Config) {
		conf.File.Compress = compress
	}
}

func WithConsoleEnable(enable bool) Option {
	return func(conf *Config) {
		conf.Console.Enable = enable
	}
}

func WithConsoleEnableColor(color bool) Option {
	return func(conf *Config) {
		conf.Console.Color = color
	}
}

func WithLokiEnable(enable bool) Option {
	return func(conf *Config) {
		conf.Loki.Enable = enable
	}
}

func WithLokiHost(host string) Option {
	return func(conf *Config) {
		conf.Loki.Host = host
	}
}

func WithLokiPort(port int) Option {
	return func(conf *Config) {
		conf.Loki.Port = port
	}
}

func WithLokiSource(source string) Option {
	return func(conf *Config) {
		conf.Loki.Source = source
	}
}

func WithLokiService(service string) Option {
	return func(conf *Config) {
		conf.Loki.Service = service
	}
}

func WithLokiJob(job string) Option {
	return func(conf *Config) {
		conf.Loki.Job = job
	}
}

func WithLokiEnvironment(environment string) Option {
	return func(conf *Config) {
		conf.Loki.Environment = environment
	}
}
