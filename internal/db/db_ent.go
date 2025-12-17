package db

import (
	"context"
	"database/sql"
	"fmt"

	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql/schema"
	_ "github.com/go-sql-driver/mysql"

	"game/ent"
)

func CreateDBEnt(dsnChuaDB, tenDB string) error {
	db, err := sql.Open("mysql", dsnChuaDB)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf(
		"CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", tenDB,
	))
	return err
}

func OpenEnt(dsnChuaDB, tenDB string) (*ent.Client, error) {
	full := fmt.Sprintf("%s%s?charset=utf8mb4&parseTime=True&loc=Local", dsnChuaDB, tenDB)
	return ent.Open(dialect.MySQL, full)
}

func RunMigrationEnt(ctx context.Context, client *ent.Client) error {
	return client.Schema.Create(
		ctx,
		schema.WithForeignKeys(true),
		schema.WithDropIndex(true),
		schema.WithDropColumn(false),
	)
}
