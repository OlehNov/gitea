// Copyright 2017 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package models

import (
	"time"

	"github.com/go-xorm/xorm"
)

// ActivityStats represets issue and pull request information.
type ActivityStats struct {
	OpenedPRs                   PullRequestList
	OpenedPRAuthorCount         int64
	MergedPRs                   PullRequestList
	MergedPRAuthorCount         int64
	OpenedIssues                IssueList
	OpenedIssueAuthorCount      int64
	ClosedIssues                IssueList
	ClosedIssueAuthorCount      int64
	UnresolvedIssues            IssueList
	PublishedReleases           []*Release
	PublishedReleaseAuthorCount int64
}

// ActivePRCount returns total active pull request count
func (stats *ActivityStats) ActivePRCount() int {
	return stats.OpenedPRCount() + stats.MergedPRCount()
}

// OpenedPRCount returns opened pull request count
func (stats *ActivityStats) OpenedPRCount() int {
	return len(stats.OpenedPRs)
}

// OpenedPRPerc returns opened pull request percents from total active
func (stats *ActivityStats) OpenedPRPerc() int {
	return int(float32(stats.OpenedPRCount()) / float32(stats.ActivePRCount()) * 100.0)
}

// MergedPRCount returns merged pull request count
func (stats *ActivityStats) MergedPRCount() int {
	return len(stats.MergedPRs)
}

// MergedPRPerc returns merged pull request percent from total active
func (stats *ActivityStats) MergedPRPerc() int {
	return int(float32(stats.MergedPRCount()) / float32(stats.ActivePRCount()) * 100.0)
}

// ActiveIssueCount returns total active issue count
func (stats *ActivityStats) ActiveIssueCount() int {
	return stats.OpenedIssueCount() + stats.ClosedIssueCount()
}

// OpenedIssueCount returns open issue count
func (stats *ActivityStats) OpenedIssueCount() int {
	return len(stats.OpenedIssues)
}

// OpenedIssuePerc returns open issue count percent from total active
func (stats *ActivityStats) OpenedIssuePerc() int {
	return int(float32(stats.OpenedIssueCount()) / float32(stats.ActiveIssueCount()) * 100.0)
}

// ClosedIssueCount returns closed issue count
func (stats *ActivityStats) ClosedIssueCount() int {
	return len(stats.ClosedIssues)
}

// ClosedIssuePerc returns closed issue count percent from total active
func (stats *ActivityStats) ClosedIssuePerc() int {
	return int(float32(stats.ClosedIssueCount()) / float32(stats.ActiveIssueCount()) * 100.0)
}

// UnresolvedIssueCount returns unresolved issue and pull request count
func (stats *ActivityStats) UnresolvedIssueCount() int {
	return len(stats.UnresolvedIssues)
}

// PublishedReleaseCount returns published release count
func (stats *ActivityStats) PublishedReleaseCount() int {
	return len(stats.PublishedReleases)
}

// FillPullRequestsForActivity returns pull request information for activity page
func FillPullRequestsForActivity(stats *ActivityStats, baseRepoID int64, fromTime time.Time) error {
	var err error
	var count int64

	// Merged pull requests
	sess := pullRequestsForActivityStatement(baseRepoID, fromTime, true)
	sess.OrderBy("pull_request.merged_unix DESC")
	stats.MergedPRs = make(PullRequestList, 0)
	if err = sess.Find(&stats.MergedPRs); err != nil {
		return err
	}
	if err = stats.MergedPRs.LoadAttributes(); err != nil {
		return err
	}

	// Merged pull request authors
	sess = pullRequestsForActivityStatement(baseRepoID, fromTime, true)
	if _, err = sess.Select("count(distinct issue.poster_id) as `count`").Table("pull_request").Get(&count); err != nil {
		return err
	}
	stats.MergedPRAuthorCount = count

	// Opened pull requests
	sess = pullRequestsForActivityStatement(baseRepoID, fromTime, false)
	sess.OrderBy("issue.created_unix ASC")
	stats.OpenedPRs = make(PullRequestList, 0)
	if err = sess.Find(&stats.OpenedPRs); err != nil {
		return err
	}
	if err = stats.OpenedPRs.LoadAttributes(); err != nil {
		return err
	}

	// Opened pull request authors
	sess = pullRequestsForActivityStatement(baseRepoID, fromTime, false)
	if _, err = sess.Select("count(distinct issue.poster_id) as `count`").Table("pull_request").Get(&count); err != nil {
		return err
	}
	stats.OpenedPRAuthorCount = count

	return nil
}

