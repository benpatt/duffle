package marshal

import (
	"fmt"

	"github.com/deislabs/cnab-go/bundle"
	"github.com/docker/go/canonical/json"

	"github.com/deislabs/duffle/pkg/crypto/digest"
)

func Bundle(bf *bundle.Bundle) ([]byte, string, error) {
	data, err := json.MarshalCanonical(bf)
	if err != nil {
		return nil, "", err
	}
	data = append(data, '\n') //TODO: why?

	digest, err := digest.OfBuffer(data)
	if err != nil {
		return nil, "", fmt.Errorf("cannot compute digest from bundle: %v", err)
	}

	return data, digest, nil
}
