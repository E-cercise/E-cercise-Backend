package request

import "regexp"

type EquipmentListRequest struct {
	Q           string `json:"q"`
	MuscleGroup string `json:"muscle_group"`
}

type EquipmentPostRequest struct {
	Band            string   `json:"band"`
	Color           string   `json:"color"`
	Images          []Image  `json:"images"`
	Material        string   `json:"material"`
	Model           string   `json:"model"`
	MuscleGroupUsed []string `json:"muscle_group_used"`
	Name            string   `json:"name"`
	Option          []Option `json:"option"`
	SpecialFeature  string   `json:"special_feature"`
}

type Image struct {
	ID        string `json:"id"`
	IsPrimary bool   `json:"is_primary"`
}

type Option struct {
	Avaliable int64   `json:"avaliable"`
	Price     float64 `json:"price"`
	Weight    float64 `json:"weight"`
}

func ValidateMuscleGroup(muscleGroups []string) bool {
	// Define the regex pattern
	pattern := `^(fk|bk)_[0-9]+$`
	re := regexp.MustCompile(pattern)

	// Validate each muscle group
	for _, group := range muscleGroups {
		if !re.MatchString(group) {
			return false
		}
	}
	return true
}
