package main

import (
	"fmt"
)

type MyError struct {
	Msg  string
	Code int
}

func (e *MyError) Error() string {
	return fmt.Sprintf("msg:%s,code:%d", e.Msg, e.Code)
}

func Test() (string, MyError) {
	e := MyError{
		"page not found",
		404,
	}
	return "", e
}

func main() {
	var b1 bool = true
	fmt.Printf("b1: %v\n", b1)
	var b2 bool
	b2 = false
	fmt.Printf("b2: %v\n", b2)
	fmt.Println("--------------")

	s := []int{1, 2, 3}
	s = append(s, 100)
	fmt.Printf("s: %v\n", s)
	fmt.Println("-------------")

	m := make(map[string]string)

	m["name"] = "tom"
	m["age"] = "20"
	m["email"] = "tom@gmail.com"
	fmt.Printf("m: %v\n", m)

	delete(m, "name")
	fmt.Printf("m: %v\n", m)

	fmt.Println("-----------")

	_, me := Test()
	fmt.Printf("me: %v\n", me)

}
