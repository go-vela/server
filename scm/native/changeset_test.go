// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

// TODO: Add support in Jenkins-X for changes test data
// https://github.com/jenkins-x/go-scm/blob/main/scm/driver/fake/git.go#L58-L60
// func TestNative_Changeset(t *testing.T) {
// 	// setup types
// 	want := []string{"file1.txt"}

// 	u := new(library.User)
// 	u.SetName("foo")
// 	u.SetToken("bar")

// 	r := new(library.Repo)
// 	r.SetOrg("repos")
// 	r.SetName("octocat")

// 	client, _ := NewTest("fake.com")

// 	// run test
// 	got, err := client.Changeset(u, r, "6dcb09b5b57875f334f61aebed695e2e4193db5e")

// 	if err != nil {
// 		t.Errorf("Changeset returned err: %v", err)
// 	}

// 	if !reflect.DeepEqual(got, want) {
// 		t.Errorf("Changeset is %v, want %v", got, want)
// 	}
// }
