Huffman compresser implementation in Go


### Compression (CLI Mode)
````
go run main.go -compress -input input.txt -output output.bin -password PASSWORD
````

### deCompression (CLI Mode)
````
go run main.go -decompress -input output.bin -output output.txt -password PASSWORD
````

### Web Mode
````
go run main.go --serve
````