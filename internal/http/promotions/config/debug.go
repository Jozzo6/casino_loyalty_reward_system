package config

import (
	"fmt"
	"reflect"
)

func Welcome(c *Config) {

	println("")
	println("\033[47m ◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆ \033[0m")
	println("\033[47m ◆ Running <<", c.ServiceName, ">> on port", c.Port, "◆ \033[0m")
	println("\033[47m ◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆◆ \033[0m")

	fmt.Println("\033[46m ► ENV VARIABLES \033[0m")
	s := reflect.ValueOf(c).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Println("\033[46m ►", typeOfT.Field(i).Name, "=", f.Interface(), "\033[0m")
	}
	println("")
}