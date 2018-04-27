// Copyright 2018 Google LLC All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/google/go-containerregistry/name"
)

var (
	baseImage, _ = name.NewTag("gcr.io/distroless/base:latest", name.WeakValidation)
)

func main() {
	// Parent command to which all subcommands are added.
	cmds := &cobra.Command{
		Use:   "ko",
		Short: "Rapidly iterate with Go, Containers, and Kubernetes.",
		Long:  "Long Desc", // K8s has a helper here?
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	addKubeCommands(cmds)

	if err := cmds.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}