package swagger

import (
	"encoding/json"
	"testing"

	"github.com/emicklei/go-restful"
)

type sample struct {
	id       string `swagger:"required"` // TODO
	items    []item
	rootItem item `json:"root"`
}

type item struct {
	itemName string `json:"name"`
}

// go test -v -test.run TestApi ...swagger
func TestApi(t *testing.T) {
	value := Api{Path: "/", Description: "Some Path", Operations: []Operation{}}
	compareJson(t, true, value, `{"path":"/","description":"Some Path"}`)
}

// go test -v -test.run TestModelToJsonSchema ...swagger
func TestModelToJsonSchema(t *testing.T) {
	sws := newSwaggerService(Config{})
	models := map[string]Model{}
	op := new(Operation)
	op.Nickname = "getSome"
	sws.addModelFromSampleTo(op, true, sample{items: []item{}}, models)
	_, err := json.MarshalIndent(op, " ", " ")
	if err != nil {
		t.Fatal(err.Error())
	}
}

// go test -v -test.run TestServiceToApi ...swagger
func TestServiceToApi(t *testing.T) {
	ws := new(restful.WebService)
	ws.Path("/tests")
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_XML)
	ws.Route(ws.GET("/all").To(dummy).Writes(sample{}))
	cfg := Config{
		WebServicesUrl: "http://here.com",
		ApiPath:        "/apipath",
		WebServices:    []*restful.WebService{ws}}
	sws := newSwaggerService(cfg)
	decl := sws.composeDeclaration(ws, "/tests")
	_, err := json.MarshalIndent(decl, " ", " ")
	if err != nil {
		t.Fatal(err.Error())
	}
}

func dummy(i *restful.Request, o *restful.Response) {}

// go test -v -test.run TestIssue78 ...swagger
type Response struct {
	Code  int
	Users *[]User
	Items *[]Item
}
type User struct {
	Id, Name string
}
type Item struct {
	Id, Name string
}

func TestIssue78(t *testing.T) {
	sws := newSwaggerService(Config{})
	models := map[string]Model{}
	sws.addModelFromSampleTo(&Operation{}, true, Response{Items: &[]Item{}}, models)
	model, ok := models["swagger.Response"]
	if !ok {
		t.Fatal("missing response model")
	}
	if "swagger.Response" != model.Id {
		t.Fatal("wrong model id:" + model.Id)
	}
	code, ok := model.Properties["Code"]
	if !ok {
		t.Fatal("missing code")
	}
	if "integer" != code.Type {
		t.Fatal("wrong code type:" + code.Type)
	}
	items, ok := model.Properties["Items"]
	if !ok {
		t.Fatal("missing items")
	}
	if "array" != items.Type {
		t.Fatal("wrong items type:" + items.Type)
	}
	items_items := items.Items
	if items_items == nil {
		t.Fatal("missing items->items")
	}
	ref := items_items["$ref"]
	if ref == "" {
		t.Fatal("missing $ref")
	}
	if ref != "swagger.Item" {
		t.Fatal("wrong $ref:" + ref)
	}
}

// go test -v -test.run TestIssue85 ...swagger
type Dataset struct {
	Names []string
}

func TestIssue85(t *testing.T) {
	sws := newSwaggerService(Config{})
	models := map[string]Model{}
	anon := struct{ Datasets []Dataset }{}
	sws.addModelFromSampleTo(&Operation{}, true, anon, models)
	_, ok := models["struct { Datasets ||swagger.Dataset }"]
	if !ok {
		for k, _ := range models {
			t.Logf("key:%s", k)
		}
		t.Fatal("missing anonymous model")
	}
}
