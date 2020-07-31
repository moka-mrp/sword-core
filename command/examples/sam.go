package examples

import (
	"fmt"
	"time"
)

func Sam() {
	fmt.Println("sam is a good man.")
	for i:=1; i<5; i++ {
		fmt.Println("sam阻塞中...")
		time.Sleep(1*time.Second)
	}
}
