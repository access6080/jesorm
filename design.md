# JesORM Design notes

### Goal
To design a simple database agnostic orm for a golang that alows users to also write thier own sql queries


### Desired features
- Create Table
- Find data
- Insert data
- Delete Data
- Update Data

---

## Create Table
The goal is to use golang struct tags to validate each member of the table.

```golang
    type User struct {
        id      int     `validate:"PrimaryKey"`
        name    string  `validate:"NotNull"`
        password string
        email   string
        BookId  int     `validate:"ForeignKey(Books)"`
    }
```


## Find Data
The goal is to expose a few functions that take in a struct(Table) and returns data based to conditions passed to it.
```golang
    user := User{...}
    id := 1

    jesorm.findById(&user, id)
 
    consition := Condition{...}
    jesorm.findAll(&user, condition)


    findById(*tableStruct interface{}, condition Condition) interface{} {}

    func findById[T any](*tableStruct T, condition Condition) T {}

    func findById[T any](*tableStruct *T, condition Condition) *T {
    // Create a new instance of the same type as tableStruct
    newInstance := new(T)

    // Populate newInstance with data from the database query (omitted for brevity)

    return newInstance
}

```

### Conditions
The condition struct is what allows users to tell jesorm what they need returned

```golang
     type Condition struct {
        Field    string      // The column name to apply the condition on
        Operator string      // The operator to use (e.g., "=", ">", "<", "LIKE", etc.)
        Value    interface{} // The value to compare the field against
        And      []Condition // Nested conditions for AND logic
        Or       []Condition // Nested conditions for OR logic
    }

    condition := Condition {
        Field:    "name",
        Operator: "=",
        Value:    "John Doe",
        And: []Condition{
            {
                Field:    "age",
                Operator: ">",
                Value:    30,
            },
        },
    }

    // OR 

    type ConditionParams struct {
        Field    string
        Operator string
        Value    interface{}
    }

    func NewCondition() *Condition {
        return &Condition{}
    }

    func (c *Condition) Args(params ConditionParams) *Condition {}

    func (c *Condition) And(params ConditionParams) *Condition {}

    func (c *Condition) Or(params ConditionParams) *Condition {}

    condition := jesorm.NewCondition()
    condition.Args("name", "=", "John Doe")
             .And()
             .Or()

```

## Insert Data
Same as find data

```golang

    user := User{}
    jesorm.InsertOne(user)

    var users []User
    users = ...
    jesorm.InsertMany(users)

```


## Delete Data
Same as find data
```golang

    import (
        "models"
    )
    id = 0
    jesorm.DeleteById(models.User, id)

    condition = Condition{}
    jesorm.Delete(models.User, condition)
```


## Update Data
Same as find data
```golang
    import "models"
    id := 0
    jesorm.UpdateById(models.User, id)

    condition := Condition{}
    jesorm.Update(models.User, condition)

```


## Initialize Orm and DB
``` golang

    db, err := jesorm.Init(config jesorm.Config{})

   // if you want to execute sql commands manually, jes.Init returns a db object, which exposes an Exec method that takes in a query and args. just like databse/sql package in golang


   query := "...."
   db.Exec(query)
```


## Config

```golang
    type Config struct {
        DriverName string
        DSN        string
    }


    
```