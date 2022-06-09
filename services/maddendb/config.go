package maddendb

import (
	"../services/dbutils"
)

//sets up a postgres madden from environment variables
func BuildPostgresMaintenanceFromEnvironment() (Maintenance, error) {
	config, err := dbutils.NewConfigFromEnvironment()
	if err != nil {
		return nil, err
	}
	db, err := dbutils.NewGormDb(config)
	if err != nil {
		return nil, err
	}
	return NewPostgresMaintenace(db), nil
}
