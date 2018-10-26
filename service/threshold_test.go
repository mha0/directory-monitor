package service

import (
	"github.com/mha0/directory-monitor/domain"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestShouldNotify(t *testing.T) {
	cfg := domain.DirectoryMonitorConfig{HeartbeatThresholdInHours: 168, DeadbeatThresholdInHours: 72}

	var fixture = []struct {
		lastRunStatus        domain.Status
		thisRunStatus        domain.Status
		lastNotificationTime time.Time
		shouldNotify         bool
		message              string
	}{
		{thisRunStatus: domain.INITIALIZED, shouldNotify: true, message: "Just INITIALIZED from scratch"},
		{lastRunStatus: domain.INITIALIZED, thisRunStatus: domain.OPERATIONAL, shouldNotify: true, message: "Turned OPERATIONAL (from INITIALIZED)"},
		{lastRunStatus: domain.INITIALIZED, thisRunStatus: domain.WARNING, shouldNotify: true, message: "Turned WARNING (from INITIALIZED)"},

		{lastRunStatus: domain.OPERATIONAL, thisRunStatus: domain.WARNING, shouldNotify: true, message: "Turned WARNING (from OPERATIONAL)"},
		{lastRunStatus: domain.OPERATIONAL, thisRunStatus: domain.OPERATIONAL, lastNotificationTime: time.Now().AddDate(0, 0, -5), shouldNotify: false},
		{lastRunStatus: domain.OPERATIONAL, thisRunStatus: domain.OPERATIONAL, lastNotificationTime: time.Now().AddDate(0, 0, -10), shouldNotify: true, message: "Heartbeat (still OPERATIONAL)"},

		{lastRunStatus: domain.WARNING, thisRunStatus: domain.OPERATIONAL, shouldNotify: true, message: "Turned OPERATIONAL (from WARNING)"},
		{lastRunStatus: domain.WARNING, thisRunStatus: domain.WARNING, lastNotificationTime: time.Now().AddDate(0, 0, -2), shouldNotify: false},
		{lastRunStatus: domain.WARNING, thisRunStatus: domain.WARNING, lastNotificationTime: time.Now().AddDate(0, 0, -5), shouldNotify: true, message: "Deadbeat (still WARNING)"},
	}

	for _, fix := range fixture {
		t.Logf("Given lastRunStatus %v, thisRunStatus %v, lastNotificationTime %v", fix.lastRunStatus, fix.thisRunStatus, fix.lastNotificationTime)
		shouldNotify, message := ShouldNotify(fix.lastRunStatus, fix.thisRunStatus, fix.lastNotificationTime, cfg)
		assert.Equal(t, fix.shouldNotify, shouldNotify, "should be equal")
		assert.Equal(t, fix.message, message, "should be equal")
		t.Logf("\tShould return the expected shouldNotify result %v", shouldNotify)
	}
}
