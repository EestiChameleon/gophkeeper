package main

import (
	"context"
	"fmt"
	"github.com/EestiChameleon/gophkeeper/models"
	"github.com/georgysavva/scany/pgxscan"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"strconv"

	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

func main() {
	//err := migration.InitMigration()
	//if err != nil {
	//	log.Fatalln(err)
	//}
	// working code ahead

	//pool, err := pgxpool.Connect(context.Background(), "postgresql://localhost:5432/yandex_practicum_db?sslmode=disable")
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//fmt.Println(pool.Ping(context.Background()))
	//p := new(models.Pair)
	//err = GetRow(pool, "pair_by_title", p, "title", 1)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//fmt.Println("is deleted:", p.DeletedAt.Valid)

	m := map[string]struct{}{"one": {}, "two": {}}
	_, ok := m["one"]
	fmt.Println(ok)
}

func GetValue(driver *pgxpool.Pool, funcName string, dest interface{}, args ...interface{}) (err error) {
	sqlQuery := "SELECT * FROM " + funcName + "("
	for i := range args {
		if i > 0 {
			sqlQuery += ", "
		}
		sqlQuery += "$" + strconv.Itoa(i+1)
	}
	sqlQuery += ");"
	if err = driver.QueryRow(context.Background(), sqlQuery, args...).Scan(dest); err != nil {
		log.Println(err)
		return err
	}
	return
}

func GetRow(driver *pgxpool.Pool, funcName string, dest interface{}, args ...interface{}) (err error) {
	sqlQuery := "SELECT * FROM " + funcName + "("
	for i := range args {
		if i > 0 {
			sqlQuery += ", "
		}
		sqlQuery += "$" + strconv.Itoa(i+1)
	}
	sqlQuery += ")"

	if err = pgxscan.Get(context.Background(), driver, dest, sqlQuery, args...); err != nil {
		log.Println(err)
		return err
	}
	return
}

// GetAll returns a table with values from offset till limit params
func All(driver *pgxpool.Pool, funcName string, dest interface{}, args ...interface{}) (err error) {
	sqlQuery := "SELECT * FROM " + funcName + "("
	for i := range args {
		if i > 0 {
			sqlQuery += ", "
		}
		sqlQuery += "$" + strconv.Itoa(i+1)
	}
	sqlQuery += ")"
	if err = pgxscan.Select(context.Background(), driver, dest, sqlQuery, args...); err != nil {
		log.Println(err)
		return err
	}
	return
}

func GetAllUserDataLastVersion(driver *pgxpool.Pool) (*models.ActualProtoData, error) {
	var err error
	data := new(models.ActualProtoData)
	if err = All(driver, "pairs_all_last_version_by_user_id", &data.Pairs, 1); err != nil {
		return nil, err
	}
	if err = All(driver, "texts_all_last_version_by_user_id", &data.Texts, 1); err != nil {
		return nil, err
	}
	if err = All(driver, "bins_all_last_version_by_user_id", &data.Bins, 1); err != nil {
		return nil, err
	}

	if err = All(driver, "cards_all_last_version_by_user_id", &data.Cards, 1); err != nil {
		return nil, err
	}

	return data, nil
}
