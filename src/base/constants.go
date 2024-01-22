package base

type Collection string

const ConfigFile string = "config.yaml"
const PasswordCost int = 12
const FileIdPathParam string = "identifier"
const LimitQueryParam string = "limit"
const SkipQueryParam string = "skip"

const (
	Users         Collection = "users"
	Files         Collection = "files"
	FilesMetadata Collection = "files_metadata"
)
