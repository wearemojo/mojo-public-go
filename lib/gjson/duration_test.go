package gjson_test

import (
	"encoding/json/v2"
	"testing"
	"time"

	"github.com/matryer/is"
	"github.com/wearemojo/mojo-public-go/lib/gjson"
)

// A bare time.Duration cannot be carried through json/v2 at all, and the
// failure takes the whole enclosing value with it, so these cover the round
// trip rather than just the type in isolation.
func TestDuration(t *testing.T) {
	type payload struct {
		ChargeIn gjson.Duration  `json:"charge_in"`
		Optional *gjson.Duration `json:"optional"`
		Other    string          `json:"other"`
	}

	t.Run("marshals as nanoseconds, matching encoding/json v1", func(t *testing.T) {
		is := is.New(t)

		data, err := json.Marshal(payload{ChargeIn: gjson.Duration(90 * time.Second), Other: "x"})

		is.NoErr(err)
		is.Equal(`{"charge_in":90000000000,"optional":null,"other":"x"}`, string(data))
	})

	t.Run("round trips", func(t *testing.T) {
		is := is.New(t)

		optional := gjson.Duration(time.Minute)
		want := payload{ChargeIn: gjson.Duration(48 * time.Hour), Optional: &optional, Other: "x"}

		data, err := json.Marshal(want)
		is.NoErr(err)

		out, err := gjson.Unmarshal[payload](data)
		is.NoErr(err)
		is.Equal(want.ChargeIn, out.ChargeIn)
		is.Equal(*want.Optional, *out.Optional)
		is.Equal(48*time.Hour, out.ChargeIn.Duration())
	})

	t.Run("zero value is not special-cased", func(t *testing.T) {
		is := is.New(t)

		data, err := json.Marshal(payload{})

		is.NoErr(err)
		is.Equal(`{"charge_in":0,"optional":null,"other":""}`, string(data))
	})
}
