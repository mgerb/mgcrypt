package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"golang.org/x/crypto/pbkdf2"
)

type flagsType struct {
	key     string
	decrypt bool
	output  string
}

var flags = &flagsType{}

func init() {

	// read flags
	keyPtr := flag.String("k", "", "Encryption key - this is your password")
	decryptPtr := flag.Bool("d", false, "Decrypt")
	outputPtr := flag.String("o", "", "Output file")

	flag.Parse()

	if *keyPtr == "" {
		fmt.Println("Invalid Arguments: Please specify a key with -k")
		os.Exit(0)
	}

	// make sure input file exists
	if len(flag.Args()) != 1 {
		fmt.Println("Invalid Arguments: Please specify one file input")
		os.Exit(0)
	}

	flags.key = *keyPtr
	flags.decrypt = *decryptPtr
	flags.output = *outputPtr
}

func main() {

	data, err := readFile(flag.Args()[0])

	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	var output []byte

	if flags.decrypt {
		output, err = decrypt([]byte(flags.key), data)
	} else {
		output, err = encrypt([]byte(flags.key), data)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	if flags.output != "" {
		writeFile(flags.output, output)
	} else {
		fmt.Printf("%s", output)
	}
}

// get the first 128 bits of the sha256 checksum of the file
func getEncryptionSalt(filedata []byte) []byte {
	sum := sha256.Sum256(filedata)
	return sum[:16]
}

// get the first 128 bits from the file as the salt
func getDecriptionSalt(filedata []byte) []byte {
	return filedata[:16]
}

// generate AES key based on key input and salt from file
func generateKey(key, salt []byte) []byte {
	return pbkdf2.Key(key, salt, 4096, 32, sha256.New)
}

func readFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func writeFile(filename string, data []byte) error {
	return ioutil.WriteFile(filename, data, 0644)
}

func encrypt(inputKey, data []byte) ([]byte, error) {
	salt := getEncryptionSalt(data)
	key := generateKey(inputKey, salt)

	// Empty array of 16 + data length
	// Include the IV at the beginning
	ciphertext := make([]byte, aes.BlockSize+len(data))

	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return ciphertext, err
	}

	// Slice of first 16 bytes
	iv := ciphertext[:aes.BlockSize]

	// Write 16 rand bytes to fill iv
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return ciphertext, err
	}

	// Return an encrypted stream
	stream := cipher.NewCFBEncrypter(block, iv)

	// Encrypt bytes from plaintext to ciphertext
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	// add the salt to the beginning of the ciphertext
	return append(salt, ciphertext...), nil
}

func decrypt(inputKey, data []byte) ([]byte, error) {
	salt := getDecriptionSalt(data)
	key := generateKey(inputKey, salt)

	// strip the salt from the file
	data = data[16:]

	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return data, err
	}

	// Before even testing the decryption,
	// if the text is too small, then it is incorrect
	if len(data) < aes.BlockSize {
		return data, errors.New("Text is too short")
	}

	// Get the 16 byte IV
	iv := data[:aes.BlockSize]

	// Remove the IV from the ciphertext
	data = data[aes.BlockSize:]

	// Return a decrypted stream
	stream := cipher.NewCFBDecrypter(block, iv)

	// Decrypt bytes from ciphertext
	stream.XORKeyStream(data, data)

	return data, nil
}
