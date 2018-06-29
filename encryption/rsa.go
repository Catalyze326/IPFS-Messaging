package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func rsaEvery() {
	kp, err := CreateKeyPair()
	if err != nil {
		fmt.Printf("%v", err)
	}
	savePublicPEMKey("keys/rsa_key.pub", kp.PublicKey)
	savePEMKey("keys/rsa_key.pri", kp)
	if err != nil {
		fmt.Printf("%v", err)
	}
	//encrypts or decrypts a file with rsa
	//the first thing passed in is the file that is to be encrypted/decrypted
	//the second is the output file
	//and the third is true if you are decrypting a file and false if you are encrypting
	encDec("keys/aes.key", "keys/aes.key.enc", false)
	deleteFile("keys/aes.key")
	encDec("keys/aes.key.enc", "keys/aes.key", true)
}

//creates the rsa keys it only returns one, but the public key is generated from the private key
func CreateKeyPair() (*rsa.PrivateKey, error) {
	//4096 bit key. You can change this, but this is the most secure for rsa
	size := 4096

	priv, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		log.Fatalf("Failed to generate %d-bit key", size)
		return nil, err
	}

	return priv, err

}

// func Encrypt(in []byte, pub rsa.PublicKey) ([]byte, error) {
// 	sha1 := sha1.New()
// 	out, err := rsa.EncryptOAEP(sha1, rand.Reader, &pub, in, nil)
// 	if err != nil {
// 		log.Fatalf("Failed to encrypt message %v", err)
// 		return nil, err
// 	}
// 	return out, nil
// }

// func Decrypt(ciphertext []byte, priv *rsa.PrivateKey) ([]byte, error) {
//
// 	sha1 := sha1.New()
//
// 	out, err := rsa.DecryptOAEP(sha1, rand.Reader, priv, ciphertext, nil)
//
// 	if err != nil {
// 		log.Fatalf("Failed to decrypt message %v", err)
// 		return nil, err
// 	}
//
// 	return out, nil
// }

//write the public key to a file as a pem file
func savePublicPEMKey(fileName string, pubkey rsa.PublicKey) {
	asn1Bytes, err := asn1.Marshal(pubkey)
	checkError(err)

	var pemkey = &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	pemfile, err := os.Create(fileName)
	checkError(err)
	defer pemfile.Close()

	err = pem.Encode(pemfile, pemkey)
	checkError(err)
}

//save the private key to file as a pem file
func savePEMKey(fileName string, key *rsa.PrivateKey) {
	outFile, err := os.Create(fileName)
	checkError(err)
	defer outFile.Close()

	var privateKey = &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	err = pem.Encode(outFile, privateKey)
	checkError(err)
}

//check error function that the savePEMKey and savePublicPEMKey need
func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

//Encrypt and decrypt the file with RSA
func encDec(inFile string, outFile string, decrypt bool) {
	//the rsa file path
	keyFile := "keys/rsa_key.pri"
	label := ""

	flag.Parse()

	// Read the input file
	in, err := ioutil.ReadFile(inFile)
	if err != nil {
		log.Fatalf("input file: %s", err)
	}

	// Read the private key
	pemData, err := ioutil.ReadFile(keyFile)
	if err != nil {
		log.Fatalf("read key file: %s", err)
	}

	// Extract the PEM-encoded data block
	block, _ := pem.Decode(pemData)
	if block == nil {
		log.Fatalf("bad key data: %s", "not PEM-encoded")
	}
	if got, want := block.Type, "RSA PRIVATE KEY"; got != want {
		log.Fatalf("unknown key type %q, want %q", got, want)
	}

	// Decode the RSA private key
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatalf("bad private key: %s", err)
	}

	var out []byte
	if decrypt {
		if label == "" {
			label = outFile
		}
		// Decrypt the data
		out, err = rsa.DecryptOAEP(sha1.New(), rand.Reader, priv, in, []byte(label))
		if err != nil {
			log.Fatalf("decrypt: %s", err)
		}
	} else {
		if label == "" {
			label = inFile
		}
		out, err = rsa.EncryptOAEP(sha1.New(), rand.Reader, &priv.PublicKey, in, []byte(label))
		if err != nil {
			log.Fatalf("encrypt: %s", err)
		}
	}

	// Write data to output file
	if err := ioutil.WriteFile(outFile, out, 0600); err != nil {
		log.Fatalf("write output: %s", err)
	}
}
