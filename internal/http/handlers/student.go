package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

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