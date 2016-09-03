package main

import "os"

func fileExists(filename string) bool {
	prefix := "./"
	prefixBuffer := make([]byte, len(filename)+2)
	prefixBufferPointer := copy(prefixBuffer, prefix)
	copy(prefixBuffer[prefixBufferPointer:], filename)

	prefixed := string(prefixBuffer)

	_, err := os.Stat(prefixed)
	// fmt.Println(prefixed, "?", err == nil)
	return err == nil
}
