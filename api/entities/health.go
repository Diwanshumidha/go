package entities

type HealthPost struct {
	Message string `json:"message" binding:"required"`
}

type HealthParams struct {
	Quantity int `uri:"quantity" binding:"required"`
}
