package pkg_test

import (
	"testing"

	"github.com/project-safari/zebra/cmd/herd/pkg"
	"github.com/stretchr/testify/assert"
)

func TestGroupLabel(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	labels := pkg.CreateLabels()

	// test for generating group label based on address.
	assert.NotNil(pkg.GroupLabels(labels, pkg.Addresses()))

	// test to see if group is created for given address.
	groupTest := pkg.GroupLabels(labels, "Mexico")
	assert.True(groupTest.MatchEqual("group", "Mexico"))
}
