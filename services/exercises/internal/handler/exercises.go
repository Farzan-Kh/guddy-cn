package service

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/Farzan-kh/guddy-cn/exercises/internal/db"
	"github.com/Farzan-kh/guddy-cn/exercises/internal/models"
)

// GetExercises godoc
// @Summary      List exercises
// @Description  Get all exercises
// @Tags         exercises
// @Produce      json
// @Param        id			query      int  	false	"Exercise ID"
// @Param        muscle   	query      string  	false  	"Target Muscle(s)"
// @Param        equipment  query      string  	false	"Equipment required for the Exercise"
// @Param        name   	query      string  	false  	"Name of the Exercise"
// @Param		 limit		query		int		false	"Limit"
// @Param		 offset		query		int		false	"Offset"
// @Success      200	{array}  models.Exercise
// @Failure      500
// @Router       /api/exercises [get]
func GetExercises(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /api/exercises endpoint called")

	query_params := r.URL.Query()
	name := query_params.Get("name")
	equipment := query_params.Get("equipment")
	muscle := query_params.Get("muscle")
	id := query_params.Get("id")
	limitParam := query_params.Get("limit")
	offsetParam := query_params.Get("offset")

	// Parse limit/offset with defaults and validation
	limit := 50
	offset := 0
	if limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}
	if offsetParam != "" {
		if o, err := strconv.Atoi(offsetParam); err == nil && o >= 0 {
			offset = o
		}
	}

	params := db.GetExercisesParams{
		Name:       nil,
		Equipment:  nil,
		Muscle:     nil,
		ExerciseID: nil,
		Offset:     offset,
		Limit:      limit,
	}
	if name != "" {
		params.Name = &name
	}
	if equipment != "" {
		params.Equipment = db.EquipmentT("Other")
	}
	if muscle != "" {
		params.Muscle = &muscle
	}
	if id != "" {
		params.ExerciseID = &id
	}

	log.Printf("Fetching exercises with params: %+v", params)
	exercisesRows, err := db.Queriez.GetExercises(r.Context(), params)
	if err != nil {
		log.Printf("Couldn't Fetch exercises from db: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	exercises := models.ExerciseFromRows(exercisesRows)

	exercises_json, err := json.MarshalIndent(exercises, "", "  ")
	if err != nil {
		log.Printf("Error at Marshaling exercises objects: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(exercises_json)
}
