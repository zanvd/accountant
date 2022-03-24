package framework

import (
	"database/sql"
)

func InitDatabase(config *Config) (*sql.DB, error) {
	return sql.Open(
		"mysql",
		config.Database.Username+":"+config.Database.Password+
			"@("+config.Database.Host+":"+config.Database.Port+")/"+config.Database.Name,
	)
}
