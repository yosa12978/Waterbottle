package controllers

import (
	"bufio"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/yosa12978/waterbottle/pkg/db"
)

func DoMigrate(args []string) error {
	filepart := strings.Split(args[0], "_")
	filename := fmt.Sprintf("%s_%s_%s.sql", filepart[0], filepart[1], args[1])

	db, err := db.InitDB(args[2])
	defer db.Close()
	if err != nil {
		return err
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	ready_data := strings.Split(string(data), ";")

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	for _, i := range ready_data {
		if strings.Trim(i, "") == "" {
			continue
		}
		_, err = tx.ExecContext(ctx, i)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	err = UpdateMigrationInfo(db, filepart[0], filepart[1], args[1])
	if err != nil {
		return err
	}

	log.Printf("Migration %s %s %s completed", filepart[0], filepart[1], args[1])
	return nil
}

func UpdateMigrationInfo(db *sql.DB, name string, version string, tpe string) error {
	create_stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS wb_migrations (" +
		"id INT AUTO_INCREMENT, " +
		"name VARCHAR(255) NOT NULL, " +
		"version VARCHAR(255) NOT NULL, " +
		"type VARCHAR(255) NOT NULL, " +
		"created DATE NOT NULL, " +
		"PRIMARY KEY(id)" +
		")")
	if err != nil {
		return err
	}
	create_stmt.Exec()
	create_stmt.Close()

	stmt, err := db.Prepare("INSERT INTO wb_migrations (name, version, type, created) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	stmt.Exec(name, version, tpe, time.Now())
	stmt.Close()
	return nil
}

func Help(args []string) error {
	if len(args) != 0 {
		return errors.New("unknown arguments")
	}
	data, err := os.Open(os.Getenv("WATERBOTTLE_PATH") + "/help/help.txt")
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	return nil
}
