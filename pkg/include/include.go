// Package include exists to import generators and outputs so the init
// function is run.
package include

import (
	_ "github.com/leehinman/spigot/pkg/generator/aws/vpcflow"
	_ "github.com/leehinman/spigot/pkg/generator/cisco/asa"
	_ "github.com/leehinman/spigot/pkg/generator/fortinet/firewall"
	_ "github.com/leehinman/spigot/pkg/output/file"
	_ "github.com/leehinman/spigot/pkg/output/rally"
	_ "github.com/leehinman/spigot/pkg/output/s3"
	_ "github.com/leehinman/spigot/pkg/output/simulate"
	_ "github.com/leehinman/spigot/pkg/output/syslog"
)
