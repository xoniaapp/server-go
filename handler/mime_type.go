package handler

var validImageTypes = map[string]bool{
	"image/jpg":  true,
	"image/gif":  true,
	"image/jpeg": true,
	"image/png":  true,
	"image/svg":  true,
}

func isAllowedImageType(mimeType string) bool {
	_, exists := validImageTypes[mimeType]

	return exists
}

var validFileTypes = map[string]bool{
	"image/jpg":                   true,
	"image/vnd.microsoft.icon":    true,
	"image/gif":                   true,
	"image/jpeg":                  true,
	"image/png":                   true,
	"image/svg":                   true,
	"audio/mp3":                   true,
	"audio/mpeg":                  true,
	"audio/opus":                  true,
	"video/mpeg":                  true,
	"video/mp4":                   true,
	"application/json":            true,
	"application/zip":             true,
	"application/gzip":            true,
	"application/ld+json":         true,
	"application/pdf":             true,
	"application/vnd.rar":         true,
	"application/x-tar":           true,
	"application/x-7z-compressed": true,
	"application/x-bzip":          true,
	"application/x-bzip2":         true,
	"application/octet-stream":    true,
	"application/x-binary":        true,
}

func isAllowedFileType(mimeType string) bool {
	_, exists := validFileTypes[mimeType]

	return exists
}
