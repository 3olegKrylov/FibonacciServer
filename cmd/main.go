package main

import (
	Fibonacci "FibonacciServer"
	"FibonacciServer/pkg/handler"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"log"
)



func main(){

	handler.MemcacheClient = memcache.New("10.0.0.1:11211", "10.0.0.2:11211", "10.0.0.3:11212")
	err := handler.MemcacheClient.Set(&memcache.Item{Key: "1", Value: []byte("0")})
	if err != nil {
		fmt.Errorf("memcache error:", err)
	}

	err = handler.MemcacheClient.Set(&memcache.Item{Key: "2", Value: []byte("1")})
	if err != nil {
		fmt.Errorf("memcache error:", err)
	}

	handlers:= new(handler.Handler)
	srv := new(Fibonacci.Server)

	if err = srv.Run("8000", handlers.InitRoutes()); err != nil{
		log.Fatalln("error, on running server:", err.Error())
	}

	return
}
