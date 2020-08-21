package sources

// -----> IMPORTANT / DISCLAIMER <-----
// The source of this file was extracted from https://github.com/google/go-github package
//
// Copyright 2016 The go-github AUTHORS. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"encoding/json"
	"time"
)

// Source from https://github.com/google/go-github
// I extract only the necessary structures.

// Event represents a GitHub event.
type Event struct {
	Type       *string          `json:"type,omitempty"`
	Public     *bool            `json:"public,omitempty"`
	RawPayload *json.RawMessage `json:"payload,omitempty"`
	Repo       *Repository      `json:"repo,omitempty"`
	Actor      *User            `json:"actor,omitempty"`
	CreatedAt  *time.Time       `json:"created_at,omitempty"`
	ID         *string          `json:"id,omitempty"`
}

// User represents a GitHub user.
type User struct {
	Login             *string    `json:"login,omitempty"`
	ID                *int64     `json:"id,omitempty"`
	NodeID            *string    `json:"node_id,omitempty"`
	AvatarURL         *string    `json:"avatar_url,omitempty"`
	HTMLURL           *string    `json:"html_url,omitempty"`
	GravatarID        *string    `json:"gravatar_id,omitempty"`
	Name              *string    `json:"name,omitempty"`
	Company           *string    `json:"company,omitempty"`
	Blog              *string    `json:"blog,omitempty"`
	Location          *string    `json:"location,omitempty"`
	Email             *string    `json:"email,omitempty"`
	Hireable          *bool      `json:"hireable,omitempty"`
	Bio               *string    `json:"bio,omitempty"`
	PublicRepos       *int       `json:"public_repos,omitempty"`
	PublicGists       *int       `json:"public_gists,omitempty"`
	Followers         *int       `json:"followers,omitempty"`
	Following         *int       `json:"following,omitempty"`
	Type              *string    `json:"type,omitempty"`
	SiteAdmin         *bool      `json:"site_admin,omitempty"`
	TotalPrivateRepos *int       `json:"total_private_repos,omitempty"`
	OwnedPrivateRepos *int       `json:"owned_private_repos,omitempty"`
	PrivateGists      *int       `json:"private_gists,omitempty"`
	DiskUsage         *int       `json:"disk_usage,omitempty"`
	Collaborators     *int       `json:"collaborators,omitempty"`

	// API URLs
	URL               *string `json:"url,omitempty"`
	EventsURL         *string `json:"events_url,omitempty"`
	FollowingURL      *string `json:"following_url,omitempty"`
	FollowersURL      *string `json:"followers_url,omitempty"`
	GistsURL          *string `json:"gists_url,omitempty"`
	ReceivedEventsURL *string `json:"received_events_url,omitempty"`
	ReposURL          *string `json:"repos_url,omitempty"`
	StarredURL        *string `json:"starred_url,omitempty"`
	SubscriptionsURL  *string `json:"subscriptions_url,omitempty"`
}


