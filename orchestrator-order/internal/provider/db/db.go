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

func (i *connectDatabase) openConn() error {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		i.cfg.Host, i.cfg.DbConfig.Port, i.cfg.User, i.cfg.DbConfig.Password, i.cfg.Name)
	db, err := sql.Open(i.cfg.Driver, dsn)
	if err != nil {
		return fmt.Errorf("failed to open connection %v", err.Error())
	}

	i.db = db
	return nil
}

func (i *connectDatabase) Conn() *sql.DB {
	return i.db
}

func NewConnectDatabase(cfg *config.Config) (ConnectDatabase, error) {
	conn := &connectDatabase{cfg: cfg}
	if err := conn.openConn(); err != nil {
		return nil, err
	}
	return conn, nil
}
