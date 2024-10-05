package utils

import "os"

// CheckRoot memeriksa apakah program dijalankan dengan hak akses root
func CheckRoot() bool {
    return os.Geteuid() == 0
}