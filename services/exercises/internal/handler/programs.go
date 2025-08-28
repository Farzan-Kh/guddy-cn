package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/Farzan-kh/guddy-cn/exercises/internal/db"
	"github.com/Farzan-kh/guddy-cn/exercises/internal/models"
)

// GetProgram godoc
// @Summary      Get Program by ID
// @Description  Get Program By ID
// @Tags         programs
// @Produce      json
// @Param        uuid		query      string  	true	"Programs UUID"
// @Success      200	{object}  models.Program
// @Failure      500
// @Router       /api/program [get]
func GetProgram(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /api/program endpoint called")
	var program_uuid pgtype.UUID
	err := program_uuid.Scan(chi.URLParam(r, "uuid"))
	if err != nil {
		log.Printf("Error scanning UUID Value from URL: %v, uuid: %s", err, chi.URLParam(r, "uuid"))
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	log.Printf("Fetching program by ID: %s", chi.URLParam(r, "uuid"))
	programRows, err := db.Queriez.GetProgramById(r.Context(), program_uuid)
	if err != nil {
		log.Printf("Error at GETting the program from DB: %v, uuid: %s", err, chi.URLParam(r, "uuid"))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	program := models.ProgramFromRows(program_uuid.Bytes, programRows)

	program_json, err := json.Marshal(program)
	if err != nil {
		log.Printf("Error at Marshaling program object: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(program_json)
}

// GetCompleteProgram godoc
// @Summary      Get Complete Program by ID
// @Description  Get Complete Program By ID, get's all the program and related info about exercises
// @Tags         programs
// @Produce      json
// @Param        uuid		query      string  	true	"Programs UUID"
// @Success      200	{object}  models.CompleteProgram
// @Failure      500
// @Router       /api/completeProgram [get]
func GetCompleteProgram(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /api/fullProgram endpoint called")
	var program_uuid pgtype.UUID
	err := program_uuid.Scan(chi.URLParam(r, "uuid"))
	if err != nil {
		log.Printf("Error scanning UUID Value from URL: %v, uuid: %s", err, chi.URLParam(r, "uuid"))
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	log.Printf("Fetching full program by ID: %s", chi.URLParam(r, "uuid"))
	programRows, err := db.Queriez.GetFullProgramById(r.Context(), program_uuid)
	if err != nil {
		log.Printf("Error at GETting the program from DB: %v, uuid: %s", err, chi.URLParam(r, "uuid"))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	program := models.FullProgramFromRows(program_uuid.Bytes, programRows)

	program_json, err := json.Marshal(program)
	if err != nil {
		log.Printf("Error at Marshaling program object: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(program_json)
}

// GetProgram godoc
// @Summary      Create a Program
// @Description  Create a new program and return it's UUID
// @Tags         programs
// @Accept       json
// @Produce      json
// @Param        uuid		body      []models.ProgramRecord  	true	"Programs UUID"
// @Success      200	{object}  uuid.UUID
// @Failure      500
// @Router       /api/program [post]
func PostProgram(w http.ResponseWriter, r *http.Request) {
	log.Println("POST /api/program endpoint called")
	var exercises_list []db.InsertToProgramsByIdParams

	//TODO: Do something about idx order_number translation
	err := json.NewDecoder(r.Body).Decode(&exercises_list)
	if err != nil {
		log.Printf("Invalid JSON in PostProgram: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	uuid := uuid.New()
	var pg_uuid pgtype.UUID
	_ = pg_uuid.Scan(uuid.String())

	log.Printf("Inserting new program with ID: %s, num_exercises: %d", uuid.String(), len(exercises_list))
	for _, v := range exercises_list {
		v.ID = pg_uuid
		err := db.Queriez.InsertToProgramsById(r.Context(), v)
		if err != nil {
			log.Printf("Error at inserting program items: %v, program_id: %s", err, uuid.String())
			http.Error(w, "Error at inserting program items", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Add("Content-Type", "application/json")
	res := fmt.Sprintf("{\n\t\"program_id\": \"%s\"\n}", uuid)
	w.Write([]byte(res))
}
