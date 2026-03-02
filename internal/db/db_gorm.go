package db

import (
	"database/sql"
	"fmt"
	"game/internal/utils"

	_ "github.com/go-sql-driver/mysql"
	gmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	dbMySql_gorm = utils.GetenvString("MYSQL_DSN_NO_DB", "root:123456@tcp(localhost:3306)/")
	name_gorm    = utils.GetenvString("MYSQL_DB_GORM", "game1")
)

// --- Structs ---

type GDoHiem struct {
	MaDoHiem       int     `gorm:"column:ma_do_hiem;primaryKey;autoIncrement;comment:Mã độ hiếm"`
	TenDoHiem      string  `gorm:"column:ten_do_hiem;size:20;not null;uniqueIndex:uq_dohiem_ten"`
	SoLuong        int     `gorm:"column:so_luong;not null;default:0;index:ix_dohiem_so_luong"`
	MauSac         *string `gorm:"column:mau_sac;size:7;index:ix_dohiem_mau_sac"`
	SatThuongBonus float64 `gorm:"column:sat_thuong_bonus;not null;default:0;index:ix_dohiem_sat_bonus"`
	TocDoDanhBonus float64 `gorm:"column:toc_do_danh_bonus;not null;default:0;index:ix_dohiem_spd_bonus"`
	CapBac         int     `gorm:"column:cap_bac;not null;unique"`
}

func (GDoHiem) TableName() string { return "do_hiem" }

type GHe struct {
	MaHe  int     `gorm:"column:ma_he;primaryKey;autoIncrement;comment:Mã hệ"`
	TenHe string  `gorm:"column:ten_he;size:50;not null;uniqueIndex:uq_he_ten"`
	MoTa  *string `gorm:"column:mo_ta;size:255"`
}

func (GHe) TableName() string { return "he" }

type GLoaiVuKhi struct {
	MaLoaiVuKhi  int     `gorm:"column:ma_loai_vu_khi;primaryKey;autoIncrement;comment:Mã loại vũ khí"`
	TenLoaiVuKhi string  `gorm:"column:ten_loai_vu_khi;size:50;not null;uniqueIndex:uq_loaivukhi_ten"`
	MoTa         *string `gorm:"column:mo_ta;size:255"`
}

func (GLoaiVuKhi) TableName() string { return "loai_vu_khi" }

type GVuKhi struct {
	MaVuKhi        int     `gorm:"column:ma_vu_khi;primaryKey;autoIncrement;comment:Mã vũ khí"`
	TenVuKhi       string  `gorm:"column:ten_vu_khi;size:100;not null;uniqueIndex:uq_vukhi_ten"`
	SatThuongCoBan int     `gorm:"column:sat_thuong_co_ban;not null;default:0;index:ix_vukhi_loaivukhi_satthuong,priority:2"`
	TocDoDanh      float64 `gorm:"column:toc_do_danh;not null;default:1.0"`
	TamDanh        int     `gorm:"column:tam_danh;not null;default:1"`
	MoTa           *string `gorm:"column:mo_ta;size:255"`
	Version        int     `gorm:"column:version;not null;default:1"`

	// Khai báo các cột khóa ngoại
	MaLoaiVuKhi int `gorm:"column:ma_loai_vu_khi;not null;index:ix_vukhi_loaivukhi_dohiem,priority:1;index:ix_vukhi_loaivukhi_he,priority:1;index:ix_vukhi_loaivukhi_he_dohiem,priority:1;index:ix_vukhi_loaivukhi_satthuong,priority:1"`
	MaDoHiem    int `gorm:"column:ma_do_hiem;not null;index:ix_vukhi_loaivukhi_dohiem,priority:2;index:ix_vukhi_loaivukhi_he_dohiem,priority:3;index:ix_vukhi_ma_dohiem"`
	MaHe        int `gorm:"column:ma_he;not null;index:ix_vukhi_loaivukhi_he,priority:2;index:ix_vukhi_loaivukhi_he_dohiem,priority:2"`

	// Ràng buộc quan hệ: Trỏ từ Vũ Khí sang các bảng danh mục
	LoaiVuKhi GLoaiVuKhi `gorm:"foreignKey:MaLoaiVuKhi;references:MaLoaiVuKhi;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
	DoHiem    GDoHiem    `gorm:"foreignKey:MaDoHiem;references:MaDoHiem;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
	He        GHe        `gorm:"foreignKey:MaHe;references:MaHe;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"-"`
}

func (GVuKhi) TableName() string { return "vu_khi" }

// --- Database Operations ---

func CreateDBGorm() error {
	admin, err := sql.Open("mysql", dbMySql_gorm)
	if err != nil {
		return err
	}
	defer admin.Close()

	_, err = admin.Exec(fmt.Sprintf(
		"CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci",
		name_gorm,
	))
	return err
}

func OpenGorm() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s%s?charset=utf8mb4&parseTime=True&loc=Local", dbMySql_gorm, name_gorm)
	return gorm.Open(gmysql.Open(dsn), &gorm.Config{})
}

func RunMigrationGorm(db *gorm.DB) error {
	// Đợt 1: Tạo các bảng "cha" trước
	err := db.AutoMigrate(&GDoHiem{}, &GHe{}, &GLoaiVuKhi{})
	if err != nil {
		return err
	}

	// Đợt 2: Tạo bảng "con" có khóa ngoại sau cùng
	return db.AutoMigrate(&GVuKhi{})
}
