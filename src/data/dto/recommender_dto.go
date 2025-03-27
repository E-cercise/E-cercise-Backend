package dto

type RecommenderRequest struct {
	UserType    string                  `json:"user_type"`
	Gender      string                  `json:"gender"`
	Age         int                     `json:"age"`
	Weight      float64                 `json:"weight"`
	Height      float64                 `json:"height"`
	Goal        string                  `json:"goal"`
	Experience  string                  `json:"experience"`
	Preferences []RecommenderPreference `json:"preferences"`
}

type RecommenderPreference struct {
	Tag       string   `json:"tag"`
	Group     string   `json:"group"`
	MaxPrice  *float64 `json:"max_price,omitempty"`
	MinWeight *float64 `json:"min_weight,omitempty"`
}

type RecommendedResponseDTO []struct {
	ID          string  `json:"id"`
	EquipmentID string  `json:"equipment_id"`
	Name        string  `json:"name"`
	Weight      float64 `json:"weight"`
	Price       float64 `json:"price"`
	Score       float64 `json:"score"`
}
