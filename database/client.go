package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
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

func NewClient(host, port, user, password string) (*Client, error) {
	var err error
	dbport, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}
	sqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, dbport, user, password, "blockparty")
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
