// Copyright 2019 TriggerMesh, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package generate

import (
	"fmt"

	"github.com/triggermesh/tm/pkg/client"
	"github.com/triggermesh/tm/pkg/file"
	"gopkg.in/yaml.v2"
)

type Project struct {
	Name      string
	Namespace string
	Runtime   string
}

func (p *Project) Generate(clientset *client.ConfigSet) error {
	ss := NewTable()

	sample, exist := (*ss)[p.Runtime]
	if !exist {
		return fmt.Errorf("runtime %q does not exist", p.Runtime)
	}

	var buildArgs []string
	if sample.handler != "" {
		buildArgs = append(buildArgs, fmt.Sprintf("HANDLER=%s", sample.handler))
	}

	provider := file.TriggermeshProvider{
		Name:     "triggermesh",
		Registry: client.Registry,
	}

	functions := map[string]file.Function{
		"example-function": file.Function{
			Source:    sample.name,
			Runtime:   sample.template,
			Buildargs: buildArgs,
			Environment: map[string]string{
				"foo": "bar",
			},
		},
	}

	template := file.Definition{
		Service:     "triggermesh-demo",
		Description: "Sample knative service",
		Provider:    provider,
		Functions:   functions,
	}

	manifest, err := yaml.Marshal(&template)
	if err != nil {
		return err
	}
	if client.Dry {
		fmt.Printf("%s:\n---\n%s\n\n", sample.name, sample.function)
		fmt.Printf("%s:\n---\n%s\n", manifestName, manifest)
		return nil
	}
	if err := file.Write(sample.name, sample.function); err != nil {
		return fmt.Errorf("writing function to file: %s", err)
	}
	if err := file.Write(manifestName, string(manifest)); err != nil {
		return fmt.Errorf("writing manifest to file: %s", err)
	}
	return nil
}
