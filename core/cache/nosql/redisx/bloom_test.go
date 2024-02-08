package redisx

import (
	"fmt"
	"testing"
)

func TestBloom(t *testing.T) {

	client := NewRedisBloom("127.0.0.1:6379", "nhelp", nil)
	key := "data-filter"

	err := client.Reserve(key, BloomError, BloomCapacity)
	if err != nil {
		fmt.Println(err)
	}

	//向布隆过滤器添加元素
	ok, err := client.Add(key, "A")
	fmt.Println("add A = ", ok, err)

	client.Add(key, "B")
	client.Add(key, "C")
	client.Add(key, "D")

	exister, err := client.Exists(key, "C")
	fmt.Println("c bloom=", exister, err)
}
