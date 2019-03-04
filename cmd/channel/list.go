// Copyright 2018 TriggerMesh, Inc
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

package channel

import (
	"encoding/json"

	"github.com/triggermesh/tm/pkg/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Channel) List(clientset *client.ConfigSet) ([]byte, error) {
	list, err := clientset.Eventing.EventingV1alpha1().Channels(c.Namespace).List(metav1.ListOptions{})
	if err != nil {
		return []byte{}, err
	}
	// if output == "" {
	// 	table.AddRow("NAMESPACE", "CHANNEL")
	// 	for _, item := range list.Items {
	// 		table.AddRow(item.Namespace, item.Name)
	// 	}
	// 	return table.String(), err
	// }
	return json.Marshal(list)
}
