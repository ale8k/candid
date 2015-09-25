// Copyright 2015 Canonical Ltd.

package idp_test

import (
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"
	"gopkg.in/yaml.v2"

	"github.com/CanonicalLtd/blues-identity/idp"
)

type idpSuite struct{}

var _ = gc.Suite(&idpSuite{})

var identityProviderUnmarshalYAMLTests = []struct {
	about       string
	data        string
	expectValue idp.IdentityProvider
	expectError string
}{{
	about:       "Ubuntu SSO",
	data:        "type: usso",
	expectValue: idp.UbuntuSSOIdentityProvider,
}, {
	about:       "Ubuntu SSO OAuth",
	data:        "type: usso_oauth",
	expectValue: idp.UbuntuSSOOAuthIdentityProvider,
}, {
	about:       "agent",
	data:        "type: agent",
	expectValue: idp.AgentIdentityProvider,
}, {
	about:       "bad type",
	data:        "type: no-such-type",
	expectError: `unrecognised identity provider type "no-such-type"`,
}, {
	about: "keystone",
	data: `type: keystone
name: ks1
domain: openstack
description: Keystone Login
url: https://example.com/keystone`,
	expectValue: idp.KeystoneIdentityProvider(&idp.KeystoneParams{
		Name:        "ks1",
		Domain:      "openstack",
		Description: "Keystone Login",
		URL:         "https://example.com/keystone",
	}),
}, {
	about: "keystone no name",
	data: `type: keystone
domain: openstack
description: Keystone Login
url: https://example.com/keystone`,
	expectError: "cannot unmarshal keystone configuration: name not specified",
}, {
	about: "keystone no url",
	data: `type: keystone
name: ks1
domain: openstack
description: Keystone Login`,
	expectError: "cannot unmarshal keystone configuration: url not specified",
}, {
	about: "keystone_userpass",
	data: `type: keystone_userpass
name: ks1
domain: openstack
description: Keystone Userpass Login
url: https://example.com/keystone`,
	expectValue: idp.KeystoneUserpassIdentityProvider(&idp.KeystoneParams{
		Name:        "ks1",
		Domain:      "openstack",
		Description: "Keystone Userpass Login",
		URL:         "https://example.com/keystone",
	}),
}, {
	about: "keystone_token",
	data: `type: keystone_token
name: ks1
domain: openstack
description: Keystone Token Login
url: https://example.com/keystone`,
	expectValue: idp.KeystoneTokenIdentityProvider(&idp.KeystoneParams{
		Name:        "ks1",
		Domain:      "openstack",
		Description: "Keystone Token Login",
		URL:         "https://example.com/keystone",
	}),
}}

func (s *idpSuite) TestIdentityProviderUnmarshalYAML(c *gc.C) {
	for i, test := range identityProviderUnmarshalYAMLTests {
		c.Logf("%d %s", i, test.about)
		var v idp.IdentityProvider
		err := yaml.Unmarshal([]byte(test.data), &v)
		if test.expectError != "" {
			c.Assert(err, gc.ErrorMatches, test.expectError)
			continue
		}
		c.Assert(err, gc.IsNil)
		c.Assert(v, jc.DeepEquals, test.expectValue)
	}
}
