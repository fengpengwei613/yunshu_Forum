package utils

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"sync"
	"time"
	"yunshu_Forum/config"
)

type Emailreq struct {
	Aim     string
	Subject string
	Body    string
}
type HandelEerror struct {
	Email Emailreq
	Err   error
}

func SendEmail(email Emailreq) HandelEerror {
	from := config.Email.Username
	pass := config.Email.Password
	host := config.Email.Host
	port := config.Email.Port
	message := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, email.Aim, email.Subject, email.Body))
	auth := smtp.PlainAuth("", from, pass, host)
	err := smtp.SendMail(fmt.Sprintf("%s:%d", host, port), auth, from, []string{email.Aim}, message)
	errhandler := HandelEerror{
		Email: email,
		Err:   err,
	}
	if err != nil {
		fmt.Println("发送邮件失败", err)
	}
	return errhandler
}

type EmailWorker struct {
	EmailChan chan Emailreq
	Errchan   chan HandelEerror
	WaitGroup sync.WaitGroup
	PoolSize  int
}

// 创建一个邮件工作池
func NewEmailWorker(poolSize int) *EmailWorker {
	return &EmailWorker{
		EmailChan: make(chan Emailreq, 100),
		Errchan:   make(chan HandelEerror, 100),
		PoolSize:  poolSize,
	}
}

// 启动协程池
func (e *EmailWorker) Start() {
	for i := 0; i < e.PoolSize; i++ {
		go func() {
			for {
				select {
				case email := <-e.EmailChan:
					errhandler := SendEmail(email)
					fmt.Println("发送邮件")
					if errhandler.Err != nil {
						fmt.Println("发送邮件失败2:", errhandler.Err)
						e.Errchan <- errhandler
					}
					e.WaitGroup.Done()
				}
			}
		}()
	}
}

// 发送邮件
func (e *EmailWorker) SendEmail(email Emailreq) {
	e.WaitGroup.Add(1)
	e.EmailChan <- email
}

// 等待所有任务完成
func (e *EmailWorker) Wait() {

	e.WaitGroup.Wait()
	close(e.EmailChan)
	close(e.Errchan)
	time.Sleep(1 * time.Second)
}

// 监听错误并记录日志
func (e *EmailWorker) ListenError() {
	root := config.GetRootPath()
	file, err_open := os.OpenFile(root+"/log/email.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err_open != nil {
		fmt.Println("打开日志文件失败:", err_open)
		return
	}
	defer file.Close()

	logger := log.New(file, "", log.LstdFlags) // 使用独立 Logger
	for errhandler := range e.Errchan {        // 自动退出循环当通道关闭
		logger.Printf("发送失败 - 收件人:%s 主题:%s 内容:%s 错误:%v\n",
			errhandler.Email.Aim,
			errhandler.Email.Subject,
			errhandler.Email.Body,
			errhandler.Err)
		fmt.Println("已记录邮件发送错误到日志")
	}
}
