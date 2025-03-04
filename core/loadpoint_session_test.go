package core

import (
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/evcc-io/evcc/api"
	coredb "github.com/evcc-io/evcc/core/db"
	"github.com/evcc-io/evcc/mock"
	serverdb "github.com/evcc-io/evcc/server/db"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSession(t *testing.T) {
	var err error
	serverdb.Instance, err = serverdb.New("sqlite", ":memory:")
	assert.NoError(t, err)

	db, err := coredb.New("foo")
	assert.NoError(t, err)

	clock := clock.NewMock()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mm := mock.NewMockMeter(ctrl)
	me := mock.NewMockMeterEnergy(ctrl)

	type EnergyDecorator struct {
		api.Meter
		api.MeterEnergy
	}

	cm := &EnergyDecorator{Meter: mm, MeterEnergy: me}

	lp := &Loadpoint{
		clock:       clock,
		db:          db,
		chargeMeter: cm,
	}

	// create session
	me.EXPECT().TotalEnergy().Return(1.0, nil)
	lp.createSession()
	assert.NotNil(t, lp.session)

	// start charging
	lp.updateSession(func(session *coredb.Session) {
		if session.Created.IsZero() {
			session.Created = lp.clock.Now()
		}
	})
	assert.Equal(t, clock.Now(), lp.session.Created)

	// stop charging
	clock.Add(time.Hour)
	lp.chargedEnergy = 1.23 * 1e3                                   // Wh
	me.EXPECT().TotalEnergy().Return(1.0+lp.chargedEnergy/1e3, nil) // match chargedEnergy

	lp.stopSession()
	assert.NotNil(t, lp.session)
	assert.Equal(t, lp.chargedEnergy/1e3, lp.session.ChargedEnergy)
	assert.Equal(t, clock.Now(), lp.session.Finished)

	s, err := db.Sessions()
	assert.NoError(t, err)
	assert.Len(t, s, 1)
	t.Log(s)

	// stop charging - 2nd leg
	clock.Add(time.Hour)
	lp.chargedEnergy *= 2
	me.EXPECT().TotalEnergy().Return(3.0, nil) // doesn't match chargedEnergy

	lp.stopSession()
	assert.NotNil(t, lp.session)
	assert.Equal(t, 3.0-1.0, lp.session.ChargedEnergy) // expect actual meter energy delta
	assert.Equal(t, clock.Now(), lp.session.Finished)

	s, err = db.Sessions()
	assert.NoError(t, err)
	assert.Len(t, s, 1)
	t.Log(s)
}
