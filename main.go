package main

import (
	"fmt"
	"log"

	"github.com/icodeologist/url_shortner/database"
	"gorm.io/gorm"
)

// repetitive actions written with simple funcitons
func printDBData(db *gorm.DB) {
	var users []database.URL

	res := db.Find(&users)

	if res.Error != nil {
		log.Fatal(res.Error)
	}

	for _, u := range users {
		fmt.Printf("ID %v Longurl %v  clicks %v\n", u.ID, u.LongUrl)
	}

}

func main() {
	db, err := database.ConnectToPSQL()
	if err != nil {
		log.Fatal(err)
	}
	printDBData(db)
}
