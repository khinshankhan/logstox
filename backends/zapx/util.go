package zapx

import (
	"github.com/khinshankhan/logstox"
)

func firstNonEmptyString(vals ...string) (string, bool) {
	return logstox.FirstNonZero("", vals...)
}
