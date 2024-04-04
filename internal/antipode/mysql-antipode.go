package antipode

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// It assumes by default that the table where the queries will be executed has
// exactly two columns called k and value
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

	// Connect to MySQL database
	db, err := sql.Open("mysql", m.dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	jsonAntiObj, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	// Prepare the statement and execute the query
	query := fmt.Sprintf("INSERT INTO %s VALUES (?, ?)", m.table)
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(key, jsonAntiObj)

	return err
}

func (m MySQL) read(ctx context.Context, key string) (AntiObj, error) {

	// Connect to MySQL database
	db, err := sql.Open("mysql", m.dsn)
	if err != nil {
		return AntiObj{}, err
	}
	defer db.Close()

	// Query the database for the value associated with the key
	var value []byte
	query := fmt.Sprintf("SELECT value FROM %s WHERE k = ?", m.table)
	err = db.QueryRow(query, key).Scan(&value)
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
