package repository_test

import (
	"context"
	"github.com/stretchr/testify/require"
	"go-labs/internal/repository"
	"go-labs/internal/testutils"
	"go-labs/internal/testutils/factory"
	"go-labs/internal/util"
	"testing"
	"time"
)

func TestNewUserRepo(t *testing.T) {
	var (
		ctx    = context.Background()
		dbConn = testutils.DB()
	)

	t.Cleanup(func() {
		err := dbConn.Close()
		require.NoError(t, err)
	})

	testUser := factory.NewUser()

	userRepo, err := repository.NewUserTestRepo(ctx, dbConn, testUser)
	require.NoError(t, err)

	user, err := userRepo.FindByID(ctx, testUser.ID)
	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, testUser.ID, user.ID)

	user, err = userRepo.FindByEmail(ctx, testUser.Email)
	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, testUser.Email, user.Email)

	user, err = userRepo.FindByEmail(ctx, "test@noemail.com")
	require.EqualError(t, err, repository.ErrRecordNotFound.Error())
	require.Nil(t, user)

	testUser.Name = "Morty"
	err = userRepo.Update(ctx, testUser)
	require.NoError(t, err)

	user, err = userRepo.FindByID(ctx, testUser.ID)
	require.NoError(t, err)
	require.Equal(t, testUser.Email, user.Email)
	require.WithinDuration(t, time.Now(), util.PtrToTime(user.UpdatedAt), time.Second)

	err = userRepo.Delete(ctx, user.ID)
	require.NoError(t, err)
}
