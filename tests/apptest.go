package tests

import (
	"github.com/revel/revel/testing"
)

//AppTest is the test structure for the controllers of the application
type AppTest struct {
	testing.TestSuite
}

//Before has setup function
func (t *AppTest) Before() {
	println("Set up")
}

//TestThatIndexPageWorks is the test for index page
func (t *AppTest) TestThatIndexPageWorks() {
	t.Get("/")
	t.AssertOk()
	t.AssertContentType("text/html; charset=utf-8")
}

//After is the teardown for test
func (t *AppTest) After() {
	println("Tear down")
}
