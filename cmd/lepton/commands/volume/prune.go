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

package volume

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"go.farcloser.world/lepton/cmd/lepton/helpers"
	"go.farcloser.world/lepton/leptonic/services/containerd"
	"go.farcloser.world/lepton/pkg/api/options"
	"go.farcloser.world/lepton/pkg/cmd/volume"
)

func pruneCommand() *cobra.Command {
	volumePruneCommand := &cobra.Command{
		Use:           "prune [flags]",
		Short:         "Remove all unused local volumes",
		Args:          cobra.NoArgs,
		RunE:          pruneAction,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	volumePruneCommand.Flags().BoolP("all", "a", false, "Remove all unused volumes, not just anonymous ones")
	volumePruneCommand.Flags().BoolP("force", "f", false, "Do not prompt for confirmation")
	return volumePruneCommand
}

func pruneOptions(cmd *cobra.Command, _ []string) (*options.VolumePrune, error) {
	all, err := cmd.Flags().GetBool("all")
	if err != nil {
		return nil, err
	}

	force, err := cmd.Flags().GetBool("force")
	if err != nil {
		return nil, err
	}

	return &options.VolumePrune{
		All:   all,
		Force: force,
	}, nil
}

func pruneAction(cmd *cobra.Command, args []string) error {
	globalOptions, err := helpers.ProcessRootCmdFlags(cmd)
	if err != nil {
		return err
	}

	opts, err := pruneOptions(cmd, args)
	if err != nil {
		return err
	}

	if !opts.Force {
		var confirm string
		msg := "This will remove all local volumes not used by at least one container."
		msg += "\nAre you sure you want to continue? [y/N] "
		fmt.Fprintf(cmd.OutOrStdout(), "WARNING! %s", msg)
		fmt.Fscanf(cmd.InOrStdin(), "%s", &confirm)

		if strings.ToLower(confirm) != "y" {
			return nil
		}
	}

	cli, ctx, cancel, err := containerd.NewClient(cmd.Context(), globalOptions.Namespace, globalOptions.Address)
	if err != nil {
		return err
	}
	defer cancel()

	return volume.Prune(ctx, cli, cmd.OutOrStdout(), globalOptions, opts)
}
