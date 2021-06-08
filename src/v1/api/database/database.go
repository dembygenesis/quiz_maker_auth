package database

import (
	"database/sql"
	"fmt"
	string2 "github.com/dembygenesis/quiz_maker_auth/src/app/common/utilities/string_utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

// https://github.com/jmoiron/sqlx]'
var dBInstance *sqlx.DB
var DBInstancePublic *sqlx.DB

type ClassDatabase struct {
	Instance *sqlx.DB
}



type UserListDisplay struct {
	ID           uint   `json:"id" db:"id"`
	FirstName    string `json:"firstname" db:"firstname"`
	LastName     string `json:"lastname" db:"lastname"`
	Email        string `json:"email" db:"email"`
	MobileNumber string `json:"mobile_number" db:"mobile_number"`
	Role         string `json:"role" db:"role"`
	BankType     string `json:"bank_type" db:"bank_type"`
	BankNo       string `json:"bank_no" db:"bank_no"`
	Address      string `json:"address" db:"address"`
	Birthday     string `json:"birthday" db:"birthday"`
	Gender       string `json:"gender" db:"gender"`
	M88Account   string `json:"m88_account" db:"m88_account"`
}


func init() {
	EstablishConnection()
}

func ValidEntry(v string, c string, t string) (bool, error) {
	var count int
	var hasEntry bool

	sql := `
		SELECT COUNT(*) FROM ` + t + `
		WHERE 1 = 1
			AND ` + c + ` = ?
	`

	err := DBInstancePublic.Get(&count, sql, v)

	if err != nil {
		return hasEntry, err
	}

	if count > 0 {
		hasEntry = true
	} else {
		hasEntry = false
	}

	return hasEntry, err
}


func GetPaginationDetails(
	sql string,
	count int,
	page int,
	rowLimit int,
	pageLimit int,
) (string, []int, int, int, int, int, int) {

	var pages []int

	pageStart := 0
	pageEnd := pageLimit
	totalPages := int(math.Ceil(float64(count) / float64(rowLimit)))

	if page > totalPages || totalPages == 1 {
		page = 0
	}

	var rowsPerPage int

	if rowLimit > count {
		rowsPerPage = count
	} else {
		rowsPerPage = rowLimit
	}

	var offset int

	if count != 0 {
		if page == 0 {
			offset = 0
		} else if totalPages == 0 {
			offset = 0
		} else {
			if page >= totalPages {
				offset = 0
			} else {
				offset = page * rowLimit
			}
		}
	} else {
		offset = 0
	}

	for !(page >= 0 && page <= pageEnd) {
		pageStart = pageStart + pageLimit
		pageEnd = pageEnd + pageLimit
	}

	for i := pageStart; i <= pageEnd; i++ {
		if i <= totalPages - 1 {
			pages = append(pages, i)
		} else {
			if len(pages) > 0 {
				previousPage := pages[0] - 1

				if previousPage > 1 {
					pages = append([]int{previousPage}, pages...)
				}
			}
		}
	}

	sql = sql + " LIMIT replace_limit OFFSET replace_offset"

	sql = strings.Replace(sql, "replace_limit", strconv.Itoa(rowsPerPage), -1)
	sql = strings.Replace(sql, "replace_offset", strconv.Itoa(offset), -1)

	if len(pages) == 0 {
		pages = append(pages, 0)
	}

	// Handle result count
	resultCount := 0

	// Check if it has next page
	hasNextPage := false

	for _, ele := range pages {
		if ele > page {
			hasNextPage = true
			break
		}
	}

	if hasNextPage == true {
		resultCount = rowLimit
	} else {
		resultCount = count - offset
	}

	return sql, pages, rowsPerPage, offset, page, count, resultCount
}

func GetQueryCount(
	sql string,
	args ...interface{},
) (int, error) {

	var (
		count int
	)

	sql = "SELECT COUNT(*) FROM (" + sql + ") AS a"

	err := DBInstancePublic.Get(&count, sql, args...)

	if err != nil {
		return count, err
	}

	return count, nil
}

// Sets the global variable db instance
func EstablishConnection() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbDatabase := os.Getenv("DB_DATABASE")
	dbPort := os.Getenv("DB_PORT")

	// Connect to MYSQL and execute queries
	connString := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":"+ dbPort +")/" + dbDatabase + "?parseTime=true"
	db, err := sqlx.Open("mysql", connString)

	if err != nil {
		fmt.Println("Error establishing database connection")
		panic(err.Error())
	}

	maxConnections, _ := strconv.Atoi(os.Getenv("DB_DATABASE"))

	db.SetMaxOpenConns(maxConnections)

	dBInstance = db
	DBInstancePublic = db

	testConnection()
}

// Performs a simple query to see if the connection succeeded
func testConnection() {
	_, err := dBInstance.Query("SELECT 5 AS test")

	if err != nil {
		fmt.Println("Error establishing database connection")
		panic(err.Error())
	}
}

func GetDynamicQuery(s string) ([][]string, error) {
	var results [][]string

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbDatabase := os.Getenv("DB_DATABASE")
	dbPort := os.Getenv("DB_PORT")

	db, err := sql.Open("mysql", dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":"+ dbPort +")/" + dbDatabase + "?parseTime=true")
	defer db.Close()

	if err != nil {
		panic(err)
	}

	rows, err := db.Query(s)

	if err != nil {
		return results, err
	}

	cols, err := rows.Columns()
	if err != nil {
		return results, err
	}

	// Result is your slice string.
	rawResult := make([][]byte, len(cols))
	result := make([]string, len(cols))
	dest := make([]interface{}, len(cols)) // A temporary interface{} slice

	for i, _ := range rawResult {
		dest[i] = &rawResult[i] // Put pointers to each string in the interface slice
	}

	for rows.Next() {
		err = rows.Scan(dest...)
		if err != nil {
			return results, err
		}

		container := make([]string, len(cols))

		for i, raw := range rawResult {
			if raw == nil {
				result[i] = "\\N"
			} else {
				result[i] = string(raw)
			}

			container = append(container, result[i])
		}

		results = append(results, container)
	}

	string2.Dump(results)
	fmt.Println("-==============================")


	panic("GG")

	return results, nil
}