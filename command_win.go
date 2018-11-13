// +build windows

package goman

import (
	"os"
	"time"

	"github.com/govenue/mousetrap"
)

var preExecHookFn = preExecHook

func preExecHook(c *Command) {
	if MousetrapHelpText != "" && mousetrap.StartedByExplorer() {
		c.Print(MousetrapHelpText)
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}
}
