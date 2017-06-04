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

// TestAddProvider checks adding new provider to version struct
func (p *ProviderSuite) TestAddProvider(c *C) {

	p1 := NewMetadataProvider(
		"virtualbox",
		"checksum1",
		"sha256",
		"http://localhost:8080/virtualbox.box",
	)

	p2 := NewMetadataProvider(
		"vmware",
		"checksum2",
		"sha256",
		"http://localhost:8080/vmware.box",
	)

	p3 := NewMetadataProvider(
		"virtualbox",
		"checksum3",
		"sha256",
		"http://localhost:8080/virtualbox.box",
	)

	v1 := NewMetadataVersion("1.0.0", make(VMetadataProvidersList, 0))

	v1.AddProvider(p1)
	v1.AddProvider(p2)

	c.Assert(2, Equals, v1.CountProviders())

	v1.AddProvider(p3)

	c.Assert(2, Equals, v1.CountProviders())
}
