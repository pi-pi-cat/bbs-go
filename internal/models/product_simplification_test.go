package models

import (
	"reflect"
	"testing"
)

func TestModelsDoesNotRegisterRemovedCommunityGrowthTables(t *testing.T) {
	removed := map[string]struct{}{
		"Favorite":      {},
		"TaskConfig":    {},
		"UserTaskEvent": {},
		"UserTaskLog":   {},
		"Badge":         {},
		"UserBadge":     {},
		"LevelConfig":   {},
		"UserExpLog":    {},
		"CheckIn":       {},
		"UserFollow":    {},
		"UserFeed":      {},
		"Link":          {},
		"UserScoreLog":  {},
	}

	for _, model := range Models {
		name := reflect.Indirect(reflect.ValueOf(model)).Type().Name()
		if _, ok := removed[name]; ok {
			t.Fatalf("Models should not register removed community growth model %s", name)
		}
	}
}
