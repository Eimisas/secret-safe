package cobra

import (
	"fmt"
	"secretSafe"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets a secret from your predifined safe",
	Run: func(cmd *cobra.Command, args []string) {
		v := secretSafe.File(encodingKey, secretsPath())
		key := args[0]
		value, err := v.Get(key)
		if err != nil {
			fmt.Println("Error while opening the file. Check if the key / encoding key is correct and the Vault file is in the right directory")
		} else {
			fmt.Printf("%s = %s\n", key, value)
		}
	},
}

func init() {
	RootCmd.AddCommand(getCmd)
}
