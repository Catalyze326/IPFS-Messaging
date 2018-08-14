package main

import (
	"crypto/aes"
	"crypto/cipher"
	cryptorand "crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	mathrand "math/rand"
	"os"
	"time"
)

func main() {
	//make the folder for all the keys
	os.Mkdir("keys", os.ModePerm)

	//generates aes 256 bit key
	//	writeFile((RandStringRunes(32)), "keys/aes.key")

	//encrypts files with aes if the last thing passed in is a false
	//in this case the first object passed in will be the file location of the file to be encrypted
	//the second one will be empty "" because this is where the message would go

	//ths is the other way to do it. THis represents encrypting a message so the third object is true
	//the first is the file which we are not using so it is blank
	//the second is the message to be encrypted
	//the messaes are currently not written to files, but that can be changed if need be
	//	primaryAES("", "Hello Please Encrypt Me", true)

	//generates rsa keys the right way to encrypt and sign keys
	genRSA()

	//either encrypts or decrypts a rsa encrypted file based on whether true or false is passed to it
	//the first file it takes in is the infile and the second is the outfile so if it is decrypting
	//than it needs the encrypted file name first and the decrypted file name second and vise
	//virsa when it is encrypting
	//True means decrypting --False means Decrypting
	encDec("keys/aes.key", "keys/aes.key.enc", false)
	deleteFile("keys/aes.key")
	//decrypt with the same func from a line earlier
	encDec("keys/aes.key.enc", "keys/aes.key", true)

	//	signRSA("keys/aes.key")
	//signs the encrypted aes key with rsa and verrifies it to prove it is coming from a
	//reliable source

}

//delete a file
//this function needs the isError func right below it
func deleteFile(path string) {
	// delete file
	var err = os.Remove(path)
	if isError(err) {
		return
	}

	fmt.Println("==> done deleting file")
}

//func used by the delete func.
func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}
	return (err != nil)
}

//Write a string to a file just takes in a string and a file name and writes it out
func writeFile(text string, fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()

	fmt.Fprintf(file, text)
}

//reads out a file and I use it to read intoa string
func readFile(file string) string {
	b, err := ioutil.ReadFile(file) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	return string(b)
}

//the first string is if you want a file put in the file location
//the second string is a message
//the third bool is true if it is a message and false if it is a file
func primaryAES(fileToEncrypt string, messageToEncrypt string, message bool) []byte {
	var plaintext []byte
	//if it is a message than take in the message as plaintext if it is a file than pull it out of the file and save it as a []byte

	if message == false {
		encryptFile := readFile(fileToEncrypt)
		plaintext = []byte(encryptFile)
	} else {
		plaintext = []byte(messageToEncrypt)
	}
	//print the orraginal text before it was turned into []byte
	//That is what the array of bytes was before I fixed it

	fmt.Println(string(plaintext))

	//generates key for cypher
	key := []byte(readFile("keys/aes.key"))

	//catch error for the message
	ciphertext, err := encrypt(key, plaintext)
	if err != nil {
		log.Fatal(err)
	}

	//catch error for the cypher
	//normally because it is the wrong length password
	//if you are ever changing this make sure you always keep a 32 char password
	fmt.Printf("%0x\n", ciphertext)

	//writes the encrypted text to a file if it came from a file
	//When you pull this back out you might need to cast it back to a []byte to get it to work
	if message == false {
		writeFile(string(ciphertext), (fileToEncrypt + ".enc"))
	}
	return ciphertext

}

//encrypt the message
func encrypt(key, text []byte) ([]byte, error) {

	//catch error
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//turns it into a base64 string so that it can be encrypted in aes
	b := base64.StdEncoding.EncodeToString(text)

	//this encrypts it with aes now that it is in base 64 it can be encrypted into aes
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	//catch error
	if _, err := io.ReadFull(cryptorand.Reader, iv); err != nil {
		return nil, err
	}

	//this turns it into aes block
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))

	//return cyper
	return ciphertext, nil
}

