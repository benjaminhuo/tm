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

import (
	"fmt"
	"time"

	v1alpha1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	"github.com/triggermesh/tm/pkg/client"
	"github.com/triggermesh/tm/pkg/printer"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/duration"
	"knative.dev/pkg/apis"
)

func (tr *TaskRun) GetTable(list *v1alpha1.TaskRunList) printer.Table {
	table := printer.Table{
		Headers: []string{
			"Namespace",
			"Name",
			"Age",
			"Ready",
			"Reason",
		},
		Rows: make([][]string, 0, len(list.Items)),
	}

	for _, item := range list.Items {
		table.Rows = append(table.Rows, tr.Row(&item))
	}
	return table
}

func (tr *TaskRun) Row(item *v1alpha1.TaskRun) []string {
	name := item.Name
	namespace := item.Namespace
	age := duration.HumanDuration(time.Since(item.GetCreationTimestamp().Time))
	ready := fmt.Sprintf("%v", item.Status.GetCondition(apis.ConditionSucceeded).IsTrue())
	reason := ""
	if !item.Status.GetCondition(apis.ConditionSucceeded).IsTrue() {
		reason = item.Status.GetCondition(apis.ConditionSucceeded).Message
	}

	row := []string{
		namespace,
		name,
		age,
		ready,
		reason,
	}

	return row
}

func (tr *TaskRun) List(clientset *client.ConfigSet) (*v1alpha1.TaskRunList, error) {
	return clientset.TektonPipelines.TektonV1alpha1().TaskRuns(tr.Namespace).List(metav1.ListOptions{})
}
