package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

// Define a struct for the Huffman tree node
type Node struct {
	Symbol    byte
	Frequency int
	Left      *Node
	Right     *Node
}

// Define a type for Huffman code
type HuffmanCode map[byte]string

// Function to build Huffman tree
func buildHuffmanTree(frequencies map[byte]int) *Node {
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
func encode(data []byte, codes HuffmanCode) string {
	var encoded strings.Builder
	for _, symbol := range data {
		encoded.WriteString(codes[symbol])
	}
	return encoded.String()
}

// Function to write encoded data to file
func writeEncodedData(encodedData string, outputFileName string, password string, codes HuffmanCode) error {
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Write password to file
	fmt.Fprintln(outputFile, password)

	// Separator between password and encoded data
	fmt.Fprintln(outputFile, "----DATA----")

	// Write Huffman codes to file
	for symbol, code := range codes {
		fmt.Fprintf(outputFile, "%c:%s\n", symbol, code)
	}

	// Separator between codes and encoded data
	fmt.Fprintln(outputFile, "----END CODES----")

	// Write encoded data to file
	_, err = outputFile.WriteString(encodedData)
	return err
}

// Function to read Huffman codes from file
func readCodes(scanner *bufio.Scanner) (HuffmanCode, error) {
	codes := make(HuffmanCode)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "DATA" {
			break
		}
		if len(line) > 2 && line[1] == ':' {
			codes[line[0]] = line[2:]
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return codes, nil
}

// Function to decode data using Huffman codes
func decode(encodedData string, codes HuffmanCode) []byte {
	var decoded strings.Builder
	code := ""
	for _, bit := range encodedData {
		code += string(bit)
		if symbol, ok := reverseLookup(codes, code); ok {
			decoded.WriteByte(symbol)
			code = ""
		}
	}
	return []byte(decoded.String())
}

// Function to reverse lookup symbol from Huffman codes
func reverseLookup(codes HuffmanCode, code string) (byte, bool) {
	for symbol, c := range codes {
		if c == code {
			return symbol, true
		}
	}
	return 0, false
}

// Function to read encoded data from file
func readEncodedData(inputFileName string) (string, string, HuffmanCode, error) {
	inputFile, err := os.Open(inputFileName)
	if err != nil {
		return "", "", nil, err
	}
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)

	// Read password from the first line
	scanner.Scan()
	password := scanner.Text()

	// Read until the separator "----DATA----"
	for scanner.Scan() {
		if scanner.Text() == "----DATA----" {
			break
		}
	}

	// Read Huffman codes
	codes := make(HuffmanCode)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "----END CODES----" {
			break
		}
		if len(line) > 2 && line[1] == ':' {
			codes[line[0]] = line[2:]
		}
	}

	// Read the rest of the data
	var encodedData string
	for scanner.Scan() {
		encodedData += scanner.Text() + "\n"
	}

	if err := scanner.Err(); err != nil {
		return "", "", nil, err
	}

	return encodedData, password, codes, nil
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
		inputFile, err := os.Open(*inputFileName)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer inputFile.Close()

		data, err := io.ReadAll(inputFile)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// Calculate symbol frequencies
		frequencies := make(map[byte]int)
		for _, symbol := range data {
			frequencies[symbol]++
		}

		// Build Huffman tree
		root := buildHuffmanTree(frequencies)

		// Generate Huffman codes
		codes := generateCodes(root)

		// Encode input data
		encodedData := encode(data, codes)

		// Write encoded data to output file
		err = writeEncodedData(encodedData, *outputFileName, *password, codes)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("Compression successful. Output written to", *outputFileName)
	} else if *decompress {
		// Read input file
		inputFile, err := os.Open(*inputFileName)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer inputFile.Close()

		scanner := bufio.NewScanner(inputFile)

		// Read password
		scanner.Scan()
		// password := scanner.Text()

		if *password == "" {
			fmt.Println("Error: Password is required for decompression.")
			return
		}

		// Read password from file
		_, storedPassword, _, err := readEncodedData(*inputFileName)
		if err != nil {
			fmt.Println("Error reading password from file:", err)
			return
		}

		// Compare passwords
		fmt.Printf("password : %s \n", *password)
		fmt.Printf("storedPassword : %s \n", storedPassword)
		if *password != storedPassword {
			fmt.Println("Error: Incorrect password.")
			return
		}

		// Read Huffman codes
		codes, err := readCodes(scanner)
		if err != nil {
			fmt.Println("Error reading Huffman codes:", err)
			return
		}

		// Read encoded
		// data
		encodedData, _, _, err := readEncodedData(*inputFileName)
		if err != nil {
			fmt.Println("Error reading encoded data:", err)
			return
		}

		// Decode data
		decodedData := decode(encodedData, codes)

		// Write decoded data to output file
		err = os.WriteFile(*outputFileName, decodedData, 0644)
		if err != nil {
			fmt.Println("Error writing decoded data:", err)
			return
		}

		fmt.Println("Decompression successful. Output written to", *outputFileName)
	} else {
		fmt.Println("Please specify either -compress or -decompress.")
	}
}
