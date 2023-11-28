package assert

import (
	"reflect"
	"testing"
)

func Error(t testing.TB, err error) {
	t.Helper()
	if err == nil {
		t.Errorf("wanted error but got nil")
	}
}

func NoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("got error but did not want one, %+v", err)
	}
}

func Equal(t testing.TB, got, want any) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v, wanted %+v", got, want)
	}
}

func NotEqual(t testing.TB, a, b any) {
	t.Helper()
	if reflect.DeepEqual(a, b) {
		t.Errorf("%+v and %+v are equal, but should not be", a, b)
	}
}

func True(t testing.TB, got bool) {
	t.Helper()
	if !got {
		t.Errorf("got %t, wanted %t", got, true)
	}
}

func False(t testing.TB, got bool) {
	t.Helper()
	if got {
		t.Errorf("got %t, wanted %t", got, false)
	}
}
