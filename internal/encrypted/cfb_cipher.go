package encrypted

import "crypto/cipher"

const (
	keyLen  = 32
	keyIter = 4096
)

// CFBCipher is a CFB Cipher.
type CFBCipher struct{}

type Cipher interface {
	// Encrypter returns encrypting Stream.
	Encrypter(block cipher.Block, iv []byte) cipher.Stream

	// Decrypter returns decrypting Stream.
	Decrypter(block cipher.Block, iv []byte) cipher.Stream
}

var _ Cipher = (*CFBCipher)(nil)

// Encrypter implements Cipher.
func (c CFBCipher) Encrypter(block cipher.Block, iv []byte) cipher.Stream {
	return cipher.NewCFBEncrypter(block, iv)
}

// Decrypter implements Cipher.
func (c CFBCipher) Decrypter(block cipher.Block, iv []byte) cipher.Stream {
	return cipher.NewCFBDecrypter(block, iv)
}
