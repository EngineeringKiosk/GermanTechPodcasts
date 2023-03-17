package io

const (
	YAMLExtension = ".yml"
	JSONExtension = ".json"
)

func GetImageExtensions() []string {
	return []string{
		".jpg",
		".jpeg",
		".png",
	}
}
