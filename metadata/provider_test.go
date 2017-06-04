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

// TestFilterProvides checks FilterProvider function
func (p *ProviderSuite) TestFilterProvider(c *C) {

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
		"custom",
		"checksum3",
		"sha256",
		"http://localhost:8080/custom.box",
	)

	v1 := NewMetadataVersion("1.0.0", make(VMetadataProvidersList, 0))

	v1.AddProvider(p1)
	v1.AddProvider(p2)
	v1.AddProvider(p3)

	v1.FilterProvider(p2, isEqualProviders)

	v2 := NewMetadataVersion("1.0.0", make(VMetadataProvidersList, 0))

	v2.AddProvider(p1)
	v2.AddProvider(p3)

	v2.FilterProvider(p2, isEqualProviders)

	c.Assert(1, Equals, v1.CountProviders())
	c.Assert(0, Equals, v2.CountProviders())
}

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

// TestRemoveProvider checks removing provider from the version struct
func (p *ProviderSuite) TestRemoveProvider(c *C) {

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

	v1 := NewMetadataVersion("1.0.0", make(VMetadataProvidersList, 0))

	v1.AddProvider(p1)
	v1.AddProvider(p2)

	v1.RemoveProvider(p1)
	v1.RemoveProvider(p2)

	c.Assert(0, Equals, v1.CountProviders())
}