// Repository represents a GitHub repository.
type Repository struct {
	ID               *int64           `json:"id,omitempty"`
	NodeID           *string          `json:"node_id,omitempty"`
	Owner            *User            `json:"owner,omitempty"`
	Name             *string          `json:"name,omitempty"`
	FullName         *string          `json:"full_name,omitempty"`
	Description      *string          `json:"description,omitempty"`
	Homepage         *string          `json:"homepage,omitempty"`

	DefaultBranch    *string          `json:"default_branch,omitempty"`
	MasterBranch     *string          `json:"master_branch,omitempty"`
	HTMLURL          *string          `json:"html_url,omitempty"`
	CloneURL         *string          `json:"clone_url,omitempty"`
	GitURL           *string          `json:"git_url,omitempty"`
	MirrorURL        *string          `json:"mirror_url,omitempty"`
	SSHURL           *string          `json:"ssh_url,omitempty"`
	SVNURL           *string          `json:"svn_url,omitempty"`
	Language         *string          `json:"language,omitempty"`
	Fork             *bool            `json:"fork,omitempty"`
	ForksCount       *int             `json:"forks_count,omitempty"`
	NetworkCount     *int             `json:"network_count,omitempty"`
	OpenIssuesCount  *int             `json:"open_issues_count,omitempty"`
	StargazersCount  *int             `json:"stargazers_count,omitempty"`
	SubscribersCount *int             `json:"subscribers_count,omitempty"`
	WatchersCount    *int             `json:"watchers_count,omitempty"`
	Size             *int             `json:"size,omitempty"`
	AutoInit         *bool            `json:"auto_init,omitempty"`
	Parent           *Repository      `json:"parent,omitempty"`
	Source           *Repository      `json:"source,omitempty"`

	Permissions      *map[string]bool `json:"permissions,omitempty"`
	AllowRebaseMerge *bool            `json:"allow_rebase_merge,omitempty"`
	AllowSquashMerge *bool            `json:"allow_squash_merge,omitempty"`
	AllowMergeCommit *bool            `json:"allow_merge_commit,omitempty"`
	Topics           []string         `json:"topics,omitempty"`

	// Additional mutable fields when creating and editing a repository
	Private           *bool   `json:"private,omitempty"`
	HasIssues         *bool   `json:"has_issues,omitempty"`
	HasWiki           *bool   `json:"has_wiki,omitempty"`
	HasPages          *bool   `json:"has_pages,omitempty"`
	HasProjects       *bool   `json:"has_projects,omitempty"`
	HasDownloads      *bool   `json:"has_downloads,omitempty"`
	LicenseTemplate   *string `json:"license_template,omitempty"`
	GitignoreTemplate *string `json:"gitignore_template,omitempty"`
	Archived          *bool   `json:"archived,omitempty"`

	// Creating an organization repository. Required for non-owners.
	TeamID *int64 `json:"team_id,omitempty"`

	// API URLs
	URL              *string `json:"url,omitempty"`
	ArchiveURL       *string `json:"archive_url,omitempty"`
	AssigneesURL     *string `json:"assignees_url,omitempty"`
	BlobsURL         *string `json:"blobs_url,omitempty"`
	BranchesURL      *string `json:"branches_url,omitempty"`
	CollaboratorsURL *string `json:"collaborators_url,omitempty"`
	CommentsURL      *string `json:"comments_url,omitempty"`
	CommitsURL       *string `json:"commits_url,omitempty"`
	CompareURL       *string `json:"compare_url,omitempty"`
	ContentsURL      *string `json:"contents_url,omitempty"`
	ContributorsURL  *string `json:"contributors_url,omitempty"`
	DeploymentsURL   *string `json:"deployments_url,omitempty"`
	DownloadsURL     *string `json:"downloads_url,omitempty"`
	EventsURL        *string `json:"events_url,omitempty"`
	ForksURL         *string `json:"forks_url,omitempty"`
	GitCommitsURL    *string `json:"git_commits_url,omitempty"`
	GitRefsURL       *string `json:"git_refs_url,omitempty"`
	GitTagsURL       *string `json:"git_tags_url,omitempty"`
	HooksURL         *string `json:"hooks_url,omitempty"`
	IssueCommentURL  *string `json:"issue_comment_url,omitempty"`
	IssueEventsURL   *string `json:"issue_events_url,omitempty"`
	IssuesURL        *string `json:"issues_url,omitempty"`
	KeysURL          *string `json:"keys_url,omitempty"`
	LabelsURL        *string `json:"labels_url,omitempty"`
	LanguagesURL     *string `json:"languages_url,omitempty"`
	MergesURL        *string `json:"merges_url,omitempty"`
	MilestonesURL    *string `json:"milestones_url,omitempty"`
	NotificationsURL *string `json:"notifications_url,omitempty"`
	PullsURL         *string `json:"pulls_url,omitempty"`
	ReleasesURL      *string `json:"releases_url,omitempty"`
	StargazersURL    *string `json:"stargazers_url,omitempty"`
	StatusesURL      *string `json:"statuses_url,omitempty"`
	SubscribersURL   *string `json:"subscribers_url,omitempty"`
	SubscriptionURL  *string `json:"subscription_url,omitempty"`
	TagsURL          *string `json:"tags_url,omitempty"`
	TreesURL         *string `json:"trees_url,omitempty"`
	TeamsURL         *string `json:"teams_url,omitempty"`
}

func (e *Event) ParsePayload() (payload interface{}, err error) {
	switch *e.Type {
	case "CreateEvent":
		payload = &CreateEvent{}
	case "ForkEvent":
		payload = &ForkEvent{}
	case "IssuesEvent":
		payload = &IssuesEvent{}
	case "PullRequestEvent":
		payload = &PullRequestEvent{}
	case "PullRequestReviewEvent":
		payload = &PullRequestReviewEvent{}
	case "PullRequestReviewCommentEvent":
		payload = &PullRequestReviewCommentEvent{}
	case "PushEvent":
		payload = &PushEvent{}
	case "ReleaseEvent":
		payload = &ReleaseEvent{}
	case "WatchEvent":
		payload = &WatchEvent{}
	}
	err = json.Unmarshal(*e.RawPayload, &payload)
	return payload, err
}


