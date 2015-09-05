package core

import (
	. "github.com/eris-ltd/eth-client/Godeps/_workspace/src/github.com/eris-ltd/common/go/log"
)

var logger *Logger

func init() {
	logger = AddLogger("mintx-core")
}
