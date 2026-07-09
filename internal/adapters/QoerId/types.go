package adapters


type KycTier3NinRequest struct {
	NIN                string `json:"nin" binding:"required"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	DOB       string `json:"dob" binding:"required"`
}

type KycTier3AddressRequest struct {
	CustomerReference         string `json:"customerReference" binding:"required"`
	Street             string `json:"street" binding:"required"`
	LgaName            string `json:"lgaName" binding:"required"`
	StateName          string `json:"stateName" binding:"required"`
	City               string `json:"city" binding:"required"`
	Landmark           string `json:"landmark" binding:"required"`
	ApplicantFirstName string `json:"applicant_firstname" binding:"required"`
	ApplicantLastName  string `json:"applicant_lastname" binding:"required"`
	ApplicantPhone     string `json:"applicant_phone" binding:"required"`
	ApplicantDOB       string `json:"applicant_dob" binding:"required"`
}

type KycTier3NinResponse struct {
	ID int `json:"id"`
	Applicant struct {
		Firstname string `json:"firstname"`
		Lastname string `json:"lastname"`
		Phone string `json:"phone"`
		Gender string `json:"gender"`
	}
	Status struct{
		State string `json:"state"`
		Status string `json:"status"`
	}

}
type KycTier3AddressResponse struct {
		ID int `json:"id"`
	Applicant struct {
		Firstname string `json:"firstname"`
		Middlename string `json:"middlename"`
		Lastname string `json:"lastname"`
		Phone string `json:"phone"`
		Gender string `json:"gender"`
		DOB string `json:"dob"`
	}

	Status struct{
		State string `json:"state"`
		Status string `json:"status"`
	}
}

