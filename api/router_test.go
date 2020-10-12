package api

import (
	"fmt"
	"log"
	"testing"
)

func TestGetAllFilePaths(t *testing.T) {
	filenames, err := GetAllFilePaths("../assets/html")
	if err != nil{
		log.Fatal(err)
		return
	}
	fmt.Println(filenames);
}

func TestLoadTemplate(t *testing.T) {
	template, err := LoadTemplate();
	if err != nil {
		fmt.Println("[Testing] Fail to load template")
	}
	println(template)
}