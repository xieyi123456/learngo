package main

import (
	"fmt"
	"learngo/service"
)

func main() {
	fmt.Println("===== JSON差异比较服务演示 =====")

	// 创建JSON差异比较服务实例
	diffService := service.NewJSONDiffService()

	// 示例JSON字符串
	json1 := `{"name":"Alice","age":30,"address":{"city":"New York"},"hobbies":["reading","swimming"]}`
	json2 := `{"name":"Alice","age":31,"address":{"city":"Boston","zip":"02108"},"hobbies":["reading","cycling"],"email":"alice@example.com"}`

	fmt.Println("原始JSON 1:")
	fmt.Println(json1)
	fmt.Println("\n原始JSON 2:")
	fmt.Println(json2)

	// 比较JSON
	diff, err := diffService.CompareJSON(json1, json2)
	if err != nil {
		fmt.Printf("比较JSON时出错: %v\n", err)
		return
	}

	fmt.Println("\nJSON差异结果:")
	// 输出差异结果
	if len(diff.Added) > 0 {
		fmt.Println("新增:")
		for _, path := range diff.Added {
			fmt.Printf("  + %s\n", path)
		}
	}

	if len(diff.Removed) > 0 {
		fmt.Println("移除:")
		for _, path := range diff.Removed {
			fmt.Printf("  - %s\n", path)
		}
	}

	if len(diff.Changed) > 0 {
		fmt.Println("变更:")
		for path, change := range diff.Changed {
			fmt.Printf("  * %s: %s\n", path, change)
		}
	}

	// 如果没有差异
	if len(diff.Added) == 0 && len(diff.Removed) == 0 && len(diff.Changed) == 0 {
		fmt.Println("两个JSON完全相同")
	}
}
