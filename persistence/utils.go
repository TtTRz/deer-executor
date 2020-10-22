package persistence

import (
	"bufio"
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"io/ioutil"
)

type DigitalSignPEM struct {
	PublicKey 			*rsa.PublicKey
	PrivateKey 			*rsa.PrivateKey
	PublicKeyRaw 		[]byte
	PrivateKeyRaw 		[]byte
}

/* SHA256 */

func SHA256String(body string) ([]byte, error) {
	return SHA256Bytes([]byte(body))
}

func SHA256Bytes(body []byte) ([]byte, error) {
	buf := bytes.NewBuffer(body)
	return SHA256Streams([]io.Reader{buf})
}

func SHA256Streams(streams []io.Reader) ([]byte, error) {
	hash := sha256.New()
	for _, stream := range streams {
		if _, err := io.Copy(hash, stream); err != nil {
			return nil, err
		}
	}
	ret := hash.Sum(nil)
	//hex.EncodeToString(ret)
	return ret, nil
}

/* RSA Sign */

func RSA2048SignString(body string, privateKey *rsa.PrivateKey) ([]byte, error) {
	return RSA2048SignBytes([]byte(body), privateKey)
}

func RSA2048SignBytes(body []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	hash, err := SHA256Bytes(body)
	if err != nil { return nil, err }
	return RSA2048Sign(hash, privateKey)
}

func RSA2048Sign(hash []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	sign, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash)
	if err != nil {
		return nil, err
	}
	return sign, nil
}

/* RSA Verify */

func RSA2048VerifyString(body string, sign []byte, publicKey *rsa.PublicKey) error {
	return RSA2048VerifyBytes([]byte(body), sign, publicKey)
}

func RSA2048VerifyBytes(body []byte, sign []byte, publicKey *rsa.PublicKey) error {
	hash, err := SHA256Bytes(body)
	if err != nil { return err }
	return RSA2048Verify(hash, sign, publicKey)
}

func RSA2048Verify(hash []byte, sign []byte, publicKey *rsa.PublicKey) error {
	return rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash, sign)
}

func ReadAndParsePublicKey(cert []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(cert)
	if block == nil {
		return nil, errors.New("read public key error")
	}
	//publicKey, err := x509.ParseCertificate(block.Bytes)
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return publicKey.(*rsa.PublicKey), nil
	//return publicKey.PublicKey.(*rsa.PublicKey), nil
}

func ReadAndParsePrivateKey(key []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errors.New("read private key error")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func ReadPemFile(filename string) ([]byte, error) {
	pemBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return pemBytes, nil
}

func GetDigitalPEMFromFile (publicKeyFile string, privateKeyFile string) (*DigitalSignPEM, error) {
	publicKey, err := ReadPemFile(publicKeyFile)
	if err != nil {
		return nil, err
	}
	privateKey, err := ReadPemFile(privateKeyFile)
	if err != nil {
		return nil, err
	}
	dp := GetDigitalPEM(publicKey, privateKey)
	return dp, nil
}

func GetDigitalPEM (publicKey []byte, privateKey []byte) *DigitalSignPEM {
	dp := DigitalSignPEM{}
	if publicKey != nil {
		dp.PublicKeyRaw = publicKey
		pk, err := ReadAndParsePublicKey(publicKey)
		dp.PublicKey = pk
		if err != nil {
			dp.PublicKey = nil
		}
	}
	if privateKey != nil {
		dp.PrivateKeyRaw = privateKey
		pk, err := ReadAndParsePrivateKey(privateKey)
		dp.PrivateKey = pk
		if err != nil {
			dp.PrivateKey = nil
		}
	}
	return &dp
}

func Gets(reader io.Reader) string {
	buf := make([]byte, 16, 16)
	rd := bufio.NewReader(reader)
	for {
		t, err := rd.ReadByte()
		if err != nil {
			return string(buf)
		}
		buf = append(buf, t)
		if t == 13 || t == 0 {
			break
		}
	}
	return string(buf)
}