package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinksContainer(t *testing.T) {
	container := NewLinksContainer()
	fullUrl, shortUrl := "fullUrl", "sh"

	//add
	container.AddLink(fullUrl, shortUrl)

	//success found
	actual, found := container.GetByShort(shortUrl)
	assert.Equal(t, fullUrl, actual)
	assert.Equal(t, true, found)

	//fail not found
	actual, found = container.GetByShort("unexistentstr")
	assert.Equal(t, "", actual)
	assert.Equal(t, false, found)

	//success found
	actual, found = container.GetByFull(fullUrl)
	assert.Equal(t, shortUrl, actual)
	assert.Equal(t, true, found)

	//fail not found
	actual, found = container.GetByFull("unexistendstr2")
	assert.Equal(t, "", actual)
	assert.Equal(t, false, found)
}
