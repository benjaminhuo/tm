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

package taskrun

// TaskRun represents tekton TaskRun object
type TaskRun struct {
	Name             string
	Namespace        string
	Params           []string
	PipelineResource Resource
	Registry         string
	RegistrySecret   string
	Source           Git
	Task             Resource
	Timeout          string
	Wait             bool
}

type Resource struct {
	Name  string
	Owned bool
}

type Git struct {
	URL      string
	Revision string
}

type registryAuths struct {
	Auths registry
}

type credentials struct {
	Username string
	Password string
}

type registry map[string]credentials
