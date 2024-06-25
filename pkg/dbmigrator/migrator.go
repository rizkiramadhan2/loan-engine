package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type migrations struct {
	files []string
}

type dbMigration struct {
	ID          int64     `db:"id"`
	Name        string    `db:"name"`
	CreatedTime time.Time `db:"created_time"`
}

var (
	db      *sqlx.DB
	sqlPath = ""
)

func main() {
	args := os.Args[1:]
	ctx := context.Background()

	sqlPath = os.Getenv("DBM_SQL_PATH")

	switch true {
	case len(args) >= 1 && args[0] == "create":
		createMigrationScript()
	case len(args) >= 1 && args[0] == "exec":
		runMigration(ctx)
	default:
		fmt.Println("invalid command")
	}
}

func createChangelogTable() error {
	q := `CREATE TABLE IF NOT EXISTS public.db_migration (
		"id" bigserial NOT NULL,
		"name" VARCHAR NOT NULL,
		"created_time" Timestamp Without Time Zone,
		PRIMARY KEY ( "id" ),
		CONSTRAINT "unique_db_migration_name" UNIQUE( "name" )
	);
	`

	_, err := db.ExecContext(context.Background(), q)
	return err
}

func runMigration(ctx context.Context) {
	// initialize config
	_ = godotenv.Load()

	masterDB := ""

	fmt.Println("1. Dev")
	fmt.Println("2. Staging")
	fmt.Println("3. Production")
	fmt.Print("Choose Env : ")

	envSelection := ""
	_, err := fmt.Scan(&envSelection)
	if err != nil {
		log.Fatal(err)
	}

	masterDBKey := ""
	switch envSelection {
	case "2":
		log.Println("using staging connection")
		masterDBKey = "DBM_STAG_MASTER_DB"
	case "3":
		log.Println("using production connection")
		masterDBKey = "DBM_PROD_MASTER_DB"
	default:
		log.Println("using development connection")
		masterDBKey = "DBM_DEV_MASTER_DB"
	}

	masterDB = os.Getenv(masterDBKey)
	if masterDB == "" {
		log.Printf("%s env not set for %v. Please input the DSN: ", masterDBKey, envSelection)
		_, err := fmt.Scan(&masterDB)
		if err != nil {
			log.Fatal(err)
		}
	}

	db, err = sqlx.Connect("postgres", masterDB)
	if err != nil {
		log.Fatal("Could not get Database connection :" + err.Error())
		return
	}

	// INIT CHANGELOG TABLE
	err = createChangelogTable()
	if err != nil {
		log.Fatalf("create changelog table failed err=%v", err)
	}

	if sqlPath == "" {
		log.Printf("SQL Files directory is empty. Please enter the path: ")
		_, err := fmt.Scan(&sqlPath)
		if err != nil {
			log.Fatal(err)
		}
	}

	files, err := ioutil.ReadDir(sqlPath)
	if err != nil {
		log.Fatal(err)
	}

	m := getMigrationsMetadata(files)
	for _, f := range m.files {
		ok, err := isFileMigrated(f)
		if err != nil || ok {
			log.Printf("skipping %s err=%v", f, err)
			continue
		}

		err = executeSQLFile(filepath.Join(sqlPath, f))
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%s executed", f)
	}
}

func executeSQLFile(sqlFile string) error {
	file, err := ioutil.ReadFile(sqlFile)
	if err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		tx.Rollback()
	}()
	for _, q := range strings.Split(string(file), ";") {
		q := strings.TrimSpace(q)
		if q == "" {
			continue
		}
		if _, err := tx.Exec(q); err != nil {
			return err
		}
	}

	if _, err = tx.Exec(db.Rebind(`INSERT INTO db_migration (name, created_time) VALUES (?, now())`), filepath.Base(sqlFile)); err != nil {
		return err
	}

	return tx.Commit()
}

func isFileMigrated(name string) (bool, error) {
	data := dbMigration{}
	err := db.Get(&data, `SELECT id FROM db_migration WHERE name = $1`, name)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	if err == sql.ErrNoRows {
		return false, nil
	}

	return true, nil
}

func createMigrationScript() {
	var name string

	fmt.Print("SQL Migration File Name: ")
	_, err := fmt.Scan(&name)
	if err != nil {
		log.Fatal(err)
	}

	// sanitize
	prefix := time.Now().Format("20060102150405")
	name = strings.ReplaceAll(name, ".sql", "")

	fileName := fmt.Sprintf("%v.%s.sql", prefix, name)
	_, err = os.Create(sqlPath + fileName)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%s created\n", fileName)
}

func getMigrationsMetadata(files []fs.FileInfo) migrations {
	m := migrations{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if filepath.Ext(file.Name()) != ".sql" {
			continue
		}

		check := strings.Split(file.Name(), ".")
		if len(check) < 3 {
			continue
		}

		m.files = append(m.files, file.Name())
	}

	sort.Strings(m.files)

	return m
}
