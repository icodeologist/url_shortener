package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

type User struct {
	ID       uint   `gorm:"primarykey" json:"id"`
	Username string `gorm:"unique;not null;" json:"username"`
	Email    string `gorm:"unique;not null;" json:"email"`
	Password string `gorm:"unique;not null;"`
	// one to many relation ship with urls
	URLs []URL
}

type URL struct {
	ID      uint   `gorm:"unique;primarykey;autoincrement:true" json:"urlid"`
	LongUrl string `json:"longurl"`
	ShortID string `gorm:"column:shortid;type:varchar(256);default:0" json:"shortid"`
	UserID  uint   `json:"user_id"` // this is the foreighkey
	User    User   `gorm:"foreignKey:UserID"`
}

func ConnectToPSQL() (*gorm.DB, error) {
	user := os.Getenv("USER")
	password := os.Getenv("PASSWORD")
	dbname := os.Getenv("NAME")

	dsn := fmt.Sprintf("user=%v password=%v dbname=%v host=localhost port=5432 sslmode=disable", user, password, dbname)
	//connect to postgres
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&User{})

	return db, nil
}

func DeleteByID(db *gorm.DB, uID uint) {
	result := db.Delete(&User{}, uID)
	if result.Error != nil {
		log.Fatal("Error deleting the record")
	}
	if result.RowsAffected == 0 {
		log.Fatal("Id was not found")
	}
	fmt.Printf("User %v deleted\n", uID)
}
