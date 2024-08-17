package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	// "encoding/base64"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Define a struct for the Huffman tree node
type Node struct {
	Symbol    rune
	Frequency int
	Left      *Node
	Right     *Node
}

// Define a type for Huffman code
type HuffmanCode map[rune]string

// Function to build Huffman tree
func buildHuffmanTree(frequencies map[rune]int) *Node {
	var nodes []*Node
	for symbol, freq := range frequencies {
		nodes = append(nodes, &Node{Symbol: symbol, Frequency: freq})
	}
	for len(nodes) > 1 {
		// Sort nodes by frequency
		sortNodes(nodes)
		// Combine two nodes with lowest frequency
		left, right := nodes[0], nodes[1]
		parent := &Node{
			Frequency: left.Frequency + right.Frequency,
			Left:      left,
			Right:     right,
		}
		nodes = append(nodes[2:], parent)
	}
	return nodes[0]
}

// Function to sort nodes by frequency
func sortNodes(nodes []*Node) {
	for i := 1; i < len(nodes); i++ {
		j := i
		for j > 0 && nodes[j].Frequency < nodes[j-1].Frequency {
			nodes[j], nodes[j-1] = nodes[j-1], nodes[j]
			j--
		}
	}
}

// Function to generate Huffman codes
func generateCodes(root *Node) HuffmanCode {
	codes := make(HuffmanCode)
	generateCodeRec(root, "", codes)
	return codes
}

// Recursive function to generate Huffman codes
func generateCodeRec(node *Node, code string, codes HuffmanCode) {
	if node == nil {
		return
	}
	if node.Left == nil && node.Right == nil {
		codes[node.Symbol] = code
		return
	}
	generateCodeRec(node.Left, code+"0", codes)
	generateCodeRec(node.Right, code+"1", codes)
}

// Function to encode data using Huffman codes
func encode(data string, codes HuffmanCode) string {
	var encoded strings.Builder
	for _, symbol := range data {
		encoded.WriteString(codes[symbol])
	}
	return encoded.String()
}

// Function to derive encryption key
func deriveKey(password []byte, salt []byte) ([]byte, []byte) {
	if salt == nil {
		salt = make([]byte, 16)
		if _, err := rand.Read(salt); err != nil {
			panic(err)
		}
	}
	return argon2.IDKey(password, salt, 1, 64*1024, 4, 32), salt
}

// Function to encrypt data
func encrypt(plaintext []byte, password string) ([]byte, error) {
	key, salt := deriveKey([]byte(password), nil)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return append(salt, ciphertext...), nil
}

// Function to decrypt data
func decrypt(ciphertext []byte, password string) ([]byte, error) {
	salt, ciphertext := ciphertext[:16], ciphertext[16:]
	key, _ := deriveKey([]byte(password), salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

// Function to write encoded data to file
func writeEncodedData(encodedData string, outputFileName string, password string, codes HuffmanCode) error {
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Encrypt the entire content
	content := fmt.Sprintf("%s\n----DATA----\n", encodedData)
	for symbol, code := range codes {
		content += fmt.Sprintf("%d:%s\n", symbol, code)
	}
	content += "----END CODES----\n"

	encryptedContent, err := encrypt([]byte(content), password)
	if err != nil {
		return err
	}

	// Write encrypted content to file
	_, err = outputFile.Write(encryptedContent)
	return err
}

// Function to read Huffman codes from file
func readCodes(scanner *bufio.Scanner) (HuffmanCode, error) {
	codes := make(HuffmanCode)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "----END CODES----" {
			break
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			symbol, _ := strconv.Atoi(parts[0])
			codes[rune(symbol)] = parts[1]
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return codes, nil
}

// Function to decode data using Huffman codes
func decode(encodedData string, codes HuffmanCode) string {
	var decoded strings.Builder
	code := ""
	for _, bit := range encodedData {
		code += string(bit)
		if symbol, ok := reverseLookup(codes, code); ok {
			decoded.WriteRune(symbol)
			code = ""
		}
	}
	return strings.ReplaceAll(decoded.String(), "\x00", "\n")
}

// Function to reverse lookup symbol from Huffman codes
func reverseLookup(codes HuffmanCode, code string) (rune, bool) {
	for symbol, c := range codes {
		if c == code {
			return symbol, true
		}
	}
	return 0, false
}

// Function to read encoded data from file
func readEncodedData(inputFileName string, password string) (string, HuffmanCode, error) {
	inputFile, err := os.ReadFile(inputFileName)
	if err != nil {
		return "", nil, err
	}

	// Decrypt the entire content
	decryptedContent, err := decrypt(inputFile, password)
	if err != nil {
		return "", nil, err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(decryptedContent)))

	// Read encoded data
	var encodedData string
	scanner.Scan()
	encodedData = scanner.Text()

	// Read until the separator "----DATA----"
	for scanner.Scan() {
		if scanner.Text() == "----DATA----" {
			break
		}
	}

	// Read Huffman codes
	codes, err := readCodes(scanner)
	if err != nil {
		return "", nil, err
	}

	return encodedData, codes, nil
}

func main() {
	compress := flag.Bool("compress", false, "Compress the input file")
	decompress := flag.Bool("decompress", false, "Decompress the input file")
	inputFileName := flag.String("input", "input.txt", "Input file name")
	outputFileName := flag.String("output", "output.bin", "Output file name")
	password := flag.String("password", "", "Password for accessing the compressed data")
	flag.Parse()

	if *compress {
		// Read input data
		data, err := os.ReadFile(*inputFileName)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// Replace newline characters with a special character
		content := strings.ReplaceAll(string(data), "\n", "\x00")

		// Calculate symbol frequencies
		frequencies := make(map[rune]int)
		for _, symbol := range content {
			frequencies[symbol]++
		}

		// Build Huffman tree
		root := buildHuffmanTree(frequencies)

		// Generate Huffman codes
		codes := generateCodes(root)

		// Encode input data
		encodedData := encode(content, codes)

		// Write encoded data to output file
		err = writeEncodedData(encodedData, *outputFileName, *password, codes)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("Compression successful. Output written to", *outputFileName)
	} else if *decompress {
		// Read input file
		if *password == "" {
			fmt.Println("Error: Password is required for decompression.")
			return
		}

		// Read encoded data from file
		encodedData, codes, err := readEncodedData(*inputFileName, *password)
		if err != nil {
			if err.Error() == "cipher: message authentication failed" {
				fmt.Println("Error reading encoded data : Invalid password")
				return
			}
			fmt.Println("Error reading encoded data:", err)
			return
		}

		// Decode data
		decodedData := decode(encodedData, codes)

		// Write decoded data to output file
		err = os.WriteFile(*outputFileName, []byte(decodedData), 0644)
		if err != nil {
			fmt.Println("Error writing decoded data:", err)
			return
		}

		fmt.Println("Decompression successful. Output written to", *outputFileName)
	} else {
		fmt.Println("Please specify either -compress or -decompress.")
	}
}
