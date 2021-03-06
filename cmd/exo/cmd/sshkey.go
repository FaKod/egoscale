package cmd

import (
	"github.com/spf13/cobra"
)

// sshkeyCmd represents the sshkey command
var sshkeyCmd = &cobra.Command{
	Use:   "sshkey",
	Short: "SSH keys pair management",
}

func init() {
	RootCmd.AddCommand(sshkeyCmd)
}
