package snake_case

import (
	"encoding/json"

	"github.com/joaopandolfi/blackwhale/handlers/conjson"
	"github.com/joaopandolfi/blackwhale/handlers/conjson/transform"
)

// JsonMarshal - Marshal struct to snake_case
func JsonMarshal(v interface{}) ([]byte, error) {
	marshaler := conjson.NewMarshaler(v, transform.ConventionalKeys())
	return json.MarshalIndent(marshaler, "", " ")
}

// JsonUnmarshal - Unmarshal json in snake_case to GolangStruct
func JsonUnmarshal(b []byte, v interface{}) error {
	return json.Unmarshal(
		b,
		conjson.NewUnmarshaler(v, transform.ConventionalKeys()),
	)
}
