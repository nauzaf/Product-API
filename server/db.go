package server

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func connectPostgres(dbURL string) (*sqlx.DB, error) {
	if dbURL == "" {
		return nil, fmt.Errorf("unspecified dburl")
	}

	var db *sqlx.DB
	parsedURL, err := url.Parse(dbURL)
	if err != nil {
		panic(fmt.Errorf("connection string parsing: %w", err))
	}

	var maxIdleConns, maxOpenConns int64
	queryPart := parsedURL.Query()
	if maxIdleConnsStr := queryPart.Get("max_idle_conns"); maxIdleConnsStr != "" {
		queryPart.Del("max_idle_conns")
		maxIdleConns, err = strconv.ParseInt(maxIdleConnsStr, 10, 32)
		if err != nil {
			panic(fmt.Errorf("parse error: %w", err))
		}
	}
	if maxOpenConnsStr := queryPart.Get("max_open_conns"); maxOpenConnsStr != "" {
		queryPart.Del("max_open_conns")
		maxOpenConns, err = strconv.ParseInt(maxOpenConnsStr, 10, 32)
		if err != nil {
			panic(fmt.Errorf("parse error: %w", err))
		}
	}
	if maxIdleConns == 0 {
		maxIdleConns = 2
	}
	if maxOpenConns == 0 {
		maxOpenConns = 8
	}

	parsedURL.RawQuery = queryPart.Encode()
	dbURL = parsedURL.String()

	for {
		db, err = sqlx.Connect("postgres", dbURL)
		if err == nil {
			break
		}
		if !strings.Contains(err.Error(), "connect: connection refused") {
			panic(fmt.Errorf("db connection: %w", err))
		}
		const retryDuration = 5 * time.Second
		time.Sleep(retryDuration)

	}
	if db != nil {
		db.SetMaxIdleConns(int(maxIdleConns))
		db.SetMaxOpenConns(int(maxOpenConns))
	}
	return db, nil
}
