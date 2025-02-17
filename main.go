package main

import "C"
import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"tinygroupcache"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

// 在当前节点创建缓存空间scores
func createGroup() *tinygroupcache.Group {
	return tinygroupcache.NewGroup("scores", 2<<10, tinygroupcache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
}

// startCacheServer 用来启动缓存服务器：创建 HTTPPool，添加节点信息，注册到缓存空间中，启动 HTTP 服务
func startCacheServerHTTP(addr string, addrs []string, cacheG *tinygroupcache.Group) {
	peers := tinygroupcache.NewHTTPPool(addr)
	peers.Set(addrs...)
	cacheG.RegisterPeers(peers)
	log.Println("tinygroupcache is running at", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

func startCacheServerGrpc(addr string, addrs []string, cacheG *tinygroupcache.Group) {
	peers := tinygroupcache.NewGrpcPool(addr)
	peers.Set(addrs...)
	cacheG.RegisterPeers(peers)
	log.Println("tinygroupcache is running at", addr)
	peers.Run()
}

// startAPIServer 用来启动一个 API 服务，与用户进行交互
func startAPIServer(apiAddr string, cacheG *tinygroupcache.Group) {
	http.Handle("/api", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			view, errG := cacheG.Get(key)
			if errG != nil {
				http.Error(w, errG.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			_, errW := w.Write(view.ByteSlice())
			if errW != nil {
				return
			}

		}))
	log.Println("fontend server is running at", apiAddr)
	log.Fatal(http.ListenAndServe(apiAddr[7:], nil))

}

func startGRPCServer() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "tinygroupcache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	addrMap := map[int]string{
		8001: ":8001",
		8002: ":8002",
		8003: ":8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	cacheG := createGroup()
	if api {
		go startAPIServer(apiAddr, cacheG)
	}
	startCacheServerGrpc(addrMap[port], addrs, cacheG)
}

func startHTTPServer() {
	var port int
	var api bool
	flag.IntVar(&port, "port", 8001, "tinygroupcache server port")
	flag.BoolVar(&api, "api", false, "Start a api server?")
	flag.Parse()

	apiAddr := "http://localhost:9999"
	// 三个端口，相当于三个真实缓存节点
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}

	var addrs []string
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	cacheG := createGroup()
	if api {
		go startAPIServer(apiAddr, cacheG)
	}
	startCacheServerHTTP(addrMap[port], []string(addrs), cacheG)
}

func main() {
	startGRPCServer()
}
