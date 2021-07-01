package handler

import (
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

var MemcacheClient *memcache.Client

func (h *Handler) get(c *gin.Context) {

	since := c.Query("since")
	before := c.Query("before")

	if since == "" || before == "" {
		c.JSON(401, "since or before is invalid")
		return
	}

	x, err := strconv.Atoi(since)
	if err != nil {
		c.JSON(401, "since is invalid")
		return
	}

	y, err := strconv.Atoi(before)
	if err != nil {
		c.JSON(401, "since or before is invalid")
		return
	}

	if x > y || x < 1 || y < 1 {
		c.JSON(401, "since or before is invalid")
		return
	}
	arr, err := findByCash(x, y)

	c.JSON(200, gin.H{
		"fibonacci": arr,
	})

}

func fibonacci(n int) int {
	it, err := MemcacheClient.Get(string(n))
	if err != nil {
		fn := make(map[int]int)
		fn[0] = 0
		fn[1] = 1
		for i := 0; i <= n-1; i++ {
			var f int
			if i < 2{
				f = fn[i]
			}else{
				f = fn[i-1] + fn[i-2]
			}
			fn[i] = f
		}
		MemcacheClient.Set(&memcache.Item{Key: fmt.Sprintf("%i", n), Value: []byte(string(fn[n-1]))})
		return fn[n-1]
	}

	res, err := strconv.Atoi(string(it.Value))
	if err != nil {
		fmt.Errorf("memcache error", err)
	}
	return res

}

func findByCash(since int, before int) ([]int, error) {
	//предполагается что если в кэше лежит наиболее
	_, err := MemcacheClient.Get(string(before))

	if err != nil {
		fibonacci(before)
	}

	arr := make([]int, before-since+1, before-since+1)
	fibvalue := before
	for i := before-since+1; i >= 1; i--  {
		arr[i-1] = getCash(fibvalue)
		fibvalue--
	}

	return arr, nil
}

func getCash(index int) int {
	it, err := MemcacheClient.Get(string(index))
	if err != nil {
		return fibonacci(index)
	}
	res, err := strconv.Atoi(string(it.Value))
	if err != nil {
		log.Fatalln("not int value:", err)
	}

	return res
}
