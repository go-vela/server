// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

// TODO: investigate ways to test this
//
// Right now the "fake" client does not support webhooks:
// https://github.com/jenkins-x/go-scm/tree/main/scm/driver/fake
// I could get around this by testing a specific client type i.e. github, bitbucket etc
// However, testing it that way would take some refactoring to how my test client works.
// func Test_client_ProcessWebhook(t *testing.T) {
// 	client, _ := NewTest("fake.com")

// 	type args struct {
// 		request *http.Request
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    *types.Webhook
// 		wantErr bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			got, err := client.ProcessWebhook(test.args.request)
// 			if (err != nil) != test.wantErr {
// 				t.Errorf("client.ProcessWebhook() error = %v, wantErr %v", err, test.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(got, test.want) {
// 				t.Errorf("client.ProcessWebhook() = %v, want %v", got, test.want)
// 			}
// 		})
// 	}
// }
