package main

import (
	_ "bookstore/internal/store" // internal/store 将自身注册到 factory 中
	"bookstore/server"
	"bookstore/store/factory"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// 测试添加：curl -X POST -H "Content-Type:application/json" -d '{"id": "978-7-111-55842-2", "name": "The Go Programming Language", "authors":["Alan A.A.Donovan", "Brian W. Kergnighan"],"press": "Pearson Education"}' http://localhost:8080/book
// 测试查询：curl -X GET -H "Content-Type:application/json" localhost:8080/book/978-7-111-55842-2
// 测试查询全部：curl -X GET -H "Content-Type:application/json" localhost:8080/book
// 测试更新：curl -X POST -H 'Content-Type: application/json' -d '{"id": "1", "name": "golang", "press": "hha"}' http://localhost:8080/book
// 测试删除：curl -X GET -H "Content-Type:application/json" localhost:8080/book/978-7-111-55842-2
func main() {
	// 创建图书数据存储模块实例
	store, err := factory.New("mem")
	if err != nil {
		panic(err)
	}

	// 创建 http 服务实例
	bookStoreServer := server.NewBookStoreServer(":8080", store)

	// 运行 http 服务
	errChan, err := bookStoreServer.ListenAndServe()
	if err != nil {
		log.Println("web server start failed: ", err)
		return
	}
	log.Println("web server start ok")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	// 监视来自 errChan 以及 c 的事件
	select {
	case err = <-errChan:
		log.Println("web server run failed: ", err)
		return
	case <-c:
		log.Println("bookstore program is exiting...")
		ctx, cf := context.WithTimeout(context.Background(), time.Second)
		defer cf()
		err = bookStoreServer.Shutdown(ctx)
	}

	if err != nil {
		log.Println("bookstore program exit error: ", err)
		return
	}
	log.Println("bookstore program exit ok")
}
