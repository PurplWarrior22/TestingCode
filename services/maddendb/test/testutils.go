package test

import (
	"fmt"

	"../services/dbutils"
	"../services/utilities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//testing cleanup/setup utilities

var (
	db *gorm.DB
)

//new db instance for internal test cleanup/build
func buildTestDbHook() error {
	username, err := utilities.GetEnvOrError(dbutils.USERNAME_ENV, fmt.Sprintf("%s is required", dbutils.USERNAME_ENV))
	if err != nil {
		return err
	}
	password, err := utilities.GetEnvOrError(dbutils.PASSWORD_ENV, fmt.Sprintf("%s is required", dbutils.PASSWORD_ENV))
	if err != nil {
		return err
	}
	host, err := utilities.GetEnvOrError(dbutils.HOST_ENV, fmt.Sprintf("%s is required", dbutils.HOST_ENV))
	if err != nil {
		return err
	}
	port, err := utilities.GetEnvOrError(dbutils.PORT_ENV, fmt.Sprintf("%s is required", dbutils.PORT_ENV))
	if err != nil {
		return err
	}
	dbName := utilities.GetEnvStringOrDefault(dbutils.DB_NAME_ENV, dbutils.DB_NAME_DEFAULT)
	sslMode := utilities.GetEnvStringOrDefault(dbutils.SSLMODE_ENV, dbutils.SSLMODE_DEFAULT)
	connString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, username, password, dbName, port, sslMode)
	db, err = gorm.Open(postgres.Open(connString), &gorm.Config{})
	return err
}
