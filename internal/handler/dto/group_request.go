package dto

type CreateGroupRequest struct {
	Name    string `json:"name" binding:"required"`
	Members []int  `json:"members" binding:"required"` // list user_id
	CreatedBy int  `json:"created_by"`
}

type GroupResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
