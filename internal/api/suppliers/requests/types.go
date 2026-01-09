package requests

type CreateSupplierRequest struct {
	CompanyName       string `json:"company_name" binding:"required"`
	ContactPersonName string `json:"contact_person_name" binding:"required"`
	ContactNumber     string `json:"contact_number" binding:"required"`
}

type UpdateSupplierRequest struct {
	CompanyName       string `json:"company_name" binding:"required"`
	ContactPersonName string `json:"contact_person_name" binding:"required"`
	ContactNumber     string `json:"contact_number" binding:"required"`
}
