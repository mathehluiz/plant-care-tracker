package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"time"
)

var db ConnectionStorer

func TestMain(m *testing.M) {
	mock, err := StartMock()
	if err != nil {
		log.Fatalln("cannot start mock: ", err)
	}

	db = mock

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestGet(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		purpose string
		wantRes string
		key     string
		wantErr error
		mock    func(t *testing.T)
	}{
		{
			"should get the result as expected",
			"foo",
			"awesome-key",
			nil,
			func(t *testing.T) {
				err := db.Set(ctx, 10*time.Second, "awesome-key", "foo")
				assert.NoError(t, err)
			},
		},
		{
			"should return error, bc dont have the key",
			"",
			"not-set-key",
			redis.Nil,
			func(t *testing.T) {},
		},
	}

	for _, tt := range cases {
		t.Run(tt.purpose, func(t *testing.T) {
			tt.mock(t)

			res, err := db.Get(ctx, tt.key)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantRes, res)
		})
	}
}

func TestGetIncludingKey(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		purpose     string
		wantRes     string
		mustInclude string
		wantErr     error
		mock        func(t *testing.T)
	}{
		{
			"should get the result as expected",
			`[{"id": 2},{"id": 1}]`,
			"foo",
			nil,
			func(t *testing.T) {
				err := db.Set(ctx, 10*time.Second, "hash@bar@foo", `{"id": 1}`)
				assert.NoError(t, err)

				err = db.Set(ctx, 10*time.Second, "hash1@foo", `{"id": 2}`)
				assert.NoError(t, err)
			},
		},
		{
			"should get the result as expected",
			`[{"id": 3}]`,
			"hash3",
			nil,
			func(t *testing.T) {
				err := db.Set(ctx, 10*time.Second, "hash3@bar@foo", `{"id": 3}`)
				assert.NoError(t, err)

				err = db.Set(ctx, 10*time.Second, "hash4@foo", `{"id": 4}`)
				assert.NoError(t, err)
			},
		},
		{
			"should return error, bc dont have the key",
			"",
			"not-set-key",
			redis.Nil,
			func(t *testing.T) {},
		},
	}

	for _, tt := range cases {
		t.Run(tt.purpose, func(t *testing.T) {
			tt.mock(t)

			res, err := db.GetIncludingKey(ctx, tt.mustInclude)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantRes, res)
		})
	}
}

func TestGetKeys(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		purpose     string
		wantRes     []string
		mustInclude string
		wantErr     error
		mock        func(t *testing.T)
	}{
		{
			"should get the result as expected",
			[]string{"key1@Topic:foo", "key2@Topic:foo@Topic:bar"},
			"Topic:foo",
			nil,
			func(t *testing.T) {
				err := db.Set(ctx, 10*time.Second, "key1@Topic:foo", ``)
				assert.NoError(t, err)

				err = db.Set(ctx, 10*time.Second, "key2@Topic:foo@Topic:bar", ``)
				assert.NoError(t, err)
			},
		},
		{
			"should get the result as expected",
			[]string{"key2@Topic:foo@Topic:bar@Address:addr1@Hash:txhash"},
			"@Topic:foo*@Hash:txhash",
			nil,
			func(t *testing.T) {
				err := db.Delete(ctx, "key1@Topic:foo")
				assert.NoError(t, err)

				err = db.Delete(ctx, "key2@Topic:foo@Topic:bar")
				assert.NoError(t, err)

				err = db.Set(ctx, 10*time.Second, "key2@Topic:foo@Topic:bar@Address:addr1@Hash:txhash", ``)
				assert.NoError(t, err)
			},
		},
		{
			"should get the result as expected",
			[]string{"key2@Topic:foo@Topic:bar@Address:addr1@Hash:txhash"},
			"@Topic:foo*@Address:addr*@Hash:txhash",
			nil,
			func(t *testing.T) {
				err := db.Delete(ctx, "key2@Topic:foo@Topic:bar@Address:addr1@Hash:txhash")
				assert.NoError(t, err)

				err = db.Set(ctx, 10*time.Second, "key2@Topic:foo@Topic:bar@Address:addr1@Hash:txhash", ``)
				assert.NoError(t, err)
			},
		},
		{
			"should get the result as expected",
			[]string{"key2@Topic:foo@Topic:bar@Address:addr1@Hash:txhash"},
			"@Topic:foo*@Address:addr",
			nil,
			func(t *testing.T) {
				err := db.Delete(ctx, "key2@Topic:foo@Topic:bar@Address:addr1@Hash:txhash")
				assert.NoError(t, err)

				err = db.Set(ctx, 10*time.Second, "key2@Topic:foo@Topic:bar@Address:addr1@Hash:txhash", ``)
				assert.NoError(t, err)
			},
		},
		{
			"should return error, bc dont have the key",
			nil,
			"not-set-key",
			redis.Nil,
			func(t *testing.T) {},
		},
	}

	for _, tt := range cases {
		t.Run(tt.purpose, func(t *testing.T) {
			tt.mock(t)

			res, err := db.GetKeys(ctx, tt.mustInclude)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantRes, res)
		})
	}
}

func TestSet(t *testing.T) {
	ctx := context.Background()

	cases := []struct {
		purpose  string
		wantRes  string
		duration time.Duration
		key      string
		wantErr  error
	}{
		{
			"should set the result as expected",
			"foo",
			time.Minute * 1,
			"awesome-key",
			nil,
		},
	}

	for _, tt := range cases {
		t.Run(tt.purpose, func(t *testing.T) {
			err := db.Set(ctx, tt.duration, tt.key, tt.wantRes)
			assert.NoError(t, err)

			res, err := db.Get(ctx, tt.key)
			assert.NoError(t, tt.wantErr, err)
			assert.Equal(t, tt.wantRes, res)
		})
	}
}
