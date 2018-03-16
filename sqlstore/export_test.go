// Copyright 2018 Canonical Ltd.

package sqlstore

import (
	"time"

	"golang.org/x/net/context"

	"github.com/CanonicalLtd/blues-identity/meeting"
)

var PutAtTime = func(ctx context.Context, s meeting.Store, id, address string, now time.Time) error {
	return s.(*meetingStore).put(id, address, now)
}