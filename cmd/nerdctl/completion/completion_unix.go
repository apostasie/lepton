//go:build unix

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

package completion

import (
	"github.com/spf13/cobra"

	"github.com/containerd/nerdctl/v2/cmd/nerdctl/helpers"
	"github.com/containerd/nerdctl/v2/leptonic/services/containerd"
	"github.com/containerd/nerdctl/v2/pkg/infoutil"
)

func NetworkDrivers(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	return []string{"bridge", "macvlan", "ipvlan"}, cobra.ShellCompDirectiveNoFileComp
}

func IPAMDrivers(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	return []string{"default", "host-local", "dhcp"}, cobra.ShellCompDirectiveNoFileComp
}

func SnapshotterNames(cmd *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	globalOptions, err := helpers.ProcessRootCmdFlags(cmd)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	client, ctx, cancel, err := containerd.NewClient(cmd.Context(), globalOptions.Namespace, globalOptions.Address)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	defer cancel()

	snapshotterPlugins, err := infoutil.GetSnapshotterNames(ctx, client)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	return snapshotterPlugins, cobra.ShellCompDirectiveNoFileComp
}
