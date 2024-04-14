package structures

type Config struct {
	AppMode      AppMode          `yaml:"app_mode"`
	Listen       Listener         `yaml:"listen"`
	Storage      StorageConfig    `yaml:"storage"`
	TestStorage  StorageConfig    `yaml:"test_storage"`
	DockerServer DockerRESTServer `yaml:"docker_rest_server"`
}

type AppMode struct {
	IsTest bool `yaml:"test"`
}

type Listener struct {
	BindIP string `yaml:"bind_ip"`
	Port   string `yaml:"port"`
}

type DockerRESTServer struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type StorageConfig struct {
	Host     string `yaml:"host" env:"DB_HOST" env-default:"my_db"`
	Port     rune   `yaml:"port" env:"DB_PORT" env-default:"5432"`
	Database string `yaml:"database" env:"DB_NAME" env-default:"avito_tech"`
	Username string `yaml:"username" env:"DB_USER" env-default:"admin"`
	Password string `yaml:"password" env:"DB_PASSWORD" env-default:"root"`
}
