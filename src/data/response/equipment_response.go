package response

import (
	"github.com/E-cercise/E-cercise/src/model"
	"github.com/google/uuid"
)

type EquipmentDetailComparisonResponse struct {
	Equipments []EquipmentDetail `json:"equipments"`
}

type EquipmentDetail struct {
	ID              uuid.UUID         `json:"id"`
	Name            string            `json:"name"`
	Brand           string            `json:"brand"`
	Color           string            `json:"color"`
	Category        string            `json:"category"`
	Description     string            `json:"description"`
	Material        string            `json:"material"`
	Model           string            `json:"model"`
	Options         []Option          `json:"options"`
	AdditionalField []AdditionalField `json:"additional_fields"`
}

type EquipmentsResponse struct {
	Equipments []Equipment `json:"equipments"`
}

type Equipment struct {
	ID               uuid.UUID `json:"ID"`
	Name             string    `json:"name"`
	Price            float64   `json:"price"`
	ImagePath        string    `json:"image_path"`
	MuscleGroupUsed  []string  `json:"muscle_group_used"`
	RemainingProduct *int64    `json:"remaining_product,omitempty"`
}

type CategoriesResponse struct {
	Categories []Category `json:"categories"`
}

type Category struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type EquipmentDetailResponse struct {
	Brand            string            `json:"brand"`
	Color            string            `json:"color"`
	Category         string            `json:"category"`
	Description      string            `json:"description"`
	Material         string            `json:"material"`
	Model            string            `json:"model"`
	MuscleGroupUsed  []string          `json:"muscle_group_used"`
	Name             string            `json:"name"`
	Options          []Option          `json:"options"`
	AdditionalFields []AdditionalField `json:"additional_fields"`
	Features         []Feature         `json:"features"`
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
	Name      string  `json:"name"`
	Available int     `json:"available"`
	Price     float64 `json:"price"`
	Weight    float64 `json:"weight"`
	Images    []Image `json:"images"`
}

type Feature struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

func FormatEquipmentDetailResponse(equipment *model.Equipment) *EquipmentDetailResponse {

	var muscleGroupUsed []string
	for _, msg := range equipment.MuscleGroups {
		muscleGroupUsed = append(muscleGroupUsed, msg.ID)
	}

	var opts []Option
	for _, opt := range equipment.EquipmentOptions {

		var imgs []Image
		for _, img := range opt.Images {
			newImg := Image{
				ID:        img.ID.String(),
				Url:       img.CloudinaryPath,
				IsPrimary: img.IsPrimary,
			}
			imgs = append(imgs, newImg)
		}

		newOpt := Option{
			ID:        opt.ID.String(),
			Name:      opt.Name,
			Available: opt.RemainingProducts,
			Price:     opt.Price,
			Weight:    opt.Weight,
			Images:    imgs,
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
		Brand:            equipment.Brand,
		Color:            equipment.Color,
		Category:         equipment.Category,
		Description:      equipment.Description,
		Material:         equipment.Material,
		Model:            equipment.Model,
		MuscleGroupUsed:  muscleGroupUsed,
		Name:             equipment.Name,
		Options:          opts,
		AdditionalFields: attributes,
		Features:         feats,
	}

	return &resp
}
