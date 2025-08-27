package main

import (
	"fmt"
	"math/rand"
)

func main() {
	// 我们可以使用rand对象生成随机数。我们应该为rand对象提供一些种子，以使生成的数量不同。
	// 如果我们不提供种子，那么编译器将始终产生相同的结。
    // 测试下来 go 库函数随机性还是挺大的
	fmt.Println("A number from 1-100", rand.Intn(100))
}
