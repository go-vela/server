// SPDX-License-Identifier: Apache-2.0

package constants

// Allowed repo events. NOTE: these can NOT change order. New events must be added at the end.
const (
	AllowPushBranch = 1 << iota // 00000001 = 1
	AllowPushTag                // 00000010 = 2
	AllowPullOpen               // 00000100 = 4
	AllowPullEdit               // ...
	AllowPullSync
	_ // AllowPullAssigned - Not Implemented
	_ // AllowPullMilestoned - Not Implemented
	AllowPullLabel
	_ // AllowPullLocked - Not Implemented
	_ // AllowPullReady - Not Implemented
	AllowPullReopen
	_ // AllowPullReviewRequest - Not Implemented
	_ // AllowPullClosed - Not Implemented
	AllowDeployCreate
	AllowCommentCreate
	AllowCommentEdit
	AllowSchedule
	AllowPushDeleteBranch
	AllowPushDeleteTag
	AllowPullUnlabel
)
