package log

import (
	"bytes"
	"fmt"
	"log"
	"testing"
)

func TestLog(t *testing.T) {
	// log 默认输出到 stderr。
	// log 没有日志级别，会自动加上换行。
	// Printf 格式化输出。
	// Println 多个 args 会以空格分割输出，Print 不会。
	// 含有 Panic Fatal 函数。

	// 默认增加的信息只有时间，可以设置时间格式。
	log.Printf("log msg %d %s", 10, "abc")
	log.Printf("log msg1")

	log.Println(3, "abc")
	log.Print("ab", 3)
	log.Print("ab1", 4)

}

func TestLogFlags(t *testing.T) {
	user := struct {
		Name string
		Age  int
	}{
		"lxl", 36,
	}

	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
	log.Printf("login: name: %s, age: %d", user.Name, user.Age)
}

func TestLogCustomLogger(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	// 第一个参数可以设置为 io.MultiWriter,将日志输出多个 w 中。
	logger := log.New(buf, "", log.Lshortfile|log.LstdFlags)

	user := struct {
		Name string
		Age  int
	}{
		"lxl", 36,
	}
	logger.Printf("login: name: %s, age: %d", user.Name, user.Age)
	fmt.Printf(buf.String())
}
