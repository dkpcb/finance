package data

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitializeDatabase(db *gorm.DB) {
	// テーブルをドロップして再作成
	db.Exec("DROP TABLE IF EXISTS trade_histories")
	db.Exec("DROP TABLE IF EXISTS reference_prices")

	db.AutoMigrate(&TradeHistory{})
	db.AutoMigrate(&ReferencePrice{})
}

func DBConnection() (*gorm.DB, error) {
	dsn := GetDBConfig()
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("DB Error: %w", err))
	}
	return db, nil
}

func GetDBConfig() string {
	user := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	hostname := os.Getenv("DB_HOSTNAME")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_DBNAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, hostname, port, dbname) + "?charset=utf8mb4&parseTime=True&loc=Local"
	return dsn
}

// func DBConnection() *sql.DB {

// 	dsnWithoutDB := getDSNWithoutDB()
// 	createDatabaseIfNotExists(dsnWithoutDB)

// 	dsn := GetDBConfig()
// 	var err error
// 	dialector := mysql.Open(dsn)
// 	option := &gorm.Config{}

// 	// DB接続
// 	if err = dbConnect(dialector, option, 10); err != nil {
// 		log.Fatalln(err)
// 	}

// 	CreateTable_reference_prices(db)
// 	CreateTable_trade_history(db)
// 	sqlDB, err := db.DB()
// 	if err != nil {
// 		panic(fmt.Errorf("DB Error: %w", err))
// 	}
// 	return sqlDB
// }

// // DSN without DB name
// func getDSNWithoutDB() string {

// 	// err := godotenv.Load()
// 	// if err != nil {
// 	// 	log.Fatalf("Error loading .env file: %v", err)
// 	// }

// 	user := os.Getenv("DB_USERNAME")
// 	password := os.Getenv("DB_PASSWORD")
// 	hostname := os.Getenv("DB_HOSTNAME")
// 	port := os.Getenv("DB_PORT")

// 	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local", user, password, hostname, port)
// 	fmt.Printf("DSN without DB: %s\n", dsn)
// 	return dsn
// }

// // データベースが存在しない場合は作成する
// func createDatabaseIfNotExists(dsn string) {
// 	db, err := sql.Open("mysql", dsn)
// 	if err != nil {
// 		log.Fatalf("Error connecting to the database: %v", err)
// 	}
// 	defer db.Close()

// 	dbname := os.Getenv("DB_DBNAME")
// 	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbname))
// 	if err != nil {
// 		log.Fatalf("Error creating database: %v", err)
// 	}
// 	log.Printf("Database %s created or already exists.", dbname)
// }

// func dbConnect(dialector gorm.Dialector, config gorm.Option, count uint) (err error) {
// 	// countで指定した回数リトライする
// 	for count > 1 {
// 		if db, err = gorm.Open(dialector, config); err != nil {
// 			time.Sleep(time.Second * 2)
// 			count--
// 			log.Printf("retry... count:%v\n", count)
// 			continue
// 		}
// 		break
// 	}
// 	// エラーを返す
// 	return err
// }

// // DBのdsnを取得する
// func GetDBConfig() string {

// 	// err := godotenv.Load()
// 	// if err != nil {
// 	// 	log.Fatalf("Error loading .env file: %v", err)
// 	// }

// 	user := os.Getenv("DB_USERNAME")
// 	password := os.Getenv("DB_PASSWORD")
// 	hostname := os.Getenv("DB_HOSTNAME")
// 	port := os.Getenv("DB_PORT")
// 	dbname := os.Getenv("DB_DBNAME")

// 	// 環境変数の値を出力してデバッグ
// 	fmt.Printf("DB_USERNAME: %s\n", user)
// 	fmt.Printf("DB_PASSWORD: %s\n", password)
// 	fmt.Printf("DB_HOSTNAME: %s\n", hostname)
// 	fmt.Printf("DB_PORT: %s\n", port)
// 	fmt.Printf("DB_DBNAME: %s\n", dbname)

// 	if user == "" || password == "" || hostname == "" || port == "" || dbname == "" {
// 		log.Fatal("One or more environment variables are not set: DB_USERNAME, DB_PASSWORD, DB_HOSTNAME, DB_PORT, DB_DBNAME")
// 	}

// 	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, hostname, port, dbname)
// 	fmt.Printf("DSN: %s\n", dsn)

// 	return dsn
// }

// ReferencePrice型のテーブルを作成する
func CreateTable_reference_prices(db *gorm.DB) {
	db.AutoMigrate(&ReferencePrice{})
}

// TradeHistory型のテーブルを作成する
func CreateTable_trade_history(db *gorm.DB) {
	db.AutoMigrate(&TradeHistory{})
}