// CreateEvent represents a created repository, branch, or tag.
// The Webhook event name is "create".
//
// Note: webhooks will not receive this event for created repositories.
// Additionally, webhooks will not receive this event for tags if more
// than three tags are pushed at once.
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#createevent
type CreateEvent struct {
	Ref *string `json:"ref,omitempty"`
	// RefType is the object that was created. Possible values are: "repository", "branch", "tag".
	RefType      *string `json:"ref_type,omitempty"`
	MasterBranch *string `json:"master_branch,omitempty"`
	Description  *string `json:"description,omitempty"`

	// The following fields are only populated by Webhook events.
	PusherType   *string       `json:"pusher_type,omitempty"`
	Repo         *Repository   `json:"repository,omitempty"`
	Sender       *User         `json:"sender,omitempty"`
}


// ForkEvent is triggered when a user forks a repository.
// The Webhook event name is "fork".
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#forkevent
type ForkEvent struct {
	// Forkee is the created repository.
	Forkee *Repository `json:"forkee,omitempty"`

	// The following fields are only populated by Webhook events.
	Repo         *Repository   `json:"repository,omitempty"`
	Sender       *User         `json:"sender,omitempty"`
}

// IssuesEvent is triggered when an issue is assigned, unassigned, labeled,
// unlabeled, opened, closed, or reopened.
// The Webhook event name is "issues".
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#issuesevent
type IssuesEvent struct {
	// Action is the action that was performed. Possible values are: "assigned",
	// "unassigned", "labeled", "unlabeled", "opened", "closed", "reopened", "edited".
	Action   *string `json:"action,omitempty"`
	Issue    *Issue  `json:"issue,omitempty"`
	Assignee *User   `json:"assignee,omitempty"`
	Repo         *Repository   `json:"repository,omitempty"`
	Sender       *User         `json:"sender,omitempty"`
}

// Issue represents a GitHub issue on a repository.
//
// Note: As far as the GitHub API is concerned, every pull request is an issue,
// but not every issue is a pull request. Some endpoints, events, and webhooks
// may also return pull requests via this struct. If PullRequestLinks is nil,
// this is an issue, and if PullRequestLinks is not nil, this is a pull request.
// The IsPullRequest helper method can be used to check that.
type Issue struct {
	ID               *int64            `json:"id,omitempty"`
	Number           *int              `json:"number,omitempty"`
	State            *string           `json:"state,omitempty"`
	Locked           *bool             `json:"locked,omitempty"`
	Title            *string           `json:"title,omitempty"`
	Body             *string           `json:"body,omitempty"`
	User             *User             `json:"user,omitempty"`
	Assignee         *User             `json:"assignee,omitempty"`
	Comments         *int              `json:"comments,omitempty"`
	ClosedAt         *time.Time        `json:"closed_at,omitempty"`
	CreatedAt        *time.Time        `json:"created_at,omitempty"`
	UpdatedAt        *time.Time        `json:"updated_at,omitempty"`
	ClosedBy         *User             `json:"closed_by,omitempty"`
	URL              *string           `json:"url,omitempty"`
	HTMLURL          *string           `json:"html_url,omitempty"`
	CommentsURL      *string           `json:"comments_url,omitempty"`
	EventsURL        *string           `json:"events_url,omitempty"`
	LabelsURL        *string           `json:"labels_url,omitempty"`
	RepositoryURL    *string           `json:"repository_url,omitempty"`
	Repository       *Repository       `json:"repository,omitempty"`
	Assignees        []*User           `json:"assignees,omitempty"`
	NodeID           *string           `json:"node_id,omitempty"`
}

