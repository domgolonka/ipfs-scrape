package database

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"os"
	"strconv"
)

type Client struct {
	*sql.DB
}

const ipfsSchema = `
CREATE TABLE IF NOT EXISTS ipfs (
	image TEXT,
	description TEXT,
	name TEXT,
	cid TEXT NOT NULL,
        constraint ipfs_pk
        unique (description, name, image, cid)
);`

func NewClient() (*Client, error) {
	godotenv.Load()

	dbuser := os.Getenv("db_username")
	dbpass := os.Getenv("db_password")
	dbhost := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	var err error
	dbport, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}
	sqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbhost, dbport, dbuser, dbpass, "blockparty")
	db, err := sql.Open("postgres", sqlConn)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	_, err = db.Exec(ipfsSchema)
	if err != nil {
		return nil, err
	}

	return &Client{db}, nil
}
