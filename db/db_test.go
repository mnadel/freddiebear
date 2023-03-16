package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkQueryText(t *testing.B) {
	bearDB, err := NewDB()
	assert.Nil(t, err, "cannot create db")

	for i := 0; i < t.N; i++ {
		_, err = bearDB.QueryText("2022")
		assert.Nil(t, err, "error searching text")
	}
}

func BenchmarkQueryTitlesExact(t *testing.B) {
	bearDB, err := NewDB()
	assert.Nil(t, err, "cannot create db")

	for i := 0; i < t.N; i++ {
		_, err = bearDB.QueryTitles("2022", true)
		assert.Nil(t, err, "error searching titles exact")
	}
}

func BenchmarkQueryTitlesFuzzy(t *testing.B) {
	bearDB, err := NewDB()
	assert.Nil(t, err, "cannot create db")

	for i := 0; i < t.N; i++ {
		_, err = bearDB.QueryTitles("2022", false)
		assert.Nil(t, err, "error searching titles fuzzy")
	}
}
