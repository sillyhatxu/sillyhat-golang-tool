// +build appengine gopherjs

package sillyhat_logrus

import (
	"io"
)

func checkIfTerminal(w io.Writer) bool {
	return true
}
