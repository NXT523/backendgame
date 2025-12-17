package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func dsnWithDBRaw(dsnNoDB, dbname string) string {
	return fmt.Sprintf("%s%s?charset=utf8mb4&parseTime=True&loc=Local", dsnNoDB, dbname)
}

// CreateDBRaw: tạo DB nếu chưa có
func CreateDBRaw(dsnNoDB, dbname string) error {
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

// OpenRaw: mở *sql.DB
func OpenRaw(dsnNoDB, dbname string) (*sql.DB, error) {
	dsn := dsnWithDBRaw(dsnNoDB, dbname)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(10)
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// RunMigrationRaw: tạo bảng bằng SQL thuần
func RunMigrationRaw(db *sql.DB) error {
	stmts := []string{
		// 1) do_hiem
		`CREATE TABLE IF NOT EXISTS do_hiem (
			ma_do_hiem INT NOT NULL AUTO_INCREMENT,
			ten_do_hiem VARCHAR(20) NOT NULL,
			so_luong INT NOT NULL DEFAULT 0,
			mau_sac VARCHAR(7),
			sat_thuong_bonus DOUBLE NOT NULL DEFAULT 0,
			toc_do_danh_bonus DOUBLE NOT NULL DEFAULT 0,
			PRIMARY KEY (ma_do_hiem),
			UNIQUE KEY uq_dohiem_ten (ten_do_hiem),

			-- index cần thiết
			KEY ix_dohiem_so_luong (so_luong),
			KEY ix_dohiem_mau_sac (mau_sac)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		// 2) he
		`CREATE TABLE IF NOT EXISTS he (
			ma_he INT NOT NULL AUTO_INCREMENT,
			ten_he VARCHAR(50) NOT NULL,
			mo_ta VARCHAR(255),
			PRIMARY KEY (ma_he),
			UNIQUE KEY uq_he_ten (ten_he)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		// 3) loai_vu_khi
		`CREATE TABLE IF NOT EXISTS loai_vu_khi (
			ma_loai INT NOT NULL AUTO_INCREMENT,
			ten_loai VARCHAR(50) NOT NULL,
			mo_ta VARCHAR(255),
			PRIMARY KEY (ma_loai),
			UNIQUE KEY uq_loaivukhi_ten (ten_loai)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,

		// 4) vu_khi
		`CREATE TABLE IF NOT EXISTS vu_khi (
			ma_vu_khi INT NOT NULL AUTO_INCREMENT,
			ten_vu_khi VARCHAR(100) NOT NULL,
			sat_thuong_co_ban INT NOT NULL DEFAULT 0,
			toc_do_danh DOUBLE NOT NULL DEFAULT 1.00,
			tam_danh INT NOT NULL DEFAULT 1,
			mo_ta VARCHAR(255),
			ma_loai INT NOT NULL,
			ma_do_hiem INT NOT NULL,
			ma_he INT NOT NULL,
			PRIMARY KEY (ma_vu_khi),

			UNIQUE KEY uq_vukhi_ten (ten_vu_khi),

			-- composite index cần thiết
			KEY ix_vukhi_loai_dohiem (ma_loai, ma_do_hiem),
			KEY ix_vukhi_loai_he (ma_loai, ma_he),
			KEY ix_vukhi_loai_he_dohiem (ma_loai, ma_he, ma_do_hiem),
			KEY ix_vukhi_loai_satthuong (ma_loai, sat_thuong_co_ban),
			KEY ix_vukhi_ma_dohiem (ma_do_hiem),

			CONSTRAINT fk_vukhi_loai FOREIGN KEY (ma_loai) REFERENCES loai_vu_khi(ma_loai)
				ON UPDATE RESTRICT ON DELETE RESTRICT,
			CONSTRAINT fk_vukhi_dohiem FOREIGN KEY (ma_do_hiem) REFERENCES do_hiem(ma_do_hiem)
				ON UPDATE RESTRICT ON DELETE RESTRICT,
			CONSTRAINT fk_vukhi_he FOREIGN KEY (ma_he) REFERENCES he(ma_he)
				ON UPDATE RESTRICT ON DELETE RESTRICT
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;`,
	}

	for _, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			return err
		}
	}
	return nil
}
