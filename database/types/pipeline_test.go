// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"
	"reflect"
	"testing"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
)

func TestDatabase_Pipeline_Compress(t *testing.T) {
	// setup tests
	tests := []struct {
		name     string
		failure  bool
		level    int
		pipeline *Pipeline
		want     []byte
	}{
		{
			name:     "compression level -1",
			failure:  false,
			level:    constants.CompressionNegOne,
			pipeline: &Pipeline{Data: testPipelineData()},
			want:     []byte{120, 156, 100, 203, 65, 14, 3, 33, 8, 133, 225, 61, 167, 120, 23, 112, 209, 173, 183, 33, 14, 99, 201, 168, 24, 156, 216, 235, 55, 214, 164, 73, 83, 118, 252, 124, 208, 20, 31, 106, 45, 226, 65, 244, 50, 191, 196, 35, 1, 103, 225, 105, 30, 81, 216, 179, 16, 208, 11, 223, 167, 121, 141, 56, 44, 93, 226, 68, 67, 124, 106, 146, 177, 116, 64, 227, 42, 17, 46, 135, 14, 2, 0, 173, 156, 191, 129, 198, 45, 253, 7, 118, 109, 249, 223, 173, 144, 172, 86, 110, 199, 71, 175, 9, 251, 22, 82, 81, 132, 231, 94, 246, 251, 59, 0, 0, 255, 255, 33, 108, 56, 191},
		},
		{
			name:     "compression level 0",
			failure:  false,
			level:    constants.CompressionZero,
			pipeline: &Pipeline{Data: testPipelineData()},
			want:     []byte{120, 1, 0, 186, 0, 69, 255, 10, 118, 101, 114, 115, 105, 111, 110, 58, 32, 49, 10, 10, 119, 111, 114, 107, 101, 114, 58, 10, 32, 32, 102, 108, 97, 118, 111, 114, 58, 32, 108, 97, 114, 103, 101, 10, 32, 32, 112, 108, 97, 116, 102, 111, 114, 109, 58, 32, 100, 111, 99, 107, 101, 114, 10, 10, 115, 101, 114, 118, 105, 99, 101, 115, 58, 10, 32, 32, 45, 32, 110, 97, 109, 101, 58, 32, 114, 101, 100, 105, 115, 10, 32, 32, 32, 32, 105, 109, 97, 103, 101, 58, 32, 114, 101, 100, 105, 115, 10, 10, 115, 116, 101, 112, 115, 58, 10, 32, 32, 45, 32, 110, 97, 109, 101, 58, 32, 112, 105, 110, 103, 10, 32, 32, 32, 32, 105, 109, 97, 103, 101, 58, 32, 114, 101, 100, 105, 115, 10, 32, 32, 32, 32, 99, 111, 109, 109, 97, 110, 100, 115, 58, 10, 32, 32, 32, 32, 32, 32, 45, 32, 114, 101, 100, 105, 115, 45, 99, 108, 105, 32, 45, 104, 32, 114, 101, 100, 105, 115, 32, 112, 105, 110, 103, 10, 1, 0, 0, 255, 255, 33, 108, 56, 191},
		},
		{
			name:     "compression level 1",
			failure:  false,
			level:    constants.CompressionOne,
			pipeline: &Pipeline{Data: testPipelineData()},
			want:     []byte{120, 1, 100, 204, 189, 173, 195, 48, 12, 69, 225, 158, 83, 220, 5, 84, 188, 150, 219, 8, 210, 181, 30, 97, 253, 129, 52, 156, 245, 131, 36, 85, 144, 83, 127, 56, 114, 211, 195, 214, 84, 252, 137, 60, 150, 159, 116, 21, 224, 232, 249, 94, 174, 232, 217, 27, 5, 216, 61, 95, 199, 242, 161, 168, 171, 156, 116, 145, 160, 223, 86, 24, 42, 64, 194, 204, 131, 10, 103, 181, 16, 0, 176, 145, 27, 21, 206, 106, 33, 18, 23, 247, 23, 220, 54, 219, 175, 3, 128, 178, 198, 200, 179, 190, 245, 171, 244, 153, 166, 210, 13, 233, 31, 206, 106, 129, 109, 179, 201, 51, 0, 0, 255, 255, 33, 108, 56, 191},
		},
		{
			name:     "compression level 2",
			failure:  false,
			level:    constants.CompressionTwo,
			pipeline: &Pipeline{Data: testPipelineData()},
			want:     []byte{120, 94, 100, 204, 65, 170, 3, 33, 12, 198, 241, 125, 78, 241, 93, 192, 197, 219, 122, 155, 160, 25, 95, 24, 53, 18, 7, 123, 253, 226, 12, 20, 74, 179, 75, 248, 253, 67, 75, 124, 170, 245, 136, 63, 162, 151, 249, 41, 30, 9, 56, 42, 47, 243, 136, 202, 94, 132, 128, 81, 249, 58, 204, 91, 68, 182, 116, 138, 19, 77, 241, 165, 73, 230, 214, 1, 157, 155, 68, 184, 100, 157, 4, 0, 218, 184, 124, 14, 52, 47, 25, 95, 112, 104, 47, 191, 110, 135, 201, 90, 227, 158, 111, 189, 39, 60, 79, 67, 170, 138, 240, 255, 44, 184, 243, 119, 0, 0, 0, 255, 255, 33, 108, 56, 191},
		},
		{
			name:     "compression level 3",
			failure:  false,
			level:    constants.CompressionThree,
			pipeline: &Pipeline{Data: testPipelineData()},
			want:     []byte{120, 94, 100, 204, 65, 170, 3, 33, 12, 198, 241, 125, 78, 241, 93, 192, 197, 219, 122, 155, 160, 25, 95, 24, 53, 18, 7, 123, 253, 226, 12, 20, 74, 179, 75, 248, 253, 67, 75, 124, 170, 245, 136, 63, 162, 151, 249, 41, 30, 9, 56, 42, 47, 243, 136, 202, 94, 132, 128, 81, 249, 58, 204, 91, 68, 182, 116, 138, 19, 77, 241, 165, 73, 230, 214, 1, 157, 155, 68, 184, 100, 157, 4, 0, 218, 184, 124, 14, 52, 47, 25, 95, 112, 104, 47, 191, 110, 135, 201, 90, 227, 158, 111, 189, 39, 60, 79, 67, 170, 138, 240, 255, 44, 184, 243, 119, 0, 0, 0, 255, 255, 33, 108, 56, 191},
		},
		{
			name:     "compression level 4",
			failure:  false,
			level:    constants.CompressionFour,
			pipeline: &Pipeline{Data: testPipelineData()},
			want:     []byte{120, 94, 100, 203, 65, 14, 3, 33, 8, 133, 225, 61, 167, 120, 23, 112, 209, 173, 183, 33, 14, 99, 201, 168, 24, 156, 216, 235, 55, 214, 164, 73, 83, 118, 252, 124, 208, 20, 31, 106, 45, 226, 65, 244, 50, 191, 196, 35, 1, 103, 225, 105, 30, 81, 216, 179, 16, 208, 11, 223, 167, 121, 141, 56, 44, 93, 226, 68, 67, 124, 106, 146, 177, 116, 64, 227, 42, 17, 46, 135, 14, 2, 0, 173, 156, 191, 129, 198, 45, 253, 7, 118, 109, 249, 223, 173, 144, 172, 86, 110, 199, 71, 175, 9, 251, 22, 82, 81, 132, 231, 94, 246, 251, 59, 0, 0, 255, 255, 33, 108, 56, 191},
		},
		{
			name:     "compression level 5",
			failure:  false,
			level:    constants.CompressionFive,
			pipeline: &Pipeline{Data: testPipelineData()},
			want:     []byte{120, 94, 100, 203, 65, 14, 3, 33, 8, 133, 225, 61, 167, 120, 23, 112, 209, 173, 183, 33, 14, 99, 201, 168, 24, 156, 216, 235, 55, 214, 164, 73, 83, 118, 252, 124, 208, 20, 31, 106, 45, 226, 65, 244, 50, 191, 196, 35, 1, 103, 225, 105, 30, 81, 216, 179, 16, 208, 11, 223, 167, 121, 141, 56, 44, 93, 226, 68, 67, 124, 106, 146, 177, 116, 64, 227, 42, 17, 46, 135, 14, 2, 0, 173, 156, 191, 129, 198, 45, 253, 7, 118, 109, 249, 223, 173, 144, 172, 86, 110, 199, 71, 175, 9, 251, 22, 82, 81, 132, 231, 94, 246, 251, 59, 0, 0, 255, 255, 33, 108, 56, 191},
		},
		{
			name:     "compression level 6",
			failure:  false,
			level:    constants.CompressionSix,
			pipeline: &Pipeline{Data: testPipelineData()},
			want:     []byte{120, 156, 100, 203, 65, 14, 3, 33, 8, 133, 225, 61, 167, 120, 23, 112, 209, 173, 183, 33, 14, 99, 201, 168, 24, 156, 216, 235, 55, 214, 164, 73, 83, 118, 252, 124, 208, 20, 31, 106, 45, 226, 65, 244, 50, 191, 196, 35, 1, 103, 225, 105, 30, 81, 216, 179, 16, 208, 11, 223, 167, 121, 141, 56, 44, 93, 226, 68, 67, 124, 106, 146, 177, 116, 64, 227, 42, 17, 46, 135, 14, 2, 0, 173, 156, 191, 129, 198, 45, 253, 7, 118, 109, 249, 223, 173, 144, 172, 86, 110, 199, 71, 175, 9, 251, 22, 82, 81, 132, 231, 94, 246, 251, 59, 0, 0, 255, 255, 33, 108, 56, 191},
		},
		{
			name:     "compression level 7",
			failure:  false,
			level:    constants.CompressionSeven,
			pipeline: &Pipeline{Data: testPipelineData()},
			want:     []byte{120, 218, 100, 203, 65, 14, 3, 33, 8, 133, 225, 61, 167, 120, 23, 112, 209, 173, 183, 33, 14, 99, 201, 168, 24, 156, 216, 235, 55, 214, 164, 73, 83, 118, 252, 124, 208, 20, 31, 106, 45, 226, 65, 244, 50, 191, 196, 35, 1, 103, 225, 105, 30, 81, 216, 179, 16, 208, 11, 223, 167, 121, 141, 56, 44, 93, 226, 68, 67, 124, 106, 146, 177, 116, 64, 227, 42, 17, 46, 135, 14, 2, 0, 173, 156, 191, 129, 198, 45, 253, 7, 118, 109, 249, 223, 173, 144, 172, 86, 110, 199, 71, 175, 9, 251, 22, 82, 81, 132, 231, 94, 246, 251, 59, 0, 0, 255, 255, 33, 108, 56, 191},
		},
		{
			name:     "compression level 8",
			failure:  false,
			level:    constants.CompressionEight,
			pipeline: &Pipeline{Data: testPipelineData()},
			want:     []byte{120, 218, 100, 203, 65, 14, 3, 33, 8, 133, 225, 61, 167, 120, 23, 112, 209, 173, 183, 33, 14, 99, 201, 168, 24, 156, 216, 235, 55, 214, 164, 73, 83, 118, 252, 124, 208, 20, 31, 106, 45, 226, 65, 244, 50, 191, 196, 35, 1, 103, 225, 105, 30, 81, 216, 179, 16, 208, 11, 223, 167, 121, 141, 56, 44, 93, 226, 68, 67, 124, 106, 146, 177, 116, 64, 227, 42, 17, 46, 135, 14, 2, 0, 173, 156, 191, 129, 198, 45, 253, 7, 118, 109, 249, 223, 173, 144, 172, 86, 110, 199, 71, 175, 9, 251, 22, 82, 81, 132, 231, 94, 246, 251, 59, 0, 0, 255, 255, 33, 108, 56, 191},
		},
		{
			name:     "compression level 9",
			failure:  false,
			level:    constants.CompressionNine,
			pipeline: &Pipeline{Data: testPipelineData()},
			want:     []byte{120, 218, 100, 203, 65, 14, 3, 33, 8, 133, 225, 61, 167, 120, 23, 112, 209, 173, 183, 33, 14, 99, 201, 168, 24, 156, 216, 235, 55, 214, 164, 73, 83, 118, 252, 124, 208, 20, 31, 106, 45, 226, 65, 244, 50, 191, 196, 35, 1, 103, 225, 105, 30, 81, 216, 179, 16, 208, 11, 223, 167, 121, 141, 56, 44, 93, 226, 68, 67, 124, 106, 146, 177, 116, 64, 227, 42, 17, 46, 135, 14, 2, 0, 173, 156, 191, 129, 198, 45, 253, 7, 118, 109, 249, 223, 173, 144, 172, 86, 110, 199, 71, 175, 9, 251, 22, 82, 81, 132, 231, 94, 246, 251, 59, 0, 0, 255, 255, 33, 108, 56, 191},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.pipeline.Compress(test.level)

			if test.failure {
				if err == nil {
					t.Errorf("Compress for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("Compress for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(test.pipeline.Data, test.want) {
				t.Errorf("Compress for %s is %v, want %v", test.name, string(test.pipeline.Data), string(test.want))
			}
		})
	}
}

func TestDatabase_Pipeline_Decompress(t *testing.T) {
	// setup tests
	tests := []struct {
		name     string
		failure  bool
		pipeline *Pipeline
		want     []byte
	}{
		{
			name:     "compression level -1",
			failure:  false,
			pipeline: &Pipeline{Data: []byte{120, 156, 100, 203, 65, 14, 3, 33, 8, 133, 225, 61, 167, 120, 23, 112, 209, 173, 183, 33, 14, 99, 201, 168, 24, 156, 216, 235, 55, 214, 164, 73, 83, 118, 252, 124, 208, 20, 31, 106, 45, 226, 65, 244, 50, 191, 196, 35, 1, 103, 225, 105, 30, 81, 216, 179, 16, 208, 11, 223, 167, 121, 141, 56, 44, 93, 226, 68, 67, 124, 106, 146, 177, 116, 64, 227, 42, 17, 46, 135, 14, 2, 0, 173, 156, 191, 129, 198, 45, 253, 7, 118, 109, 249, 223, 173, 144, 172, 86, 110, 199, 71, 175, 9, 251, 22, 82, 81, 132, 231, 94, 246, 251, 59, 0, 0, 255, 255, 33, 108, 56, 191}},
			want:     testPipelineData(),
		},
		{
			name:     "compression level 0",
			failure:  false,
			pipeline: &Pipeline{Data: []byte{120, 1, 0, 186, 0, 69, 255, 10, 118, 101, 114, 115, 105, 111, 110, 58, 32, 49, 10, 10, 119, 111, 114, 107, 101, 114, 58, 10, 32, 32, 102, 108, 97, 118, 111, 114, 58, 32, 108, 97, 114, 103, 101, 10, 32, 32, 112, 108, 97, 116, 102, 111, 114, 109, 58, 32, 100, 111, 99, 107, 101, 114, 10, 10, 115, 101, 114, 118, 105, 99, 101, 115, 58, 10, 32, 32, 45, 32, 110, 97, 109, 101, 58, 32, 114, 101, 100, 105, 115, 10, 32, 32, 32, 32, 105, 109, 97, 103, 101, 58, 32, 114, 101, 100, 105, 115, 10, 10, 115, 116, 101, 112, 115, 58, 10, 32, 32, 45, 32, 110, 97, 109, 101, 58, 32, 112, 105, 110, 103, 10, 32, 32, 32, 32, 105, 109, 97, 103, 101, 58, 32, 114, 101, 100, 105, 115, 10, 32, 32, 32, 32, 99, 111, 109, 109, 97, 110, 100, 115, 58, 10, 32, 32, 32, 32, 32, 32, 45, 32, 114, 101, 100, 105, 115, 45, 99, 108, 105, 32, 45, 104, 32, 114, 101, 100, 105, 115, 32, 112, 105, 110, 103, 10, 1, 0, 0, 255, 255, 33, 108, 56, 191}},
			want:     testPipelineData(),
		},
		{
			name:     "compression level 1",
			failure:  false,
			pipeline: &Pipeline{Data: []byte{120, 1, 100, 204, 189, 173, 195, 48, 12, 69, 225, 158, 83, 220, 5, 84, 188, 150, 219, 8, 210, 181, 30, 97, 253, 129, 52, 156, 245, 131, 36, 85, 144, 83, 127, 56, 114, 211, 195, 214, 84, 252, 137, 60, 150, 159, 116, 21, 224, 232, 249, 94, 174, 232, 217, 27, 5, 216, 61, 95, 199, 242, 161, 168, 171, 156, 116, 145, 160, 223, 86, 24, 42, 64, 194, 204, 131, 10, 103, 181, 16, 0, 176, 145, 27, 21, 206, 106, 33, 18, 23, 247, 23, 220, 54, 219, 175, 3, 128, 178, 198, 200, 179, 190, 245, 171, 244, 153, 166, 210, 13, 233, 31, 206, 106, 129, 109, 179, 201, 51, 0, 0, 255, 255, 33, 108, 56, 191}},
			want:     testPipelineData(),
		},
		{
			name:     "compression level 2",
			failure:  false,
			pipeline: &Pipeline{Data: []byte{120, 94, 100, 204, 65, 170, 3, 33, 12, 198, 241, 125, 78, 241, 93, 192, 197, 219, 122, 155, 160, 25, 95, 24, 53, 18, 7, 123, 253, 226, 12, 20, 74, 179, 75, 248, 253, 67, 75, 124, 170, 245, 136, 63, 162, 151, 249, 41, 30, 9, 56, 42, 47, 243, 136, 202, 94, 132, 128, 81, 249, 58, 204, 91, 68, 182, 116, 138, 19, 77, 241, 165, 73, 230, 214, 1, 157, 155, 68, 184, 100, 157, 4, 0, 218, 184, 124, 14, 52, 47, 25, 95, 112, 104, 47, 191, 110, 135, 201, 90, 227, 158, 111, 189, 39, 60, 79, 67, 170, 138, 240, 255, 44, 184, 243, 119, 0, 0, 0, 255, 255, 33, 108, 56, 191}},
			want:     testPipelineData(),
		},
		{
			name:     "compression level 3",
			failure:  false,
			pipeline: &Pipeline{Data: []byte{120, 94, 100, 204, 65, 170, 3, 33, 12, 198, 241, 125, 78, 241, 93, 192, 197, 219, 122, 155, 160, 25, 95, 24, 53, 18, 7, 123, 253, 226, 12, 20, 74, 179, 75, 248, 253, 67, 75, 124, 170, 245, 136, 63, 162, 151, 249, 41, 30, 9, 56, 42, 47, 243, 136, 202, 94, 132, 128, 81, 249, 58, 204, 91, 68, 182, 116, 138, 19, 77, 241, 165, 73, 230, 214, 1, 157, 155, 68, 184, 100, 157, 4, 0, 218, 184, 124, 14, 52, 47, 25, 95, 112, 104, 47, 191, 110, 135, 201, 90, 227, 158, 111, 189, 39, 60, 79, 67, 170, 138, 240, 255, 44, 184, 243, 119, 0, 0, 0, 255, 255, 33, 108, 56, 191}},
			want:     testPipelineData(),
		},
		{
			name:     "compression level 4",
			failure:  false,
			pipeline: &Pipeline{Data: []byte{120, 94, 100, 203, 65, 14, 3, 33, 8, 133, 225, 61, 167, 120, 23, 112, 209, 173, 183, 33, 14, 99, 201, 168, 24, 156, 216, 235, 55, 214, 164, 73, 83, 118, 252, 124, 208, 20, 31, 106, 45, 226, 65, 244, 50, 191, 196, 35, 1, 103, 225, 105, 30, 81, 216, 179, 16, 208, 11, 223, 167, 121, 141, 56, 44, 93, 226, 68, 67, 124, 106, 146, 177, 116, 64, 227, 42, 17, 46, 135, 14, 2, 0, 173, 156, 191, 129, 198, 45, 253, 7, 118, 109, 249, 223, 173, 144, 172, 86, 110, 199, 71, 175, 9, 251, 22, 82, 81, 132, 231, 94, 246, 251, 59, 0, 0, 255, 255, 33, 108, 56, 191}},
			want:     testPipelineData(),
		},
		{
			name:     "compression level 5",
			failure:  false,
			pipeline: &Pipeline{Data: []byte{120, 94, 100, 203, 65, 14, 3, 33, 8, 133, 225, 61, 167, 120, 23, 112, 209, 173, 183, 33, 14, 99, 201, 168, 24, 156, 216, 235, 55, 214, 164, 73, 83, 118, 252, 124, 208, 20, 31, 106, 45, 226, 65, 244, 50, 191, 196, 35, 1, 103, 225, 105, 30, 81, 216, 179, 16, 208, 11, 223, 167, 121, 141, 56, 44, 93, 226, 68, 67, 124, 106, 146, 177, 116, 64, 227, 42, 17, 46, 135, 14, 2, 0, 173, 156, 191, 129, 198, 45, 253, 7, 118, 109, 249, 223, 173, 144, 172, 86, 110, 199, 71, 175, 9, 251, 22, 82, 81, 132, 231, 94, 246, 251, 59, 0, 0, 255, 255, 33, 108, 56, 191}},
			want:     testPipelineData(),
		},
		{
			name:     "compression level 6",
			failure:  false,
			pipeline: &Pipeline{Data: []byte{120, 156, 100, 203, 65, 14, 3, 33, 8, 133, 225, 61, 167, 120, 23, 112, 209, 173, 183, 33, 14, 99, 201, 168, 24, 156, 216, 235, 55, 214, 164, 73, 83, 118, 252, 124, 208, 20, 31, 106, 45, 226, 65, 244, 50, 191, 196, 35, 1, 103, 225, 105, 30, 81, 216, 179, 16, 208, 11, 223, 167, 121, 141, 56, 44, 93, 226, 68, 67, 124, 106, 146, 177, 116, 64, 227, 42, 17, 46, 135, 14, 2, 0, 173, 156, 191, 129, 198, 45, 253, 7, 118, 109, 249, 223, 173, 144, 172, 86, 110, 199, 71, 175, 9, 251, 22, 82, 81, 132, 231, 94, 246, 251, 59, 0, 0, 255, 255, 33, 108, 56, 191}},
			want:     testPipelineData(),
		},
		{
			name:     "compression level 7",
			failure:  false,
			pipeline: &Pipeline{Data: []byte{120, 218, 100, 203, 65, 14, 3, 33, 8, 133, 225, 61, 167, 120, 23, 112, 209, 173, 183, 33, 14, 99, 201, 168, 24, 156, 216, 235, 55, 214, 164, 73, 83, 118, 252, 124, 208, 20, 31, 106, 45, 226, 65, 244, 50, 191, 196, 35, 1, 103, 225, 105, 30, 81, 216, 179, 16, 208, 11, 223, 167, 121, 141, 56, 44, 93, 226, 68, 67, 124, 106, 146, 177, 116, 64, 227, 42, 17, 46, 135, 14, 2, 0, 173, 156, 191, 129, 198, 45, 253, 7, 118, 109, 249, 223, 173, 144, 172, 86, 110, 199, 71, 175, 9, 251, 22, 82, 81, 132, 231, 94, 246, 251, 59, 0, 0, 255, 255, 33, 108, 56, 191}},
			want:     testPipelineData(),
		},
		{
			name:     "compression level 8",
			failure:  false,
			pipeline: &Pipeline{Data: []byte{120, 218, 100, 203, 65, 14, 3, 33, 8, 133, 225, 61, 167, 120, 23, 112, 209, 173, 183, 33, 14, 99, 201, 168, 24, 156, 216, 235, 55, 214, 164, 73, 83, 118, 252, 124, 208, 20, 31, 106, 45, 226, 65, 244, 50, 191, 196, 35, 1, 103, 225, 105, 30, 81, 216, 179, 16, 208, 11, 223, 167, 121, 141, 56, 44, 93, 226, 68, 67, 124, 106, 146, 177, 116, 64, 227, 42, 17, 46, 135, 14, 2, 0, 173, 156, 191, 129, 198, 45, 253, 7, 118, 109, 249, 223, 173, 144, 172, 86, 110, 199, 71, 175, 9, 251, 22, 82, 81, 132, 231, 94, 246, 251, 59, 0, 0, 255, 255, 33, 108, 56, 191}},
			want:     testPipelineData(),
		},
		{
			name:     "compression level 9",
			failure:  false,
			pipeline: &Pipeline{Data: []byte{120, 218, 100, 203, 65, 14, 3, 33, 8, 133, 225, 61, 167, 120, 23, 112, 209, 173, 183, 33, 14, 99, 201, 168, 24, 156, 216, 235, 55, 214, 164, 73, 83, 118, 252, 124, 208, 20, 31, 106, 45, 226, 65, 244, 50, 191, 196, 35, 1, 103, 225, 105, 30, 81, 216, 179, 16, 208, 11, 223, 167, 121, 141, 56, 44, 93, 226, 68, 67, 124, 106, 146, 177, 116, 64, 227, 42, 17, 46, 135, 14, 2, 0, 173, 156, 191, 129, 198, 45, 253, 7, 118, 109, 249, 223, 173, 144, 172, 86, 110, 199, 71, 175, 9, 251, 22, 82, 81, 132, 231, 94, 246, 251, 59, 0, 0, 255, 255, 33, 108, 56, 191}},
			want:     testPipelineData(),
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.pipeline.Decompress()

			if test.failure {
				if err == nil {
					t.Errorf("Decompress for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("Decompress for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(test.pipeline.Data, test.want) {
				t.Errorf("Decompress for %s is %v, want %v", test.name, string(test.pipeline.Data), string(test.want))
			}
		})
	}
}

func TestDatabase_Pipeline_Nullify(t *testing.T) {
	// setup types
	var p *Pipeline

	want := &Pipeline{
		ID:       sql.NullInt64{Int64: 0, Valid: false},
		RepoID:   sql.NullInt64{Int64: 0, Valid: false},
		Commit:   sql.NullString{String: "", Valid: false},
		Flavor:   sql.NullString{String: "", Valid: false},
		Platform: sql.NullString{String: "", Valid: false},
		Ref:      sql.NullString{String: "", Valid: false},
		Type:     sql.NullString{String: "", Valid: false},
		Version:  sql.NullString{String: "", Valid: false},
	}

	// setup tests
	tests := []struct {
		pipeline *Pipeline
		want     *Pipeline
	}{
		{
			pipeline: testPipeline(),
			want:     testPipeline(),
		},
		{
			pipeline: p,
			want:     nil,
		},
		{
			pipeline: new(Pipeline),
			want:     want,
		},
	}

	// run tests
	for _, test := range tests {
		got := test.pipeline.Nullify()

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Nullify is %v, want %v", got, test.want)
		}
	}
}

func TestDatabase_Pipeline_ToAPI(t *testing.T) {
	// setup types
	want := new(api.Pipeline)

	want.SetID(1)
	want.SetRepo(testRepo().ToAPI())
	want.SetCommit("48afb5bdc41ad69bf22588491333f7cf71135163")
	want.SetFlavor("large")
	want.SetPlatform("docker")
	want.SetRef("refs/heads/main")
	want.SetType(constants.PipelineTypeYAML)
	want.SetVersion("1")
	want.SetExternalSecrets(false)
	want.SetInternalSecrets(false)
	want.SetServices(true)
	want.SetStages(false)
	want.SetSteps(true)
	want.SetTemplates(false)
	want.SetData(testPipelineData())

	// run test
	got := testPipeline().ToAPI()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ToAPI is %v, want %v", got, want)
	}
}

func TestDatabase_Pipeline_Validate(t *testing.T) {
	// setup tests
	tests := []struct {
		failure  bool
		pipeline *Pipeline
	}{
		{
			failure:  false,
			pipeline: testPipeline(),
		},
		{ // no commit set for pipeline
			failure: true,
			pipeline: &Pipeline{
				ID:      sql.NullInt64{Int64: 1, Valid: true},
				RepoID:  sql.NullInt64{Int64: 1, Valid: true},
				Ref:     sql.NullString{String: "refs/heads/main", Valid: true},
				Type:    sql.NullString{String: constants.PipelineTypeYAML, Valid: true},
				Version: sql.NullString{String: "1", Valid: true},
			},
		},
		{ // no ref set for pipeline
			failure: true,
			pipeline: &Pipeline{
				ID:      sql.NullInt64{Int64: 1, Valid: true},
				RepoID:  sql.NullInt64{Int64: 1, Valid: true},
				Commit:  sql.NullString{String: "48afb5bdc41ad69bf22588491333f7cf71135163", Valid: true},
				Type:    sql.NullString{String: constants.PipelineTypeYAML, Valid: true},
				Version: sql.NullString{String: "1", Valid: true},
			},
		},
		{ // no repo_id set for pipeline
			failure: true,
			pipeline: &Pipeline{
				ID:      sql.NullInt64{Int64: 1, Valid: true},
				Commit:  sql.NullString{String: "48afb5bdc41ad69bf22588491333f7cf71135163", Valid: true},
				Ref:     sql.NullString{String: "refs/heads/main", Valid: true},
				Type:    sql.NullString{String: constants.PipelineTypeYAML, Valid: true},
				Version: sql.NullString{String: "1", Valid: true},
			},
		},
		{ // no type set for pipeline
			failure: true,
			pipeline: &Pipeline{
				ID:      sql.NullInt64{Int64: 1, Valid: true},
				RepoID:  sql.NullInt64{Int64: 1, Valid: true},
				Commit:  sql.NullString{String: "48afb5bdc41ad69bf22588491333f7cf71135163", Valid: true},
				Ref:     sql.NullString{String: "refs/heads/main", Valid: true},
				Version: sql.NullString{String: "1", Valid: true},
			},
		},
		{ // no version set for pipeline
			failure: true,
			pipeline: &Pipeline{
				ID:     sql.NullInt64{Int64: 1, Valid: true},
				RepoID: sql.NullInt64{Int64: 1, Valid: true},
				Commit: sql.NullString{String: "48afb5bdc41ad69bf22588491333f7cf71135163", Valid: true},
				Ref:    sql.NullString{String: "refs/heads/main", Valid: true},

				Type: sql.NullString{String: constants.PipelineTypeYAML, Valid: true},
			},
		},
	}

	// run tests
	for _, test := range tests {
		err := test.pipeline.Validate()

		if test.failure {
			if err == nil {
				t.Errorf("Validate should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("Validate returned err: %v", err)
		}
	}
}

func TestDatabase_PipelineFromAPI(t *testing.T) {
	// setup types
	want := &Pipeline{
		ID:              sql.NullInt64{Int64: 1, Valid: true},
		RepoID:          sql.NullInt64{Int64: 1, Valid: true},
		Commit:          sql.NullString{String: "48afb5bdc41ad69bf22588491333f7cf71135163", Valid: true},
		Flavor:          sql.NullString{String: "large", Valid: true},
		Platform:        sql.NullString{String: "docker", Valid: true},
		Ref:             sql.NullString{String: "refs/heads/main", Valid: true},
		Type:            sql.NullString{String: constants.PipelineTypeYAML, Valid: true},
		Version:         sql.NullString{String: "1", Valid: true},
		ExternalSecrets: sql.NullBool{Bool: false, Valid: true},
		InternalSecrets: sql.NullBool{Bool: false, Valid: true},
		Services:        sql.NullBool{Bool: true, Valid: true},
		Stages:          sql.NullBool{Bool: false, Valid: true},
		Steps:           sql.NullBool{Bool: true, Valid: true},
		Templates:       sql.NullBool{Bool: false, Valid: true},
		Data:            testPipelineData(),
	}

	p := new(api.Pipeline)

	p.SetID(1)
	p.SetRepo(testRepo().ToAPI())
	p.SetCommit("48afb5bdc41ad69bf22588491333f7cf71135163")
	p.SetFlavor("large")
	p.SetPlatform("docker")
	p.SetRef("refs/heads/main")
	p.SetType(constants.PipelineTypeYAML)
	p.SetVersion("1")
	p.SetExternalSecrets(false)
	p.SetInternalSecrets(false)
	p.SetServices(true)
	p.SetStages(false)
	p.SetSteps(true)
	p.SetTemplates(false)
	p.SetData(testPipelineData())

	// run test
	got := PipelineFromAPI(p)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("PipelineFromAPI is %v, want %v", got, want)
	}
}

// testPipeline is a test helper function to create a Pipeline
// type with all fields set to a fake value.
func testPipeline() *Pipeline {
	return &Pipeline{
		ID:              sql.NullInt64{Int64: 1, Valid: true},
		RepoID:          sql.NullInt64{Int64: 1, Valid: true},
		Commit:          sql.NullString{String: "48afb5bdc41ad69bf22588491333f7cf71135163", Valid: true},
		Flavor:          sql.NullString{String: "large", Valid: true},
		Platform:        sql.NullString{String: "docker", Valid: true},
		Ref:             sql.NullString{String: "refs/heads/main", Valid: true},
		Type:            sql.NullString{String: constants.PipelineTypeYAML, Valid: true},
		Version:         sql.NullString{String: "1", Valid: true},
		ExternalSecrets: sql.NullBool{Bool: false, Valid: true},
		InternalSecrets: sql.NullBool{Bool: false, Valid: true},
		Services:        sql.NullBool{Bool: true, Valid: true},
		Stages:          sql.NullBool{Bool: false, Valid: true},
		Steps:           sql.NullBool{Bool: true, Valid: true},
		Templates:       sql.NullBool{Bool: false, Valid: true},
		Data:            testPipelineData(),

		Repo: *testRepo(),
	}
}

// testPipelineData is a test helper function to create the
// content for the Data field for the Pipeline type.
func testPipelineData() []byte {
	return []byte(`
version: 1

worker:
  flavor: large
  platform: docker

services:
  - name: redis
    image: redis

steps:
  - name: ping
    image: redis
    commands:
      - redis-cli -h redis ping
`)
}
