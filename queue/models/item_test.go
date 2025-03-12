// SPDX-License-Identifier: Apache-2.0

package models

import (
	"reflect"
	"testing"

	api "github.com/go-vela/server/api/types"
)

func TestTypes_ToItem(t *testing.T) {
	// setup types
	booL := false
	num16 := int32(1)
	num64 := int64(1)
	str := "foo"
	e := new(api.Events)

	b := &api.Build{
		ID: &num64,
		Repo: &api.Repo{
			ID: &num64,
			Owner: &api.User{
				ID:     &num64,
				Name:   &str,
				Token:  &str,
				Active: &booL,
				Admin:  &booL,
			},
			Org:         &str,
			Name:        &str,
			FullName:    &str,
			Link:        &str,
			Clone:       &str,
			Branch:      &str,
			Timeout:     &num16,
			Visibility:  &str,
			Private:     &booL,
			Trusted:     &booL,
			Active:      &booL,
			AllowEvents: e,
		},
		Number:   &num64,
		Parent:   &num64,
		Event:    &str,
		Status:   &str,
		Error:    &str,
		Enqueued: &num64,
		Created:  &num64,
		Started:  &num64,
		Finished: &num64,
		Deploy:   &str,
		Clone:    &str,
		Source:   &str,
		Title:    &str,
		Message:  &str,
		Commit:   &str,
		Sender:   &str,
		Author:   &str,
		Branch:   &str,
		Ref:      &str,
		BaseRef:  &str,
	}
	want := &Item{
		Build: &api.Build{
			ID: &num64,
			Repo: &api.Repo{
				ID: &num64,
				Owner: &api.User{
					ID:     &num64,
					Name:   &str,
					Token:  &str,
					Active: &booL,
					Admin:  &booL,
				},
				Org:         &str,
				Name:        &str,
				FullName:    &str,
				Link:        &str,
				Clone:       &str,
				Branch:      &str,
				Timeout:     &num16,
				Visibility:  &str,
				Private:     &booL,
				Trusted:     &booL,
				Active:      &booL,
				AllowEvents: e,
			},
			Number:   &num64,
			Parent:   &num64,
			Event:    &str,
			Status:   &str,
			Error:    &str,
			Enqueued: &num64,
			Created:  &num64,
			Started:  &num64,
			Finished: &num64,
			Deploy:   &str,
			Clone:    &str,
			Source:   &str,
			Title:    &str,
			Message:  &str,
			Commit:   &str,
			Sender:   &str,
			Author:   &str,
			Branch:   &str,
			Ref:      &str,
			BaseRef:  &str,
		},
		ItemVersion: ItemVersion,
	}

	// run test
	got := ToItem(b)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ToItem is %v, want %v", got, want)
	}
}