//decrypt the message that we just encrypted
func decrypt(key, text []byte) ([]byte, error) {
	//this creates a new aes 256 key
	block, err := aes.NewCipher(key)

	//catch error
	if err != nil {
		return nil, err
	}
	//catch error
	// I don't know when this would happen, but the text gets much larger when it is encrypted so if it were to be smaller something would have gone wrong
	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short\n")
	}

	//finding the size of the aes encrypted block and creating the var iv and setting text to that value
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]

	//Decrypt the block of data that we created in the encrypt function
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)
	//decode the message. THe reason it is base64 is because we already set it to base 64 in the encryption process
	//so now when we are decrypting it we have to keep that in mind
	data, err := base64.StdEncoding.DecodeString(string(text))
	//catch error
	if err != nil {
		return nil, err
	}
	//return the decrypted data and not an error
	return data, nil
}

//these are charectors that can be included in the alpanumaric string that creates the key for the aes encryption
var letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!@#$%^&*()-=_+[]{}|;':<>?,./`~")

//this Function and the one following it are used to create the key that encrypts the data that is randomly generated with time
//This creates the seed for the random number for the aes key
func init() {
	mathrand.Seed(time.Now().UnixNano())
}

//this created the random number string we need for the aes key. We want a 32 character key, but I made it so that it can take a
// any length key so it takes in the length but we are using 32 because we are using aes 256 bit encryption
func RandStringRunes(n int) string {
	//creates a slice how many chars long you you say which in our case is 32
	b := make([]rune, n)
	//takes slice b and itorates through it and fills it with random chars from letterRunes
	for i := range b {
		b[i] = letterRunes[mathrand.Intn(len(letterRunes))]
	}
	//returnes the filled slice
	return string(b)
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
		out, err = rsa.DecryptOAEP(sha1.New(), cryptorand.Reader, priv, in, []byte(label))
		if err != nil {
			log.Fatalf("decrypt: %s", err)
		}
	} else {
		if label == "" {
			label = inFile
		}
		out, err = rsa.EncryptOAEP(sha1.New(), cryptorand.Reader, &priv.PublicKey, in, []byte(label))
		if err != nil {
			log.Fatalf("encrypt: %s", err)
		}
	}

	// Write data to output file
	if err := ioutil.WriteFile(outFile, out, 0600); err != nil {
		log.Fatalf("write output: %s", err)
	}
}

//generates the rsa public and private keys
//please do not try to do this a different way because the signing will not work unless you
//create the rsa keys such that the private key is defined as a rsa key and the public key
//is ideantified as a public key. I don't know why the rsa libray for go works this way, but
//it does so don't change this to any other way to gen rsa keys
func genRSA() {
	//this is the line that actaully generates the private key.
	priv, err := rsa.GenerateKey(cryptorand.Reader, 1024)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = priv.Validate()
	if err != nil {
		fmt.Println("Validation failed.", err)
	}

	// Get der format. priv_der []byte
	priv_der := x509.MarshalPKCS1PrivateKey(priv)

	// pem.Block
	// blk pem.Block
	priv_blk := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   priv_der,
	}

	// Resultant private key in PEM format.
	// priv_pem string
	priv_pem := string(pem.EncodeToMemory(&priv_blk))

	//	fmt.Printf(priv_pem)
	writeFile(priv_pem, "keys/rsa_key.pri")

	// Public Key generation from the private key
	pub := priv.PublicKey
	pub_der, err := x509.MarshalPKIXPublicKey(&pub)
	if err != nil {
		fmt.Println("Failed to get der format for PublicKey.", err)
		return
	}

	//yes this is supposed to be public key. I don't really know why it works this way, but it does
	pub_blk := pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   pub_der,
	}
	pub_pem := string(pem.EncodeToMemory(&pub_blk))
	//	fmt.Printf(pub_pem)
	writeFile(pub_pem, "keys/rsa_key.pub")
}
