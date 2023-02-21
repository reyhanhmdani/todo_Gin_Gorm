package database

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
	"todoGin/config"
	"todoGin/model/entity"
)

func DatabaseInit(ctx context.Context, cfg *config.Config) (*gorm.DB, error) {

	fmt.Printf("%+v\n ", cfg)
	//
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		cfg.DBUsername,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logger.Info,
				Colorful:      true,
			},
		),
	})
	if err != nil {
		panic("Cannot Connect to database")
		//return nil, err
	}
	err = db.AutoMigrate(&entity.Todolist{})
	if err != nil {
		logrus.Error(err)
	}

	//logrus.Info("Database Migrated")

	logrus.Info("Connect to Database")
	return db, err
}
