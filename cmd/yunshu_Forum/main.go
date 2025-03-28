package main

import (
	"fmt"
	_ "time"
	"yunshu_Forum/config"
	"yunshu_Forum/internal/utils"
)

func main() {
	config.InitConfig()
	utils.DBinit()
	utils.OssInit()
	utils.RedisInit()
	testDb()
	//printall()

}

func testDb() {
	//mysql主数据库测试
	db := utils.GetMasterDB()
	var count int
	db.Table("users").Count(&count)
	fmt.Println("users表总行数：", count)
	//mysql从数据库测试
	db_slave := utils.GetSlaveDB()
	db_slave.Table("users").Count(&count)
	fmt.Println("users表总行数：", count)
	//OSS测试
	name, err_oss := utils.UploadFileToOss("test.txt", "hello world")
	fmt.Println(name)
	if err_oss != nil {
		fmt.Println(err_oss)
	}
	url := utils.GetUrl("test.txt")
	fmt.Println("文件URL：", url)
	err_del := utils.DeleteFileFromOss("test.txt")
	if err_del != nil {
		fmt.Println(err_del)
	}
	//redis测试
	redis := utils.GetRedis()
	redis.Set("test", "hello world", 0)
	val, _ := redis.Get("test").Result()
	fmt.Println("Redis测试：", val)
}
func printall() {
	fmt.Println("服务器监听端口号：", config.Server.Addr)
	fmt.Println("运行模式：", config.Server.Mode)
	fmt.Println("数据库主库DSN：", config.Database.Master.Dsn)
	fmt.Println("数据库从库DSN：", config.Database.Slave.Dsn)
	fmt.Println("数据库最大空闲连接数：", config.Database.Maxidleconn)
	fmt.Println("数据库最大打开连接数：", config.Database.Maxopenconn)
	fmt.Println("Redis地址：", config.Redis.Addr)
	fmt.Println("Redis密码：", config.Redis.Password)
	fmt.Println("Redis数据库：", config.Redis.DB)
	fmt.Println("Redis连接池大小：", config.Redis.Poolsize)
	fmt.Println("JWT密钥：", config.JWT.Secret)
	fmt.Println("JWT过期时间：", config.JWT.Expire)
	fmt.Println("OSS配置ID：", config.Oss.AccessKeyId)
	fmt.Println("OSS配置Secret：", config.Oss.AccessKeySecret)
	fmt.Println("OSS配置Bucket：", config.Oss.Bucket)
	fmt.Println("OSS配置Url：", config.Oss.Url)
	fmt.Println("Email配置提供商：", config.Email.Privider)
	fmt.Println("Email配置SMTP地址：", config.Email.Host)
	fmt.Println("Email配置SMTP端口：", config.Email.Port)
	fmt.Println("Email配置账号：", config.Email.Username)
	fmt.Println("Email配置密码：", config.Email.Password)
}
