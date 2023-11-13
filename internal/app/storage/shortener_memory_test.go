package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinksMapping(t *testing.T) {
	container := NewLinksMapping()
	fullURL, shortURL := "fullUrl", "sh"

	//add
	container.AddLink(fullURL, shortURL)

	//success found
	actual, found := container.GetByShort(shortURL)
	assert.Equal(t, fullURL, actual)
	assert.Equal(t, true, found)

	//fail not found
	actual, found = container.GetByShort("unexistentstr")
	assert.Equal(t, "", actual)
	assert.Equal(t, false, found)

	//success found
	actual, found = container.GetByFull(fullURL)
	assert.Equal(t, shortURL, actual)
	assert.Equal(t, true, found)

	//fail not found
	actual, found = container.GetByFull("unexistendstr2")
	assert.Equal(t, "", actual)
	assert.Equal(t, false, found)
}
