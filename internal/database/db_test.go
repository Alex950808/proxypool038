package database

import (
	"fmt"
	"log"
	"testing"
)

func TestConnect(t *testing.T) {
	//t.SkipNow()
	//connect()
	InitTables()
}

func TestGetAllProxies(t *testing.T) {
	connect()
	proxies := GetAllProxies();
	fmt.Println(proxies.Len())
	fmt.Println(proxies[0])
}

func TestDeleteProxyList(t *testing.T) {
	connect()
	if err := DB.Delete(&Proxy{},"id > ?",1); err != nil{
		log.Print("Delete failed", err)
	}

}
