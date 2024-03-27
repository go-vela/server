// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"reflect"
	"testing"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/types/library"
)

func TestTypes_ToItem(t *testing.T) {
	// setup types
	booL := false
	num := 1
	num64 := int64(num)
	str := "foo"
	e := new(api.Events)

	b := &library.Build{
		ID:       &num64,
		RepoID:   &num64,
		Number:   &num,
		Parent:   &num,
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
	r := &api.Repo{
		ID: &num64,
		Owner: &library.User{
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
		Timeout:     &num64,
		Visibility:  &str,
		Private:     &booL,
		Trusted:     &booL,
		Active:      &booL,
		AllowPull:   &booL,
		AllowPush:   &booL,
		AllowDeploy: &booL,
		AllowTag:    &booL,
		AllowEvents: e,
	}
	want := &Item{
		Build: &library.Build{
			ID:       &num64,
			RepoID:   &num64,
			Number:   &num,
			Parent:   &num,
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
		Repo: &api.Repo{
			ID: &num64,
			Owner: &library.User{
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
			Timeout:     &num64,
			Visibility:  &str,
			Private:     &booL,
			Trusted:     &booL,
			Active:      &booL,
			AllowPull:   &booL,
			AllowPush:   &booL,
			AllowDeploy: &booL,
			AllowTag:    &booL,
			AllowEvents: e,
		},
		ItemVersion: ItemVersion,
	}

	// run test
	got := ToItem(b, r)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ToItem is %v, want %v", got, want)
	}
}
