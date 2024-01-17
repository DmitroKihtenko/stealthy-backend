package api

type HealthcheckResponse struct {
	Status string `json:"status" example:"ok"`
} //@name HealthcheckResponse

type TokenResponse struct {
	Token string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.ey"`
} //@name TokenResponse

type SignUpRequest struct {
	Username string `json:"username" validate:"required,username" example:"john_doe"`
	Password string `json:"password" validate:"required,password" example:"p@ssw0rd"`
} //@name SignUpRequest

type SignInRequest struct {
	SignUpRequest
} //@name SignInRequest

type AddFileResponse struct {
	FileId string `json:"file_id" bson:"file_id" validate:"required" example:"YTE1YzhmMjMtYTEwMi00ZmQ0LTk1ZWUtZmM4ZDAyMjc3MmNm"`
} //@name AddFileResponse

type ErrorResponse struct {
	Summary string `json:"summary" validate:"required" example:"Invalid authorization token"`
	Detail  any    `json:"detail"`
} //@name ErrorResponse

type UserResponse struct {
	Username string `json:"username" validate:"required,username" example:"john_doe"`
} //@name UserResponse

type FileMetadataListResponse struct {
	Records []*FileMetadata `json:"records" validate:"required,records"`
	Total   int64           `json:"total" validate:"gte=0" example:"10"`
} //@name FileMetadataListResponse

type PaginationQueryParameters struct {
	Skip  int64 `validate:"gte=0" query:"skip" example:"3" default:"0"`
	Limit int64 `validate:"gte=1" query:"limit" example:"20" default:"20"`
} //@name PaginationQueryParameters
