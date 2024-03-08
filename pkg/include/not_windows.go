// Package include exists to import generators and outputs so the init
// function is run.

//go:build !windows
package include

import (
	_ "github.com/leehinman/spigot/pkg/output/syslog"
)
