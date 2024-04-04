package antipode

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type MySQL struct {
	dsn       string
	datastore string
	table     string
}

func CreateMySQL(host string, port string, user string, password string, datastore string, table string) MySQL {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, datastore)

	return MySQL{dsn, datastore, table}
}

func (m MySQL) write(ctx context.Context, key string, obj AntiObj) error {

	db, err := sql.Open("mysql", m.dsn)
	if err != nil {
		log.Fatalf("impossible to create the connection: %s", err)
	}
	defer db.Close()

	jsonAntiObj, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("INSERT INTO %s VALUES ( %s, %s )", m.table, key, jsonAntiObj)

	_, err = db.Query(query)

	if err != nil {
		return err
	}

	return err
}

func (m MySQL) read(ctx context.Context, key string) (AntiObj, error) {

	// Connect to MySQL database
	db, err := sql.Open("mysql", "root:password@tcp(mysql1:3306)/mydatabase")
	if err != nil {
		return AntiObj{}, err
	}
	defer db.Close()

	// Query the database for the value associated with the key "key1"
	var value []byte
	query := fmt.Sprintf("SELECT value FROM %s WHERE key = ?", m.table)
	err = db.QueryRow(query, "key1").Scan(&value)
	if err != nil {
		return AntiObj{}, err
	}

	var antiObj AntiObj
	err = json.Unmarshal(value, &antiObj)
	if err != nil {
		return AntiObj{}, err
	}

	return antiObj, err
}

func (m MySQL) barrier(ctx context.Context, lineage []WriteIdentifier, datastoreID string) error {
	//
	//TO-DO
	//
	return nil
}
