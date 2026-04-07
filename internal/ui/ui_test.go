package ui

import (
	"fmt"
	"testing"

	"charm.land/huh/v2"
	"github.com/matryer/is"
)

func TestSelectProfileNoProfiles(t *testing.T) {
	is := is.New(t)

	_, err := SelectProfile(nil, nil, nil)

	is.True(err != nil)
	is.Equal(err.Error(), "There are no available profiles")
}

func TestIsAborted(t *testing.T) {
	is := is.New(t)

	is.True(IsAborted(huh.ErrUserAborted))
	is.True(IsAborted(fmt.Errorf("wrapped: %w", huh.ErrUserAborted)))
	is.True(!IsAborted(fmt.Errorf("other error")))
}
