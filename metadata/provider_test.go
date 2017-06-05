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

// NewVersion returns prepared VMetadataVersion struct
func (p *ProviderSuite) NewVersion() (*VMetadataVersion) {
	return NewMetadataVersion("1.0.0", VMetadataProvidersList{
		NewMetadataProvider(
			"virtualbox",
			"ea8f7fabe99ccaf221f5d61af5bb51cd96fbd675",
			"sha256",
			"http://localhost:8080/virtualbox.box",
		),
		NewMetadataProvider(
			"docker",
			"d25de0ec16c8def969d926c1f56977a8085fdaac",
			"sha256",
			"http://localhost:8080/docker.box",
		),
		NewMetadataProvider(
			"hyper-v",
			"40ef470f82420505ae1cdb6982f2e49cd4a1e83d",
			"sha256",
			"http://localhost:8080/hyper-v.box",
		),
	})
}

// TestFilterProvides checks FilterProvider function
func (p *ProviderSuite) TestFilterProviders(c *C) {
	version := p.NewVersion()

	custom_provider := NewMetadataProvider(
		"custom",
		"9288216bccd3c4399b56b0978d3db7fbfc35376b",
		"sha256",
		"http://localhost:8080/custom.box",
	)

	version.FilterProvider(custom_provider, isEqualProviders)
	c.Assert(0, Equals, version.CountProviders())

	version.AddProvider(custom_provider)
	version.FilterProvider(custom_provider, isEqualProviders)
	c.Assert(1, Equals, version.CountProviders())
}

// TestAddProvider checks adding new provider to version struct
func (p *ProviderSuite) TestAddProvider(c *C) {

	version := p.NewVersion()

	custom_provider := NewMetadataProvider(
		"custom",
		"9288216bccd3c4399b56b0978d3db7fbfc35376b",
		"sha256",
		"http://localhost:8080/custom.box",
	)

	count := version.CountProviders()

	version.AddProvider(custom_provider)
	version.AddProvider(custom_provider)

	c.Assert(count + 1, Equals, version.CountProviders())
}

// TestRemoveProvider checks removing provider from the version struct
func (p *ProviderSuite) TestRemoveProvider(c *C) {
	version := p.NewVersion()

	custom_provider := NewMetadataProvider(
		"custom",
		"9288216bccd3c4399b56b0978d3db7fbfc35376b",
		"sha256",
		"http://localhost:8080/custom.box",
	)

	count := version.CountProviders()

	version.AddProvider(custom_provider)
	version.RemoveProvider(custom_provider)

	c.Assert(count, Equals, version.CountProviders())
}
