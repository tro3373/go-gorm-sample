package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type SampleTable struct {
	Id         string `gorm:"column:id"`
	Code       string `gorm:"column:code"`
	Name       string `gorm:"column:name"`
	UpdateFlag string `gorm:"column:update_flag"`
}

func (t SampleTable) String() string {
	return fmt.Sprintf("%#+v", t)
}

func sub(args []string) error {
	config, err := NewConfig(args)
	if err != nil {
		return err
	}
	// log.Info("==> Config:", config)
	log.Info("==> DSN: ", config.ToDsn())
	db, err := gorm.Open(mysql.Open(config.ToDsn()), &gorm.Config{})
	if err != nil {
		return err
	}

	// Raw SQLを実行する
	res := db.Raw("SELECT * FROM sample_table WHERE update_flag = ?", 0)
	if res.Error != nil {
		return res.Error
	}
	var deletes []SampleTable
	res = res.Scan(&deletes)
	if res.Error != nil {
		return res.Error
	}
	for i, d := range deletes {
		fmt.Println(i, d)
	}
	return nil
}

func main() {
	if err := sub(os.Args[1:]); err != nil {
		log.Error("Failed to execute")
		panic(err)
	}
}
