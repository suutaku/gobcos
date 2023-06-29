package client

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suutaku/gobcos/pkg/option"
)

func TestNewClient(t *testing.T) {
	c := NewClient(
		option.WithHost("https://119.29.91.75:17779"),
		option.WithAppId("d74aa1c13b2142fabab5"),
		option.WithSignUser("12SSS36SSS0237af5132184a5a95afdb8c2ace21eb"),
		option.WithAppKey("132CF1AF0CFA217FE9823EAE79EEAA8F758021788417E029D21619AB5A3809BD"),
	)
	require.NotNil(t, c)
}
