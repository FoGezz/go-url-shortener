package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinksMapping(t *testing.T) {
	container := NewLinksMapping()
	fullURL, shortURL := "fullUrl", "sh"

	//add
	container.AddLink(context.Background(), fullURL, shortURL)

	//success found
	actual, found := container.GetByShort(context.Background(), shortURL)
	assert.Equal(t, fullURL, actual)
	assert.Equal(t, true, found)

	//fail not found
	actual, found = container.GetByShort(context.Background(), "unexistentstr")
	assert.Equal(t, "", actual)
	assert.Equal(t, false, found)

	//success found
	actual, found = container.GetByFull(context.Background(), fullURL)
	assert.Equal(t, shortURL, actual)
	assert.Equal(t, true, found)

	//fail not found
	actual, found = container.GetByFull(context.Background(), "unexistendstr2")
	assert.Equal(t, "", actual)
	assert.Equal(t, false, found)
}
