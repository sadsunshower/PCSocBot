package main

import (
	"os"
	"testing"

	"github.com/bwmarrin/discordgo"
	com "github.com/unswpcsoc/PCSocBot/commands"
)

// signal testing
// https://stackoverflow.com/questions/43409919/init-function-breaking-unit-tests

var _ = (func() interface{} {
	_testing = true
	return nil
}())

// Preamble

func TestMain(m *testing.M) {
	// run tests
	os.Exit(m.Run())
}

type Example struct {
	names []string
	desc  string
}

func NewExample() *Example {
	return &Example{
		names: []string{"example", "an extended command string"},
		desc:  "Example!",
	}
}

func (e *Example) Names() []string {
	return e.names
}

func (e *Example) Desc() string {
	return e.desc
}

func (e *Example) Roles() []string {
	return nil
}

func (e *Example) Channels() []string {
	return nil
}

func (e *Example) MsgHandle(ses *discordgo.Session, msg *discordgo.Message) (*com.CommandSend, error) {
	return com.NewSimpleSend(msg.ChannelID, "Pong!"), nil
}

// Tests

func TestAddCommand(t *testing.T) {
	// init router
	router := NewRouter()

	// create command
	exp := NewExample()

	// add to router
	router.AddCommand(exp, exp.Names())

	// assert single route made
	r1 := "example"
	got, ok := router.routes.leaves[r1]
	if !ok {
		t.Errorf("%s: no route made for %s\n", t.Name(), r1)
	}
	if got.command != exp {
		t.Errorf("%s: made route to %v, expected %v\n", t.Name(), got.command, exp)
	}

	// assert lengthy route made
	r2 := []string{"an", "extended", "command", "string"}
	curr := router.routes.leaves
	for _, str := range r2 {
		got, ok = curr[str]
		if !ok {
			t.Errorf("%s: no route made for %v\n", t.Name(), r2)
		}
		curr = got.leaves
	}
	if got.command != exp {
		t.Errorf("%s: made route to %v, expected %v\n", t.Name(), got.command, exp)
	}
}

func TestRoute(t *testing.T) {
	// init router
	router := NewRouter()

	// create command
	exp := NewExample()

	// add simple route manually
	expl := NewLeaf(exp)
	router.routes.leaves["example"] = expl

	// assert simple routing works
	got, ind := router.Route([]string{"example"})
	if ind == 0 {
		t.Errorf("%s: route did not find anything\n", t.Name())
	}
	if got != exp {
		t.Errorf("%s: got %v, expected %v\n", t.Name(), got, exp)
	}

	// add lengthy route manually
	exp2 := NewLeaf(exp)
	router.routes.leaves["an"] = NewLeaf(nil)
	router.routes.leaves["an"].leaves["extended"] = NewLeaf(nil)
	router.routes.leaves["an"].leaves["extended"].leaves["command"] = NewLeaf(nil)
	router.routes.leaves["an"].leaves["extended"].leaves["command"].leaves["string"] = exp2

	// assert lengthy routing works
	got, ind = router.Route([]string{"an", "extended", "command", "string"})
	if ind == 0 {
		t.Errorf("%s: route did not find anything\n", t.Name())
	}
	if got != exp {
		t.Errorf("%s: got %v, expected %v\n", t.Name(), got, exp)
	}
}