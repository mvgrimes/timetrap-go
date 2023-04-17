package common

import "github.com/olebedev/when/rules"

var All = []rules.Rule{
	SlashMDY(rules.Override),
}
