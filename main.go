package main

import (
	"bytedance_SimpleDouyin/dao"
	"bytedance_SimpleDouyin/router"
	"fmt"
)

func main() {
	r := router.InitRouter()
	err := r.Run(":8080")
	if err != nil {
		fmt.Println("run err")
		return
	}
	defer dao.CloseDB()

}
