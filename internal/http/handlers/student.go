package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/SOMBIT4/students-api-using-GO/internal/storage"
	"github.com/SOMBIT4/students-api-using-GO/internal/types"
	"github.com/SOMBIT4/students-api-using-GO/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
	   // This function returns a handler that responds with a welcome message.
   return func(w http.ResponseWriter, r *http.Request) {
     slog.Info("creating new student handler")

	 var student types.Student

	err:=  json.NewDecoder(r.Body).Decode(&student)

	if errors.Is(err, io.EOF){
		response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("request body is empty")))
		return
	}

	if err != nil {
		response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
		return 
	}
	//request validation
    if err:=  validator.New().Struct(student); err!=nil{

		validateErrs := err.(validator.ValidationErrors)
		response.WriteJson(w, http.StatusBadRequest,response.ValidationError(validateErrs))
		return
	}

	lastid, err:= storage.CreateStudent(
		student.Name,
		student.Email,
		student.Age,
	)
  slog.Info("user created successfullly",slog.String("userId",fmt.Sprint(lastid)))
	if err != nil {	
		response.WriteJson(w, http.StatusInternalServerError, err)
		return
	}

	  response.WriteJson(w, http.StatusCreated, map[string]int64{"id":lastid})
   }
}

func GetbyId(storage storage.Storage) http.HandlerFunc {
	// This function returns a handler that responds with a welcome message.
	return func(w http.ResponseWriter, r *http.Request) {
		id:= r.PathValue("id")
		slog.Info("getting student by id handler",slog.String("method", r.Method), slog.String("id",id))


		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		student, err := storage.GetStudentById(intId)
		if err != nil {
			slog.Error("failed to get student by id", slog.String("id", id))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, student)
	}
}

func GetList(storage storage.Storage) http.HandlerFunc {
	// This function returns a handler that responds with a welcome message.
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("getting all students")

		students, err := storage.GetStudents()
		if err != nil {
			slog.Error("failed to get students", slog.String("method", r.Method), slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, students)
	}
}

func UpdateById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("updating student by id handler", slog.String("method", r.Method), slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		var student types.Student
		if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}
		defer r.Body.Close()

		student.ID = intId // enforce path param ID over body ID

		if err := storage.UpdateStudentById(intId, student); err != nil {
			slog.Error("failed to update student by id", slog.String("id", id), slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, map[string]string{"message": "student updated successfully"})
	}
}
