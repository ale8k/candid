package mgostore_test

import (
	"os"
	"testing"

	qt "github.com/frankban/quicktest"
	aclstore "github.com/juju/aclstore/v2"
	"github.com/juju/mgotest"
	"golang.org/x/net/context"
	errgo "gopkg.in/errgo.v1"

	"github.com/CanonicalLtd/candid/meeting"
	"github.com/CanonicalLtd/candid/store"
	"github.com/CanonicalLtd/candid/store/mgostore"
	"github.com/CanonicalLtd/candid/store/storetest"
)

func TestKeyValueStore(t *testing.T) {
	c := qt.New(t)
	defer c.Done()

	storetest.TestKeyValueStore(c, func(c *qt.C) store.ProviderDataStore {
		return newFixture(c).backend.ProviderDataStore()
	})
}

func TestStore(t *testing.T) {
	c := qt.New(t)
	defer c.Done()

	storetest.TestStore(c, func(c *qt.C) store.Store {
		return newFixture(c).backend.Store()
	})
}

func TestMeetingStore(t *testing.T) {
	c := qt.New(t)
	defer c.Done()

	storetest.TestMeetingStore(c, func(c *qt.C) meeting.Store {
		return newFixture(c).backend.MeetingStore()
	}, mgostore.PutAtTime)
}

func TestACLStore(t *testing.T) {
	c := qt.New(t)
	defer c.Done()

	storetest.TestACLStore(c, func(c *qt.C) aclstore.ACLStore {
		return newFixture(c).backend.ACLStore()
	})
}

func TestRootKeyStore(t *testing.T) {
	c := qt.New(t)
	defer c.Done()

	f := newFixture(c)

	rks := f.backend.BakeryRootKeyStore()

	ctx := context.Background()

	key, id, err := rks.RootKey(ctx)
	c.Assert(err, qt.Equals, nil)

	key2, err := rks.Get(ctx, id)
	c.Assert(err, qt.Equals, nil)

	c.Assert(key2, qt.DeepEquals, key)
}

type fixture struct {
	backend store.Backend
	db      *mgotest.Database
	connStr string
}

func newFixture(c *qt.C) *fixture {
	db, err := mgotest.New()
	if errgo.Cause(err) == mgotest.ErrDisabled {
		c.Skip("mgotest disabled")
	}
	c.Assert(err, qt.Equals, nil)
	backend, err := mgostore.NewBackend(db.Database)
	if err != nil {
		db.Close()
		c.Fatal(err)
	}
	c.Assert(err, qt.Equals, nil)
	c.Defer(backend.Close)

	connStr := os.Getenv("MGOCONNECTIONSTRING")
	if connStr == "" {
		connStr = "localhost"
	}
	return &fixture{
		db:      db,
		backend: backend,
		connStr: connStr,
	}
}
