package gen

import (
	"go/types"
	"strings"

	"github.com/dimes/dihedral/typeutil"
)

// SanitizeName returns a name that can be used as a Go identifier
// from the given name. The generated identifier should not clash
// with any other types.Named, unless the underlying named object
// is the same
func SanitizeName(name *types.Named) string {
	sanitized := typeutil.IDFromNamed(name)
	sanitized = strings.ReplaceAll(sanitized, "/", "_")
	sanitized = strings.ReplaceAll(sanitized, ".", "_")
	sanitized = strings.ReplaceAll(sanitized, "-", "_")
	return sanitized
}
