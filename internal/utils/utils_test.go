package utils_test

import (
	"context"
	"errors"
	"github.com/etiennecoutaud/curlme-controller/internal/utils"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"reflect"
	"testing"
)

var curlmeAnnotationKey = utils.GetCurlmeAnnotationKey()

func TestGetAnnotationValue(t *testing.T) {
	type Test struct {
		Input          map[string]string
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
				"x-k8s.io/curl-me":  "joke=curl-a-joke.herokuapp.com",
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
		"foo=127.0.0.1:8080",
	}

	NonValidTests := []string{
		"joke:curl-a-joke.herokuapp",
		"joke =http://curl-a-joke.herokuapp.com",

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
	type Output struct {
		key   string
		value string
		err   error
	}

	tests := []struct {
		input  string
		output Output
	}{
		{
			input: "foo:bar",
			output: Output{
				key:   "",
				value: "",
				err:   errors.New("annotation value not well format, expect value=curl-url"),
			},
		},
		{
			input: "foo=bar",
			output: Output{
				key:   "foo",
				value: "bar",
				err:   nil,
			},
		},
		{
			input: "foo=bar.com",
			output: Output{
				key:   "foo",
				value: "bar.com",
				err:   nil,
			},
		},
		{
			input: "foo =bar.com",
			output: Output{
				key:   "",
				value: "",
				err:   errors.New("annotation value not well format, expect value=curl-url"),
			},
		},
	}

	for _, test := range tests {
		key, value, err := utils.SplitAnnotationValue(test.input)
		if key != test.output.key || value != test.output.value || !reflect.DeepEqual(test.output.err, err) {
			//(test.output.err == nil && err != nil) && (err.Error() != test.output.err.Error()) {
			t.Errorf("Expected: %s, %s, %v, Got: %s, %s, %v", test.output.key, test.output.value, test.output.err, key, value, err)
		}
	}
}

func TestContainsAnnotation(t *testing.T) {
	tests := []struct {
		clientSet      kubernetes.Interface
		expectedResult bool
	}{
		{
			clientSet: fake.NewSimpleClientset(&v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "default",
					Annotations: map[string]string{
						"x-k8s.io/curl-me-that": "foo",
					},
				},
			}),
			expectedResult: true,
		},
		{
			clientSet: fake.NewSimpleClientset(&v1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test",
					Namespace: "default",
					Annotations: map[string]string{
						"foo": "bar",
					},
				},
			}),
			expectedResult: false,
		},
	}

	for _, test := range tests {
		cm, err := test.clientSet.CoreV1().ConfigMaps("default").Get(context.TODO(), "test", metav1.GetOptions{})
		if err != nil {
			t.Error(err)
		}
		res := utils.ContainsAnnotation(cm)
		if test.expectedResult != res {
			t.Errorf("Expected %v, Got : %v", test.expectedResult, res)
		}
	}
}
