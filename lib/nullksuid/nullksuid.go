package nullksuid

import (
	"github.com/cuvva/cuvva-public-go/lib/ksuid"
	"github.com/wearemojo/mojo-public-go/lib/ptr"
)

// P takes an empty string and returns a pointer to a KSUID if the string isn't
// null. It exists purely because the mongo driver doesn't support nullable
// ksuids.
func P(in *string) (k *ksuid.ID) {
	if in != nil {
		k = ptr.P(ksuid.MustParse(*in))
	}
	return
}
