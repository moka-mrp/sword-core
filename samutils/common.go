package samutils

import (
	"fmt"
	"os"
)

func P(i ...interface{} ) {
	for  _,val:=range i{
		fmt.Printf("%+v\r\n",val)
	}

}


func Pd(i ...interface{} ) {
	for  _,val:=range i{
		fmt.Printf("%+v\r\n",val)
	}
	os.Exit(0)
}

