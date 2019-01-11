package set

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/triggermesh/tm/pkg/client"
)

var r Route
var c Credentials

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set resource parameters",
}

// NewSetCmd returns "Set" cobra CLI command with its subcommands
func NewSetCmd(clientset *client.ConfigSet) *cobra.Command {
	setCmd.AddCommand(cmdSetRoutes(clientset))
	setCmd.AddCommand(cmdSetRegistryCreds(clientset))
	setCmd.AddCommand(cmdSetPullSecret(clientset))
	return setCmd
}

func cmdSetRoutes(clientset *client.ConfigSet) *cobra.Command {
	setRoutesCmd := &cobra.Command{
		Use:   "route",
		Short: "Configure service route",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := r.SetPercentage(args, clientset); err != nil {
				log.Fatalln(err)
			}
			fmt.Println("Routes successfully updated")
		},
	}

	setRoutesCmd.Flags().StringSliceVarP(&r.Revisions, "revisions", "r", []string{}, "Set traffic percentage for revision")
	setRoutesCmd.Flags().StringSliceVarP(&r.Configs, "configurations", "c", []string{}, "Set traffic percentage for configuration")
	return setRoutesCmd
}

func cmdSetRegistryCreds(clientset *client.ConfigSet) *cobra.Command {
	setRegistryCredsCmd := &cobra.Command{
		Use:   "registry-auth",
		Short: "Create secret with registry credentials",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := c.SetRegistryCreds(args, clientset); err != nil {
				log.Fatalln(err)
			}
			fmt.Println("Registry credentials set")
		},
	}

	setRegistryCredsCmd.Flags().StringVar(&c.Host, "registry", "", "Registry host address")
	setRegistryCredsCmd.Flags().StringVar(&c.Username, "username", "", "Registry username")
	setRegistryCredsCmd.Flags().StringVar(&c.Password, "password", "", "Registry password")
	return setRegistryCredsCmd
}

func cmdSetPullSecret(clientset *client.ConfigSet) *cobra.Command {
	setPullSecretCmd := &cobra.Command{
		Use:   "pull-secret",
		Short: "Image pull secret for service account",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if err := c.SetPullSecret(args, clientset); err != nil {
				log.Fatalln(err)
			}
			fmt.Println("Image pull secret created")
		},
	}

	setPullSecretCmd.Flags().StringVar(&c.Host, "registry", "", "Registry host address")
	setPullSecretCmd.Flags().StringVar(&c.Username, "username", "", "Registry username")
	setPullSecretCmd.Flags().StringVar(&c.Password, "password", "", "Registry password")
	setPullSecretCmd.Flags().StringVar(&c.Password, "email", "", "User email")
	return setPullSecretCmd
}
