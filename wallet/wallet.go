package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"math/big"
	"nomadcoin/utils"
	"os"
)

const Filename string = "root.wallet"

type wallet struct {
	privateKey *ecdsa.PrivateKey
	Address    string
}

var w *wallet

func hasWalletFile() bool {
	_, err := os.Stat(Filename)
	return !os.IsNotExist(err)
}

func createPrivKey() *ecdsa.PrivateKey {
	priveKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleErr(err)
	return priveKey
}

func saveKey(key *ecdsa.PrivateKey) {
	bytes, err := x509.MarshalECPrivateKey(key)
	utils.HandleErr(err)
	//0644 = read and write
	err = os.WriteFile("root.wallet", bytes, 0644)
	utils.HandleErr(err)

}

func restoreKey() *ecdsa.PrivateKey {
	keyAsBytes, err := os.ReadFile(Filename)
	utils.HandleErr(err)
	key, err := x509.ParseECPrivateKey(keyAsBytes)
	utils.HandleErr(err)
	return key
}
func encodeBigIntsByte(a, b []byte) string {
	bigIntBytes := append(a, b...)
	return fmt.Sprintf("%x", bigIntBytes)
}

func getAddressFromPrivKey(key *ecdsa.PrivateKey) string {
	x := key.X.Bytes()
	y := key.Y.Bytes()
	return encodeBigIntsByte(x, y)
}

func Sign(payload string, wallet wallet) string {
	payloadAsByte, err := hex.DecodeString(payload)
	utils.HandleErr(err)
	r, s, err := ecdsa.Sign(rand.Reader, wallet.privateKey, payloadAsByte)
	utils.HandleErr(err)
	signature := encodeBigIntsByte(r.Bytes(), s.Bytes())
	return signature
}

func restoreBigInt(payload string) (*big.Int, *big.Int, error) {
	sigBytes, err := hex.DecodeString(payload)
	if err != nil {
		return nil, nil, err
	}
	utils.HandleErr(err)
	firstHalfBytes := sigBytes[:len(sigBytes)/2]
	secondHalfBytes := sigBytes[len(sigBytes)/2:]

	bigR := &big.Int{}
	bigS := &big.Int{}

	bigR.SetBytes(firstHalfBytes)
	bigS.SetBytes(secondHalfBytes)
	return bigR, bigS, nil
}

func Verify(signature, payload, address string) bool {
	r, s, err := restoreBigInt(signature)
	utils.HandleErr(err)
	x, y, err := restoreBigInt(address)
	utils.HandleErr(err)
	publickKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}
	payloadBytes, err := hex.DecodeString(payload)
	utils.HandleErr(err)
	ok := ecdsa.Verify(&publickKey, payloadBytes, r, s)
	return ok
}

func Wallet() *wallet {
	if w == nil {
		w = &wallet{}
		//has a wallet already?
		if hasWalletFile() {
			w.privateKey = restoreKey()
		} else {
			key := createPrivKey()
			saveKey(key)
			w.privateKey = key
		}
		//yes -> restore from file
		// no -> create prv key,save to file
		w.Address = getAddressFromPrivKey(w.privateKey)
	}
	return w
}
