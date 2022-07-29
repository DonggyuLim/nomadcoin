package wallet

import (
	"io/fs"
	"os"
)

const Filename string = "root.wallet"

//File interface
type fileLayer interface {
	hasWalletFile() bool
	writeFile(name string, data []byte, perm fs.FileMode) error
	readFile(name string) ([]byte, error)
}
type layer struct{}

func (layer) hasWalletFile() bool {
	_, err := os.Stat(Filename)
	return !os.IsNotExist(err)
}

func (layer) writeFile(name string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(name, data, perm)
}

func (layer) readFile(filename string) ([]byte, error) {
	keyAsBytes, err := os.ReadFile(filename)
	return keyAsBytes, err
}

var files fileLayer = layer{}