func pullRequestsForActivityStatement(baseRepoID int64, fromTime time.Time, merged bool) *xorm.Session {
	sess := x.Where("pull_request.base_repo_id=?", baseRepoID).
		Join("INNER", "issue", "pull_request.issue_id = issue.id")

	if merged {
		sess.And("pull_request.has_merged = ?", true)
		sess.And("pull_request.merged_unix >= ?", fromTime.Unix())
	} else {
		sess.And("issue.is_closed = ?", false)
		sess.And("issue.created_unix >= ?", fromTime.Unix())
	}

	return sess
}

// FillIssuesForActivity returns issue information for activity page
func FillIssuesForActivity(stats *ActivityStats, baseRepoID int64, fromTime time.Time) error {
	var err error
	var count int64

	// Closed issues
	sess := issuesForActivityStatement(baseRepoID, fromTime, true, false)
	sess.OrderBy("issue.updated_unix DESC")
	stats.ClosedIssues = make(IssueList, 0)
	if err = sess.Find(&stats.ClosedIssues); err != nil {
		return err
	}

	// Closed issue authors
	sess = issuesForActivityStatement(baseRepoID, fromTime, true, false)
	if _, err = sess.Select("count(distinct issue.poster_id) as `count`").Table("issue").Get(&count); err != nil {
		return err
	}
	stats.ClosedIssueAuthorCount = count

	// New issues
	sess = issuesForActivityStatement(baseRepoID, fromTime, false, false)
	sess.OrderBy("issue.created_unix ASC")
	stats.OpenedIssues = make(IssueList, 0)
	if err = sess.Find(&stats.OpenedIssues); err != nil {
		return err
	}

	// Opened issue authors
	sess = issuesForActivityStatement(baseRepoID, fromTime, false, false)
	if _, err = sess.Select("count(distinct issue.poster_id) as `count`").Table("issue").Get(&count); err != nil {
		return err
	}
	stats.OpenedIssueAuthorCount = count

	return nil
}

// FillUnresolvedIssuesForActivity returns unresolved issue and pull request information for activity page
func FillUnresolvedIssuesForActivity(stats *ActivityStats, baseRepoID int64, fromTime time.Time, issues, prs bool) error {
	// Check if we need to select anything
	if !issues && !prs {
		return nil
	}
	sess := issuesForActivityStatement(baseRepoID, fromTime, false, true)
	if !issues || !prs {
		sess.And("issue.is_pull = ?", prs)
	}
	sess.OrderBy("issue.updated_unix DESC")
	stats.UnresolvedIssues = make(IssueList, 0)
	return sess.Find(&stats.UnresolvedIssues)
}

func issuesForActivityStatement(baseRepoID int64, fromTime time.Time, closed, unresolved bool) *xorm.Session {
	sess := x.Where("issue.repo_id = ?", baseRepoID).
		And("issue.is_closed = ?", closed)

	if !unresolved {
		sess.And("issue.is_pull = ?", false)
		sess.And("issue.created_unix >= ?", fromTime.Unix())
	} else {
		sess.And("issue.created_unix < ?", fromTime.Unix())
		sess.And("issue.updated_unix >= ?", fromTime.Unix())
	}

	return sess
}

// FillReleasesForActivity returns release information for activity page
func FillReleasesForActivity(stats *ActivityStats, baseRepoID int64, fromTime time.Time) error {
	var err error
	var count int64

	// Published releases list
	sess := releasesForActivityStatement(baseRepoID, fromTime)
	sess.OrderBy("release.created_unix DESC")
	stats.PublishedReleases = make([]*Release, 0)
	if err = sess.Find(&stats.PublishedReleases); err != nil {
		return err
	}

	// Published releases authors
	sess = releasesForActivityStatement(baseRepoID, fromTime)
	if _, err = sess.Select("count(distinct release.publisher_id) as `count`").Table("release").Get(&count); err != nil {
		return err
	}
	stats.PublishedReleaseAuthorCount = count

	return nil
}

func releasesForActivityStatement(baseRepoID int64, fromTime time.Time) *xorm.Session {
	return x.Where("release.repo_id = ?", baseRepoID).
		And("release.is_draft = ?", false).
		And("release.created_unix >= ?", fromTime.Unix())
}
