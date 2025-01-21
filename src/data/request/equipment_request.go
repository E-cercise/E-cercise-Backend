package request

import "regexp"

type EquipmentListRequest struct {
	Q           string `json:"q"`
	MuscleGroup string `json:"muscle_group"`
}

type EquipmentPostRequest struct {
	Band            string            `json:"band"`
	Color           string            `json:"color"`
	Images          []Image           `json:"images"`
	Material        string            `json:"material"`
	Model           string            `json:"model"`
	MuscleGroupUsed []string          `json:"muscle_group_used"`
	Name            string            `json:"name"`
	Option          []Option          `json:"option"`
	SpecialFeature  string            `json:"special_feature"`
	AdditionalField []AdditionalField `json:"additional_field"`
}

type AdditionalField struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
type Image struct {
	ID        string `json:"id"`
	IsPrimary bool   `json:"is_primary"`
}

type Option struct {
	Available int     `json:"available"`
	Price     float64 `json:"price"`
	Weight    float64 `json:"weight"`
}

func ValidateMuscleGroup(muscleGroups []string) bool {
	// Define the regex pattern
	pattern := `^(ft|bk)_[0-9]+$`
	re := regexp.MustCompile(pattern)

	// Validate each muscle group
	for _, group := range muscleGroups {
		if !re.MatchString(group) {
			return false
		}
	}
	return true
}

type EquipmentPutRequest struct {
	AdditionalField *AdditionalFieldPut `json:"additional_field,omitempty"`
	Band            *string             `json:"band,omitempty"`
	Color           *string             `json:"color,omitempty"`
	Images          *Images             `json:"images,omitempty"`
	Material        *string             `json:"material,omitempty"`
	Model           *string             `json:"model,omitempty"`
	MuscleGroupUsed []string            `json:"muscle_group_used,omitempty"`
	Name            *string             `json:"name,omitempty"`
	Option          *OptionPut          `json:"option,omitempty"`
	SpecialFeature  *string             `json:"special_feature,omitempty"`
}

type AdditionalFieldPut struct {
	Created []AdditionalFieldCreated `json:"created,omitempty"`
	Deleted []string                 `json:"deleted,omitempty"`
	Updated []AdditionalFieldUpdated `json:"updated,omitempty"`
}

type AdditionalFieldCreated struct {
	Key   *string `json:"key,omitempty"`
	Value *string `json:"value,omitempty"`
}

type AdditionalFieldUpdated struct {
	ID    *string `json:"id,omitempty"`
	Key   *string `json:"key,omitempty"`
	Value *string `json:"value,omitempty"`
}

type Images struct {
	DeletedID []DeletedID `json:"deleted_id"`
	UploadID  []UploadID  `json:"upload_id"`
}

type DeletedID struct {
	ID        *string `json:"id,omitempty"`
	IsPrimary *bool   `json:"is_primary,omitempty"`
	URL       *string `json:"url,omitempty"`
}

type UploadID struct {
	ID        *string `json:"id,omitempty"`
	IsPrimary *bool   `json:"is_primary,omitempty"`
}

type OptionPut struct {
	Added   []Added         `json:"added"`
	Created []OptionCreated `json:"created"`
	Deleted []string        `json:"deleted"`
	Updated []OptionUpdated `json:"updated"`
}

type Added struct {
	Available *int64   `json:"available,omitempty"`
	Price     *float64 `json:"price,omitempty"`
	Weight    *float64 `json:"weight,omitempty"`
}

type OptionCreated struct {
	Available *int64   `json:"available,omitempty"`
	Price     *float64 `json:"price,omitempty"`
	Weight    *float64 `json:"weight,omitempty"`
}

type OptionUpdated struct {
	Available *int64   `json:"available,omitempty"`
	ID        *string  `json:"id,omitempty"`
	Price     *float64 `json:"price,omitempty"`
	Weight    *float64 `json:"weight,omitempty"`
}
