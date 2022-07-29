package wallet

import (
	"crypto/x509"
	"encoding/hex"
	"io/fs"
	"nomadcoin/utils"
	"reflect"
	"testing"
)

const (
	testPayload = "0b46d0aa25e3ea7655ac1c3c5b8265e9df715b6284a02596bf915ca9df40edee"
	testKey     = "30770201010420d29c888773d87d3a7f2a333522dbfe4107fdd5f94a3e673e34340dd6524b4541a00a06082a8648ce3d030107a144034200044b3e8a1d593ac1975c653f68f854fd6c7e681b75af5eddbb280df449d3c519a466ff5c1f7ffee33b240e25b1c158d71b596841d1834eebd94bacfd6e2b4ff61b"
	testSig     = "c5791e964cb7b12fd0479f4616233d185868239e7094a967d74f5997fbaee4956a62dbfbf90e1927dac5278c3439363b6544ef8d60bdf2cb6397b26051c13517"
)

func makeTestWallet() *wallet {
	w := &wallet{}
	b, _ := hex.DecodeString(testKey)
	key, _ := x509.ParseECPrivateKey(b)
	w.privateKey = key
	w.Address = getAddressFromPrivKey(key)
	return w
}
func TestSign(t *testing.T) {
	s := Sign("", *makeTestWallet())
	_, err := hex.DecodeString(s)
	if err != nil {
		t.Errorf("Sign() should return a hex encoded string,got %s", s)
	}
}

func TestVerify(t *testing.T) {
	type test struct {
		input string
		ok    bool
	}
	tests := []test{
		{testPayload, true},
		{"0b26d0aa25e3ea7655ac1c3c5b8265e9df715b6284a02596bf915ca9df40edee", false},
	}
	for _, tc := range tests {
		w := makeTestWallet()
		ok := Verify(testSig, tc.input, w.Address)
		if ok != tc.ok {
			t.Error("Verify() could not verify testSignature and testPayload")
		}
	}

}

func TestRestoreBigInts(t *testing.T) {
	_, _, err := restoreBigInt("xx")
	if err == nil {
		t.Error("restoreBigInts should return error when payload is not hex")
	}
}

type fakeLayer struct {
	fakeHasWalletFile func() bool
}

func (f fakeLayer) hasWalletFile() bool {
	return f.fakeHasWalletFile()
}

func (fakeLayer) writeFile(name string, data []byte, perm fs.FileMode) error {
	return nil
}

func (fakeLayer) readFile(filename string) ([]byte, error) {
	return utils.ToBytes(makeTestWallet().privateKey), nil
}

func TestWallet(t *testing.T) {
	t.Run("Wallet is created", func(t *testing.T) {
		files = fakeLayer{
			fakeHasWalletFile: func() bool { return false },
		}
		w := Wallet()
		if reflect.TypeOf(w) != reflect.TypeOf(&wallet{}) {
			t.Error("New wallet should return a new wallet instance")
		}
	})
	t.Run("Wallet is restored", func(t *testing.T) {
		files = fakeLayer{
			fakeHasWalletFile: func() bool { return true },
		}
		w = nil
		if reflect.TypeOf(w) != reflect.TypeOf(&wallet{}) {
			t.Error("New wallet should return a new wallet instance")
		}
	})
}
