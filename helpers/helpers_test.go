package helpers_test

import (
	"reflect"
	"testing"

	// "github.com/access6080/jesorm/helpers"
	"github.com/access6080/jesorm/helpers"
	"github.com/access6080/jesorm/structures"
)

func TestTestEqualFunc(t *testing.T) {
	a := []string{"ForeignKey(Account.id)", "NotNull", "Unique"}
	b := []string{"ForeignKey(Account.id)", "NotNull", "Unique"}

	equal := structures.TestEq(a, b)

	if !equal {
		t.Errorf("Expected slices to be equal, but they were not")
	}
}

func TestCompareModels(t *testing.T) {
	type User struct {
		Name       string `type:"VARCHAR(255)" constraints:"PrimaryKey"`
		Id         int
		Registered bool `constraints:"ForeignKey(Messages.id)"`
		Account    int  `type:"INTEGER" constraints:"ForeignKey(Account.id), NotNull, Unique"`
	}
	oldColumns := structures.CreateColumn(User{})
	newColumns := structures.CreateColumn(User{})

	oldModel := make(map[string][]structures.Column)
	oldModel["User"] = oldColumns

	newModel := make(structures.ModelMap)
	newModel["User"] = newColumns

	migrate, _ := helpers.CompareModels(oldModel, newModel)

	if migrate {
		t.Errorf("Expected models to be equal, but they were not")
	}

	type NewUser struct {
		Name       string `type:"VARCHAR(255)" constraints:"PrimaryKey"`
		Id         int    `type:"INT(2)"`
		Registered bool   `constraints:"ForeignKey(Account.id)"`
		Account    int    `type:"INTEGER" constraints:"ForeignKey(Account.id), NotNull, Unique"`
	}

	newColumns = structures.CreateColumn(NewUser{})
	newModel = make(structures.ModelMap)
	newModel["User"] = newColumns

	migrate, _ = helpers.CompareModels(oldModel, newModel)

	if !migrate {
		t.Errorf("Expected models not to be equal, but they were")
	}
}

func TestMigrateColumnsReturned(t *testing.T) {
	type User struct {
		Name       string `type:"VARCHAR(255)" constraints:"PrimaryKey"`
		Id         int
		Registered bool `constraints:"ForeignKey(Account.id)"`
		Account    int  `type:"INTEGER" constraints:"ForeignKey(Account.id), NotNull, Unique"`
	}
	oldColumns := structures.CreateColumn(User{})
	oldModel := make(map[string][]structures.Column)
	oldModel["User"] = oldColumns

	type NewUser struct {
		Name       string `type:"VARCHAR(255)" constraints:"PrimaryKey"`
		Id         int    `type:"INT(2)"`
		Registered bool   `constraints:"ForeignKey(Account.id)"`
		Account    int    `type:"INTEGER" constraints:"ForeignKey(Account.id), NotNull, Unique"`
	}

	newColumns := structures.CreateColumn(NewUser{})
	newModel := make(structures.ModelMap)
	newModel["User"] = newColumns

	_, migrateMap := helpers.CompareModels(oldModel, newModel)

	if migrateMap == nil {
		t.Errorf("Migrate map should not be nil")
	}

	returnedMap := structures.ModelMap{
		"User": []structures.Column{
			{
				Name:    "Id",
				Gotype:  "int",
				Sqltype: "INT(2)",
			},
		},
	}

	if !reflect.DeepEqual(migrateMap, returnedMap) {
		t.Errorf("Migrate map and returned map should be equal, but they were not")
	}
}
