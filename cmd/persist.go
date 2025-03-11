package main

import (
	"os"

	"github.com/adrg/xdg"
	"github.com/gurgeous/vectro/internal"
	"gopkg.in/yaml.v3"
)

const statePath = "vectro/state.yml"

type state struct {
	Version int      `yaml:"version"`
	Stack   []string `yaml:"stack"`
	History []string `yaml:"history"`
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
	if err = yaml.Unmarshal(data, &state); err != nil {
		return // ignore
	}
	if state.Version != 1 {
		return // ignore
	}

	c.SetStackString(state.Stack)
	c.SetHistory(state.History)
}

// save calculator state. bail if we get any kind of error
func Save(c *internal.Calculator) {
	state := state{Version: 1, Stack: c.GetStackString(), History: c.GetHistory()}
	data, err := yaml.Marshal(state)
	if err != nil {
		panic(err)
	}

	path, err := xdg.ConfigFile(statePath)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(path, data, 0600)
	if err != nil {
		panic(err)
	}
}
