package main

import (
	"bytes"
	"os"
	"reflect"
	"testing"
)

func TestBuildHuffmanTree(t *testing.T) {
	frequencies := map[rune]int{
		'a': 5,
		'b': 2,
		'c': 1,
		'd': 3,
	}
	root := buildHuffmanTree(frequencies)
	if root == nil {
		t.Error("Expected non-nil root node")
	}
	if root != nil && root.Frequency != 11 {
		t.Errorf("Expected root frequency to be 11, got %d", root.Frequency)
	}
}

func TestGenerateCodes(t *testing.T) {
	root := &Node{
		Frequency: 11,
		Left: &Node{
			Symbol:    'a',
			Frequency: 5,
		},
		Right: &Node{
			Frequency: 6,
			Left: &Node{
				Symbol:    'b',
				Frequency: 2,
			},
			Right: &Node{
				Frequency: 4,
				Left: &Node{
					Symbol:    'c',
					Frequency: 1,
				},
				Right: &Node{
					Symbol:    'd',
					Frequency: 3,
				},
			},
		},
	}
	expectedCodes := HuffmanCode{
		'a': "0",
		'b': "10",
		'c': "110",
		'd': "111",
	}
	codes := generateCodes(root)
	if !reflect.DeepEqual(codes, expectedCodes) {
		t.Errorf("Expected codes %v, got %v", expectedCodes, codes)
	}
}

func TestEncode(t *testing.T) {
	codes := HuffmanCode{
		'a': "0",
		'b': "10",
		'c': "110",
		'd': "111",
	}
	data := "abcd"
	expected := "010110111"
	encoded := encode(data, codes)
	if encoded != expected {
		t.Errorf("Expected encoded data to be %s, got %s", expected, encoded)
	}
}

func TestDecode(t *testing.T) {
	codes := HuffmanCode{
		'a': "0",
		'b': "10",
		'c': "110",
		'd': "111",
	}
	encodedData := "01011011"
	expected := "abc"
	decoded := decode(encodedData, codes)
	if decoded != expected {
		t.Errorf("Expected decoded data to be %s, got %s", expected, decoded)
	}
}

func TestEncryptDecrypt(t *testing.T) {
	plaintext := []byte("Hello, World!")
	password := "testpassword"

	encrypted, err := encrypt(plaintext, password)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	decrypted, err := decrypt(encrypted, password)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("Decrypted text doesn't match original plaintext")
	}
}

func TestWriteAndReadEncodedData(t *testing.T) {
	tempFile := "test_output.bin"
	defer os.Remove(tempFile)

	encodedData := "01011011"
	password := "testpassword"
	codes := HuffmanCode{
		'a': "0",
		'b': "10",
		'c': "110",
		'd': "111",
	}

	err := writeEncodedData(encodedData, tempFile, password, codes)
	if err != nil {
		t.Fatalf("Failed to write encoded data: %v", err)
	}

	readEncodedData, readCodes, err := readEncodedData(tempFile, password)
	if err != nil {
		t.Fatalf("Failed to read encoded data: %v", err)
	}

	if readEncodedData != encodedData {
		t.Errorf("Read encoded data doesn't match original. Expected %s, got %s", encodedData, readEncodedData)
	}

	if !reflect.DeepEqual(codes, readCodes) {
		t.Errorf("Read codes don't match original. Expected %v, got %v", codes, readCodes)
	}
}
