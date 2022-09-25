// Copyright 2022 Gleb Otochkin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package sampleapp

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func getEmployeesPG(db *sql.DB) (EmpData, error) {
	employees := EmpData{}
	t1 := time.Now()
	rows, err := db.Query("SELECT employee_id,first_name,last_name,hire_date,manager_id FROM employees ORDER BY 4 DESC LIMIT 10")
	if err != nil {
		return employees, fmt.Errorf("an employees rows scan error: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var (
			employee_id int64
			first_name  string
			last_name   string
			hire_date   time.Time
			manager_id  int64
		)
		err := rows.Scan(&employee_id, &first_name, &last_name, &hire_date, &manager_id)
		if err != nil {
			return employees, fmt.Errorf("an employees rows scan error: %v", err)
		}
		employees.Employees = append(employees.Employees, Employee{Employee_id: employee_id, First_name: first_name, Last_name: last_name, Hire_date: hire_date, Manager_id: manager_id})
	}
	//
	employees.GetResponseTime = time.Since(t1).String()
	return employees, nil
}
func PostEmployeePG(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	//
	if err := r.ParseForm(); err != nil {
		log.Printf("PostEmployee: failed to parse input: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	first_name := r.FormValue("fname")
	last_name := r.FormValue("lname")
	hire_date := r.FormValue("hdate")
	manager_id, errint := strconv.Atoi(r.FormValue("mgrid"))
	if errint != nil {
		log.Printf("func PostEmployee: unable to convert manager_id to int: %v", errint)
	}
	//Insert data
	insertEmp := "INSERT INTO employees(first_name, last_name, hire_date, manager_id) VALUES( $1, $2, TO_DATE($3,'mm-dd-yyyy'), $4)"
	_, err := db.Exec(insertEmp, first_name, last_name, hire_date, manager_id)
	if err != nil {
		log.Printf("func PostEmployee: unable to save employee: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	fmt.Fprintf(w, "The employee %s is successfully added!", first_name)

}
func checkDBObjectPG(db *sql.DB, dbname string, objname string) (int, error) {
	//
	//defer elapsedTime(time.Now(), "chekObject")
	var cnt int
	err := db.QueryRow("select count(*) from information_schema.tables where table_schema=? and table_name=?", dbname, objname).Scan(&cnt)
	if err != nil {
		//return 0, fmt.Errorf("DB.QueryRow: %v", err)
		return -1, err
	}
	if cnt > 0 {
		return cnt, nil
	}
	return cnt, nil
}

func initDBPG(db *sql.DB) error {
	//Create table if it doesn't exist
	var errddl error
	createEmp := `CREATE TABLE IF NOT EXISTS employees (
		employee_id SERIAL NOT NULL,
		first_name VARCHAR(50) NOT NULL,
		last_name VARCHAR(50) NOT NULL,
		hire_date DATE NOT NULL,
		manager_id BIGINT,
		PRIMARY KEY (employee_id)
	);`
	if os.Getenv("DBVERSION") == "" {
		//
		chk, chkerr := checkDBObjectPG(db, os.Getenv("DBNAME"), "employee")
		if chkerr != nil {
			log.Fatal(chkerr)
		}
		if chk == 0 {
			//
			errddl = execStmt(createEmp)
			if errddl != nil {
				log.Fatalf("Unable to create object: %s", errddl)
			}
		}
	}
	return errddl
}
