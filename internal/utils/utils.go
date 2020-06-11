package utils

import (
	"errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"regexp"
	"strings"
)

const curlmeAnnotationKey = "x-k8s.io/curl-me-that"
const annotationValueRegex = `^\S+\=(http:\/\/www\.|https:\/\/www\.|http:\/\/|https:\/\/)?[a-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,5}(:[0-9]{1,5})?(\/.*)?$`

// ContainsAnnotation check if ContainsAnnotation is present in annotation
func ContainsAnnotation(cm metav1.Object) bool {
	return GetAnnotationValue(cm.GetAnnotations()) != ""
}

// GetAnnotationValue find if annotation curlmeAnnotationKey is present in annotation map
func GetAnnotationValue(as map[string]string) string {
	for key, value := range as {
		if key == curlmeAnnotationKey {
			return value
		}
	}
	return ""
}

// CompareAnnotation in case of update value to avoid resync when message is updated
func CompareAnnotation(old map[string]string, new map[string]string) bool {
	return GetAnnotationValue(old) == GetAnnotationValue(new)
}

// VerifyValueFormat find if annotation value respect regex
func VerifyValueFormat(v string) error {
	matched, err := regexp.Match(annotationValueRegex, []byte(v))
	if err != nil {
		return err
	}
	if !matched {
		return errors.New("annotation value not well format, expect value=curl-url")
	}
	return nil
}

// SplitAnnotationValue split annotation format foo=bar into string foo and bar
func SplitAnnotationValue(value string) (string, string, error) {
	err := VerifyValueFormat(value)
	if err != nil {
		return "", "", err
	}
	split := strings.Split(value, "=")
	return split[0], split[1], err
}

// GetCurlmeAnnotationKey return const use for testing package
func GetCurlmeAnnotationKey() string {
	return curlmeAnnotationKey
}
