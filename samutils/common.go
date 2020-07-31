package samutils

import (
	"fmt"
	"os"
)

func P(i interface{}) {
	fmt.Printf("%+v\r\n",i)
}


func Pd(i interface{}) {
	fmt.Printf("%+v\r\n",i)
	os.Exit(0)
}

