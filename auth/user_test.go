package auth_test

import (
	"context"
	"testing"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/auth"
	"github.com/stretchr/testify/assert"
)

const userType = "user"

//nolint:funlen
func TestUser(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	all, e := auth.NewPriv("", true, true, true, true)
	assert.Nil(e)
	assert.NotNil(all)
	admin := &auth.Role{"admin", []*auth.Priv{all}}

	writeAll, e := auth.NewPriv("", true, false, true, true)
	assert.Nil(e)
	assert.NotNil(writeAll)
	writer := &auth.Role{"writer", []*auth.Priv{writeAll}}

	readAll, e := auth.NewPriv("", false, true, false, false)
	assert.Nil(e)
	assert.NotNil(all)

	rwOne, e := auth.NewPriv("eden", true, true, true, true)
	assert.Nil(e)
	assert.NotNil(rwOne)

	user := &auth.Role{"user", []*auth.Priv{readAll, rwOne}}

	godKey, err := auth.Generate()
	assert.Nil(err)
	assert.NotNil(godKey)

	ctx := context.Background()

	uType := auth.UserType()
	god, ok := uType.New().(*auth.User)
	assert.True(ok)

	god.Name = "almighty"
	god.ID = "00000000000001"
	god.Type = userType
	assert.NotNil(god.Validate(ctx))

	god.Key = godKey
	assert.NotNil(god.Validate(ctx))

	god.Role = admin
	assert.NotNil(god.Validate(ctx))

	god.PasswordHash = auth.HashPassword("youhaveachoice")
	assert.Nil(god.Validate(ctx))

	luciferKey, err := auth.Generate()
	assert.Nil(err)
	assert.NotNil(luciferKey)

	lucifer := new(auth.User)
	lucifer.Name = "lucifer"
	lucifer.ID = "00000000000002"
	lucifer.Type = userType
	lucifer.Key = luciferKey
	lucifer.Role = writer
	lucifer.PasswordHash = auth.HashPassword("youwillbetempted")

	adamKey, err := auth.Generate()
	assert.Nil(err)
	assert.NotNil(adamKey)

	adam := new(auth.User)
	adam.Name = "adam"
	adam.ID = "00000000000003"
	adam.Type = userType
	adam.Key = adamKey
	adam.Role = user
	adam.PasswordHash = auth.HashPassword("iloveeve")

	eveKey, err := auth.Generate()
	assert.Nil(err)
	assert.NotNil(eveKey)

	eve := new(auth.User)
	eve.Name = "eve"
	eve.ID = "00000000000004"
	eve.Type = userType
	eve.Key = eveKey
	eve.Role = user
	eve.PasswordHash = auth.HashPassword("iloveadam")

	token, err := godKey.Sign([]byte(auth.SharedSecret))
	assert.Nil(err)
	assert.NotEmpty(token)

	assert.Nil(god.Authenticate(string(token)))
	assert.Nil(god.AuthenticatePassword("youhaveachoice"))
	assert.NotNil(lucifer.Authenticate(string(token)))
	assert.NotNil(lucifer.AuthenticatePassword("youhaveachoice"))

	assert.True(god.Create("universe"))
	assert.True(god.Read("universe"))
	assert.True(god.Delete("universe"))
	assert.True(god.Update("universe"))
	assert.True(god.Write("universe"))
	assert.False(lucifer.Read("universe"))
	assert.True(lucifer.Write("universe"))
	assert.False(adam.Write("something"))
	assert.False(adam.Delete("something"))
	assert.False(eve.Create("universe"))
	assert.False(eve.Delete("universe"))
	assert.True(adam.Create("eden"))
	assert.True(eve.Write("eden"))
	assert.True(eve.Update("eden"))
	assert.False(eve.Update("universe"))

	assert.True(eve.Write("eden"))
	assert.True(eve.Update("eden"))
	assert.False(eve.Update("universe"))

	newUser := auth.NewUser("eve", "eve123", eveKey, zebra.Labels{})
	assert.NotNil(newUser)
	assert.Nil(newUser.Validate(context.Background()))
}
