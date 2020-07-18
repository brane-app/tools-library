package middleware

import (
	"github.com/imonke/monkebase"
	"github.com/imonke/monketype"

	"os"
	"testing"
)

const (
	nick  = "foobar"
	email = "foo@bar.com"
)

var (
	user  monketype.User
	token string
)

func TestMain(main *testing.M) {
	monkebase.Connect(os.Getenv("MONKEBASE_CONNECTION"))
	user = monketype.NewUser(nick, "", email)

	var err error
	if err = monkebase.WriteUser(user.Map()); err != nil {
		panic(err)
	}

	if token, _, err = monkebase.CreateToken(user.ID); err != nil {
		panic(err)
	}

	var result int = main.Run()
	monkebase.DeleteUser(user.ID)
	os.Exit(result)
}
