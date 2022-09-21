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
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
)

var (
	httpPort = flag.String("port", ":8080", "Listen port")
	db       *sql.DB
)

func renderTmpl(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	//
	// data := EmpData{}

	// data.Employees = append(data.Employees, Employee{Employee_id: 1, First_name: "Gleb", Last_name: "Otochkin", Hire_date: time.Now(), Manager_id: 1})
	// data.Employees = append(data.Employees, Employee{Employee_id: 2, First_name: "David", Last_name: "Finch", Hire_date: time.Now(), Manager_id: 1})
	// data.Employees = append(data.Employees, Employee{Employee_id: 3, First_name: "Harry", Last_name: "Windsor", Hire_date: time.Now(), Manager_id: 1})
	//t := time.Now().String
	// Get the table data
	data, err := getEmployees(db)
	if err != nil {
		log.Printf("func getEmployees: failed to get data: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	indexTmpl := template.Must(template.New("index").Parse(indexFile))
	err = indexTmpl.Execute(w, data)
	if err != nil {
		log.Printf("func renderTmpl: failed to render template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func dbConnect() *sql.DB {
	//
	var (
		db  *sql.DB
		err error
	)
	if os.Getenv("DB_VERSION") == "" {
		db, err = connectPostgres()
		if err != nil {
			log.Fatalf("connectPostgres: Unable to connect to the database %s", err)
		}
		// if err := initDBPG(db); err != nil {
		// 	log.Fatalf("initDBPG unable to create table: %s", err)
		// }
	}

	return db
}
func checkDBObjectPG(dbname string, objname string) (int, error) {
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
func execStmt(tdll string) error {
	//
	_, err := db.Exec(tdll)
	return err
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
	if os.Getenv("DB_VERSION") == "" {
		//
		chk, chkerr := checkDBObjectPG(os.Getenv("DBNAME"), "employee")
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

type Employee struct {
	Employee_id int64
	First_name  string
	Last_name   string
	Hire_date   time.Time
	Manager_id  int64
}
type EmpData struct {
	Employees []Employee
}

func getEmployees(db *sql.DB) (EmpData, error) {
	employees := EmpData{}
	rows, err := db.Query("SELECT employee_id,first_name,last_name,hire_date,manager_id FROM employees ORDER BY 4 DESC LIMIT 10")
	if err != nil {
		return employees, fmt.Errorf("An employees rows scan error: %v", err)
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
			return employees, fmt.Errorf("An employees rows scan error: %v", err)
		}
		employees.Employees = append(employees.Employees, Employee{Employee_id: employee_id, First_name: first_name, Last_name: last_name, Hire_date: hire_date, Manager_id: manager_id})
	}
	return employees, nil
}
func PostEmployee(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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

func RunApp(w http.ResponseWriter, r *http.Request) {
	//
	switch r.Method {
	case http.MethodGet:
		renderTmpl(w, r, dbConnect())
	case http.MethodPost:
		//
		PostEmployee(w, r, dbConnect())
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// ######################################################
// This section can be moved to the connect_pg.go as one of the option for connection
//
func configurePool(db *sql.DB) {

	// Maximum number of connections in idle connection pool.
	db.SetMaxIdleConns(3)

	// Maximum number of open connections to the database.
	db.SetMaxOpenConns(10)

	// Maximum time (in seconds) that a connection can remain open.
	db.SetConnMaxLifetime(1800 * time.Second)

}

func connectPostgres() (*sql.DB, error) {
	// Connection parameters
	// Here is example how to use environment variable for that, but it is not secure
	// Better to use integrations with secret stores
	var (
		dbPort string
		dbName = os.Getenv("DBNAME")
		dbUser = os.Getenv("DBUSER")
		dbPwd  = os.Getenv("DBPASS")
		dbHost = os.Getenv("DBHOST")
	)
	//Default port for Postgres
	if os.Getenv("DBPORT") == "" {
		dbPort = "5432"
	} else {
		dbPort = os.Getenv("DBPORT")
	}
	//create URI for the database connection
	dbURI := fmt.Sprintf("host=%s user=%s password=%s port=%s database=%s", dbHost, dbUser, dbPwd, dbPort, dbName)

	//Create a connection pool
	dbPool, err := sql.Open("pgx", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql open pool error: %v", err)
	}

	//Configure the pool
	configurePool(dbPool)

	return dbPool, nil

}

// ######################################################
//

func main() {
	flag.Parse()
	//http.HandleFunc("/", renderTmpl)
	http.HandleFunc("/", RunApp)
	http.Handle("/test", http.FileServer(http.Dir("./html")))
	log.Printf("Listening on port %s", *httpPort)
	log.Fatal(http.ListenAndServe(*httpPort, nil))
}

var indexFile = `
<html lang="en">
<head>
  <head>
  	<title>Sample DB app</title>
  	<meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.2.1/dist/css/bootstrap.min.css" rel="stylesheet">
	<script src="https://cdn.jsdelivr.net/npm/bootstrap@5.2.1/dist/css/bootstrap.min.js"></script>
	<script src="https://ajax.googleapis.com/ajax/libs/jquery/3.6.0/jquery.min.js"></script>
</head>
<body>
  <div class="container">
    <div class="page-title">
      <h1>Sample App</h1>
      <p class="lead">Sample DB app with Go</p>
      <hr />
    </div>
	<h3>Last 10 added employees</h3>
	<div class="table-responsive">          
	<table class="table">
	  <thead>
		<tr>
		  <th>Employee_id</th>
		  <th>First_name</th>
		  <th>Last_name</th>
		  <th>Hire_date</th>
		  <th>Manager_id</th>
		</tr>
	  </thead>	  
	  <tbody>
	    {{ range .Employees }}
		<tr>
		  <td>{{ .Employee_id }}</td>
		  <td>{{ .First_name }}</td>
		  <td>{{ .Last_name }}</td>
		  <td>{{ .Hire_date }}</td>
		  <td>{{ .Manager_id }}</td>
		</tr>
		{{ end}}
	  </tbody>
	</table>
	</div>
	<h3>Add employee</h3>
	<form class="form-horizontal" method="post">
		<div class="form-group">
			<label class="control-label col-sm-1" for="fname">First Name:</label>
			<div class="col-sm-10">
			<input type="text" class="form-control" id="fname" placeholder="Enter First Name">
			</div>
		</div>
		<div class="form-group">
			<label class="control-label col-sm-1" for="lname">Last Name:</label>
			<div class="col-sm-10">
			<input type="text" class="form-control" id="lname" placeholder="Enter Last Name">
			</div>
		</div>
		<div class="form-group">
			<label class="control-label col-sm-1" for="hdate">Hire Date:</label>
			<div class="col-sm-10">
			<input type="text" class="form-control" id="hdate" placeholder="mm-dd-yyyy">
			</div>
		</div>
		<div class="form-group">
			<label class="control-label col-sm-1" for="mgrid">Manager ID:</label>
			<div class="col-sm-10">
			<input type="text" class="form-control" id="mgrid" placeholder="Enter Manager ID">
			</div>
		</div>
		<div class="form-group">
			<div class="col-sm-offset-2 col-sm-10">
			<button type="button" class="btn btn-default" id="submitEmployee">Submit</button>
			</div>
		</div>
	</form>
	<script>
    function postemployee(fname,lname,hdate,mgrid) {
        var xhr = new XMLHttpRequest();
        xhr.onreadystatechange = function () {
            if (this.readyState == 4) {
                window.location.reload();
            }
        };
        xhr.open("POST", "/", true);
        xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
        xhr.send("fname=" + fname + "&lname=" + lname + "&hdate=" + hdate + "&mgrid=" + mgrid);
    }
    document.getElementById("submitEmployee").addEventListener("click", function () {
        postemployee(fname.value,lname.value,hdate.value,mgrid.value);
    });
	</script>
</body>
</html>
`
