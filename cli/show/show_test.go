package show_test

import (
	"testing"

	"github.com/project-safari/zebra/cli/show"
	"github.com/stretchr/testify/assert"
)

func TestNewZebraCommand(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	cmd := show.NewZebra()
	assert.NotNil(cmd)
}

func TestNew(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	cmd := show.New()
	assert.NotNil(cmd)
}

func TestNewNetCmd(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	netC := show.NewNetCmd()
	assert.NotNil(netC)
}

func TestNewSrvCmd(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	srvC := show.NewSrvCmd()
	assert.NotNil(srvC)
}

func TestNewDcCmd(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	dcC := show.NewDCCmd()
	assert.NotNil(dcC)
}
