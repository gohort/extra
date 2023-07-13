package extra

import (
	"encoding/json"
	"testing"

	old "github.com/gohort/extra"
)

var testUserOld = &UserOld{
	Username: "username",
	Topics: &UserTopicsList{
		Topics: []Topic{
			{
				Id: 12345,
				Slug: "salt1",
			},
			{
				Id: 4311,
				Slug: "salt2",
			},
			{
				Id: 234572,
				Slug: "salt3",
			},
			{
				Id: 45923546,
				Slug: "salt4",
			},
			{
				Id: 961236,
				Slug: "salt5",
			},
			{
				Id: 384596,
				Slug: "salt6",
			},
		},
		MoreTopicsURL: "morestuffhere!",
	},
	X:old.Any{
		"extra": "thing",
		"here": 123,
	},
}

var testUser = &User{
	Username: "username",
	Topics: &UserTopicsList{
		Topics: []Topic{
			{
				Id: 12345,
				Slug: "salt1",
			},
			{
				Id: 4311,
				Slug: "salt2",
			},
			{
				Id: 234572,
				Slug: "salt3",
			},
			{
				Id: 45923546,
				Slug: "salt4",
			},
			{
				Id: 961236,
				Slug: "salt5",
			},
			{
				Id: 384596,
				Slug: "salt6",
			},
		},
		MoreTopicsURL: "morestuffhere!",
	},
	X: Map{
		"extra": "thing",
		"here": 123,
	},
}

func BenchmarkMarshalWithMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		data, err := MarshalWithMap(testUser, testUser.X)
		_ = data
		_ = err
	}
}

func BenchmarkMarshalOld(b *testing.B) {
	for i := 0; i < b.N; i++ {
		data, err := old.Marshal(testUserOld, &testUserOld.X)
		_ = data
		_ = err
	}
}

func BenchmarkMarshalVanilla(b *testing.B) {
	for i := 0; i < b.N; i++ {
		data, err := json.Marshal(testUser)
		_ = data
		_ = err
	}
}