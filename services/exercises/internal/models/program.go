package models

import (
	"strings"

	"github.com/google/uuid"

	"github.com/Farzan-kh/guddy-cn/exercises/internal/db"
)

type CompleteProgram struct {
	UUID      uuid.UUID `json:"uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	Exercises []ProgramExercise
}

type Program struct {
	UUID      uuid.UUID `json:"uuid" example:"123e4567-e89b-12d3-a456-426614174000"`
	Exercises []ProgramRecord
}

type ProgramRecord struct {
	ExerciseId int `json:"exerciseId" example:"12"`
	Idx        int `json:"idx" example:"1"`
	Sets       int `json:"sets" example:"3"`
	Reps       int `json:"reps" example:"10"`
}

type ProgramExercise struct {
	Exercise Exercise
	Idx      int `json:"idx" example:"1"`
	Sets     int `json:"sets" example:"3"`
	Reps     int `json:"reps" example:"10"`
}

func FullProgramFromRows(uuid uuid.UUID, rows []db.GetFullProgramByIdRow) *CompleteProgram {
	exercises := make([]ProgramExercise, 0, len(rows))
	for _, row := range rows {
		exerciseNames := strings.Split(string(row.NamesGrouped), ",")
		muscles := strings.Split(string(row.MusclesGrouped), ",")
		visuals := strings.Split(string(row.MusclesGrouped), ",")

		exercise := ProgramExercise{
			Idx:  int(row.Idx),
			Sets: int(row.Sets),
			Reps: int(row.Reps),
			Exercise: Exercise{
				Names:     exerciseNames,
				Equipment: row.Equipment.EquipmentT,
				Muscles:   muscles,
				Visuals:   visuals,
			},
		}

		exercises = append(exercises, exercise)
	}

	return &CompleteProgram{
		uuid,
		exercises,
	}
}

func ProgramFromRows(uuid uuid.UUID, rows []db.GetProgramByIdRow) *Program {
	exercises := make([]ProgramRecord, 0, len(rows))
	for _, row := range rows {
		exercise := ProgramRecord{
			Idx:        int(row.Idx),
			Sets:       int(row.Sets),
			Reps:       int(row.Reps),
			ExerciseId: int(row.ExerciseID),
		}

		exercises = append(exercises, exercise)
	}

	return &Program{
		uuid,
		exercises,
	}
}
