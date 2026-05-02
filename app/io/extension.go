package io

const (
	YAMLExtension    = ".yml"
	YAMLExtensionAlt = ".yaml"
	JSONExtension    = ".json"
)

func GetYAMLExtensions() []string {
	return []string{YAMLExtension, YAMLExtensionAlt}
}

func GetImageExtensions() []string {
	return []string{
		".jpg",
		".jpeg",
		".png",
	}
}
