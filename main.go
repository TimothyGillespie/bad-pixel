package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

var badPixelBytes []byte

var (
	dsn string
	db  *gorm.DB
)

type PageHitLogEntry struct {
	gorm.Model
	IP         string
	UserAgent  string
	Host       string
	RequestUri string
	CreatedAt  time.Time
}

func main() {

	waitSec := os.Getenv("BAD_PIXEL_WAIT_SECONDS")

	if waitSec != "" {
		waitTime, err := strconv.Atoi(waitSec)
		if err != nil {
			panic(err)
		}

		time.Sleep(time.Duration(int64(waitTime) * time.Second.Nanoseconds()))
	}

	dbUsername := os.Getenv("BAD_PIXEL_DB_USER")
	dbPassword := os.Getenv("BAD_PIXEL_DB_PASSWORD")
	dbURI := os.Getenv("BAD_PIXEL_DB_URI")
	dbPort := os.Getenv("BAD_PIXEL_DB_PORT")
	dbDatabase := os.Getenv("BAD_PIXEL_DB_DATABASE")

	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUsername, dbPassword, dbURI, dbPort, dbDatabase)
	db, _ = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	db.Migrator().AutoMigrate(&PageHitLogEntry{})

	badPixelFile, err := os.Open("bad-pixel.png")

	badPixelBytes, err = io.ReadAll(badPixelFile)
	if err != nil {
		fmt.Println(err)
	}

	badPixelFile.Close()

	http.HandleFunc("/", HelloHandler)
	fmt.Println("Server started at port 8080")
	http.ListenAndServe(":8080", nil)
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	var ip string

	passThrougIp := r.Header.Get("X-Real-IP")

	if passThrougIp != "" {
		ip = passThrougIp
	} else {
		ip = r.RemoteAddr
	}

	db.Create(&PageHitLogEntry{
		IP:         ip,
		Host:       r.Host,
		UserAgent:  r.UserAgent(),
		RequestUri: r.RequestURI,
	})

	w.Write(badPixelBytes)
	w.Header().Set("Content-Type", "image/png")
}
