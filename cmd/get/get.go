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

package get

import (
	"encoding/json"

	"github.com/gosuri/uitable"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/triggermesh/tm/pkg/client"
	yaml "gopkg.in/yaml.v2"
)

var (
	log    *logrus.Logger
	table  *uitable.Table
	output string
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Retrieve resources from k8s cluster",
}

func NewGetCmd(clientset *client.ClientSet, log *logrus.Logger, encode *string) *cobra.Command {
	getCmd.AddCommand(cmdListBuild(clientset))
	getCmd.AddCommand(cmdListBuildtemplates(clientset))
	getCmd.AddCommand(cmdListConfigurations(clientset))
	getCmd.AddCommand(cmdListRevision(clientset))
	getCmd.AddCommand(cmdListRoute(clientset))
	getCmd.AddCommand(cmdListService(clientset))

	table = uitable.New()
	table.Wrap = true
	table.MaxColWidth = 50
	output = *encode

	return getCmd
}

func format(v interface{}) (string, error) {
	switch output {
	case "json":
		o, err := json.MarshalIndent(v, "", "    ")
		return string(o), err
	case "yaml":
		o, err := yaml.Marshal(v)
		return string(o), err
	}
	return "", nil
}
