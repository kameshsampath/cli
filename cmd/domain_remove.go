package cmd

import (
	"errors"
	"fmt"

	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"os"

	"github.com/spf13/cobra"
)

var domainRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm", "delete", "destroy"},
	Short:   "Remove a domain",
	Example: "civo domain remove DOMAIN/DOMAIN_ID",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		domain, err := client.FindDNSDomain(args[0])
		if err != nil {
			if errors.Is(err, civogo.ZeroMatchesError) {
				utility.Error("sorry this domain (%s) does not exist in your account", args[0])
				os.Exit(1)
			}
			if errors.Is(err, civogo.MultipleMatchesError) {
				utility.Error("sorry we found more than one domain with that name in your account", args[0])
				os.Exit(1)
			}
		}

		if utility.UserConfirmedDeletion("domain", defaultYes) == true {

			_, err = client.DeleteDNSDomain(domain)

			ow := utility.NewOutputWriterWithMap(map[string]string{"ID": domain.ID, "Name": domain.Name})

			switch outputFormat {
			case "json":
				ow.WriteSingleObjectJSON()
			case "custom":
				ow.WriteCustomOutput(outputFields)
			default:
				fmt.Printf("The domain called %s with ID %s was deleted\n", utility.Green(domain.Name), utility.Green(domain.ID))
			}
		} else {
			fmt.Println("Operation aborted")
		}
	},
}
