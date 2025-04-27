package crypto

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPublicKeyEncode(t *testing.T) {
	privKey := GeneratePrivateKey()
	pubKey := privKey.PublicKey

	buf := &bytes.Buffer{}
	assert.Nil(t, SerializePublicKey(&pubKey).Encode(buf))

	decodedPk := &SerializablePublicKey{}
	assert.Nil(t, decodedPk.Decode(buf))

	pubKey2 := DecodePublicKey(decodedPk)
	assert.Equal(t, pubKey, *pubKey2)
}

func TestSignatureEncodeDecode(t *testing.T) {
	privKey := GeneratePrivateKey()
	data := []byte("Hello world")

	sig, err := SignBytes(privKey, data)
	assert.Nil(t, err)
	assert.True(t, sig.Verify(&privKey.PublicKey, data))

	buf := &bytes.Buffer{}
	assert.Nil(t, sig.Encode(buf))

	sig2 := &Signature{}
	assert.Nil(t, sig2.Decode(buf))
	assert.Equal(t, sig, sig2)
}
