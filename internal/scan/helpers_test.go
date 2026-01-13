// internal/scan/helpers_test.go
package scan

type testErr struct{}

func (e *testErr) Error() string {
	return "fake error"
}

func fakeErr() error {
	return &testErr{}
}
