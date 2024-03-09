package api

type User struct {
	Username     string `json:"username" validate:"required,username"`
	PasswordHash string `json:"password_hash" bson:"password_hash"`
}

type FileMetadata struct {
	Identifier string `json:"identifier" bson:"identifier" validate:"required" example:"YTE1YzhmMjMtYTEwMi00ZmQ0LTk1ZWUtZmM4ZDAyMjc3MmNm"`
	Name       string `json:"name" validate:"required,filename" example:"my_image.png"`
	Username   string `json:"username" validate:"required,username" example:"john_doe"`
	Size       int64  `json:"size" validate:"required,gt=0" example:"12894"`
	Mimetype   string `json:"mimetype" validate:"required" example:"image/png"`
	Creation   int64  `json:"creation" validate:"required" example:"1699651187"`
	Expiration int64  `json:"expiration" validate:"required" example:"1699644399"`
} //@name FileMetadata

type FileData struct {
	Identifier string `json:"identifier" bson:"identifier" validate:"required"`
	Data       []byte `json:"data" validate:"required"`
}
