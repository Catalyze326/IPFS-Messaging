//genAES() genRSA() encAESKey() and decryptAESKey() must be run at least once and they must be run in the order I displayed them in
//Without these nothing will work

package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	cryptorand "crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	mathrand "math/rand"
	"os"
	"time"
)

func main() {
	//generates aes key
	writeFile((RandStringRunes(32)), "keys/aes.key")
	//generates rsa keys and then encrypts the aes key with that key it just generated
	rsaEvery()
}

//
// func genAES() {
// 	output3, err3 := exec.Command("wsl", "openssl", "rand", "-out", "keys/aes.key", "32").CombinedOutput()
// 	fmt.Println(string(output3))
// 	if err3 != nil {
// 		os.Stderr.WriteString(err3.Error())
// 	}
// }
//
// func encAESKey() {
// 	output6, err6 := exec.Command("wsl", "openssl", "enc", "-d", "aes-256-cbc", "-in", "keys/aes.key", "-out", "keys/aes.key.enc", "-kfile", "keys/aes.key").CombinedOutput()
// 	fmt.Println(string(output6))
// 	if err6 != nil {
// 		os.Stderr.WriteString(err6.Error())
// 	}
//
// 	output4, err4 := exec.Command("wsl", "openssl", "rsautl", "-encrypt", "-inkey", "keys/rsa_key.pub", "-pubin", "-in", "keys/aes.key", "-out", "keys/aes.key.enc").CombinedOutput()
// 	fmt.Println(string(output4))
// 	if err4 != nil {
// 		os.Stderr.WriteString(err4.Error())
// 	}
//
// 	deleteFile("keys/aes.key")
// }
//
// func decryptAESKey() {
// 	output5, err5 := exec.Command("wsl", "openssl", "rsautl", "-decrypt", "-inkey", "keys/rsa_key.pri", "-in", "keys/aes.key.enc", "-out", "keys/aes.key").CombinedOutput()
// 	fmt.Println(string(output5))
// 	if err5 != nil {
// 		os.Stderr.WriteString(err5.Error())
// 	}
// }

//function to take input from a file
func takeInput(message string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(message)
	text, _ := reader.ReadString('\n')
	fmt.Println("")
	return text
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
