package httptest

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/gavv/httpexpect/v2"
	"github.com/gin-gonic/gin"
)

func TestIdKeys(t *testing.T) {
	want := Responses{
		{Key: "id", Value: 0, Type: "ge"},
	}
	t.Run("Test id keys", func(t *testing.T) {
		idKeys := IdKeys()
		if !reflect.DeepEqual(want, idKeys) {
			t.Errorf("IdKeys want %+v but get %+v", want, idKeys)
		}
	})
}

func TestHttpTest(t *testing.T) {
	engine := gin.New()
	// Add /example route via handler function to the gin instance
	handler := GinHandler(engine)
	// Create httpexpect instance
	e := httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(handler),
			Jar:       httpexpect.NewJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})
	pageKeys := Responses{{Key: "message", Value: "OK"}, {Key: "status", Value: 200}, {Key: "data", Value: Response{Key: "message", Value: "pong"}}}
	value := e.GET("/example").Expect().Status(http.StatusOK).JSON()

	Test(value, pageKeys)
}

func TestHttpTestArray(t *testing.T) {
	engine := gin.New()
	// Add /example route via handler function to the gin instance
	handler := GinHandler(engine)
	// Create httpexpect instance
	e := httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(handler),
			Jar:       httpexpect.NewJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})
	pageKeys := []string{"1", "2"}
	value := e.GET("/array").Expect().Status(http.StatusOK).JSON()
	Test(value, pageKeys)
}

func TestHttpScan(t *testing.T) {
	engine := gin.New()
	// Add /example route via handler function to the gin instance
	handler := GinHandler(engine)
	// Create httpexpect instance
	e := httpexpect.WithConfig(httpexpect.Config{
		Client: &http.Client{
			Transport: httpexpect.NewBinder(handler),
			Jar:       httpexpect.NewJar(),
		},
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})
	pageKeys := Responses{{Key: "message", Value: ""}}
	obj := e.GET("/example").Expect().Status(http.StatusOK).JSON().Object()

	Scan(obj, pageKeys)
	x := pageKeys.GetString("data.message")
	if x != "pong" {
		t.Errorf("Scan want get pong but get %s", x)
	}
}
