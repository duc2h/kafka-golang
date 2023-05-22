package configs

import (
	"fmt"
	"strconv"
)

type Config struct {
	Postgresql *PostgreSQL `yaml:"postgresql" mapstructure:"postgresql"`
}

// PostgreSQL is settings of a PostgreSQL server. It contains almost same fields as mysql.Config,
// but with some different field names and tags.
type PostgreSQL struct {
	Host     string `yaml:"host" mapstructure:"host"`
	Port     string `yaml:"port" mapstructure:"port"`
	Database string `yaml:"database" mapstructure:"database"`
	Username string `yaml:"username" mapstructure:"username"`
	Password string `yaml:"password" mapstructure:"password"`
	SSLMode  string `yaml:"sslmode" mapstructure:"sslmode"`
}

func (m *PostgreSQL) FormatDSN() string {
	port, _ := strconv.Atoi(m.Port) //nolint:errcheck
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", m.Host, port, m.Username, m.Password, m.Database, m.SSLMode)
	return dsn
}
