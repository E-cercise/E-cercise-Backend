package helper

import "github.com/E-cercise/E-cercise/src/model"

func FindPrimaryImageFromEquipment(equipment model.Equipment) *model.Image {
	for _, image := range equipment.Images {
		if image.IsPrimary {
			return &image
		}
	}
	return nil
}