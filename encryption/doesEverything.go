package main

import (
	"os"
	"fmt"
	"io/ioutil"
	"os/exec"
)

func main() {
deleteFile("keys/id_rsa")
deleteFile("keys/id_rsa.pub")



cmd("","", "", "", "", "", "","", "", "", "", "openssl", "rand", "-out", "keys/rsaPasswd.key", "32")
cmd("","","", "", "","","","ssh-keygen", "-t", "rsa", "-b", "4096", "-N", readFile("keys/rsaPasswd.key"), "-f", "keys/id_rsa")

fmt.Println("The file to be encrypted is " + readFile("secretfile.txt"))

cmd("", "","", "", "", "", "","", "", "", "", "openssl", "rand", "-out", "keys/aes.key", "32")
fmt.Println("The aes key to be encrypted is " + readFile("keys/aes.key"))

cmd("", "","", "", "","","","", "openssl", "aes-256-cbc", "-in", "secretfile.txt", "-out", "secretfile.txt.enc", "-pass", "file:keys/aes.key")
fmt.Println("The encrypted file is " + readFile("secretfile.txt.enc"))

cmd("openssl", "rsautl", "-encrypt", "-oaep", "-pubin", "-inkey", "<(ssh-keygen", "", "-f", "recipients-key.pub", "-m", "PKCS8)" , "-in","keys/aes.key","-out","keys/aes.key.enc")
fmt.Println("The encrypted aes key is " + readFile("aes.key.enc"))

deleteFile("keys/aes.key")
cmd("","", "", "", "", "","openssl", "rsautl", "-decrypt", "-oaep", "-inkey", "keys/id_rsa", "-in", "aes.key.enc", "-out", "aes.key")
fmt.Println("The unencrypted aes key is " + readFile("aes.key"))
secretKey := readFile("file:" + "keys/aes.key")

cmd("", "", "","", "", "", "", "openssl", "aes-256-cbc", "-d", "-in", "secretfile.txt.enc", "-out", "secretfile.txt", "-pass", secretKey)
fmt.Println("The decrypted file is " + readFile("secretfile.txt"))




}




func deleteFile(path string) {
	// delete file
	var err = os.Remove(path)
	if isError(err) { return }

	fmt.Println("==> done deleting file")
}

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return (err != nil)
}

//run ipfs commands from go
func cmd(one string, two string, three string, four string, five string, six string, seven string, eight string, nine string, ten string, eleven string, twelve string, thirteen string, fourteen string, fifteen string, sixteen string) {
	output, err := exec.Command(one, two, three, four, five, six, seven, eight, nine, ten, eleven, twelve, thirteen, fourteen, fifteen, sixteen).CombinedOutput()
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
	fmt.Println(string(output))
}

//reads out a file and I use it to read intoa string
func readFile(file string) string {
	b, err := ioutil.ReadFile(file) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	return string(b)
}
