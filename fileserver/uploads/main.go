package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
    // Buka file untuk dibaca
    file, err := os.Open("config_pew.json")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer file.Close()

    // Buat scanner untuk membaca file baris per baris
    scanner := bufio.NewScanner(file)
    lineCount := 0

    for scanner.Scan() {
        lineCount++
    }

    // Periksa jika ada kesalahan saat membaca file
    if err := scanner.Err(); err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Printf("Jumlah baris dalam file: %d\n", lineCount)
}
