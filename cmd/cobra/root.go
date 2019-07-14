package cobra

import (
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "secret",
	Short: "Secret is a secret like a password or api key",
}
var encodingKey string
var saveDirectory string

func init() {
	RootCmd.PersistentFlags().StringVarP(&encodingKey, "key", "k", "", "The key to use when encoding and decoding secrets (empty by default)")
	RootCmd.PersistentFlags().StringVarP(&saveDirectory, "dir", "", "/Desktop/Go/src/secretSafe/db", "The key to use when encoding and decoding secrets (empty by default)")

}

func secretsPath() string {
	home, _ := homedir.Dir()
	//! This next line should be set to home for universal version
	//! Or anywhere else you want to save your Vault file
	return filepath.Join(home+saveDirectory, ".secrets")
}
