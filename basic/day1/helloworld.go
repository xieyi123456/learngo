package main

import "fmt"

func main() {
	fmt.Println("qwe\tqwe")
	fmt.Println("---------------")
	fmt.Println("aaa\naaa")
	fmt.Println("---------------")
	fmt.Println("abcdef\rqqq") //从当前行最前面开始输出，覆盖掉之前的内容
	fmt.Println("---------------")
	fmt.Println("qwe\\qwe")
	fmt.Println("---------------")
	fmt.Println("qwe\"qwe")
}
