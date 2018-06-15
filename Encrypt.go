package main

//I renamed the two rands because they interfeared with eachother
import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	cryptorand "crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	mathrand "math/rand"
	"os"
	"time"
	"os/exec"
)


func main() {
	//generates rsa encryption key
	//working on making ir random
	//it was already random
	//gen()

	key := []byte(RandStringRunes(32))
	plaintext := []byte(takeInput("\n\nWhat would you like your message to be?\n"))

	//catch error for the message
	//	fmt.Print(plaintext)
	ciphertext, err := encrypt(key, plaintext)
	if err != nil {
		log.Fatal(err)
	}

	//catch error for the cypher
	//normally because it is the wrong length password
	//if you are ever changing this make sure you always keep a 32 char password
	fmt.Printf("%0x\n", ciphertext)
	result, err := decrypt(key, ciphertext)
	if err != nil {
		log.Fatal(err)
	}
	//if no errors are cought than it just prints out the encrypted result
	fmt.Printf("%s\n", result)

	//	fmt.Println(SignPKCS1v15(io.reader, private.key, , cyphertext))

}
//run ipfs commands from go 
func cmd() {
    output, err := exec.Command("ipfs", "daemon").CombinedOutput()
        if err != nil {
          os.Stderr.WriteString(err.Error())
    }
    fmt.Println(string(output))
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

//take inputs from the user
func takeInput(message string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(message)
	text, _ := reader.ReadString('\n')
	return text
}

// See alternate IV creation from ciphertext below
//var iv = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}
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
