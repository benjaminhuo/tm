/*
Copyright (c) 2018 TriggerMesh, Inc

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

package deploy

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/triggermesh/tm/pkg/client"
)

var s Service
var b Build
var bt Buildtemplate

// NewDeployCmd returns deploy cobra command and its subcommands
func NewDeployCmd(clientset *client.ConfigSet) *cobra.Command {
	var file string
	deployCmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy knative resource",
		Run: func(cmd *cobra.Command, args []string) {
			if err := fromYAML(file, clientset); err != nil {
				log.Fatal(err)
			}
		},
	}

	deployCmd.Flags().StringVarP(&file, "file", "f", "serverless.yaml", "Deploy functions defined in yaml")

	deployCmd.AddCommand(cmdDeployService(clientset))
	deployCmd.AddCommand(cmdDeployBuild(clientset))
	deployCmd.AddCommand(cmdDeployBuildTemplate(clientset))

	return deployCmd
}

func cmdDeployService(clientset *client.ConfigSet) *cobra.Command {
	deployServiceCmd := &cobra.Command{
		Use:     "service",
		Aliases: []string{"services", "svc"},
		Short:   "Deploy knative service",
		Args:    cobra.ExactArgs(1),
		Example: "tm -n default deploy service foo --build-template kaniko --from-image gcr.io/google-samples/hello-app:1.0",
		Run: func(cmd *cobra.Command, args []string) {
			s.Name = args[0]
			if err := s.DeployService(clientset); err != nil {
				log.Fatal(err)
			}
		},
	}

	// kept for back compatibility
	deployServiceCmd.Flags().StringVar(&s.Source, "from-path", "", "Local file path to deploy")
	deployServiceCmd.Flags().StringVar(&s.Source, "from-image", "", "Image to deploy")
	deployServiceCmd.Flags().StringVar(&s.Source, "from-source", "", "Git source URL to deploy")

	deployServiceCmd.Flags().StringVar(&s.Source, "source", "s", "Service source to deploy: local folder with sources, git repository or docker image")
	deployServiceCmd.Flags().StringVar(&s.Revision, "revision", "master", "Git revision (branch, tag, commit SHA or ref)")
	deployServiceCmd.Flags().BoolVar(&s.Wait, "wait", false, "Wait for successful service deployment")
	deployServiceCmd.Flags().StringVar(&s.Buildtemplate, "build-template", "", "Build template to use with service")
	deployServiceCmd.Flags().StringVar(&s.ResultImageTag, "tag", "latest", "Image tag to build")
	deployServiceCmd.Flags().StringVar(&s.PullPolicy, "image-pull-policy", "Always", "Image pull policy")
	deployServiceCmd.Flags().StringVar(&s.RunRevision, "run-revision", "", "Revision name to run service on")
	deployServiceCmd.Flags().StringVar(&s.Definition, "definition", "", "Path to function definition yaml file (serverless framework format)")
	deployServiceCmd.Flags().StringSliceVar(&s.BuildArgs, "build-argument", []string{}, "Buildtemplate arguments")
	deployServiceCmd.Flags().StringSliceVarP(&s.Labels, "label", "l", []string{}, "Service labels")
	deployServiceCmd.Flags().StringSliceVarP(&s.Env, "env", "e", []string{}, "Environment variables of the service, eg. `--env foo=bar`")

	return deployServiceCmd
}

func cmdDeployBuildTemplate(clientset *client.ConfigSet) *cobra.Command {
	deployBuildTemplateCmd := &cobra.Command{
		Use:     "buildtemplate",
		Aliases: []string{"buildtemplates", "bldtmpl"},
		Short:   "Deploy knative build template",
		Example: "tm -n default deploy buildtemplate -f https://raw.githubusercontent.com/triggermesh/nodejs-runtime/master/knative-build-template.yaml",
		Run: func(cmd *cobra.Command, args []string) {
			if err := bt.DeployBuildTemplate(clientset); err != nil {
				log.Fatal(err)
			}
		},
	}

	// kept for back compatibility
	deployBuildTemplateCmd.Flags().StringVar(&bt.File, "from-url", "", "Build template yaml URL")
	deployBuildTemplateCmd.Flags().StringVar(&bt.File, "from-file", "", "Local file path to deploy")

	deployBuildTemplateCmd.Flags().StringVar(&bt.File, "file", "f", "Build template yaml URL")
	deployBuildTemplateCmd.Flags().StringVar(&bt.RegistryCreds, "credentials", "", "Name of registry credentials to use in buildtemplate")

	return deployBuildTemplateCmd
}

func cmdDeployBuild(clientset *client.ConfigSet) *cobra.Command {
	deployBuildCmd := &cobra.Command{
		Use:     "build",
		Aliases: []string{"builds"},
		Args:    cobra.ExactArgs(1),
		Short:   "Deploy knative build",
		Example: "tm deploy build foo-builder --source git-repo --buildtemplate kaniko --args IMAGE=knative-local-registry:5000/foo-image",
		Run: func(cmd *cobra.Command, args []string) {
			b.Name = args[0]
			if err := b.DeployBuild(clientset); err != nil {
				log.Fatal(err)
			}
		},
	}

	deployBuildCmd.Flags().StringVar(&b.Source, "source", "", "Git URL to get sources from")
	deployBuildCmd.Flags().StringVar(&b.Revision, "revision", "master", "Git source revision")
	deployBuildCmd.Flags().StringVar(&b.Buildtemplate, "buildtemplate", "", "Buildtemplate name to use with build")
	deployBuildCmd.Flags().StringVar(&b.Step, "step", "", "Build step (container) to run on provided source")
	deployBuildCmd.Flags().StringVar(&b.Image, "image", "", "Image for build step")
	deployBuildCmd.Flags().StringSliceVar(&b.Command, "command", []string{}, "Build step (container) command")
	deployBuildCmd.Flags().StringSliceVar(&b.Args, "args", []string{}, "Build arguments")
	deployBuildCmd.MarkFlagRequired("source")

	return deployBuildCmd
}
