package utils

//使用gorm连接数据库
import (
	_ "database/sql"
	"fmt"
	"yunshu_Forum/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// 定义全局数据库连接对象
type DB struct {
	Master *gorm.DB
	Slave  *gorm.DB
}

var db *DB

func DBinit() {
	var err error
	db = &DB{}
	// 初始化主库连接
	db.Master, err = gorm.Open("mysql", config.Database.Master.Dsn)
	if err != nil {
		fmt.Println("主库连接失败: ", err)
		return
	}
	// 初始化从库连接
	db.Slave, err = gorm.Open("mysql", config.Database.Slave.Dsn)
	if err != nil {
		fmt.Println("从库连接失败: ", err)
		return
	}
	// 设置数据库连接池
	db.Master.DB().SetMaxIdleConns(config.Database.Maxidleconn)
	db.Master.DB().SetMaxOpenConns(config.Database.Maxopenconn)
	//db.Master.DB().SetConnMaxLifetime(config.Database.Maxlifetime*time.Second)
	//db.Master.DB().SetMaxIdleConns(config.Database.Maxidleconn*time.Second)
	db.Slave.DB().SetMaxIdleConns(config.Database.Maxidleconn)
	db.Slave.DB().SetMaxOpenConns(config.Database.Maxopenconn)
	//db.Slave.DB().SetConnMaxLifetime(config.Database.Maxlifetime*time.Second)
	//db.Slave.DB().SetMaxIdleConns(config.Database.Maxidleconn*time.Second)
}

func GetMasterDB() *gorm.DB {
	return db.Master
}

func GetSlaveDB() *gorm.DB {
	return db.Slave
}
