package dto

type Tier2Request struct {
	BVN       string `json:"bvn" binding:"required"`
	FirstName string `json:"firstname" binding:"required"`
	LastName  string `json:"lastname" binding:"required"`
}

type InitiateTier2Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Url string `json:"url"`
	Ref string `json:"ref"`
}

type BvnRetrievalResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
    FirstName string `json:"firstName"`
    LastName string `json:"lastName"`
    DateOfBirth string `json:"dateOfBirth"`
    Gender string `json:"gender"`
    MobileNumber string `json:"mobileNumber"`
    Email string `json:"email"`
    Address string `json:"address"`
    State string `json:"state"`
    Lga string `json:"lga"`
    City string `json:"city"`
    Landmark string `json:"landmark"`
	Bvn string `json:"bvn"`
	Nin string `json:"nin"`
	

}

type Tier3NinRequest struct {
	NIN                string `json:"nin" binding:"required"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	DOB       string `json:"dob" binding:"required"`
}

type Tier3AddressRequest struct{
	
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

type Tier3NINResponse struct {
	Status string `json:"status"`
	Message string `json:"message"`
}

type Tier3AddressResponse struct {
	Status string `json:"status"`
	Message string `json:"message"`
	FirstName string `json:"firstName"`
	LastName string `json:"lastName"`
	MiddleName string `json:"middleName"`
	DateOfBirth string `json:"dateOfBirth"`
	Gender string `json:"gender"`
	MobileNumber string `json:"mobileNumber"`
	Email string `json:"email"`
	Address string `json:"address"`
	State string `json:"state"`
	Lga string `json:"lga"`
	City string `json:"city"`
	Landmark string `json:"landmark"`
	Bvn string `json:"bvn"`
	Nin string `json:"nin"`
}

type Tier3Request struct {
	NIN       string `json:"nin" binding:"required"`
	Street    string `json:"street" binding:"required"`
	LgaName   string `json:"lgaName" binding:"required"`
	StateName string `json:"stateName" binding:"required"`
	City      string `json:"city" binding:"required"`
	Landmark  string `json:"landmark" binding:"required"`
	FirstName string `json:"firstName" binding:"required"`
	LastName  string `json:"lastName" binding:"required"`
	Phone     string `json:"phone" binding:"required"`
	DOB       string `json:"dob" binding:"required"`
}

type Tier3Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}



