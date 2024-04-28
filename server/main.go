package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

type StudentService struct {
	DB *sql.DB
}

type Student struct {
	ID      int    `json:"id"`
	Code    string `json:"code"`
	Name    string `json:"name"`
	Program string `json:"program"`
}

type RPCRequest struct {
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
	ID     int           `json:"id"`
}

type RPCResponse struct {
	Result interface{} `json:"result"`
	Error  string      `json:"error,omitempty"`
	ID     int         `json:"id"`
}

func (s *StudentService) HandleRequest(req RPCRequest, res *RPCResponse) error {
	switch req.Method {
	case "InsertStudent":
		if len(req.Params) != 3 {
			res.Error = "Invalid number of parameters"
			return nil
		}
		code, ok := req.Params[0].(string)
		if !ok {
			res.Error = "Code parameter must be a string"
			return nil
		}
		name, ok := req.Params[1].(string)
		if !ok {
			res.Error = "Name parameter must be a string"
			return nil
		}
		program, ok := req.Params[2].(string)
		if !ok {
			res.Error = "Program parameter must be a string"
			return nil
		}
		s.insertStudent(code, name, program)
		res.Result = "Student inserted successfully"
	case "DisplayStudents":
		students, err := s.displayStudents()
		if err != nil {
			res.Error = err.Error()
			return nil
		}
		res.Result = students
	default:
		res.Error = "Method not found"
	}
	return nil
}

func (s *StudentService) insertStudent(code, name, program string) {
	log.Println("Inserting student record ...")
	insertStudentSQL := `INSERT INTO student(code, name, program) VALUES (?, ?, ?)`
	statement, err := s.DB.Prepare(insertStudentSQL)
	if err != nil {
		log.Println("Error preparing statement:", err)
		return
	}
	defer statement.Close()

	_, err = statement.Exec(code, name, program)
	if err != nil {
		log.Println("Error inserting student:", err)
		return
	}
	log.Println("Student inserted successfully")
}

func (s *StudentService) displayStudents() ([]Student, error) {
	log.Println("Displaying students ...")
	rows, err := s.DB.Query("SELECT * FROM student ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var students []Student
	for rows.Next() {
		var student Student
		err := rows.Scan(&student.ID, &student.Code, &student.Name, &student.Program)
		if err != nil {
			return nil, err
		}
		students = append(students, student)
	}
	log.Println("Students displayed successfully")
	return students, nil
}

func main() {
	os.Remove("sqlite-database.db") // I delete the file to avoid duplicated records.

	log.Println("Creating sqlite-database.db...")
	file, err := os.Create("sqlite-database.db") // Create SQLite file
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("sqlite-database.db created")

	sqliteDatabase, err := sql.Open("sqlite3", "./sqlite-database.db") // Open the created SQLite File
	if err != nil {
		log.Fatal(err)
	}
	defer sqliteDatabase.Close() // Defer Closing the database
	createTable(sqliteDatabase)  // Create Database Tables

	// INSERT RECORDS
	studentService := &StudentService{DB: sqliteDatabase}
	studentService.insertStudent("0001", "Liana Kim", "Bachelor")
	studentService.insertStudent("0002", "Glen Rangel", "Bachelor")
	studentService.insertStudent("0003", "Martin Martins", "Master")
	studentService.insertStudent("0004", "Alayna Armitage", "PHD")
	studentService.insertStudent("0005", "Marni Benson", "Bachelor")
	studentService.insertStudent("0006", "Derrick Griffiths", "Master")
	studentService.insertStudent("0007", "Leigh Daly", "Bachelor")
	studentService.insertStudent("0008", "Marni Benson", "PHD")
	studentService.insertStudent("0009", "Klay Correa", "Bachelor")

	// DISPLAY INSERTED RECORDS
	displayStudents(sqliteDatabase)

	// Setup RPC
	rpc.Register(studentService)
	rpc.HandleHTTP()

	// Start RPC server
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("Error starting RPC server:", err)
	}
	log.Println("RPC server started on port 1234")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req RPCRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		var res RPCResponse
		err = studentService.HandleRequest(req, &res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(res)
	})

	log.Fatal(http.Serve(listener, nil))
}

func createTable(db *sql.DB) {
	createStudentTableSQL := `CREATE TABLE student (
		"idStudent" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"code" TEXT,
		"name" TEXT,
		"program" TEXT		
	  );` // SQL Statement for Create Table

	log.Println("Create student table...")
	statement, err := db.Prepare(createStudentTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	log.Println("student table created")
}

func displayStudents(db *sql.DB) {
	log.Println("Displaying students ...")
	rows, err := db.Query("SELECT * FROM student ORDER BY name")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var code string
		var name string
		var program string
		err := rows.Scan(&id, &code, &name, &program)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Student: ", code, " ", name, " ", program)
	}
	log.Println("Students displayed successfully")
}
