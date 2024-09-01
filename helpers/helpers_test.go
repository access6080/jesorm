package helpers_test

import (
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
		Name  	string	`type:"VARCHAR(255)" constraints:"PrimaryKey"`
		Id 		int
		Registered bool  `constraints:"ForeignKey(Messages.id)"`
		Account		int	 `type:"INTEGER" constraints:"ForeignKey(Account.id), NotNull, Unique"`
	}
 	oldColumns := structures.CreateColumn(User{})
	newColumns := structures.CreateColumn(User{})

	oldModel := make(map[string][]structures.Column)
	oldModel["User"] = oldColumns

	newModel := make(structures.ModelMap)
	newModel["User"] = newColumns


	equal, err := helpers.CompareModels(oldModel, newModel)
	if err != nil {
        t.Fatalf("CompareModels returned an error: %v", err)
    }

    if !equal {
        t.Errorf("Expected models to be equal, but they were not")
    }

}