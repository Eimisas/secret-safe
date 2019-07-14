package cobra

import (
	"fmt"
	"secretSafe"

	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Sets a secret in your predefined safe",
	Run: func(cmd *cobra.Command, args []string) {
		v := secretSafe.File(encodingKey, secretsPath())
		key, value := args[0], args[1]
		err := v.Set(key, value)
		if err != nil {
			panic(err)
		}
		fmt.Println("Value is set successfully")
	},
}

func init() {
	RootCmd.AddCommand(setCmd)
}
