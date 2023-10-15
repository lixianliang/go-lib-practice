package feature

import (
	"fmt"
	"sync"
	"testing"
)

func TestForLoops(t *testing.T) {
	// go1.21 go test 输出不一致。
	// GOEXPERIMENT=loopvar go test 运行两者情况不一样。
	var wg sync.WaitGroup
	words := []string{"a", "b", "c"}
	for _, word := range words {
		w := word
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println("word", word)
			fmt.Println("copy word", w)
		}()
	}
	wg.Wait()
}
