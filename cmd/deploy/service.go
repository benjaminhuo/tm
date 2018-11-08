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
	"errors"
	"fmt"
	"os"
	"path"
	"regexp"
	"time"

	"github.com/triggermesh/tm/cmd/describe"

	buildv1alpha1 "github.com/knative/build/pkg/apis/build/v1alpha1"
	servingv1alpha1 "github.com/knative/serving/pkg/apis/serving/v1alpha1"
	"github.com/triggermesh/tm/pkg/client"
	corev1 "k8s.io/api/core/v1"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Knative build timeout in minutes
const (
	timeout           = 10
	uploadDoneTrigger = "/home/.done"
)

type status struct {
	domain string
	err    error
}

type Options struct {
	PullPolicy     string
	ResultImageTag string
	Buildtemplate  string
	RunRevision    string
	Env            []string
	Labels         []string
	BuildArgs      []string
	Wait           bool
}

type Repository struct {
	URL      string
	Revision string
}

type Registry struct {
	URL string
}

type Image struct {
	Source Repository
	Image  Registry
	Path   string
}

type Service struct {
	From Image
	Options
}

func (s *Service) DeployService(args []string, clientset *client.ClientSet) error {
	configuration := servingv1alpha1.ConfigurationSpec{}
	buildArguments, templateParams := getBuildArguments(fmt.Sprintf("%s/%s-%s", clientset.Registry, clientset.Namespace, args[0]), s.BuildArgs)

	if _, err := describe.BuildTemplate(s.Buildtemplate, clientset); err != nil {
		return err
	}

	switch {
	case len(s.From.Image.URL) != 0:
		configuration = s.fromImage(args)
	case len(s.From.Source.URL) != 0:
		configuration = s.fromSource(args, clientset)
		if err := updateBuildTemplate(s.Buildtemplate, templateParams, clientset); err != nil {
			return err
		}
		configuration.Build.Template = &buildv1alpha1.TemplateInstantiationSpec{
			Name:      s.Buildtemplate,
			Arguments: buildArguments,
		}
	case len(s.From.Path) != 0:
		configuration = s.fromPath(args, clientset)
		configuration.Build.Template = &buildv1alpha1.TemplateInstantiationSpec{
			Name:      s.Buildtemplate,
			Arguments: buildArguments,
		}
	}

	configuration.RevisionTemplate.ObjectMeta = metav1.ObjectMeta{
		Name: args[0],
		Annotations: map[string]string{
			"sidecar.istio.io/inject": "true",
		},
	}

	envVars := []corev1.EnvVar{
		corev1.EnvVar{
			Name:  "timestamp",
			Value: time.Now().Format("2006-01-02 15:04:05"),
		},
	}
	for k, v := range getArgsFromSlice(s.Env) {
		envVars = append(envVars, corev1.EnvVar{Name: k, Value: v})
	}
	configuration.RevisionTemplate.Spec.Container.Env = envVars
	configuration.RevisionTemplate.Spec.Container.ImagePullPolicy = corev1.PullPolicy(s.PullPolicy)

	spec := servingv1alpha1.ServiceSpec{
		RunLatest: &servingv1alpha1.RunLatestType{
			Configuration: configuration,
		},
	}

	if s.RunRevision != "" {
		spec = servingv1alpha1.ServiceSpec{
			Pinned: &servingv1alpha1.PinnedType{
				RevisionName:  s.RunRevision,
				Configuration: configuration,
			},
		}
	}

	serviceObject := servingv1alpha1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "serving.knative.dev/servingv1alpha1",
		},

		ObjectMeta: metav1.ObjectMeta{
			Name:      args[0],
			Namespace: clientset.Namespace,
			CreationTimestamp: metav1.Time{
				time.Now(),
			},
			Labels: getArgsFromSlice(s.Labels),
		},

		Spec: spec,
	}

	if err := s.createOrUpdateObject(serviceObject, clientset); err != nil {
		return err
	}

	if len(s.From.Path) != 0 {
		if _, err := os.Stat(s.From.Path); err != nil {
			return err
		}
		fmt.Println("Uploading sources")
		if err := injectSources(args, s.From.Path, clientset); err != nil {
			return err
		}
	}

	fmt.Printf("Deployment started. Run \"tm -n %s describe service %s\" to see the details\n", clientset.Namespace, args[0])

	if s.Wait {
		fmt.Print("Waiting for ready state")
		domain, err := waitService(args[0], clientset)
		if err != nil {
			return err
		}
		fmt.Printf("\nService domain: %s\n", domain)
	}

	return nil
}

