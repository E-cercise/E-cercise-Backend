package config

import (
	"fmt"
	"github.com/E-cercise/E-cercise/src/model"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
	"strconv"
)

func DatabaseConnection() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		slog.Error("Error loading .env file")
	}

	port, err := strconv.Atoi(DatabasePort)

	var (
		host     = DatabaseHost
		user     = DatabaseUsername
		password = DatabasePassword
		dbName   = DatabaseName
	)

	sqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Bangkok", host, port, user, password, dbName)

	db, err := gorm.Open(postgres.Open(sqlInfo), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	err = db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error
	if err != nil {
		panic(err)
	}

	err = migrateEnum(db)
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&model.User{}, &model.Equipment{}, &model.EquipmentOption{}, &model.EquipmentFeature{}, &model.Image{},
		&model.Attribute{}, &model.Cart{}, &model.LineEquipment{}, &model.Order{}, &model.MuscleGroup{})

	if err != nil {
		panic(err)
	}

	return db
}

func migrateEnum(db *gorm.DB) error {

	err := db.Exec(`
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'role_type') THEN
				CREATE TYPE role_type AS ENUM (
				 	'USER',
					'ADMIN'
				);
			END IF;
		END$$;
	`).Error

	if err != nil {
		return err
	}

	err = db.Exec(`
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'order_status') THEN
				CREATE TYPE order_status AS ENUM (
				 	'Placed',
					'Paid',
					'Shipped out',
					'Received'
				);
			END IF;
		END$$;
	`).Error

	if err != nil {
		return err
	}

	return nil
}
