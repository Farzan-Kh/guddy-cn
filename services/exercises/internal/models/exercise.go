package models

import (
	"strings"

	"github.com/Farzan-kh/guddy-cn/exercises/internal/db"
)

type Exercise struct {
	Id        int32         `json:"id" example:"12"`
	Names     []string      `json:"names" example:"Push Up"`
	Muscles   []string      `json:"muscles" example:"Chest, Triceps, Shoulders"`
	Equipment db.EquipmentT `json:"equipment" example:"Bodyweight"`
	Visuals   []string      `json:"visuals" example:"pushup.jpg,pushup2.jpg"`
}

func ExerciseFromRows(rows []db.GetExercisesRow) *[]Exercise {
	exercises := make([]Exercise, 0, len(rows))
	for _, v := range rows {
		names := strings.Split(string(v.NamesGrouped), ",")
		muscles := strings.Split(string(v.MusclesGrouped), ",")
		visuals := strings.Split(string(v.VisualsGrouped), ",")
		var equipment db.EquipmentT
		if v.Equipment.Valid {
			equipment = v.Equipment.EquipmentT
		} else {
			equipment = ""
		}

		e := Exercise{
			Id:        v.ID,
			Names:     names,
			Equipment: equipment,
			Muscles:   muscles,
			Visuals:   visuals,
		}

		exercises = append(exercises, e)
	}

	return &exercises
}
