package utils

import (
	"encoding/binary"
	"github.com/mix-labs/IC-Go/utils/principal"
)

func TokenId2TokenIdentifier(canisterID string, tokenID uint32) (string, error) {
	var tokenIdentifier string
	canisterBlob, err := principal.Decode(canisterID)
	if err != nil {
		return "", err
	}
	idBlob := make([]byte, 4)
	binary.BigEndian.PutUint32(idBlob, tokenID)

	data := []uint8{10, 116, 105, 100}
	data = append(data, canisterBlob...)
	data = append(data, idBlob...)
	tokenIdentifier = principal.New(data).Encode()
	return tokenIdentifier, err
}