// PullRequestEvent is triggered when a pull request is assigned, unassigned,
// labeled, unlabeled, opened, closed, reopened, or synchronized.
// The Webhook event name is "pull_request".
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#pullrequestevent
type PullRequestEvent struct {
	// Action is the action that was performed. Possible values are:
	// "assigned", "unassigned", "review_requested", "review_request_removed", "labeled", "unlabeled",
	// "opened", "closed", "reopened", "synchronize", "edited".
	// If the action is "closed" and the merged key is false,
	// the pull request was closed with unmerged commits. If the action is "closed"
	// and the merged key is true, the pull request was merged.
	Action      *string      `json:"action,omitempty"`
	Number      *int         `json:"number,omitempty"`
	PullRequest *PullRequest `json:"pull_request,omitempty"`

	// RequestedReviewer is populated in "review_requested", "review_request_removed" event deliveries.
	// A request affecting multiple reviewers at once is split into multiple
	// such event deliveries, each with a single, different RequestedReviewer.
	RequestedReviewer *User         `json:"requested_reviewer,omitempty"`
	Repo              *Repository   `json:"repository,omitempty"`
	Sender            *User         `json:"sender,omitempty"`
}

type PullRequest struct {
	ID                  *int64     `json:"id,omitempty"`
	Number              *int       `json:"number,omitempty"`
	State               *string    `json:"state,omitempty"`
	Title               *string    `json:"title,omitempty"`
	Body                *string    `json:"body,omitempty"`
	CreatedAt           *time.Time `json:"created_at,omitempty"`
	UpdatedAt           *time.Time `json:"updated_at,omitempty"`
	ClosedAt            *time.Time `json:"closed_at,omitempty"`
	MergedAt            *time.Time `json:"merged_at,omitempty"`
	User                *User      `json:"user,omitempty"`
	Merged              *bool      `json:"merged,omitempty"`
	Mergeable           *bool      `json:"mergeable,omitempty"`
	MergeableState      *string    `json:"mergeable_state,omitempty"`
	MergedBy            *User      `json:"merged_by,omitempty"`
	MergeCommitSHA      *string    `json:"merge_commit_sha,omitempty"`
	Comments            *int       `json:"comments,omitempty"`
	Commits             *int       `json:"commits,omitempty"`
	Additions           *int       `json:"additions,omitempty"`
	Deletions           *int       `json:"deletions,omitempty"`
	ChangedFiles        *int       `json:"changed_files,omitempty"`
	URL                 *string    `json:"url,omitempty"`
	HTMLURL             *string    `json:"html_url,omitempty"`
	IssueURL            *string    `json:"issue_url,omitempty"`
	StatusesURL         *string    `json:"statuses_url,omitempty"`
	DiffURL             *string    `json:"diff_url,omitempty"`
	PatchURL            *string    `json:"patch_url,omitempty"`
	CommitsURL          *string    `json:"commits_url,omitempty"`
	CommentsURL         *string    `json:"comments_url,omitempty"`
	ReviewCommentsURL   *string    `json:"review_comments_url,omitempty"`
	ReviewCommentURL    *string    `json:"review_comment_url,omitempty"`
	Assignee            *User      `json:"assignee,omitempty"`
	Assignees           []*User    `json:"assignees,omitempty"`
	MaintainerCanModify *bool      `json:"maintainer_can_modify,omitempty"`
	AuthorAssociation   *string    `json:"author_association,omitempty"`
	NodeID              *string    `json:"node_id,omitempty"`
	RequestedReviewers  []*User    `json:"requested_reviewers,omitempty"`
}

// PullRequestReviewEvent is triggered when a review is submitted on a pull
// request.
// The Webhook event name is "pull_request_review".
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#pullrequestreviewevent
type PullRequestReviewEvent struct {
	// Action is always "submitted".
	Action      *string            `json:"action,omitempty"`
	Review      *PullRequestReview `json:"review,omitempty"`
	PullRequest *PullRequest       `json:"pull_request,omitempty"`

}

// PullRequestReview represents a review of a pull request.
type PullRequestReview struct {
	ID             *int64     `json:"id,omitempty"`
	User           *User      `json:"user,omitempty"`
	Body           *string    `json:"body,omitempty"`
	SubmittedAt    *time.Time `json:"submitted_at,omitempty"`
	CommitID       *string    `json:"commit_id,omitempty"`
	HTMLURL        *string    `json:"html_url,omitempty"`
	PullRequestURL *string    `json:"pull_request_url,omitempty"`
	State          *string    `json:"state,omitempty"`
}

// PullRequestReviewCommentEvent is triggered when a comment is created on a
// portion of the unified diff of a pull request.
// The Webhook event name is "pull_request_review_comment".
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#pullrequestreviewcommentevent
type PullRequestReviewCommentEvent struct {
	// Action is the action that was performed on the comment.
	// Possible values are: "created", "edited", "deleted".
	Action      *string             `json:"action,omitempty"`
	PullRequest *PullRequest        `json:"pull_request,omitempty"`
	Comment     *PullRequestComment `json:"comment,omitempty"`
}

