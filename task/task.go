package task

import (
	"fmt"
	"github.com/robfig/cron"
)

func InitTask() {
	c := cron.New()
	spec := "0 0 0 * * ?"
	c.AddFunc(spec, func() {

		fmt.Println("执行完成")
	})
	c.Start()
	select {}
}
