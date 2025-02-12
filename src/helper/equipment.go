package helper

import (
	"github.com/E-cercise/E-cercise/src/model"
)

func FindPrimaryImageFromEquipment(equipment model.Equipment) *model.Image {
	for _, image := range equipment.EquipmentOptions[0].Images {
		if image.IsPrimary {
			return &image
		}
	}
	return nil
}

func GetMuscleGroupIDFromEquipment(equipment model.Equipment) []string {
	var muscleGroups []string
	for _, musGroup := range equipment.MuscleGroups {
		muscleGroups = append(muscleGroups, musGroup.ID)
	}
	return muscleGroups
}
