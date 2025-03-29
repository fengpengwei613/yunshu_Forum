package utils

//使用gorm连接数据库
import (
	"yunshu_Forum/config"
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
	"gorm.io/plugin/dbresolver"

)


var db *gorm.DB

func DBinit() {
	var err error
    db, err = gorm.Open(mysql.Open(config.Database.Master.Dsn), &gorm.Config{})
    if err != nil {
        panic("主库连接失败: " + err.Error())
    }

    // 注册读写分离中间件
    db.Use(dbresolver.Register(dbresolver.Config{
        Sources:  []gorm.Dialector{mysql.Open(config.Database.Master.Dsn)},
        Replicas: []gorm.Dialector{mysql.Open(config.Database.Slave.Dsn)},
        Policy:   dbresolver.RandomPolicy{},
    }))

    // 设置连接池
    sqlDB, _ := db.DB()
    sqlDB.SetMaxIdleConns(config.Database.Maxidleconn)
    sqlDB.SetMaxOpenConns(config.Database.Maxopenconn)
}

func GetDB() *gorm.DB {
	return db
}









// // 定义全局数据库连接对象
// type DB struct {
// 	Master *gorm.DB
// 	Slave  *gorm.DB
// }

// var db *DB
// func DBinit() {
// 	var err error
// 	db = &DB{}
// 	// 初始化主库连接
// 	db.Master, err = gorm.Open("mysql", config.Database.Master.Dsn)
// 	if err != nil {
// 		fmt.Println("主库连接失败: ", err)
// 		return
// 	}
// 	// 初始化从库连接
// 	db.Slave, err = gorm.Open("mysql", config.Database.Slave.Dsn)
// 	if err != nil {
// 		fmt.Println("从库连接失败: ", err)
// 		return
// 	}
// 	// 设置数据库连接池
// 	db.Master.DB().SetMaxIdleConns(config.Database.Maxidleconn)
// 	db.Master.DB().SetMaxOpenConns(config.Database.Maxopenconn)
// 	//db.Master.DB().SetConnMaxLifetime(config.Database.Maxlifetime*time.Second)
// 	//db.Master.DB().SetMaxIdleConns(config.Database.Maxidleconn*time.Second)
// 	db.Slave.DB().SetMaxIdleConns(config.Database.Maxidleconn)
// 	db.Slave.DB().SetMaxOpenConns(config.Database.Maxopenconn)
// 	//db.Slave.DB().SetConnMaxLifetime(config.Database.Maxlifetime*time.Second)
// 	//db.Slave.DB().SetMaxIdleConns(config.Database.Maxidleconn*time.Second)
// }

// func GetMasterDB() *gorm.DB {
// 	return db.Master
// }

// func GetSlaveDB() *gorm.DB {
// 	return db.Slave
// }
