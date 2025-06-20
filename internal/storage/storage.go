package storage

import "github.com/SOMBIT4/students-api-using-GO/internal/types"

type Storage interface {
	CreateStudent(name string, email string , age int) (int64, error)
	GetStudentById(id int64) (types.Student,error)
	GetStudents() ([]types.Student, error)
   UpdateStudentById(id int64, student types.Student) error
   DeleteStudentById(id int64) error

}