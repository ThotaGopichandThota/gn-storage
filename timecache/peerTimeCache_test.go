package timecache

import (
	"errors"
	"testing"
	"time"

	"github.com/ThotaGopichandThota/gn-core3/core"
	"github.com/ThotaGopichandThota/gn-core3/core/check"
	"github.com/ThotaGopichandThota/gn-storage/common"
	"github.com/ThotaGopichandThota/gn-storage/testscommon"
	"github.com/stretchr/testify/assert"
)

func TestNewPeerTimeCache_NilTimeCacheShouldErr(t *testing.T) {
	t.Parallel()

	ptc, err := NewPeerTimeCache(nil)

	assert.Equal(t, common.ErrNilTimeCache, err)
	assert.True(t, check.IfNil(ptc))
}

func TestNewPeerTimeCache_ShouldWork(t *testing.T) {
	t.Parallel()

	ptc, err := NewPeerTimeCache(&testscommon.TimeCacheStub{})

	assert.Nil(t, err)
	assert.False(t, check.IfNil(ptc))
}

func TestPeerTimeCache_Methods(t *testing.T) {
	t.Parallel()

	pid := core.PeerID("test peer id")
	unexpectedErr := errors.New("unexpected error")
	updateWasCalled := false
	hasWasCalled := false
	sweepWasCalled := false
	ptc, _ := NewPeerTimeCache(&testscommon.TimeCacheStub{
		UpsertCalled: func(key string, span time.Duration) error {
			if key != string(pid) {
				return unexpectedErr
			}

			updateWasCalled = true
			return nil
		},
		HasCalled: func(key string) bool {
			if key != string(pid) {
				return false
			}

			hasWasCalled = true
			return true
		},
		SweepCalled: func() {
			sweepWasCalled = true
		},
	})

	assert.Nil(t, ptc.Upsert(pid, time.Second))
	assert.True(t, ptc.Has(pid))
	ptc.Sweep()

	assert.True(t, updateWasCalled)
	assert.True(t, hasWasCalled)
	assert.True(t, sweepWasCalled)
}
