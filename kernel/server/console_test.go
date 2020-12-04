package server

import (
	"github.com/robfig/cron"
	"log"
	"testing"
)


func TestStartConsole(t *testing.T) {
	StartConsole(TempRegisterSchedule)
}

func TempRegisterSchedule(c *cron.Cron) {
	c.AddFunc("@every 1s",hello)
}

func hello()  {
	log.Println("hello world")
}


