package storage

import (
	"database/sql"
	"fmt"
	"os"
	"songsapi/logger"

	_ "github.com/lib/pq"
)


func GetDatabaseURL(withName bool) string {
	if withName {
		return os.Getenv("DB_URL")
	}
	return os.Getenv("POSTGRES_URL")
}


func createDBIfNotExists(dbName string) error {
	connPostgres := GetDatabaseURL(false)

	db, err := sql.Open("postgres", connPostgres)

	if err != nil {
		logger.Err.Println("can't establish connection with postgres - ", err)
		return err
	}

	defer db.Close()

	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = '%s')", dbName)
	err = db.QueryRow(query).Scan(&exists);
	if  err != nil {
		logger.Err.Println("check database existence failed - ", err)
		return err
	}

	if !exists {
		logger.Warn.Println("database doesn't exist - start db creation process...")
		if _, err := db.Exec("CREATE DATABASE " + dbName); err != nil {
			logger.Err.Println("can't create database - ", err)
			return err
		}	
	}
	
	return nil
}

func GetDBConnection() *sql.DB {
	logger.Debug.Println("getting connection to database...")
	connStr := GetDatabaseURL(true)

	if err := createDBIfNotExists(os.Getenv("DB_NAME")); err != nil {
		return nil
	}

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		logger.Err.Println("can't establish connection with postgres - ", err)
		return nil
	}

	logger.Debug.Println("connected to the database!")
	return db
}