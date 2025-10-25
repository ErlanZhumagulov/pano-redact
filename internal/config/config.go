package config

type Config struct {
	ServerPort    string
	AdminUsername string
	AdminPassword string
	UploadDir     string
	StaticDir     string
}

func Load() *Config {
	return &Config{
		ServerPort:    ":8080",
		AdminUsername: "Mor.max.c@gmail.com",
		AdminPassword: "Gf71@koki_1%",
		UploadDir:     "static/uploads",
		StaticDir:     "static/templates",
	}
}
