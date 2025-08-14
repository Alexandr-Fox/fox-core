package database

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"strconv"
	"syscall"
)

type DB struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Pass     string `yaml:"pass"`
	Name     string `yaml:"name"`
	Port     int    `yaml:"port,omitempty"`
	Schema   string `yaml:"schema,omitempty"`
	SSL      bool   `yaml:"ssl,omitempty"`
	TimeZone string `yaml:"time_zone,omitempty"`
}

func ReadDatabaseConfig() DB {
	var (
		useK8s bool
		err    error
	)
	v, found := syscall.Getenv("USE_K8S")
	if found {
		useK8s, err = strconv.ParseBool(v)
		if err != nil {
			panic(err)
		}
	} else {
		useK8s = false
	}

	if useK8s {
		host, err := os.ReadFile("/etc/app/postgres/host")
		if err != nil {
			panic(err)
		}
		user, err := os.ReadFile("/etc/app/postgres/user")
		if err != nil {
			panic(err)
		}
		pass, err := os.ReadFile("/etc/app/postgres/pass")
		if err != nil {
			panic(err)
		}
		name, err := os.ReadFile("/etc/app/postgres/name")
		if err != nil {
			panic(err)
		}
		port, err := os.ReadFile("/etc/app/postgres/port")
		if err != nil {
			port = []byte("5432")
		}
		ssl, err := os.ReadFile("/etc/app/postgres/ssl")
		if err != nil {
			ssl = []byte("false")
		}
		timeZone, err := os.ReadFile("/etc/app/postgres/time_zone")
		if err != nil {
			timeZone = []byte("Asia/Yekaterinburg")
		}
		schema, err := os.ReadFile("/etc/app/postgres/schema")
		if err != nil {
			timeZone = []byte("core")
		}
		portNumber, err := strconv.Atoi(string(port))
		if err != nil {
			portNumber = 5432
		}
		sslBool, err := strconv.ParseBool(string(ssl))
		if err != nil {
			sslBool = false
		}

		return DB{
			Host:     string(host),
			User:     string(user),
			Pass:     string(pass),
			Name:     string(name),
			Port:     portNumber,
			SSL:      sslBool,
			Schema:   string(schema),
			TimeZone: string(timeZone),
		}
	}

	var cfg DB

	data, err := os.ReadFile("database.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		panic(err)
	}

	if cfg.Port == 0 {
		cfg.Port = 5432
	}

	if cfg.Schema == "" {
		cfg.Schema = "core"
	}

	if cfg.TimeZone == "" {
		cfg.TimeZone = "Asia/Yekaterinburg"
	}

	return cfg
}

func (d DB) ToDSN() string {
	SSLMode := "disable"
	if d.SSL {
		SSLMode = "enable"
	}
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s search_path=%s",
		d.Host, d.User, d.Pass, d.Name, d.Port, SSLMode, d.TimeZone, d.Schema,
	)
}

func GetCoreDB() *gorm.DB {
	db, err := gorm.Open(postgres.Open(ReadDatabaseConfig().ToDSN()), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}
	return db
}
