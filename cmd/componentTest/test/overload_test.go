package test

import (
	"log"
	"testing"
)

func overload(args ...interface{}) {
	log.Printf("start print args: ")
	for _, arg := range args {
		log.Printf("%v ", arg)
	}
	log.Printf("end print args\n\n")
}

type User struct {
	Name string
	Age  int
}

// TestOverload golang重载测试,golang虽然不支持函数的重载，但是可以通过...来模拟重载
func TestOverload(t *testing.T) {
	overload("zhangSan")            // 单参数
	overload("zhangSan", "leeSi")   // 同类型多参数
	overload("zhangSan", 123, User{ // 不同类型多参数
		Name: "Kobe",
		Age:  18,
	})
}
