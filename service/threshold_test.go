package service

import (
	"github.com/mha0/directory-monitor/domain"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestShouldNotify(t *testing.T) {
	cfg  := domain.DirectoryMonitorConfig{HeartbeatThresholdInHours: 168, DeadbeatThresholdInHours: 72}

	var fixture = []struct {
		lastRunStatus       domain.Status
		thisRunStatus       domain.Status
		lastTransitionTime  time.Time
		notificationCounter int
		shouldNotify        bool
		message             string
	}{
		{thisRunStatus: domain.INITIALIZED, shouldNotify: true, message: "Just INITIALIZED from scratch"},
		{lastRunStatus: domain.INITIALIZED, thisRunStatus: domain.OPERATIONAL, shouldNotify: true, message: "Turned OPERATIONAL (from INITIALIZED)"},
		{lastRunStatus: domain.INITIALIZED, thisRunStatus: domain.WARNING, shouldNotify: true, message: "Turned WARNING (from INITIALIZED)"},

		{lastRunStatus: domain.OPERATIONAL, thisRunStatus: domain.WARNING, shouldNotify: true, message: "Turned WARNING (from OPERATIONAL)"},
		{lastRunStatus: domain.OPERATIONAL, thisRunStatus: domain.OPERATIONAL, lastTransitionTime: time.Now().AddDate(0, 0, -15), notificationCounter: 2, shouldNotify: false},
		{lastRunStatus: domain.OPERATIONAL, thisRunStatus: domain.OPERATIONAL, lastTransitionTime: time.Now().AddDate(0, 0, -15), notificationCounter: 1, shouldNotify: true, message: "Heartbeat (still OPERATIONAL)"},

		{lastRunStatus: domain.WARNING, thisRunStatus: domain.OPERATIONAL, shouldNotify: true, message: "Turned OPERATIONAL (from WARNING)"},
		{lastRunStatus: domain.WARNING, thisRunStatus: domain.WARNING, lastTransitionTime: time.Now().AddDate(0, 0, -7), notificationCounter: 2, shouldNotify: false},
		{lastRunStatus: domain.WARNING, thisRunStatus: domain.WARNING, lastTransitionTime: time.Now().AddDate(0, 0, -7), notificationCounter: 1, shouldNotify: true, message: "Deadbeat (still WARNING)"},
	}

	for _, fix := range fixture {
		t.Logf("Given lastRunStatus %v, thisRunStatus %v, lastTransitionTime %v, notificationCounter %v", fix.lastRunStatus, fix.thisRunStatus, fix.lastTransitionTime, fix.notificationCounter)
		shouldNotify, message := ShouldNotify(fix.lastRunStatus, fix.thisRunStatus, fix.lastTransitionTime, fix.notificationCounter, &cfg)
		assert.Equal(t, fix.shouldNotify, shouldNotify, "should be equal")
		assert.Equal(t, fix.message, message, "should be equal")
		t.Logf("\tShould return the expected shouldNotify result %v", shouldNotify)
	}
}