// PullRequestComment represents a comment left on a pull request.
type PullRequestComment struct {
	ID                  *int64     `json:"id,omitempty"`
	InReplyTo           *int64     `json:"in_reply_to_id,omitempty"`
	Body                *string    `json:"body,omitempty"`
	Path                *string    `json:"path,omitempty"`
	DiffHunk            *string    `json:"diff_hunk,omitempty"`
	PullRequestReviewID *int64     `json:"pull_request_review_id,omitempty"`
	Position            *int       `json:"position,omitempty"`
	OriginalPosition    *int       `json:"original_position,omitempty"`
	CommitID            *string    `json:"commit_id,omitempty"`
	OriginalCommitID    *string    `json:"original_commit_id,omitempty"`
	User                *User      `json:"user,omitempty"`
	CreatedAt           *time.Time `json:"created_at,omitempty"`
	UpdatedAt           *time.Time `json:"updated_at,omitempty"`
	// AuthorAssociation is the comment author's relationship to the pull request's repository.
	// Possible values are "COLLABORATOR", "CONTRIBUTOR", "FIRST_TIMER", "FIRST_TIME_CONTRIBUTOR", "MEMBER", "OWNER", or "NONE".
	AuthorAssociation *string `json:"author_association,omitempty"`
	URL               *string `json:"url,omitempty"`
	HTMLURL           *string `json:"html_url,omitempty"`
	PullRequestURL    *string `json:"pull_request_url,omitempty"`
}

// PushEvent represents a git push to a GitHub repository.
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#pushevent
type PushEvent struct {
	PushID       *int64            `json:"push_id,omitempty"`
	Head         *string           `json:"head,omitempty"`
	Ref          *string           `json:"ref,omitempty"`
	Size         *int              `json:"size,omitempty"`
	Commits      []PushEventCommit `json:"commits,omitempty"`
	Before       *string           `json:"before,omitempty"`
	DistinctSize *int              `json:"distinct_size,omitempty"`

	// The following fields are only populated by Webhook events.
	After        *string              `json:"after,omitempty"`
	Created      *bool                `json:"created,omitempty"`
	Deleted      *bool                `json:"deleted,omitempty"`
	Forced       *bool                `json:"forced,omitempty"`
	BaseRef      *string              `json:"base_ref,omitempty"`
	Compare      *string              `json:"compare,omitempty"`
	Repo         *PushEventRepository `json:"repository,omitempty"`
	HeadCommit   *PushEventCommit     `json:"head_commit,omitempty"`
	Pusher       *User                `json:"pusher,omitempty"`
	Sender       *User                `json:"sender,omitempty"`
}

// PushEventRepository represents the repo object in a PushEvent payload.
type PushEventRepository struct {
	ID              *int64              `json:"id,omitempty"`
	NodeID          *string             `json:"node_id,omitempty"`
	Name            *string             `json:"name,omitempty"`
	FullName        *string             `json:"full_name,omitempty"`
	Owner           *PushEventRepoOwner `json:"owner,omitempty"`
	Private         *bool               `json:"private,omitempty"`
	Description     *string             `json:"description,omitempty"`
	Fork            *bool               `json:"fork,omitempty"`
	CreatedAt       *Timestamp          `json:"created_at,omitempty"`
	PushedAt        *Timestamp          `json:"pushed_at,omitempty"`
	UpdatedAt       *Timestamp          `json:"updated_at,omitempty"`
	Homepage        *string             `json:"homepage,omitempty"`
	Size            *int                `json:"size,omitempty"`
	StargazersCount *int                `json:"stargazers_count,omitempty"`
	WatchersCount   *int                `json:"watchers_count,omitempty"`
	Language        *string             `json:"language,omitempty"`
	HasIssues       *bool               `json:"has_issues,omitempty"`
	HasDownloads    *bool               `json:"has_downloads,omitempty"`
	HasWiki         *bool               `json:"has_wiki,omitempty"`
	HasPages        *bool               `json:"has_pages,omitempty"`
	ForksCount      *int                `json:"forks_count,omitempty"`
	OpenIssuesCount *int                `json:"open_issues_count,omitempty"`
	DefaultBranch   *string             `json:"default_branch,omitempty"`
	MasterBranch    *string             `json:"master_branch,omitempty"`
	Organization    *string             `json:"organization,omitempty"`
	URL             *string             `json:"url,omitempty"`
	ArchiveURL      *string             `json:"archive_url,omitempty"`
	HTMLURL         *string             `json:"html_url,omitempty"`
	StatusesURL     *string             `json:"statuses_url,omitempty"`
	GitURL          *string             `json:"git_url,omitempty"`
	SSHURL          *string             `json:"ssh_url,omitempty"`
	CloneURL        *string             `json:"clone_url,omitempty"`
	SVNURL          *string             `json:"svn_url,omitempty"`
}

