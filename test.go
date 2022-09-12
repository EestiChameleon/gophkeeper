package main

import (
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	//err := migration.InitMigration()
	//if err != nil {
	//	log.Fatalln(err)
	//}
	// working code ahead
	//
	//pool, err := pgxpool.Connect(context.Background(), "postgresql://localhost:5432/yandex_practicum_db?sslmode=disable")
	//if err != nil {
	//	log.Fatalln(err)
	//}
	////fmt.Println(pool.Ping(context.Background()))
	//c := new(models.Card)
	////err = GetRow(pool, "card_by_title", c, "ctitle2", 3)
	////if err != nil {
	////	log.Fatalln(err)
	////}
	//
	//err = pool.QueryRow(context.Background(),
	//	"SELECT id, user_id, title, number, expiration_date, comment, version, deleted_at FROM gk_card "+
	//		"WHERE title = $1 and user_id = $2 ORDER BY version DESC LIMIT 1;",
	//	"ctitle2", 3).Scan(&c.ID, &c.UserID, &c.Title, &c.Number, &c.ExpirationDate, &c.Comment, &c.Version, &c.DeletedAt)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//fmt.Println("is deleted:", c.DeletedAt.Valid)
	//fmt.Println(c)

}
