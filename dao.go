package awesomeProject

import (
	"database/sql"
	"fmt"
)

type IDao interface {
	CMD(query string, dest interface{}) error
}

type Dao struct {
	DB *sql.DB
}

type DBInfo struct {
	User     string
	Password string
	DbName   string
}

func New(dbInfo *DBInfo) *Dao {
	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s", dbInfo.User, dbInfo.Password, dbInfo.DbName)
	db, _ := sql.Open("postgres", connectionString)

	return &Dao{
		DB: db,
	}
}

func (d *Dao) CMD(query string, dest interface{}) error {
	return nil
}