func (s *Service) createOrUpdateObject(serviceObject servingv1alpha1.Service, clientset *client.ClientSet) error {
	_, err := clientset.Serving.ServingV1alpha1().Services(clientset.Namespace).Create(&serviceObject)
	if k8sErrors.IsAlreadyExists(err) {
		service, err := clientset.Serving.ServingV1alpha1().Services(clientset.Namespace).Get(serviceObject.ObjectMeta.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		serviceObject.ObjectMeta.ResourceVersion = service.GetResourceVersion()
		service, err = clientset.Serving.ServingV1alpha1().Services(clientset.Namespace).Update(&serviceObject)
		return err
	}
	return err
}

func (s *Service) fromImage(args []string) servingv1alpha1.ConfigurationSpec {
	return servingv1alpha1.ConfigurationSpec{
		RevisionTemplate: servingv1alpha1.RevisionTemplateSpec{
			Spec: servingv1alpha1.RevisionSpec{
				Container: corev1.Container{
					Image: s.From.Image.URL,
				},
			},
		},
	}
}

func (s *Service) fromSource(args []string, clientset *client.ClientSet) servingv1alpha1.ConfigurationSpec {
	return servingv1alpha1.ConfigurationSpec{
		Build: &buildv1alpha1.BuildSpec{
			Source: &buildv1alpha1.SourceSpec{
				Git: &buildv1alpha1.GitSourceSpec{
					Url:      s.From.Source.URL,
					Revision: s.From.Source.Revision,
				},
			},
		},
		RevisionTemplate: servingv1alpha1.RevisionTemplateSpec{
			Spec: servingv1alpha1.RevisionSpec{
				Container: corev1.Container{
					Image: fmt.Sprintf("%s/%s-%s:%s", clientset.Registry, clientset.Namespace, args[0], s.ResultImageTag),
				},
			},
		},
	}
}

func (s *Service) fromPath(args []string, clientset *client.ClientSet) servingv1alpha1.ConfigurationSpec {
	return servingv1alpha1.ConfigurationSpec{
		Build: &buildv1alpha1.BuildSpec{
			Source: &buildv1alpha1.SourceSpec{
				Custom: &corev1.Container{
					Image:   "library/busybox",
					Command: []string{"sh"},
					Args:    []string{"-c", "while [ -z \"$(ls " + uploadDoneTrigger + ")\" ]; do sleep 1; done; sync; mv /home/" + path.Base(s.From.Path) + "/* /workspace; sync"},
				},
			},
		},
		RevisionTemplate: servingv1alpha1.RevisionTemplateSpec{
			Spec: servingv1alpha1.RevisionSpec{
				Container: corev1.Container{
					Image: fmt.Sprintf("%s/%s-%s:%s", clientset.Registry, clientset.Namespace, args[0], s.ResultImageTag),
				},
			},
		},
	}
}

func getArgsFromSlice(slice []string) map[string]string {
	m := make(map[string]string)
	for _, s := range slice {
		t := regexp.MustCompile("[:=]").Split(s, 2)
		if len(t) != 2 {
			fmt.Printf("Can't parse argument slice %s\n", s)
			continue
		}
		m[t[0]] = t[1]
	}
	return m
}

func updateBuildTemplate(name string, params []buildv1alpha1.ParameterSpec, clientset *client.ClientSet) error {
	buildTemplate, err := clientset.Build.BuildV1alpha1().BuildTemplates(clientset.Namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	// Matching new build parameters with existing to check if need to update build template
	var new bool
	for _, v := range params {
		new = true
		for _, vv := range buildTemplate.Spec.Parameters {
			if v.Name == vv.Name {
				new = false
				break
			}
		}
		if new {
			break
		}
	}

	if new {
		buildTemplate.Spec.Parameters = params
		_, err = clientset.Build.BuildV1alpha1().BuildTemplates(clientset.Namespace).Update(buildTemplate)
	}

	return err
}

func injectSources(args []string, filepath string, clientset *client.ClientSet) error {
	var latestRevision string
	for latestRevision == "" {
		service, err := clientset.Serving.ServingV1alpha1().Services(clientset.Namespace).Get(args[0], metav1.GetOptions{})
		if err != nil {
			return err
		}
		latestRevision = service.Status.LatestCreatedRevisionName
		time.Sleep(2 * time.Second)
	}

	var buildPod string
	for buildPod == "" {
		build, err := clientset.Build.BuildV1alpha1().Builds(clientset.Namespace).Get(latestRevision, metav1.GetOptions{})
		if err != nil {
			return err
		}

		if build.Status.Cluster != nil {
			buildPod = build.Status.Cluster.PodName
		}
		time.Sleep(2 * time.Second)
	}

	var sourceContainer string
	for sourceContainer == "" {
		pod, err := clientset.Core.CoreV1().Pods(clientset.Namespace).Get(buildPod, metav1.GetOptions{})
		if err != nil {
			return err
		}
		for _, v := range pod.Status.InitContainerStatuses {
			if v.Name == "build-step-custom-source" && v.State.Running != nil {
				sourceContainer = v.Name
				break
			}
		}
		time.Sleep(2 * time.Second)
	}

	c := Copy{
		Pod:         buildPod,
		Container:   sourceContainer,
		Source:      filepath,
		Destination: "/home",
	}

	if err := c.Upload(clientset); err != nil {
		return err
	}

	if _, _, err := c.RemoteExec(clientset, "touch "+uploadDoneTrigger, nil); err != nil {
		return err
	}

	return nil
}

func waitService(name string, clientset *client.ClientSet) (string, error) {
	quit := time.After(timeout * time.Minute)
	tick := time.Tick(5 * time.Second)
	for {
		select {
		case <-quit:
			return "", errors.New("Service status wait timeout")
		case <-tick:
			fmt.Print(".")
			domain, err := readyDomain(name, clientset)
			if err != nil {
				return "", err
			} else if domain != "" {
				return domain, nil
			}
		}
	}
	return "", nil
}

func readyDomain(name string, clientset *client.ClientSet) (string, error) {
	service, err := clientset.Serving.ServingV1alpha1().Services(clientset.Namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	for _, v := range service.Status.Conditions {
		if v.Status == corev1.ConditionFalse {
			return "", errors.New(v.Message)
		}
	}
	if service.Status.IsReady() {
		return service.Status.Domain, nil
	}
	return "", nil
}
