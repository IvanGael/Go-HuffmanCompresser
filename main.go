package main

import (
	"fmt"
	"io"
	"os"
	"sort"
)

// Define a struct for the Huffman tree node
type Node struct {
	Symbol      byte
	Frequency   int
	Left, Right *Node
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
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Frequency < nodes[j].Frequency
	})
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
	encoded := ""
	for _, symbol := range data {
		encoded += codes[symbol]
	}
	return encoded
}

// Function to write encoded data to file
func writeEncodedData(encodedData string, codes HuffmanCode, outputFileName string) error {
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Write Huffman codes to file for decompression
	for symbol, code := range codes {
		fmt.Fprintf(outputFile, "%c:%s\n", symbol, code)
	}

	// Write encoded data to file
	_, err = outputFile.WriteString(encodedData)
	return err
}

func main() {
	// Read input data
	inputFileName := "input.txt"
	inputFile, err := os.Open(inputFileName)
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
	outputFileName := "compressed.bin"
	err = writeEncodedData(encodedData, codes, outputFileName)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Compression successful. Output written to", outputFileName)
}
