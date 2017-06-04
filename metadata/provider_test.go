package metadata

import (
	"testing"

	. "pkg.re/check.v1"
)

// ////////////////////////////////////////////////////////////////////////////////// //

type ProviderSuite struct{}

func Test(t *testing.T) { TestingT(t) }

var _ = Suite(&ProviderSuite{})

// ////////////////////////////////////////////////////////////////////////////////// //

// TestCount checks counts of providers in the list
func (p *ProviderSuite) TestCountProviders(c *C) {
	p1 := NewMetadataProvider(
		"test",
		"checksum",
		"sha256",
		"virtualbox",
	)

	p2 := NewMetadataProvider(
		"test",
		"checksum",
		"sha256",
		"vmware",
	)

	pl1 := make(VMetadataProvidersList, 0)

	pl1 = append(pl1, p1)
	pl1 = append(pl1, p2)

	v1 := NewMetadataVersion("1.0.0", pl1)

	c.Assert(2, Equals, v1.CountProviders())
}
