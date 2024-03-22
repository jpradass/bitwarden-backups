package main

import "github.com/jpradass/bitwarden-backups/api/bitwarden"

func main() {
	bw := bitwarden.New()

	bw.ListCollections()
}
