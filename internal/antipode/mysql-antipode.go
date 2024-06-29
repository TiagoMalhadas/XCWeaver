package antipode

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// It assumes by default that the table where the queries will be executed has
// exactly two columns called k and value
type MySQL struct {
	dsn       string
	datastore string
}

func CreateMySQL(host string, port string, user string, password string, datastore string) MySQL {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, datastore)

	return MySQL{dsn, datastore}
}

func (m MySQL) write(ctx context.Context, table string, key string, obj AntiObj) error {

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
	query := fmt.Sprintf("INSERT INTO %s VALUES (?, ?)", table)
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(key, jsonAntiObj)

	return err
}

func (m MySQL) read(ctx context.Context, table string, key string) (AntiObj, error) {

	// Connect to MySQL database
	db, err := sql.Open("mysql", m.dsn)
	if err != nil {
		return AntiObj{}, err
	}
	defer db.Close()

	// Query the database for the value associated with the key
	var value []byte
	query := fmt.Sprintf("SELECT value FROM %s WHERE k = ?", table)
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

func (m MySQL) consume(context.Context, string, string, chan struct{}) (<-chan AntiObj, error) {
	return nil, nil
}

func (m MySQL) barrier(ctx context.Context, lineage []WriteIdentifier, datastoreID string) error {

	// Connect to MySQL database
	db, err := sql.Open("mysql", m.dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	for _, writeIdentifier := range lineage {
		fmt.Println("key after for: ", writeIdentifier.Key)
		if writeIdentifier.Dtstid == datastoreID {
			for {
				// Query the database for the value associated with the writeIdentifier.Key
				var value []byte
				query := fmt.Sprintf("SELECT value FROM %s WHERE k = ?", writeIdentifier.TableId)
				err = db.QueryRow(query, writeIdentifier.Key).Scan(&value)

				if !errors.Is(err, sql.ErrNoRows) && err != nil {
					return err
				} else if errors.Is(err, sql.ErrNoRows) { //the version replication process is not yet completed
					fmt.Println("replication in progress")
					continue
				} else {
					var antiObj AntiObj
					err = json.Unmarshal(value, &antiObj)
					if err != nil {
						return err
					}

					if antiObj.Version == writeIdentifier.Version { //the version replication process is already completed
						fmt.Println("replication done: ", antiObj.Version)
						break
					} else { //the version replication process is not yet completed
						fmt.Println("replication of the new version in progress")
						continue
					}
				}
			}
		}
	}
	return nil
}
