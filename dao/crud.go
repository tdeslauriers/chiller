package dao

import (
	"fmt"
	"reflect"
	"strings"
)

func InsertRecord(uri string, table string, insert interface{}) error {

	db := dbConn(uri)
	defer db.Close()

	// build query
	query := fmt.Sprintf("INSERT INTO %s (", table)

	inserts := structToMap(insert)
	keys := make([]string, len(inserts))
	values := make([]interface{}, len(inserts))
	i := 0
	for k, v := range inserts {
		keys[i] = k
		values[i] = v
		i++
	}

	query += strings.Join(keys, ", ")
	query += fmt.Sprintf(") VALUES (%s)", strings.Repeat("?, ", len(keys))[0:len(keys)*3-2])

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(values...)
	if err != nil {
		return err
	}

	return nil
}

func UpdateRecord(uri string, table string, update interface{}) error {

	db := dbConn(uri)
	defer db.Close()

	// build query
	query := fmt.Sprintf("UPDATE %s SET ", table)

	updates := structToMap(update)
	id := updates["Id"]
	delete(updates, "Id") // remove id since it is not being updated.

	keys := make([]string, len(updates))
	values := make([]interface{}, len(keys)+1) // +1 for the id, later
	i := 0
	for k, v := range updates {
		keys[i] = fmt.Sprintf("%s = ?", k)
		values[i] = v
		i++
	}

	query += strings.Join(keys, ", ")
	query += " WHERE id = ?"

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil
	}
	defer stmt.Close()

	// adding Id value to the end for 'WHERE id = ?'
	values[len(keys)] = id // index offset by one so dont need + 1

	_, err = stmt.Exec(values...)
	if err != nil {
		return err
	}

	return nil
}

func structToMap(s interface{}) map[string]interface{} {

	result := make(map[string]interface{})

	v := reflect.ValueOf(s)
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		result[t.Field(i).Name] = v.Field(i).Interface()
	}

	return result
}
