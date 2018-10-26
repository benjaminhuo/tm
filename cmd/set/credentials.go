//    Copyright 2018 TriggerMesh, Inc
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package set

import (
	"fmt"

	"github.com/triggermesh/tm/pkg/client"
	corev1 "k8s.io/api/core/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Credentials struct {
	Host     string
	Username string
	Password string
}

func (c *Credentials) SetCredentials(args []string, clientset *client.ClientSet) error {
	secret := make(map[string]string)
	secret["config.json"] = fmt.Sprintf("{\"auths\":{\"%s\":{\"username\":\"%s\",\"password\":\"%s\"}}}", c.Host, c.Username, c.Password)
	newSecret := corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      args[0],
			Namespace: clientset.Namespace,
		},
		StringData: secret,
	}
	s, err := clientset.Core.CoreV1().Secrets(clientset.Namespace).Get(args[0], metav1.GetOptions{})
	if err == nil {
		newSecret.ObjectMeta.ResourceVersion = s.ObjectMeta.ResourceVersion
		_, err = clientset.Core.CoreV1().Secrets(clientset.Namespace).Update(&newSecret)
		return err
	} else if k8sErrors.IsNotFound(err) {
		_, err = clientset.Core.CoreV1().Secrets(clientset.Namespace).Create(&newSecret)
		return err
	}
	return nil
}
