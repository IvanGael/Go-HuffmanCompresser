package main

import (
	"bufio"
	"flag"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestBuildHuffmanTree(t *testing.T) {
	// Define input symbol frequencies
	frequencies := map[rune]int{'a': 5, 'b': 10, 'c': 15}

	// Build Huffman tree
	root := buildHuffmanTree(frequencies)

	// Assert that the root node is not nil
	if root == nil {
		t.Error("Expected non-nil root node")
	}

	// Assert that all symbols are present in the tree
	for symbol := range frequencies {
		if !symbolExistsInTree(root, symbol) {
			t.Errorf("Symbol %c is missing in the Huffman tree", symbol)
		}
	}

	// Assert that the frequencies are correctly propagated throughout the tree
	assertFrequencies(t, root, frequencies)

	// Assert that the tree satisfies the Huffman coding property
	if !huffmanCodingPropertySatisfied(root) {
		t.Error("Huffman coding property is not satisfied")
	}
}

// Helper function to check if a symbol exists in the Huffman tree
func symbolExistsInTree(node *Node, symbol rune) bool {
	if node == nil {
		return false
	}
	if node.Symbol == symbol {
		return true
	}
	return symbolExistsInTree(node.Left, symbol) || symbolExistsInTree(node.Right, symbol)
}

// Helper function to assert that frequencies are correctly propagated throughout the tree
func assertFrequencies(t *testing.T, node *Node, frequencies map[rune]int) {
	if node == nil {
		return
	}
	if node.Symbol != 0 {
		if freq, ok := frequencies[node.Symbol]; ok {
			if node.Frequency != freq {
				t.Errorf("Frequency mismatch for symbol %c: expected %d, got %d", node.Symbol, freq, node.Frequency)
			}
		} else {
			t.Errorf("Symbol %c is not found in the input frequencies", node.Symbol)
		}
	}
	assertFrequencies(t, node.Left, frequencies)
	assertFrequencies(t, node.Right, frequencies)
}

// Helper function to check if the Huffman coding property is satisfied
func huffmanCodingPropertySatisfied(node *Node) bool {
	if node == nil {
		return true
	}
	if node.Left == nil && node.Right == nil {
		return true
	}
	if node.Left != nil && node.Right != nil {
		return huffmanCodingPropertySatisfied(node.Left) && huffmanCodingPropertySatisfied(node.Right)
	}
	return false
}

func TestSortNodes(t *testing.T) {
	nodes := []*Node{
		{Frequency: 5},
		{Frequency: 10},
		{Frequency: 3},
	}

	sortNodes(nodes)

	// Verify that nodes are sorted by frequency
	for i := 1; i < len(nodes); i++ {
		if nodes[i].Frequency < nodes[i-1].Frequency {
			t.Error("Nodes are not sorted correctly")
		}
	}
}

func TestGenerateCodes(t *testing.T) {
	// Create a sample Huffman tree for testing
	root := &Node{
		Left:  &Node{Symbol: 'a'},
		Right: &Node{Symbol: 'b'},
	}

	codes := generateCodes(root)

	expected := HuffmanCode{'a': "0", 'b': "1"}
	if !reflect.DeepEqual(codes, expected) {
		t.Errorf("Generated codes: %v, Expected: %v", codes, expected)
	}
}

func TestEncode(t *testing.T) {
	codes := HuffmanCode{'a': "0", 'b': "10", 'c': "110", 'd': "111"}
	data := "abcdaba"

	encoded := encode(data, codes)

	expected := "0101101110100"
	if encoded != expected {
		t.Errorf("Encoded data: %s, Expected: %s", encoded, expected)
	}
}

func TestDecode(t *testing.T) {
	codes := HuffmanCode{'a': "0", 'b': "10", 'c': "110", 'd': "111"}
	encoded := "0101101110100"

	decoded := decode(encoded, codes)

	expected := "abcdaba"
	if decoded != expected {
		t.Errorf("Decoded data: %s, Expected: %s", decoded, expected)
	}
}

func TestReadCodes(t *testing.T) {
	input := "97:0\n98:10\n99:110\n100:111\n----END CODES----\n"

	scanner := bufio.NewScanner(strings.NewReader(input))
	codes, err := readCodes(scanner)

	if err != nil {
		t.Errorf("Error reading codes: %v", err)
	}

	expected := HuffmanCode{'a': "0", 'b': "10", 'c': "110", 'd': "111"}
	if !reflect.DeepEqual(codes, expected) {
		t.Errorf("Read codes: %v, Expected: %v", codes, expected)
	}
}

func TestMain_CompressAndDecompress(t *testing.T) {
	// Test compression and decompression using sample input files and passwords
	inputFileName := "input.txt"
	outputFileName := "output.bin"
	password := "testpassword"

	// Create a sample input file
	sampleInput := "Hello, World!\nThis is a test.\nSpecial characters: £¤"
	err := os.WriteFile(inputFileName, []byte(sampleInput), 0644)
	if err != nil {
		t.Fatalf("Error creating sample input file: %v", err)
	}

	// Test compression
	compressArgs := []string{"-compress", "-input", inputFileName, "-output", outputFileName, "-password", password}
	flag.CommandLine = flag.NewFlagSet("test", flag.ExitOnError)
	os.Args = append([]string{"test"}, compressArgs...)
	main()

	// Test decompression
	decompressArgs := []string{"-decompress", "-input", outputFileName, "-output", "output.txt", "-password", password}
	flag.CommandLine = flag.NewFlagSet("test", flag.ExitOnError)
	os.Args = append([]string{"test"}, decompressArgs...)
	main()

	// Compare the original input file with the decompressed output file
	originalData, err := os.ReadFile(inputFileName)
	if err != nil {
		t.Errorf("Error reading original input file: %v", err)
	}

	decodedData, err := os.ReadFile("output.txt")
	if err != nil {
		t.Errorf("Error reading decompressed output file: %v", err)
	}

	if string(originalData) != string(decodedData) {
		t.Errorf("Decompressed data does not match original input data")
	}
}
