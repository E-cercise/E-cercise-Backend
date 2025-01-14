package request

type EquipmentListRequest struct {
	Q           string `json:"q"`
	MuscleGroup string `json:"muscle_group"`
}
