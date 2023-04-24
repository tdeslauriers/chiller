package dao

import (
	"fmt"
	"reflect"
	"strings"
)

func InsertRecord(uri string, table string, cols []string, values ...interface{}) error {

	db := dbConn(uri)
	defer db.Close()

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", table, strings.Join(cols, ", "), strings.Repeat("?, ", len(values))[0:len(values)*3-2])
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

	fields := make([]string, len(updates))
	i := 0
	for k, _ := range updates {
		fields[i] = fmt.Sprintf("%s = ?", k)
		i++
	}

	query += strings.Join(fields, ", ")
	query += fmt.Sprintln(" WHERE id = ?")

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil
	}
	defer stmt.Close()

	values := make([]interface{}, len(fields)+1)
	i = 0
	for _, v := range updates {
		values[i] = v
		i++
	}
	values[len(fields)] = id // index offset by one so dont need + 1

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
