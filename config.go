package nebula_amqp_lib

type Config struct {
	Driver   Driver `json:"driver"`
	Enabled  bool   `json:"enabled"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	password string `json:"pwd"`
}
