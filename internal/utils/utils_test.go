package utils_test

import (
	"errors"
	"github.com/etiennecoutaud/curlme-controller/internal/utils"
	"testing"
)

var curlmeAnnotationKey = utils.GetCurlmeAnnotationKey()

func TestGetAnnotationValue(t *testing.T) {
	type Test struct {
		Input map[string]string
		ExpectedResult string
	}

	tests := []Test{
		{
			Input: map[string]string{
				curlmeAnnotationKey: "foo",
			},
			ExpectedResult: "foo",
		},
		{
			Input: map[string]string{
				curlmeAnnotationKey: "joke=curl-a-joke.herokuapp.com",
			},
			ExpectedResult: "joke=curl-a-joke.herokuapp.com",
		},
		{
			Input: map[string]string{
				"x-k8s.io/curl-me": "joke=curl-a-joke.herokuapp.com",
			},
			ExpectedResult: "",
		},
		{
			Input: map[string]string{
				"x-k8s.io/curl-me": "joke=curl-a-joke.herokuapp.com",
				curlmeAnnotationKey: "joke=curl-a-joke.herokuapp.com",
			},
			ExpectedResult: "joke=curl-a-joke.herokuapp.com",
		},
	}

	for _, test := range tests {
		res := utils.GetAnnotationValue(test.Input)
		if test.ExpectedResult != res {
			t.Errorf("Expected: %s, Got: %s, For %s input", test.ExpectedResult, res, test.Input)
		}
	}
}

func TestVerifyValueFormat(t *testing.T) {

	expectedErr := errors.New("annotation value not well format, expect value=curl-url")
	validTests := []string{
		"joke=curl-a-joke.herokuapp.com",
		"joke=http://curl-a-joke.herokuapp.com",
		"joke=https://curl-a-joke.herokuapp.com",
		"joke=www.curl-a-joke.herokuapp.io",
		"joke-test=www.curl-a-joke.herokuapp.io",
		"joke1=www.curl-a-joke.herokuapp.io",
		"joke1_azerty.json=www.curl-a-joke.herokuapp.io",
	}

	NonValidTests := []string{
		"joke=curl-a-joke.herokuapp",
		"joke =http://curl-a-joke.herokuapp.com",
		"joke=https:curl-a-joke.herokuapp.com",
		"joke=www",
		"joke- test=www.curl-a-joke.herokuapp.io",
		"joke:www.curl-a-joke.herokuapp.io",
		"json=www.curl-a-joke.herokuapp.io toto=www.curl-a-joke.herokuapp.io",
	}

	for _, test := range validTests {
		res := utils.VerifyValueFormat(test)
		if res != nil {
			t.Errorf("Expected : %v, Got : %v - Input: %s", nil, res, test)
		}
	}

	for _, test := range NonValidTests {
		res := utils.VerifyValueFormat(test)
		if res == nil || res.Error() != expectedErr.Error() {
			t.Errorf("Expected: %v, Got: %v - Input: %s", expectedErr, res, test)
		}
	}
}

func TestSplitAnnotationValue(t *testing.T) {
}
