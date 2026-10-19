package acceptance

import (
	"encoding/json"
	"strconv"
	"strings"
)

// Numeric accepts the host renderer's decimal strings as well as API JSON numbers.
type Numeric float64

func (n *Numeric) UnmarshalJSON(raw []byte) error {
	s := strings.TrimSpace(string(raw))
	if strings.HasPrefix(s, `"`) {
		if e := json.Unmarshal(raw, &s); e != nil {
			return e
		}
	}
	v, e := strconv.ParseFloat(s, 64)
	if e == nil {
		*n = Numeric(v)
	}
	return e
}
