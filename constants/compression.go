// SPDX-License-Identifier: Apache-2.0

package constants

// Log Compression Levels.
const (
	// The default compression level for the compress/zlib library
	// for log data stored in the database.
	CompressionNegOne = -1

	// Enables no compression for log data stored in the database.
	//
	// This produces no compression for the log data.
	CompressionZero = 0

	// Enables the best speed for log data stored in the database.
	//
	// This produces compression for the log data the fastest but
	// has a tradeoff of producing the largest amounts of data.
	CompressionOne = 1

	// Second compression level for log data stored in the database.
	CompressionTwo = 2

	// Third compression level for log data stored in the database.
	CompressionThree = 3

	// Fourth compression level for log data stored in the database.
	CompressionFour = 4

	// Enables an even balance of speed and compression for log
	// data stored in the database.
	//
	// This produces compression for the log data with an even
	// balance of speed while producing smaller amounts of data.
	CompressionFive = 5

	// Sixth compression level for log data stored in the database.
	CompressionSix = 6

	// Seventh compression level for log data stored in the database.
	CompressionSeven = 7

	// Eighth compression level for log data stored in the database.
	CompressionEight = 8

	// Enables the best compression for log data stored in the database.
	//
	// This produces compression for the log data the slowest but
	// has a tradeoff of producing the smallest amounts of data.
	CompressionNine = 9
)
