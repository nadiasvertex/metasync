package schema

type Article struct {
	Uuid     string
	RootUuid string

	DocumentId uint32

	MepsLanguage string

	FileSignature string
	MlSignature   string
}