// PushEventRepoOwner is a basic representation of user/org in a PushEvent payload.
type PushEventRepoOwner struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
}

// PushEventCommit represents a git commit in a GitHub PushEvent.
type PushEventCommit struct {
	Message  *string       `json:"message,omitempty"`
	Author   *CommitAuthor `json:"author,omitempty"`
	URL      *string       `json:"url,omitempty"`
	Distinct *bool         `json:"distinct,omitempty"`

	// The following fields are only populated by Events API.
	SHA *string `json:"sha,omitempty"`

	// The following fields are only populated by Webhook events.
	ID        *string       `json:"id,omitempty"`
	TreeID    *string       `json:"tree_id,omitempty"`
	Timestamp *Timestamp    `json:"timestamp,omitempty"`
	Committer *CommitAuthor `json:"committer,omitempty"`
	Added     []string      `json:"added,omitempty"`
	Removed   []string      `json:"removed,omitempty"`
	Modified  []string      `json:"modified,omitempty"`
}

// CommitAuthor represents the author or committer of a commit. The commit
// author may not correspond to a GitHub User.
type CommitAuthor struct {
	Date  *time.Time `json:"date,omitempty"`
	Name  *string    `json:"name,omitempty"`
	Email *string    `json:"email,omitempty"`

	// The following fields are only populated by Webhook events.
	Login *string `json:"username,omitempty"` // Renamed for go-github consistency.
}

type Timestamp struct {
	time.Time
}

// ReleaseEvent is triggered when a release is published.
// The Webhook event name is "release".
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#releaseevent
type ReleaseEvent struct {
	// Action is the action that was performed. Possible value is: "published".
	Action  *string            `json:"action,omitempty"`
	Release *RepositoryRelease `json:"release,omitempty"`

	// The following fields are only populated by Webhook events.
	Repo         *Repository   `json:"repository,omitempty"`
	Sender       *User         `json:"sender,omitempty"`
}

// RepositoryRelease represents a GitHub release in a repository.
type RepositoryRelease struct {
	ID              *int64         `json:"id,omitempty"`
	TagName         *string        `json:"tag_name,omitempty"`
	TargetCommitish *string        `json:"target_commitish,omitempty"`
	Name            *string        `json:"name,omitempty"`
	Body            *string        `json:"body,omitempty"`
	Draft           *bool          `json:"draft,omitempty"`
	Prerelease      *bool          `json:"prerelease,omitempty"`
	CreatedAt       *Timestamp     `json:"created_at,omitempty"`
	PublishedAt     *Timestamp     `json:"published_at,omitempty"`
	URL             *string        `json:"url,omitempty"`
	HTMLURL         *string        `json:"html_url,omitempty"`
	AssetsURL       *string        `json:"assets_url,omitempty"`
	UploadURL       *string        `json:"upload_url,omitempty"`
	ZipballURL      *string        `json:"zipball_url,omitempty"`
	TarballURL      *string        `json:"tarball_url,omitempty"`
	Author          *User          `json:"author,omitempty"`
	NodeID          *string        `json:"node_id,omitempty"`
}

// WatchEvent is related to starring a repository, not watching. See this API
// blog post for an explanation: https://developer.github.com/changes/2012-09-05-watcher-api/
//
// The event’s actor is the user who starred a repository, and the event’s
// repository is the repository that was starred.
//
// GitHub API docs: https://developer.github.com/v3/activity/events/types/#watchevent
type WatchEvent struct {
	// Action is the action that was performed. Possible value is: "started".
	Action *string `json:"action,omitempty"`

	// The following fields are only populated by Webhook events.
	Repo         *Repository   `json:"repository,omitempty"`
	Sender       *User         `json:"sender,omitempty"`
}
