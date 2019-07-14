package cobra

import (
	"fmt"
	"secretSafe"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Deletes a secret from the vault",
	Run: func(cmd *cobra.Command, args []string) {
		v := secretSafe.File(encodingKey, secretsPath())

		for _, k := range args {
			err := v.Delete(k)

			if err != nil {
				fmt.Println(err)
				break
			}

			fmt.Printf("Secret %s is removed from super secret vault!\n", k)
		}

	},
}

func init() {
	RootCmd.AddCommand(deleteCmd)
}
