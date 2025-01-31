/*
   Copyright Farcloser.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package system

import (
	"github.com/spf13/cobra"

	"github.com/containerd/nerdctl/v2/cmd/nerdctl/helpers"
	"github.com/containerd/nerdctl/v2/leptonic/services/containerd"
	"github.com/containerd/nerdctl/v2/pkg/api/options"
	"github.com/containerd/nerdctl/v2/pkg/cmd/system"
	"github.com/containerd/nerdctl/v2/pkg/formatter"
)

func EventsCommand() *cobra.Command {
	shortHelp := `Get real time events from the server`
	longHelp := shortHelp + "\nNOTE: The output format is not compatible with Docker."
	var eventsCommand = &cobra.Command{
		Use:           "events",
		Args:          cobra.NoArgs,
		Short:         shortHelp,
		Long:          longHelp,
		RunE:          eventsAction,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	eventsCommand.Flags().String("format", "", "Format the output using the given Go template, e.g, '{{json .}}'")
	eventsCommand.RegisterFlagCompletionFunc("format", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{formatter.FormatJSON}, cobra.ShellCompDirectiveNoFileComp
	})
	eventsCommand.Flags().StringSliceP("filter", "f", []string{}, "Filter matches containers based on given conditions")
	return eventsCommand
}

func eventsOptions(cmd *cobra.Command, _ []string) (*options.SystemEvents, error) {
	globalOptions, err := helpers.ProcessRootCmdFlags(cmd)
	if err != nil {
		return nil, err
	}
	format, err := cmd.Flags().GetString("format")
	if err != nil {
		return nil, err
	}
	filters, err := cmd.Flags().GetStringSlice("filter")
	if err != nil {
		return nil, err
	}
	return &options.SystemEvents{
		Stdout:   cmd.OutOrStdout(),
		GOptions: globalOptions,
		Format:   format,
		Filters:  filters,
	}, nil
}

func eventsAction(cmd *cobra.Command, args []string) error {
	options, err := eventsOptions(cmd, args)
	if err != nil {
		return err
	}
	cli, ctx, cancel, err := containerd.NewClient(cmd.Context(), options.GOptions.Namespace, options.GOptions.Address)
	if err != nil {
		return err
	}
	defer cancel()
	return system.Events(ctx, cli, options)
}
