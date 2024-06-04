Huffman compresser implementation in Go


### Compression
````
go run main.go -compress -input input.txt -output output.bin -password PASSWORD
````

### deCompression
````
go run main.go -decompress -input output.bin -output output.txt -password PASSWORD
````