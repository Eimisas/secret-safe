package cobra

import (
	"fmt"
	"secretSafe"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all secrets within the file",
	Run: func(cmd *cobra.Command, args []string) {
		v := secretSafe.File(encodingKey, secretsPath())

		secrets, err := v.ReturnAll()

		if err != nil {
			fmt.Println("Error while opening the file. Check if the key is correct and the Vault file is in the right directory")
		} else {

			fmt.Println("Secrets in the encrypted file: ")
			fmt.Println("-----------------------------")

			for k := range secrets {
				fmt.Printf("%s\n", k)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
