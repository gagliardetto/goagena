## goagena

goagena converts structs to goa types.

### Install

```
go get github.com/gagliardetto/goagena
go install github.com/gagliardetto/goagena
```

### How it works

`goagena` will scan the whole package looking for structs, and will convert them by default
to `Type`. If you want to convert a struct to a `MediaType`, add `goagena:mediatype` to the comment
of the struct.

### Examples

```
goagena --pkg="github.com/gagliardetto/goagena/example"
```

Output:

```golang
package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)


var Address = Type("Address", func() {
	Attribute("cap", Integer)
	Attribute("city", String)

	Required("cap", "city")
})
var User = Type("User", func() {
	Attribute("id", UUID)
	Attribute("username", String)
	Attribute("email", String)
	Attribute("password", String)
	Attribute("indreezzo", Address) // Is optional
	Attribute("max_concurrent", Integer)
	Attribute("max_concurrentz", Integer)
	Attribute("created_at", DateTime) // Is optional
	Attribute("things", ArrayOf(String))

	Required("id", "username", "email", "password", "max_concurrent", "max_concurrentz", "things")
})
var YammeBell = Type("YammeBell", func() {
	Attribute("hello", String)

	Required("hello")
})
```


## TODO:

- [x] Support Types
- [x] Support MediaTypes
