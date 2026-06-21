package dto

type CreateStudentRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Age       int    `json:"age" binding:"required,gte=1,lte=120"`
}

type StudentResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Age       int    `json:"age"`
}

type UpdateStudentRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Age       int    `json:"age" binding:"required,gte=1,lte=120"`
}

type PatchStudentRequest struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Email     *string `json:"email"`
	Age       *int    `json:"age"`
}

type StudentFilter struct {
	Page   int
	Limit  int
	Age    int
	Search string
	Sort   string
}
