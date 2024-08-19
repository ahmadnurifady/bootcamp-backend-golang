package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"orchestrator-order/internal/config"
)

type ConnectDatabase interface {
	Conn() *sql.DB
}

type connectDatabase struct {
	db  *sql.DB
	cfg *config.Config
}

func (c *connectDatabase) openConn() error {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.cfg.Host, c.cfg.DbConfig.Port, c.cfg.User, c.cfg.DbConfig.Password, c.cfg.Name)
	db, err := sql.Open(c.cfg.Driver, dsn)
	if err != nil {
		return fmt.Errorf("failed to open connection %v", err.Error())
	}

	c.db = db
	return nil
}

func (c *connectDatabase) Conn() *sql.DB {
	return c.db
}

func NewConnectDatabase(cfg *config.Config) (ConnectDatabase, error) {
	conn := &connectDatabase{cfg: cfg}
	if err := conn.openConn(); err != nil {
		return nil, err
	}
	return conn, nil
}
