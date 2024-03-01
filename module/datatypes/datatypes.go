package datatypes

type Config struct {
	HTTP     HTTPConfig  `yaml:"http"`
	Database DBConfig    `yaml:"database"`
	Other    OtherConfig `yaml:"other"`
}

type OtherConfig struct {
	Debug   bool `yaml:"debug"`
	Dryrun  bool `yaml:"dryrun"`
	Version bool `yaml:"version"`
}

type HTTPConfig struct {
	Address string `yaml:"address"`
}

type DBConfig struct {
	Master                 DBNodeConfig `yaml:"master"`
	Replica                DBNodeConfig `yaml:"replica"`
	MaxIdleConns           int          `yaml:"maxIdleConns"`
	MaxOpenConns           int          `yaml:"maxOpenConns"`
	MaxConnLifetimeSeconds int          `yaml:"maxConnLifetimeSeconds"`
	Backend                string       `yaml:"backend"`
}
type DBNodeConfig struct {
	DSN string `yaml:"DSN"`
}
