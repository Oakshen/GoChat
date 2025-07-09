package database

import (
	"fmt"
	"gochat/internal/config"
	"gochat/internal/models/entities"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Connect 连接数据库
func Connect(cfg *config.DatabaseConfig) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// 获取底层的sql.DB对象来配置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxOpenConns(25)        // 最大打开连接数
	sqlDB.SetMaxIdleConns(10)        // 最大空闲连接数
	sqlDB.SetConnMaxLifetime(5 * 60) // 连接最大生存时间（秒）

	return nil
}

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate() error {
	err := DB.AutoMigrate(
		&entities.User{},
		&entities.Room{},
		&entities.Message{},
		&entities.RoomMember{},
		&entities.Attachment{}, // 添加附件表
	)
	if err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}

	// 添加唯一约束
	if !DB.Migrator().HasConstraint(&entities.RoomMember{}, "unique_room_user") {
		err := DB.Migrator().CreateConstraint(&entities.RoomMember{}, "unique_room_user")
		if err != nil {
			return err
		}
	}

	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}

// Close 关闭数据库连接
func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
