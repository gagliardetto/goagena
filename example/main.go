package main

//go:generate checos gen
import (
	"time"

	uuid "github.com/satori/go.uuid"
)

func main() {

}

// User is the object that defines a user
//goagena:mediatype
type User struct {
	ID uuid.UUID

	// Username is the username of user
	Username       string
	Email          string
	Password       string
	Indreezzo      *Address
	MaxConcurrent  int
	MaxConcurrentz int64
	CreatedAt      *time.Time
	Things         []string
}

type YammeBell struct {
	Hello string
}

type Address struct {
	CAP  int
	City string
}

func (a *Address) DoSomethingOnAddress() {

}

func (a *User) DoSomethingOnUser() {

}

func (a *Address) Validate() error {
	return nil
}
