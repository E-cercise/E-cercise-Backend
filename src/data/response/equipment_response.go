package response

import (
	"github.com/E-cercise/E-cercise/src/model"
	"github.com/google/uuid"
)

type EquipmentsResponse struct {
	Equipments []Equipment
}

type Equipment struct {
	ID        uuid.UUID `json:"ID"`
	Name      string    `json:"Description"`
	Price     float64   `json:"price"`
	ImagePath string    `json:"image_path"`
	//Rating float64 `json:"rating"`
}

//type Option struct {
//	RemainingProduct int64   `json:"remaining_product"`
//	Price            float64 `json:"price"`
//	Weight           float64 `json:"weight"`
//}

type EquipmentDetailResponse struct {
	Band            string            `json:"band"`
	Color           string            `json:"color"`
	Images          []Image           `json:"images"`
	Material        string            `json:"material"`
	Model           string            `json:"model"`
	MuscleGroupUsed []string          `json:"muscle_group_used"`
	Name            string            `json:"Description"`
	Option          []Option          `json:"option"`
	SpecialFeature  string            `json:"special_feature"`
	AdditionalField []AdditionalField `json:"additional_field"`
	Feature         []Feature         `json:"feature"`
}

type AdditionalField struct {
	ID    string `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}
type Image struct {
	ID        string `json:"id"`
	Url       string `json:"url"`
	IsPrimary bool   `json:"is_primary"`
}

type Option struct {
	ID        string  `json:"id"`
	Available int     `json:"available"`
	Price     float64 `json:"price"`
	Weight    float64 `json:"weight"`
}

type Feature struct {
	ID          string `json:"id"`
	Description string `json:"Description"`
}

func FormatEquipmentDetailResponse(equipment *model.Equipment) *EquipmentDetailResponse {
	var imgs []Image
	for _, img := range equipment.Images {
		newImg := Image{
			ID:        img.ID.String(),
			Url:       img.CloudinaryPath,
			IsPrimary: img.IsPrimary,
		}
		imgs = append(imgs, newImg)
	}

	var muscleGroupUsed []string
	for _, msg := range equipment.MuscleGroups {
		muscleGroupUsed = append(muscleGroupUsed, msg.ID)
	}

	var opts []Option
	for _, opt := range equipment.EquipmentOptions {
		newOpt := Option{
			ID:        opt.ID.String(),
			Available: opt.RemainingProducts,
			Price:     opt.Price,
			Weight:    opt.Weight,
		}
		opts = append(opts, newOpt)
	}

	var feats []Feature
	for _, feat := range equipment.EquipmentFeature {
		newFeat := Feature{
			ID:          feat.ID.String(),
			Description: feat.Description,
		}
		feats = append(feats, newFeat)
	}

	var attributes []AdditionalField
	for _, field := range equipment.Attribute {
		newField := AdditionalField{
			ID:    field.ID.String(),
			Key:   field.Key,
			Value: field.Value,
		}
		attributes = append(attributes, newField)
	}

	resp := EquipmentDetailResponse{
		Band:            equipment.Brand,
		Color:           equipment.Color,
		Images:          imgs,
		Material:        equipment.Material,
		Model:           equipment.Model,
		MuscleGroupUsed: muscleGroupUsed,
		Name:            equipment.Name,
		Option:          opts,
		SpecialFeature:  equipment.SpecialFeature,
		AdditionalField: attributes,
		Feature:         feats,
	}

	return &resp
}
