package deploy

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/triggermesh/tm/cmd/delete"
	"github.com/triggermesh/tm/pkg/client"
	corev1 "k8s.io/api/core/v1"
)

func TestBuildDeploy(t *testing.T) {
	configSet, err := client.NewClient("")
	assert.NoError(t, err)
	build := Build{
		Name:          "testbuild",
		Buildtemplate: "knative-go-runtime",
	}
	err = build.Deploy(&configSet)
	assert.NoError(t, err)

	err = delete.Build([]string{"testbuild"}, &configSet)
	assert.NoError(t, err)
}

func TestFromBuildTemplate(t *testing.T) {
	build := Build{}
	buildSpec := build.fromBuildtemplate("fooName", map[string]string{"foo": "bar"})

	assert.Equal(t, "fooName", buildSpec.Template.Name)
	assert.Equal(t, "foo", buildSpec.Template.Arguments[0].Name)
	assert.Equal(t, "bar", buildSpec.Template.Arguments[0].Value)
}

func TestFromBuildSteps(t *testing.T) {
	build := Build{}

	buildSpec := build.fromBuildSteps("testStep", "golang:alpine", []string{"/bin/bash"}, []string{"-c", "cat README.md"}, []corev1.Container{})

	assert.Equal(t, "testStep", buildSpec.Steps[0].Name)
	assert.Equal(t, "golang:alpine", buildSpec.Steps[0].Image)
}
