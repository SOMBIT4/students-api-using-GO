package sqlite

import (
	"database/sql"

	"github.com/SOMBIT4/students-api-using-GO/internal/config"
	_ "github.com/mattn/go-sqlite3" // Import the SQLite driver

)


type Sqlite struct {
	Db *sql.DB
}


func New(cfg *config.Config )(*Sqlite, error) {
	db,err:= sql.Open("sqlite3", cfg.StoragePath)
    if err != nil {
		return nil, err
	}
	_, err=db.Exec(`CREATE TABLE IF NOT EXISTS students (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT ,
		email TEXT,
		age INTEGER
	)`)

	if err!= nil {
		return nil, err
	}
	return &Sqlite{Db: db}, nil
}

func (s *Sqlite)CreateStudent(name string, email string , age int) (int64, error){
    
	stmt , err := s.Db.Prepare("INSERT INTO students (name, email, age) VALUES (?, ?, ?)")
	
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result , err:= stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}

	lastid, err:= result.LastInsertId()

	if err != nil {
		return 0, err
	}
	return lastid, nil
}