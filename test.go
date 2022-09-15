package main

import (
	"fmt"
	"github.com/EestiChameleon/gophkeeper/models"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	data := new(models.ActualData)
	fmt.Println(data)
}
