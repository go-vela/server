// SPDX-License-Identifier: Apache-2.0

package github

// MergeQueueBranchPrefix outputs the prefix for merge queue branches.
func (c *Client) MergeQueueBranchPrefix() string {
	return mergeQueueBranchPrefix
}
