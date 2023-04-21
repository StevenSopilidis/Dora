package main

import (
	v "github.com/stevensopilidis/dora/vault"
)

func main() {
	vault := v.InitializeVault(&v.InitializeVaultConfig{
		Host: "localhost",
		Port: 5432,
		Db:   "doradb",
		User: "dora",
		Pass: "dora123",
	})
	defer v.CloseVault(vault)
}
