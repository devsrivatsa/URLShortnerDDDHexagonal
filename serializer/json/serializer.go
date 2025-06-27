package json

import (
	"encoding/json"

	"github.com/devsrivatsa/URLShortnerDDDHexagonal/urlShortner"
	errs "github.com/pkg/errors"
)

type Redirect struct{}

func (r Redirect) Decode(input []byte) (*urlShortner.Redirect, error) {
	redirect := &urlShortner.Redirect{}
	if err := json.Unmarshal(input, redirect); err != nil {
		return nil, errs.Wrap(err, "serializer.Redirect.Decode")
	}

	return redirect, nil
}

func (r Redirect) Encode(input *urlShortner.Redirect) ([]byte, error) {
	encoded, err := json.Marshal(input)
	if err != nil {
		return nil, errs.Wrap(err, "serializer.Redirect.Encode")
	}
	return encoded, nil
}
