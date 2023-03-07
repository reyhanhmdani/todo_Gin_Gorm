package database

import (
	//"context"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
	"todoGin/config"
)

func DatabaseInit(ctx context.Context, cfg *config.Config) (*gorm.DB, error) {

	fmt.Printf("%+v\n", cfg)
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
	//err = db.AutoMigrate(&entity.Todolist{})
	//if err != nil {
	//	logrus.Error(err)
	//}
	//
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	//ping database to make sure connection is established successfully
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, err
	}

	//logrus.Info("Database Migrated")

	logrus.Info("Connect to Database")
	return db, err
}

// migrate -database "mysql://Raihan:Pastibisa@tcp(localhost:3306)/Gin_todo" -path db/migrations up

//func Migrate(db *gorm.DB) error {
//	logrus.Info("running database migration")
//
//	sqlDB, err := db.DB()
//	if err != nil {
//		return err
//	}
//
//	driver, err := mysqlMigration.WithInstance(sqlDB, &mysqlMigration.Config{})
//	if err != nil {
//		return err
//	}
//
//	m, err := migrate.NewWithDatabaseInstance(
//		"file://database/migrations",
//		"mysql", driver)
//	if err != nil {
//		return err
//	}
//
//	err = m.Up()
//	if err != nil && err == migrate.ErrNoChange {
//		logrus.Info("No schema changes to apply")
//		return nil
//	}
//
//	return err
//}
