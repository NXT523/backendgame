package db

import (
	"context"
	"database/sql"
	"fmt"
	"game/ent"
	"game/internal/utils"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"
	_ "github.com/go-sql-driver/mysql"
)

var (
	dbMySql_ent = utils.GetenvString("MYSQL_DSN_NO_DB", "root:123456@tcp(localhost:3306)/")
	name_ent    = utils.GetenvString("MYSQL_DB_ENT", "game")
)

// Tạo database bằng ENT
func CreateDBEnt() error {

	db, err := sql.Open("mysql", dbMySql_ent)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf(
		"CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", name_ent,
	))
	return err
}

// Kết nối database bằng ENT
func OpenEnt() (*ent.Client, error) {
	duongdan := fmt.Sprintf("%s%s?charset=utf8mb4&parseTime=True&loc=Local", dbMySql_ent, name_ent)

	db, err := sql.Open("mysql", duongdan)
	if err != nil {
		return nil, err
	}

	// Quan trọng: Kiểm tra kết nối thực tế
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("không thể ping tới database: %v", err)
	}

	db.SetMaxOpenConns(1000)
	db.SetMaxIdleConns(100)
	db.SetConnMaxLifetime(time.Hour)

	drv := entsql.OpenDB(dialect.MySQL, db)
	return ent.NewClient(ent.Driver(drv)), nil
}

func RunMigrationEnt(ctx context.Context, client *ent.Client) error {
	return client.Schema.Create(
		ctx,
		schema.WithForeignKeys(true),
		schema.WithDropIndex(true),
		schema.WithDropColumn(false),
	)
}

func DropDBEnt(dsnChuaDB, tenDB string) error {
	db, err := sql.Open("mysql", dsnChuaDB)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf(
		"DROP DATABASE IF EXISTS %s", tenDB,
	))
	return err
}
