package snowflake

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestSnowFlake_GenerateID(t *testing.T) {
	var wt sync.WaitGroup
	for i:=0; i<5; i++{
		wt.Add(1)
		go func(num int) {
			for j:=0; j<1000; j++{
				id := SnowFlakeInstance.GenerateID()
				fmt.Println("go-", num, id)
				time.Sleep(10 * time.Millisecond)
			}
			wt.Done()
		}(i)
	}
	wt.Wait()
}
