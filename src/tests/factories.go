package tests

import (
	"github.com/Goldziher/fabricator"
	"stealthy-backend/api"
)

var FileDataFactory = fabricator.New[api.FileData](api.FileData{})
