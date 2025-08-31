package config

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

var (
	clickhouseConn driver.Conn
	once           sync.Once
)

func Clickhouse() (driver.Conn, error) {
	var err error
	once.Do(func() {
		ctx := context.Background()
		clickhouseConn, err = clickhouse.Open(&clickhouse.Options{
			Addr: []string{fmt.Sprintf("%s:%s", Config("CLICKHOUSE_HOST"), Config("CLICKHOUSE_PORT"))},
			Auth: clickhouse.Auth{
				Database: Config("CLICKHOUSE_DB"),
				Username: Config("CLICKHOUSE_USER"),
				Password: Config("CLICKHOUSE_PASSWORD"),
			},
			MaxOpenConns:    10,
			MaxIdleConns:    5,
			ConnMaxLifetime: 1 * time.Hour,
		})
		if err != nil {
			return
		}
		if pingErr := clickhouseConn.Ping(ctx); pingErr != nil {
			err = pingErr
		}
	})
	return clickhouseConn, err
}
