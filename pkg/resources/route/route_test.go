// Copyright 2020 TriggerMesh Inc.
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

package route

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/triggermesh/tm/pkg/client"
)

func TestList(t *testing.T) {
	namespace := "test-namespace"
	if ns, ok := os.LookupEnv("NAMESPACE"); ok {
		namespace = ns
	}
	routeClient, err := client.NewClient(client.ConfigPath(""))
	assert.NoError(t, err)

	r := &Route{Name: "Foo", Namespace: namespace}

	_, err = r.List(&routeClient)
	assert.NoError(t, err)
}

func TestGet(t *testing.T) {
	namespace := "test-namespace"
	if ns, ok := os.LookupEnv("NAMESPACE"); ok {
		namespace = ns
	}
	routeClient, err := client.NewClient(client.ConfigPath(""))
	assert.NoError(t, err)

	r := &Route{Name: "Foo", Namespace: namespace}
	result, err := r.Get(&routeClient)
	assert.Error(t, err)
	assert.Equal(t, "", result.Name)
}

func TestDelete(t *testing.T) {
	namespace := "test-namespace"
	if ns, ok := os.LookupEnv("NAMESPACE"); ok {
		namespace = ns
	}
	routeClient, err := client.NewClient(client.ConfigPath(""))
	assert.NoError(t, err)

	r := &Route{Name: "Foo", Namespace: namespace}
	err = r.Delete(&routeClient)
	assert.Error(t, err)
}
