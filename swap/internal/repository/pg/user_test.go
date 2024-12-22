package pg

import (
	"context"
	"testing"

	"code.emcdtech.com/b2b/swap/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestUser_AddFindOne(t *testing.T) {
	ctx := context.Background()

	defer func() {
		err := truncateAll(ctx)
		if err != nil {
			t.Fatal(err)
		}
	}()

	rep := NewUser(db)

	us := model.User{
		ID:       uuid.New(),
		Email:    "email@gmail.com",
		Language: "ru",
	}

	err := rep.Add(ctx, &us)
	require.NoError(t, err)

	user, err := rep.FindOne(ctx, &model.UserFilter{ID: &us.ID})
	require.NoError(t, err)

	require.Equal(t, &us, user)
}
