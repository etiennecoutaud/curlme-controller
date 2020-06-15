package controller_test

import (
	"github.com/etiennecoutaud/curlme-controller/internal/controller"
	"github.com/etiennecoutaud/curlme-controller/internal/fakehttpserver"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/diff"
	"k8s.io/client-go/informers"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	core "k8s.io/client-go/testing"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"net/http"
	"reflect"
	"testing"
)

const (
	noResyncPeriodFunc = 0
)

var (
	alwaysReady = func() bool { return true }
)

type fixture struct {
	t *testing.T

	client          *k8sfake.Clientset
	configMapLister []*corev1.ConfigMap
	actions         []core.Action
	objects         []runtime.Object
}

func newFixture(t *testing.T) *fixture {
	return &fixture{
		t:       t,
		objects: []runtime.Object{},
	}
}

func newConfigMap(name string, annotation map[string]string, data map[string]string) *corev1.ConfigMap {
	return &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   metav1.NamespaceDefault,
			Annotations: annotation,
		},
		Data: data,
	}
}

func (f *fixture) newController() (*controller.Controller, informers.SharedInformerFactory) {
	f.client = k8sfake.NewSimpleClientset(f.objects...)
	i := informers.NewSharedInformerFactory(f.client, noResyncPeriodFunc)

	c := controller.New(f.client, i.Core().V1().ConfigMaps())
	c.SetConfigMapSynced(alwaysReady)
	c.SetRecorder(&record.FakeRecorder{})

	for _, cm := range f.configMapLister {
		i.Core().V1().ConfigMaps().Informer().GetIndexer().Add(cm)
	}

	return c, i
}

func (f *fixture) run(cmName string) {
	f.runController(cmName, true, false)
}

func (f *fixture) runExpectError(cmName string) {
	f.runController(cmName, true, true)
}

func (f *fixture) runController(cmName string, startInformers, expectError bool) {
	c, i := f.newController()
	if startInformers {
		stopCh := make(chan struct{})
		defer close(stopCh)
		i.Start(stopCh)
	}

	err := c.SyncHandler(cmName)
	if !expectError && err != nil {
		f.t.Errorf("error syncing cm: %v", err)
	} else if expectError && err == nil {
		f.t.Errorf("expected error syncing cm, got nil")
	}

	actions := filterInformerAction(f.client.Actions())
	for i, action := range actions {
		if len(f.actions) < i+1 {
			f.t.Errorf("%d unexpected actions: %+v", len(actions)-len(f.actions), actions[i:])
			break
		}
		expectedAction := f.actions[i]
		checkAction(expectedAction, action, f.t)
	}

	if len(f.actions) > len(actions) {
		f.t.Errorf("%d additional expected actions:%+v", len(f.actions)-len(actions), f.actions[len(actions):])
	}

}

func filterInformerAction(actions []core.Action) []core.Action {
	ret := []core.Action{}
	for _, action := range actions {
		if len(action.GetNamespace()) == 0 &&
			action.Matches("watch", "configmaps") ||
			action.Matches("list", "configmaps") {
			continue
		}
		ret = append(ret, action)
	}
	return ret
}

func checkAction(expected, actual core.Action, t *testing.T) {
	if !(expected.Matches(actual.GetVerb(), actual.GetResource().Resource)) {
		t.Errorf("Expected\n\t%#v\ngot\n\t%#v", expected, actual)
		return
	}

	if reflect.TypeOf(actual) != reflect.TypeOf(expected) {
		t.Errorf("Action has wrong type. Expected: %t. Got: %t", expected, actual)
		return
	}

	switch a := actual.(type) {
	case core.UpdateActionImpl:
		e, _ := expected.(core.UpdateActionImpl)
		expObject := e.GetObject()
		object := a.GetObject()
		if !reflect.DeepEqual(expObject, object) {
			t.Errorf("Action %s %s has wrong object\nDiff:\n %s",
				a.GetVerb(), a.GetResource().Resource, diff.ObjectGoPrintSideBySide(expObject, object))
		}
	default:
		t.Errorf("Uncaptured Action %s %s, you should explicitly add a case to capture it",
			actual.GetVerb(), actual.GetResource().Resource)
	}
}

func (f *fixture) expectUpdateCMStatusAction(cm *corev1.ConfigMap) {
	action := core.NewUpdateAction(schema.GroupVersionResource{Resource: "configmaps"}, cm.Namespace, cm)
	f.actions = append(f.actions, action)
}

func getKey(cm *corev1.ConfigMap, t *testing.T) string {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(cm)
	if err != nil {
		t.Errorf("Unexpected error getting key for foo %v: %v", cm.Name, err)
		return ""
	}
	return key
}

func TestUpdateCMEmptyData(t *testing.T) {
	f := newFixture(t)
	key := "joke"
	body := "this is an awesome joke !"
	fakeHTTPServer := fakehttpserver.New(body, http.StatusOK)

	initialCm := newConfigMap("test",
		map[string]string{
			"x-k8s.io/curl-me-that": key + "=" + fakeHTTPServer.GetServerAddr()},
		nil)
	expectedCm := newConfigMap("test",
		map[string]string{
			"x-k8s.io/curl-me-that": key + "=" + fakeHTTPServer.GetServerAddr()},
		map[string]string{
			key: body,
		})
	f.configMapLister = append(f.configMapLister, initialCm)
	f.objects = append(f.objects, initialCm)

	f.expectUpdateCMStatusAction(expectedCm)
	f.run(getKey(initialCm, t))

}

func TestUpdateCMData(t *testing.T) {
	f := newFixture(t)
	key := "joke"
	body := "this is an awesome joke !"
	fakeHTTPServer := fakehttpserver.New(body, http.StatusOK)

	initialCm := newConfigMap("test",
		map[string]string{
			"x-k8s.io/curl-me-that": key + "=" + fakeHTTPServer.GetServerAddr()},
		map[string]string{
			"foo": "bar",
		})
	expectedCm := newConfigMap("test",
		map[string]string{
			"x-k8s.io/curl-me-that": key + "=" + fakeHTTPServer.GetServerAddr()},
		map[string]string{
			"foo": "bar",
			key:   body,
		})
	f.configMapLister = append(f.configMapLister, initialCm)
	f.objects = append(f.objects, initialCm)

	f.expectUpdateCMStatusAction(expectedCm)
	f.run(getKey(initialCm, t))
}
