package main

import (
	"os"

	"github.com/adrg/xdg"
	"github.com/gurgeous/vectro/internal"
	"gopkg.in/yaml.v3"
)

const statePath = "vectro/state.yml"

type state struct {
	Version int
	Stack   []string
	History []string
}

// load calculator state. bail if we get any kind of error
func Load(c *internal.Calculator) {
	path, err := xdg.ConfigFile(statePath)
	if err != nil {
		return // ignore
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return // ignore
	}

	var state state
	yaml.Unmarshal(data, &state)
	if state.Version != 1 {
		return // ignore
	}
	c.SetState(state.Stack, state.History)
}

// save calculator state. bail if we get any kind of error
func Save(c *internal.Calculator) {
	stack, history := c.GetState()

	data, err := yaml.Marshal(state{Version: 1, Stack: stack, History: history})
	if err != nil {
		panic(err)
	}

	path, err := xdg.ConfigFile(statePath)
	if err != nil {
		return // ignore
	}
	if err = os.WriteFile(path, data, 0600); err != nil {
		return // ignore
	}
}
