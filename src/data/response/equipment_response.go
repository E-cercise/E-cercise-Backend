package response

import "github.com/google/uuid"

type EquipmentsResponse struct {
	Equipments []Equipment
}

type Equipment struct {
	ID        uuid.UUID `json:"ID"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	ImagePath string    `json:"image_path"`
	//Rating float64 `json:"rating"`
}

//type Option struct {
//	RemainingProduct int64   `json:"remaining_product"`
//	Price            float64 `json:"price"`
//	Weight           float64 `json:"weight"`
//}
