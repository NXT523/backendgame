package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	gmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DoHiem
type GDoHiem struct {
	ID             int     `gorm:"column:ma_do_hiem;primaryKey;autoIncrement;comment:Mã độ hiếm"`
	TenDoHiem      string  `gorm:"column:ten_do_hiem;size:20;not null;uniqueIndex:uq_dohiem_ten"`
	SoLuong        int     `gorm:"column:so_luong;not null;default:0;index:ix_dohiem_so_luong"`
	MauSac         string  `gorm:"column:mau_sac;size:7;index:ix_dohiem_mau_sac"`
	SatThuongBonus float64 `gorm:"column:sat_thuong_bonus;not null;default:0"`
	TocDoDanhBonus float64 `gorm:"column:toc_do_danh_bonus;not null;default:0"`
}

func (GDoHiem) TableName() string { return "do_hiem" }

// He
type GHe struct {
	ID    int     `gorm:"column:ma_he;primaryKey;autoIncrement;comment:Mã hệ"`
	TenHe string  `gorm:"column:ten_he;size:50;not null;uniqueIndex:uq_he_ten"`
	MoTa  *string `gorm:"column:mo_ta;size:255"`
}

func (GHe) TableName() string { return "he" }

// LoaiVuKhi
type GLoaiVuKhi struct {
	ID      int     `gorm:"column:ma_loai;primaryKey;autoIncrement"`
	TenLoai string  `gorm:"column:ten_loai;size:50;not null;uniqueIndex:uq_loaivukhi_ten"`
	MoTa    *string `gorm:"column:mo_ta;size:255"`
}

func (GLoaiVuKhi) TableName() string { return "loai_vu_khi" }

// VuKhi
type GVuKhi struct {
	ID             int     `gorm:"column:ma_vu_khi;primaryKey;autoIncrement;comment:Mã vũ khí"`
	TenVuKhi       string  `gorm:"column:ten_vu_khi;size:100;not null;uniqueIndex:uq_vukhi_ten"`
	SatThuongCoBan int     `gorm:"column:sat_thuong_co_ban;not null;default:0;index:ix_vukhi_loai_satthuong,priority:2"`
	TocDoDanh      float64 `gorm:"column:toc_do_danh;not null;default:1.00"`
	TamDanh        int     `gorm:"column:tam_danh;not null;default:1"`
	MoTa           *string `gorm:"column:mo_ta;size:255"`

	MaLoai   int `gorm:"column:ma_loai;not null;index:ix_vukhi_loai_dohiem,priority:1;index:ix_vukhi_loai_he,priority:1;index:ix_vukhi_loai_he_dohiem,priority:1;index:ix_vukhi_loai_satthuong,priority:1"`
	MaDoHiem int `gorm:"column:ma_do_hiem;not null;index:ix_vukhi_loai_dohiem,priority:2;index:ix_vukhi_loai_he_dohiem,priority:3;index:ix_vukhi_ma_dohiem"`
	MaHe     int `gorm:"column:ma_he;not null;index:ix_vukhi_loai_he,priority:2;index:ix_vukhi_loai_he_dohiem,priority:2"`

	Loai   GLoaiVuKhi `gorm:"foreignKey:MaLoai;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT" json:"-"`
	DoHiem GDoHiem    `gorm:"foreignKey:MaDoHiem;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT" json:"-"`
	He     GHe        `gorm:"foreignKey:MaHe;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT" json:"-"`
}

func (GVuKhi) TableName() string {
	return "vu_khi"
}

// Kết nối & migrate
func dsnWithDB(dsnNoDB, dbname string) string {
	return fmt.Sprintf("%s%s?charset=utf8mb4&parseTime=True&loc=Local", dsnNoDB, dbname)
}

func CreateDBGorm(dsnNoDB, dbname string) error {
	admin, err := sql.Open("mysql", dsnNoDB)
	if err != nil {
		return err
	}
	defer admin.Close()
	_, err = admin.Exec(fmt.Sprintf(
		"CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci",
		dbname,
	))
	return err
}

func OpenGorm(dsnNoDB, dbname string) (*gorm.DB, error) {
	dsn := dsnWithDB(dsnNoDB, dbname)
	db, err := gorm.Open(gmysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if sqlDB, err := db.DB(); err == nil {
		sqlDB.SetMaxOpenConns(50)
		sqlDB.SetMaxIdleConns(10)
		_ = sqlDB.Ping()
	}
	return db, nil
}

func RunMigrationGorm(db *gorm.DB) error {
	return db.AutoMigrate(
		&GDoHiem{},
		&GHe{},
		&GLoaiVuKhi{},
		&GVuKhi{},
	)
}
