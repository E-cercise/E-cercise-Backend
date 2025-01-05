package config

import (
	"fmt"
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

	//err = migrateEnum(db)
	//if err != nil {
	//	panic(err)
	//}
	//
	//err = db.AutoMigrate(&model.User{}, &model.MasterSkill{}, &model.UserSkill{}, &model.UserProfile{},
	//	&model.ResetPasswordToken{}, &model.TimeSheet{},
	//	&model.Epic{}, &model.OTRequest{},
	//	&model.Section{}, &model.Project{},
	//	&model.Request{}, &model.Position{},
	//	&model.TimeSheetRow{}, &model.Role{},
	//	&model.Permission{}, &model.Section{},
	//	&model.Department{}, &model.Position{},
	//	&model.Address{}, &model.Attachment{},
	//	&model.InternDetail{}, &model.PartTimeDetail{},
	//	&model.Technology{}, &model.University{},
	//	&model.Major{}, &model.Bank{},
	//	&model.UserBank{}, &model.Consent{},
	//)

	if err != nil {
		panic(err)
	}

	return db
}
