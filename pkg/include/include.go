// Package include exists to import generators and outputs so the init
// function is run.
package include

import (
	_ "github.com/elastic/spigot/pkg/generator/aws/firewall"
	_ "github.com/elastic/spigot/pkg/generator/aws/vpcflow"
	_ "github.com/elastic/spigot/pkg/generator/cef"
	_ "github.com/elastic/spigot/pkg/generator/cisco/asa"
	_ "github.com/elastic/spigot/pkg/generator/citrix/cef"
	_ "github.com/elastic/spigot/pkg/generator/clf"
	_ "github.com/elastic/spigot/pkg/generator/fortinet/firewall"
	_ "github.com/elastic/spigot/pkg/generator/gotext"
	_ "github.com/elastic/spigot/pkg/generator/winlog"
	_ "github.com/elastic/spigot/pkg/output/file"
	_ "github.com/elastic/spigot/pkg/output/rally"
	_ "github.com/elastic/spigot/pkg/output/s3"
	_ "github.com/elastic/spigot/pkg/output/simulate"
)
