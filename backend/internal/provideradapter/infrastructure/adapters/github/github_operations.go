package github

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	goGithub "github.com/google/go-github/v71/github"

	credDomain "github.com/context-space/context-space/backend/internal/credentialmanagement/domain"
)

// Define constants for operation IDs used by handlers.
const (
	operationIDGetMe                        = "get_me"
	operationIDGetIssue                     = "get_issue"
	operationIDGetIssueComments             = "get_issue_comments"
	operationIDUpdateIssue                  = "update_issue"
	operationIDCreateIssue                  = "create_issue"
	operationIDListRepositories             = "list_repositories"
	operationIDGetRepository                = "get_repository"
	operationIDCreateFile                   = "create_file"
	operationIDUpdateFile                   = "update_file"
	operationIDGetRepositoryContent         = "get_repository_content"
	operationIDListCommits                  = "list_commits"
	operationIDDeleteFile                   = "delete_file"
	operationIDGetCommit                    = "get_commit"
	operationIDListPullRequests             = "list_pull_requests"
	operationIDGetPullRequest               = "get_pull_request"
	operationIDCreatePullRequest            = "create_pull_request"
	operationIDListRepositoryIssues         = "list_repository_issues"
	operationIDUpdatePullRequest            = "update_pull_request"
	operationIDDeleteRef                    = "delete_ref"
	operationIDGetRef                       = "get_ref"
	operationIDCreateRef                    = "create_ref"
	operationIDMergePullRequest             = "merge_pull_request"
	operationIDCreatePullRequestReview      = "create_pull_request_review"
	operationIDListPullRequestReviews       = "list_pull_request_reviews"
	operationIDGetGitTree                   = "get_git_tree"
	operationIDGetBlob                      = "get_blob"
	operationIDCreateBlob                   = "create_blob"
	operationIDCreateRepositoryFromTemplate = "create_repository_from_template"
	operationIDStarRepository               = "star_repository"
	operationIDUnstarRepository             = "unstar_repository"
)

// Parameter schema structs for GitHub operations

// GetIssueParams defines parameters for getting an issue
type GetIssueParams struct {
	Owner       string `mapstructure:"owner" validate:"required"`
	Repo        string `mapstructure:"repo" validate:"required"`
	IssueNumber int    `mapstructure:"issue_number" validate:"required,gt=0"`
}

// GetIssueCommentsParams defines parameters for listing issue comments
type GetIssueCommentsParams struct {
	Owner       string `mapstructure:"owner" validate:"required"`
	Repo        string `mapstructure:"repo" validate:"required"`
	IssueNumber int    `mapstructure:"issue_number" validate:"required,gt=0"`
	Page        int    `mapstructure:"page" validate:"omitempty,gte=1" default:"1"`
	PerPage     int    `mapstructure:"per_page" validate:"omitempty,gte=1,lte=100" default:"30"`
}

// UpdateIssueParams defines parameters for updating an issue
type UpdateIssueParams struct {
	Owner       string   `mapstructure:"owner" validate:"required"`
	Repo        string   `mapstructure:"repo" validate:"required"`
	IssueNumber int      `mapstructure:"issue_number" validate:"required,gt=0"`
	Title       string   `mapstructure:"title" validate:"omitempty"`
	Body        string   `mapstructure:"body" validate:"omitempty"`
	State       string   `mapstructure:"state" validate:"omitempty,oneof=open closed"`
	Labels      []string `mapstructure:"labels" validate:"omitempty"`
	Assignees   []string `mapstructure:"assignees" validate:"omitempty"`
	Milestone   int      `mapstructure:"milestone" validate:"omitempty,gte=0"`
}

// CreateIssueParams defines parameters for creating an issue
type CreateIssueParams struct {
	Owner     string   `mapstructure:"owner" validate:"required"`
	Repo      string   `mapstructure:"repo" validate:"required"`
	Title     string   `mapstructure:"title" validate:"required"`
	Body      string   `mapstructure:"body" validate:"omitempty"`
	Labels    []string `mapstructure:"labels" validate:"omitempty"`
	Assignees []string `mapstructure:"assignees" validate:"omitempty"`
	Milestone *int     `mapstructure:"milestone" validate:"omitempty,gte=0"` // Pointer to handle optional 0 value vs. not provided
}

// ListRepositoriesParams defines parameters for listing repositories for the authenticated user.
type ListRepositoriesParams struct {
	Visibility  string `mapstructure:"visibility" validate:"omitempty,oneof=all public private"`         // Default: all
	Affiliation string `mapstructure:"affiliation" validate:"omitempty"`                                 // Comma-separated: owner,collaborator,organization_member. Default: owner,collaborator,organization_member
	Type        string `mapstructure:"type" validate:"omitempty,oneof=all owner public private member"`  // Default: owner
	Sort        string `mapstructure:"sort" validate:"omitempty,oneof=created updated pushed full_name"` // Default: created
	Direction   string `mapstructure:"direction" validate:"omitempty,oneof=asc desc"`                    // Default: desc
	Page        int    `mapstructure:"page" validate:"omitempty,gte=1" default:"1"`
	PerPage     int    `mapstructure:"per_page" validate:"omitempty,gte=1,lte=100" default:"30"`
}

// GetRepositoryParams defines parameters for getting a repository.
type GetRepositoryParams struct {
	Owner string `mapstructure:"owner" validate:"required"`
	Repo  string `mapstructure:"repo" validate:"required"`
}

// CreateFileParams defines parameters for creating a file.
type CreateFileParams struct {
	Owner          string `mapstructure:"owner" validate:"required"`
	Repo           string `mapstructure:"repo" validate:"required"`
	Path           string `mapstructure:"path" validate:"required"`    // Path to the file
	Message        string `mapstructure:"message" validate:"required"` // Commit message
	Content        string `mapstructure:"content" validate:"required"` // Raw file content as a string
	Branch         string `mapstructure:"branch" validate:"omitempty"` // Branch name. Default: repository's default branch
	CommitterName  string `mapstructure:"committer_name" validate:"omitempty"`
	CommitterEmail string `mapstructure:"committer_email" validate:"omitempty"`
	AuthorName     string `mapstructure:"author_name" validate:"omitempty"`
	AuthorEmail    string `mapstructure:"author_email" validate:"omitempty"`
}

// UpdateFileParams defines parameters for updating a file.
type UpdateFileParams struct {
	Owner          string `mapstructure:"owner" validate:"required"`
	Repo           string `mapstructure:"repo" validate:"required"`
	Path           string `mapstructure:"path" validate:"required"`    // Path to the file
	Message        string `mapstructure:"message" validate:"required"` // Commit message
	Content        string `mapstructure:"content" validate:"required"` // Raw file content as a string
	SHA            string `mapstructure:"sha" validate:"required"`     // Blob SHA of the file being replaced
	Branch         string `mapstructure:"branch" validate:"omitempty"` // Branch name. Default: repository's default branch
	CommitterName  string `mapstructure:"committer_name" validate:"omitempty"`
	CommitterEmail string `mapstructure:"committer_email" validate:"omitempty"`
	AuthorName     string `mapstructure:"author_name" validate:"omitempty"`
	AuthorEmail    string `mapstructure:"author_email" validate:"omitempty"`
}

// GetRepositoryContentParams defines parameters for getting repository content.
type GetRepositoryContentParams struct {
	Owner string `mapstructure:"owner" validate:"required"`
	Repo  string `mapstructure:"repo" validate:"required"`
	Path  string `mapstructure:"path" validate:"required"` // Path to the content
	Ref   string `mapstructure:"ref" validate:"omitempty"` // Branch, tag, or commit SHA. Default: repository's default branch
}

// ListCommitsParams defines parameters for listing commits.
type ListCommitsParams struct {
	Owner     string     `mapstructure:"owner" validate:"required"`
	Repo      string     `mapstructure:"repo" validate:"required"`
	SHA       string     `mapstructure:"sha" validate:"omitempty"`       // SHA or branch to start listing commits from. Default: the repository's default branch
	Path      string     `mapstructure:"path" validate:"omitempty"`      // Only commits containing this file path will be returned.
	Author    string     `mapstructure:"author" validate:"omitempty"`    // GitHub login or email address by which to filter commits.
	Committer string     `mapstructure:"committer" validate:"omitempty"` // GitHub login or email address by which to filter commits.
	Since     *time.Time `mapstructure:"since" validate:"omitempty"`     // Only show commits later than or on this date.
	Until     *time.Time `mapstructure:"until" validate:"omitempty"`     // Only show commits earlier than or on this date.
	Page      int        `mapstructure:"page" validate:"omitempty,gte=1" default:"1"`
	PerPage   int        `mapstructure:"per_page" validate:"omitempty,gte=1,lte=100" default:"30"`
}

// DeleteFileParams defines parameters for deleting a file.
type DeleteFileParams struct {
	Owner          string `mapstructure:"owner" validate:"required"`
	Repo           string `mapstructure:"repo" validate:"required"`
	Path           string `mapstructure:"path" validate:"required"`    // Path to the file
	Message        string `mapstructure:"message" validate:"required"` // Commit message
	SHA            string `mapstructure:"sha" validate:"required"`     // Blob SHA of the file being deleted
	Branch         string `mapstructure:"branch" validate:"omitempty"` // Branch name. Default: repository's default branch
	CommitterName  string `mapstructure:"committer_name" validate:"omitempty"`
	CommitterEmail string `mapstructure:"committer_email" validate:"omitempty"`
	AuthorName     string `mapstructure:"author_name" validate:"omitempty"`
	AuthorEmail    string `mapstructure:"author_email" validate:"omitempty"`
}

// GetCommitParams defines parameters for getting a specific commit.
type GetCommitParams struct {
	Owner string `mapstructure:"owner" validate:"required"`
	Repo  string `mapstructure:"repo" validate:"required"`
	Ref   string `mapstructure:"ref" validate:"required"` // SHA, branch or tag name
}

// ListPullRequestsParams defines parameters for listing pull requests.
type ListPullRequestsParams struct {
	Owner     string `mapstructure:"owner" validate:"required"`
	Repo      string `mapstructure:"repo" validate:"required"`
	State     string `mapstructure:"state" validate:"omitempty,oneof=open closed all" default:"open"`
	Head      string `mapstructure:"head" validate:"omitempty"` // Filter by head user and branch name (user:ref-name)
	Base      string `mapstructure:"base" validate:"omitempty"` // Filter by base branch name
	Sort      string `mapstructure:"sort" validate:"omitempty,oneof=created updated popularity long-running" default:"created"`
	Direction string `mapstructure:"direction" validate:"omitempty,oneof=asc desc" default:"desc"`
	Page      int    `mapstructure:"page" validate:"omitempty,gte=1" default:"1"`
	PerPage   int    `mapstructure:"per_page" validate:"omitempty,gte=1,lte=100" default:"30"`
}

// GetPullRequestParams defines parameters for getting a pull request.
type GetPullRequestParams struct {
	Owner  string `mapstructure:"owner" validate:"required"`
	Repo   string `mapstructure:"repo" validate:"required"`
	Number int    `mapstructure:"number" validate:"required,gt=0"`
}

// CreatePullRequestParams defines parameters for creating a pull request.
type CreatePullRequestParams struct {
	Owner               string `mapstructure:"owner" validate:"required"`
	Repo                string `mapstructure:"repo" validate:"required"`
	Title               string `mapstructure:"title" validate:"required"`
	Head                string `mapstructure:"head" validate:"required"` // The name of the branch where your changes are implemented.
	Base                string `mapstructure:"base" validate:"required"` // The name of the branch you want the changes pulled into.
	Body                string `mapstructure:"body" validate:"omitempty"`
	MaintainerCanModify *bool  `mapstructure:"maintainer_can_modify" validate:"omitempty"` // Use pointer for optional boolean
	Draft               *bool  `mapstructure:"draft" validate:"omitempty"`                 // Use pointer for optional boolean
	Issue               *int   `mapstructure:"issue" validate:"omitempty,gt=0"`            // Associate with an issue number
}

// ListRepositoryIssuesParams defines parameters for listing repository issues.
type ListRepositoryIssuesParams struct {
	Owner     string     `mapstructure:"owner" validate:"required"`
	Repo      string     `mapstructure:"repo" validate:"required"`
	Milestone string     `mapstructure:"milestone" validate:"omitempty"` // Milestone number, 'none', or '*'
	State     string     `mapstructure:"state" validate:"omitempty,oneof=open closed all" default:"open"`
	Assignee  string     `mapstructure:"assignee" validate:"omitempty"`  // Assignee login, 'none', or '*'
	Creator   string     `mapstructure:"creator" validate:"omitempty"`   // User login
	Mentioned string     `mapstructure:"mentioned" validate:"omitempty"` // User login
	Labels    string     `mapstructure:"labels" validate:"omitempty"`    // Comma-separated list of label names
	Sort      string     `mapstructure:"sort" validate:"omitempty,oneof=created updated comments" default:"created"`
	Direction string     `mapstructure:"direction" validate:"omitempty,oneof=asc desc" default:"desc"`
	Since     *time.Time `mapstructure:"since" validate:"omitempty"`
	Page      int        `mapstructure:"page" validate:"omitempty,gte=1" default:"1"`
	PerPage   int        `mapstructure:"per_page" validate:"omitempty,gte=1,lte=100" default:"30"`
}

// UpdatePullRequestParams defines parameters for updating a pull request.
type UpdatePullRequestParams struct {
	Owner               string `mapstructure:"owner" validate:"required"`
	Repo                string `mapstructure:"repo" validate:"required"`
	Number              int    `mapstructure:"number" validate:"required,gt=0"`
	Title               string `mapstructure:"title" validate:"omitempty"`
	Body                string `mapstructure:"body" validate:"omitempty"`
	State               string `mapstructure:"state" validate:"omitempty,oneof=open closed"`
	Base                string `mapstructure:"base" validate:"omitempty"` // Change the base branch for the PR
	MaintainerCanModify *bool  `mapstructure:"maintainer_can_modify" validate:"omitempty"`
}

// DeleteRefParams defines parameters for deleting a Git reference (branch or tag).
type DeleteRefParams struct {
	Owner string `mapstructure:"owner" validate:"required"`
	Repo  string `mapstructure:"repo" validate:"required"`
	Ref   string `mapstructure:"ref" validate:"required,startswith=refs/"` // Full ref path, e.g., "refs/heads/branch-name" or "refs/tags/tag-name"
}

// GetRefParams defines parameters for getting a Git reference.
type GetRefParams struct {
	Owner string `mapstructure:"owner" validate:"required"`
	Repo  string `mapstructure:"repo" validate:"required"`
	Ref   string `mapstructure:"ref" validate:"required"` // The ref to get. Needs to be fully qualified (e.g., "refs/heads/main" or just "heads/main").
}

// CreateRefParams defines parameters for creating a Git reference.
type CreateRefParams struct {
	Owner string `mapstructure:"owner" validate:"required"`
	Repo  string `mapstructure:"repo" validate:"required"`
	Ref   string `mapstructure:"ref" validate:"required,startswith=refs/"` // Fully qualified name of the ref, e.g., "refs/heads/new-branch"
	SHA   string `mapstructure:"sha" validate:"required"`                  // SHA1 value to set this reference to
}

// MergePullRequestParams defines parameters for merging a pull request.
type MergePullRequestParams struct {
	Owner          string `mapstructure:"owner" validate:"required"`
	Repo           string `mapstructure:"repo" validate:"required"`
	Number         int    `mapstructure:"number" validate:"required,gt=0"`
	CommitTitle    string `mapstructure:"commit_title" validate:"omitempty"`                           // Title for the merge commit. Default: PR title + PR number.
	CommitMessage  string `mapstructure:"commit_message" validate:"omitempty"`                         // Body of the merge commit. Default: PR body.
	SHA            string `mapstructure:"sha" validate:"omitempty"`                                    // SHA that pull request head must match to allow merge.
	MergeMethod    string `mapstructure:"merge_method" validate:"omitempty,oneof=merge squash rebase"` // Default: merge.
	Squash         *bool  `mapstructure:"squash" validate:"omitempty"`                                 // Deprecated: Use MergeMethod instead.
	MergeCommitSHA string `mapstructure:"merge_commit_sha" validate:"omitempty"`                       // Deprecated: Use SHA instead.
}

// CreatePullRequestReviewParams defines parameters for creating a pull request review.
type CreatePullRequestReviewParams struct {
	Owner    string                         `mapstructure:"owner" validate:"required"`
	Repo     string                         `mapstructure:"repo" validate:"required"`
	Number   int                            `mapstructure:"number" validate:"required,gt=0"`
	CommitID string                         `mapstructure:"commit_id" validate:"omitempty"`                                   // The SHA of the commit that needs a review. Defaults to the head commit of the pull request.
	Body     string                         `mapstructure:"body" validate:"omitempty"`                                        // Required if event is "COMMENT" or "REQUEST_CHANGES". The body text of the review.
	Event    string                         `mapstructure:"event" validate:"omitempty,oneof=APPROVE REQUEST_CHANGES COMMENT"` // The review action. Defaults to "PENDING" if body is provided.
	Comments []*goGithub.DraftReviewComment `mapstructure:"comments" validate:"omitempty"`                                    // Use DraftReviewComment for creating review comments. Note: complex struct, may need specific input format.
}

// ListPullRequestReviewsParams defines parameters for listing pull request reviews.
type ListPullRequestReviewsParams struct {
	Owner   string `mapstructure:"owner" validate:"required"`
	Repo    string `mapstructure:"repo" validate:"required"`
	Number  int    `mapstructure:"number" validate:"required,gt=0"`
	Page    int    `mapstructure:"page" validate:"omitempty,gte=1" default:"1"`
	PerPage int    `mapstructure:"per_page" validate:"omitempty,gte=1,lte=100" default:"30"`
}

// GetGitTreeParams defines parameters for getting a Git tree.
type GetGitTreeParams struct {
	Owner     string `mapstructure:"owner" validate:"required"`
	Repo      string `mapstructure:"repo" validate:"required"`
	TreeSHA   string `mapstructure:"tree_sha" validate:"required"`   // SHA of the tree object
	Recursive *bool  `mapstructure:"recursive" validate:"omitempty"` // If non-nil and true, recursively get subtrees. Note: pointer for optional bool.
}

// GetBlobParams defines parameters for getting a Git blob.
type GetBlobParams struct {
	Owner   string `mapstructure:"owner" validate:"required"`
	Repo    string `mapstructure:"repo" validate:"required"`
	FileSHA string `mapstructure:"file_sha" validate:"required"` // SHA of the blob object
}

// CreateBlobParams defines parameters for creating a Git blob.
type CreateBlobParams struct {
	Owner    string `mapstructure:"owner" validate:"required"`
	Repo     string `mapstructure:"repo" validate:"required"`
	Content  string `mapstructure:"content" validate:"required"`                   // Content of the blob
	Encoding string `mapstructure:"encoding" validate:"omitempty" default:"utf-8"` // Encoding of the content, "utf-8" or "base64". Defaults to "utf-8".
}

// CreateRepositoryFromTemplateParams defines parameters for creating a repository from a template.
type CreateRepositoryFromTemplateParams struct {
	TemplateOwner      string `mapstructure:"template_owner" validate:"required"` // Owner of the template repository
	TemplateRepo       string `mapstructure:"template_repo" validate:"required"`  // Name of the template repository
	Name               string `mapstructure:"name" validate:"required"`           // Name of the new repository
	Owner              string `mapstructure:"owner" validate:"omitempty"`         // Owner of the new repository (organization or authenticated user if blank)
	Description        string `mapstructure:"description" validate:"omitempty"`
	Private            *bool  `mapstructure:"private" validate:"omitempty"`              // Default: false
	IncludeAllBranches *bool  `mapstructure:"include_all_branches" validate:"omitempty"` // Default: false
}

// StarRepositoryParams defines parameters for starring a repository.
type StarRepositoryParams struct {
	Owner string `mapstructure:"owner" validate:"required"`
	Repo  string `mapstructure:"repo" validate:"required"`
}

// UnstarRepositoryParams defines parameters for unstarring a repository.
type UnstarRepositoryParams struct {
	Owner string `mapstructure:"owner" validate:"required"`
	Repo  string `mapstructure:"repo" validate:"required"`
}

// OperationHandler defines a GitHub operation that can be executed
type OperationHandler func(ctx context.Context, client *goGithub.Client, params interface{}) (interface{}, error)

// OperationDefinition combines the parameter schema and handler for an operation
type OperationDefinition struct {
	Schema                interface{}      // Parameter schema (struct pointer)
	Handler               OperationHandler // Operation handler function
	PermissionIdentifiers []string
}

// Operations maps operation IDs to their definitions
type Operations map[string]OperationDefinition

// RegisterOperation registers both the parameter schema and handler for an operation
func (a *GitHubAdapter) RegisterOperation(operationID string, schema interface{}, handler OperationHandler) {
	// Register the parameter schema with the base adapter
	a.BaseAdapter.RegisterOperation(operationID, schema)

	// Store both schema and handler in our operations map
	a.operations[operationID] = OperationDefinition{
		Schema:  schema,
		Handler: handler,
	}
}

func (a *GitHubAdapter) registerOperations() {
	a.RegisterOperation(
		operationIDGetMe,
		&struct{}{},
		handleGetMe,
	)

	a.RegisterOperation(
		operationIDGetIssue,
		&GetIssueParams{},
		handleGetIssue,
	)

	a.RegisterOperation(
		operationIDGetIssueComments,
		&GetIssueCommentsParams{},
		handleGetIssueComments,
	)

	a.RegisterOperation(
		operationIDUpdateIssue,
		&UpdateIssueParams{},
		handleUpdateIssue,
	)

	a.RegisterOperation(
		operationIDCreateIssue,
		&CreateIssueParams{},
		handleCreateIssue,
	)

	a.RegisterOperation(
		operationIDListRepositories,
		&ListRepositoriesParams{},
		handleListRepositories,
	)

	a.RegisterOperation(
		operationIDGetRepository,
		&GetRepositoryParams{},
		handleGetRepository,
	)

	a.RegisterOperation(
		operationIDCreateFile,
		&CreateFileParams{},
		handleCreateFile,
	)

	a.RegisterOperation(
		operationIDUpdateFile,
		&UpdateFileParams{},
		handleUpdateFile,
	)

	a.RegisterOperation(
		operationIDGetRepositoryContent,
		&GetRepositoryContentParams{},
		handleGetRepositoryContent,
	)

	a.RegisterOperation(
		operationIDListCommits,
		&ListCommitsParams{},
		handleListCommits,
	)

	a.RegisterOperation(
		operationIDDeleteFile,
		&DeleteFileParams{},
		handleDeleteFile,
	)

	a.RegisterOperation(
		operationIDGetCommit,
		&GetCommitParams{},
		handleGetCommit,
	)

	a.RegisterOperation(
		operationIDListPullRequests,
		&ListPullRequestsParams{},
		handleListPullRequests,
	)

	a.RegisterOperation(
		operationIDGetPullRequest,
		&GetPullRequestParams{},
		handleGetPullRequest,
	)

	a.RegisterOperation(
		operationIDCreatePullRequest,
		&CreatePullRequestParams{},
		handleCreatePullRequest,
	)

	a.RegisterOperation(
		operationIDListRepositoryIssues,
		&ListRepositoryIssuesParams{},
		handleListRepositoryIssues,
	)

	a.RegisterOperation(
		operationIDUpdatePullRequest,
		&UpdatePullRequestParams{},
		handleUpdatePullRequest,
	)

	a.RegisterOperation(
		operationIDDeleteRef,
		&DeleteRefParams{},
		handleDeleteRef,
	)

	a.RegisterOperation(
		operationIDGetRef,
		&GetRefParams{},
		handleGetRef,
	)

	a.RegisterOperation(
		operationIDCreateRef,
		&CreateRefParams{},
		handleCreateRef,
	)

	a.RegisterOperation(
		operationIDMergePullRequest,
		&MergePullRequestParams{},
		handleMergePullRequest,
	)

	a.RegisterOperation(
		operationIDCreatePullRequestReview,
		&CreatePullRequestReviewParams{},
		handleCreatePullRequestReview,
	)

	a.RegisterOperation(
		operationIDListPullRequestReviews,
		&ListPullRequestReviewsParams{},
		handleListPullRequestReviews,
	)

	a.RegisterOperation(
		operationIDGetGitTree,
		&GetGitTreeParams{},
		handleGetGitTree,
	)

	a.RegisterOperation(
		operationIDGetBlob,
		&GetBlobParams{},
		handleGetBlob,
	)

	a.RegisterOperation(
		operationIDCreateBlob,
		&CreateBlobParams{},
		handleCreateBlob,
	)

	a.RegisterOperation(
		operationIDCreateRepositoryFromTemplate,
		&CreateRepositoryFromTemplateParams{},
		handleCreateRepositoryFromTemplate,
	)

	a.RegisterOperation(
		operationIDStarRepository,
		&StarRepositoryParams{},
		handleStarRepository,
	)

	a.RegisterOperation(
		operationIDUnstarRepository,
		&UnstarRepositoryParams{},
		handleUnstarRepository,
	)
}

func handleGetMe(ctx context.Context, client *goGithub.Client, _ interface{}) (interface{}, error) {
	user, resp, err := client.Users.Get(ctx, "")
	if err := handleGitHubResponse(resp, err, "get user"); err != nil {
		return nil, err
	}
	return user, nil
}

func handleGetIssue(ctx context.Context, client *goGithub.Client, params interface{}) (interface{}, error) {
	p := params.(*GetIssueParams)
	issue, resp, err := client.Issues.Get(ctx, p.Owner, p.Repo, p.IssueNumber)
	if err := handleGitHubResponse(resp, err, "get issue"); err != nil {
		return nil, err
	}
	return issue, nil
}

func handleGetIssueComments(ctx context.Context, client *goGithub.Client, params interface{}) (interface{}, error) {
	p := params.(*GetIssueCommentsParams)
	opts := &goGithub.IssueListCommentsOptions{
		ListOptions: goGithub.ListOptions{
			Page:    p.Page,
			PerPage: p.PerPage,
		},
	}

	comments, resp, err := client.Issues.ListComments(ctx, p.Owner, p.Repo, p.IssueNumber, opts)
	if err := handleGitHubResponse(resp, err, "get issue comments"); err != nil {
		return nil, err
	}
	return comments, nil
}

func handleUpdateIssue(ctx context.Context, client *goGithub.Client, params interface{}) (interface{}, error) {
	p := params.(*UpdateIssueParams)

	// Create the issue request with only provided fields
	issueRequest := &goGithub.IssueRequest{}

	// Set optional parameters if provided
	if p.Title != "" {
		issueRequest.Title = &p.Title
	}

	if p.Body != "" {
		issueRequest.Body = &p.Body
	}

	if p.State != "" {
		issueRequest.State = &p.State
	}

	// Get labels
	if len(p.Labels) > 0 {
		issueRequest.Labels = &p.Labels
	}

	// Get assignees
	if len(p.Assignees) > 0 {
		issueRequest.Assignees = &p.Assignees
	}

	// Get milestone
	if p.Milestone > 0 {
		issueRequest.Milestone = &p.Milestone
	}

	updatedIssue, resp, err := client.Issues.Edit(ctx, p.Owner, p.Repo, p.IssueNumber, issueRequest)
	if err := handleGitHubResponse(resp, err, "update issue"); err != nil {
		return nil, err
	}
	return updatedIssue, nil
}

func handleCreateIssue(ctx context.Context, client *goGithub.Client, params interface{}) (interface{}, error) {
	p := params.(*CreateIssueParams)
	issueRequest := &goGithub.IssueRequest{
		Title: &p.Title,
	}
	// Set optional parameters if provided
	if p.Body != "" {
		issueRequest.Body = &p.Body
	}
	if len(p.Labels) > 0 {
		issueRequest.Labels = &p.Labels
	}
	if len(p.Assignees) > 0 {
		issueRequest.Assignees = &p.Assignees
	}
	if p.Milestone != nil {
		issueRequest.Milestone = p.Milestone // Assign pointer directly
	}

	createdIssue, resp, err := client.Issues.Create(ctx, p.Owner, p.Repo, issueRequest)
	if err := handleGitHubResponse(resp, err, "create issue"); err != nil {
		return nil, err
	}
	return createdIssue, nil
}

// createGitHubClient creates a new GitHub client from the given credentials
func (a *GitHubAdapter) createGitHubClient(credential interface{}) (*goGithub.Client, error) {
	oauthCred, ok := credential.(*credDomain.OAuthCredential)
	if !ok {
		return nil, fmt.Errorf("invalid credential type for GitHub")
	}

	return goGithub.NewClient(nil).WithAuthToken(oauthCred.Token.AccessToken), nil
}

func handleListRepositories(ctx context.Context, client *goGithub.Client, params interface{}) (interface{}, error) {
	p := params.(*ListRepositoriesParams)
	// Use the non-deprecated ListByAuthenticatedUser method
	opts := &goGithub.RepositoryListByAuthenticatedUserOptions{
		Visibility:  p.Visibility,
		Affiliation: p.Affiliation,
		Type:        p.Type,
		Sort:        p.Sort,
		Direction:   p.Direction,
		ListOptions: goGithub.ListOptions{
			Page:    p.Page,
			PerPage: p.PerPage,
		},
	}
	// Call ListByAuthenticatedUser directly
	repos, resp, err := client.Repositories.ListByAuthenticatedUser(ctx, opts)
	if err := handleGitHubResponse(resp, err, "list authenticated user repositories"); err != nil {
		return nil, err
	}
	return repos, nil
}

func handleGetRepository(ctx context.Context, client *goGithub.Client, params interface{}) (interface{}, error) {
	p := params.(*GetRepositoryParams)
	repo, resp, err := client.Repositories.Get(ctx, p.Owner, p.Repo)
	if err := handleGitHubResponse(resp, err, "get repository"); err != nil {
		return nil, err
	}
	return repo, nil
}

func handleCreateFile(ctx context.Context, client *goGithub.Client, params interface{}) (interface{}, error) {
	p := params.(*CreateFileParams)
	opts := &goGithub.RepositoryContentFileOptions{
		Message: &p.Message,
		Content: []byte(p.Content), // Pass raw content bytes; go-github handles base64 encoding
		// Branch:  &p.Branch, // Pass nil if empty to use default branch
	}
	if p.Branch != "" {
		opts.Branch = &p.Branch
	}
	if p.CommitterName != "" || p.CommitterEmail != "" {
		opts.Committer = &goGithub.CommitAuthor{Name: &p.CommitterName, Email: &p.CommitterEmail}
	}
	if p.AuthorName != "" || p.AuthorEmail != "" {
		opts.Author = &goGithub.CommitAuthor{Name: &p.AuthorName, Email: &p.AuthorEmail}
	}
	// NOTE: Content needs to be base64 decoded before passing to API if user provides it encoded.

	contentResponse, resp, err := client.Repositories.CreateFile(ctx, p.Owner, p.Repo, p.Path, opts)
	if err := handleGitHubResponse(resp, err, "create file"); err != nil {
		return nil, err
	}
	return contentResponse, nil
}

func handleUpdateFile(ctx context.Context, client *goGithub.Client, params interface{}) (interface{}, error) {
	p := params.(*UpdateFileParams)
	opts := &goGithub.RepositoryContentFileOptions{
		Message: &p.Message,
		Content: []byte(p.Content), // Pass raw content bytes; go-github handles base64 encoding
		SHA:     &p.SHA,
		// Branch:  &p.Branch, // Pass nil if empty to use default branch
	}
	if p.Branch != "" {
		opts.Branch = &p.Branch
	}
	if p.CommitterName != "" || p.CommitterEmail != "" {
		opts.Committer = &goGithub.CommitAuthor{Name: &p.CommitterName, Email: &p.CommitterEmail}
	}
	if p.AuthorName != "" || p.AuthorEmail != "" {
		opts.Author = &goGithub.CommitAuthor{Name: &p.AuthorName, Email: &p.AuthorEmail}
	}

	contentResponse, resp, err := client.Repositories.UpdateFile(ctx, p.Owner, p.Repo, p.Path, opts)
	if err := handleGitHubResponse(resp, err, "update file"); err != nil {
		return nil, err
	}
	return contentResponse, nil
}

func handleGetRepositoryContent(ctx context.Context, client *goGithub.Client, params interface{}) (interface{}, error) {
	p := params.(*GetRepositoryContentParams)
	opts := &goGithub.RepositoryContentGetOptions{
		Ref: p.Ref,
	}
	// GetContent can return *RepositoryContent (for a file) or []*RepositoryContent (for a dir)
	fileContent, directoryContent, resp, err := client.Repositories.GetContents(ctx, p.Owner, p.Repo, p.Path, opts)
	if err := handleGitHubResponse(resp, err, "get repository content"); err != nil {
		return nil, err
	}
	// Return whichever is not nil
	if fileContent != nil {
		return fileContent, nil
	}
	return directoryContent, nil
}

func handleListCommits(ctx context.Context, client *goGithub.Client, params interface{}) (interface{}, error) {
	p := params.(*ListCommitsParams)
	opts := &goGithub.CommitsListOptions{
		SHA:    p.SHA,
		Path:   p.Path,
		Author: p.Author,
		ListOptions: goGithub.ListOptions{
			Page:    p.Page,
			PerPage: p.PerPage,
		},
	}
	// Assign Since and Until only if they are not nil
	if p.Since != nil {
		opts.Since = *p.Since
	}
	if p.Until != nil {
		opts.Until = *p.Until
	}

	commits, resp, err := client.Repositories.ListCommits(ctx, p.Owner, p.Repo, opts)
	if err := handleGitHubResponse(resp, err, "list commits"); err != nil {
		return nil, err
	}
	return commits, nil
}

func handleDeleteFile(ctx context.Context, client *goGithub.Client, params interface{}) (interface{}, error) {
	p := params.(*DeleteFileParams)
	opts := &goGithub.RepositoryContentFileOptions{
		Message: &p.Message,
		SHA:     &p.SHA,
		// Branch:  &p.Branch, // Pass nil if empty to use default branch
	}
	if p.Branch != "" {
		opts.Branch = &p.Branch
	}
	if p.CommitterName != "" || p.CommitterEmail != "" {
		opts.Committer = &goGithub.CommitAuthor{Name: &p.CommitterName, Email: &p.CommitterEmail}
	}
	if p.AuthorName != "" || p.AuthorEmail != "" {
		opts.Author = &goGithub.CommitAuthor{Name: &p.AuthorName, Email: &p.AuthorEmail}
	}

	deleteResponse, resp, err := client.Repositories.DeleteFile(ctx, p.Owner, p.Repo, p.Path, opts)
	if err := handleGitHubResponse(resp, err, "delete file"); err != nil {
		// Specifically handle '404 Not Found' during delete as non-fatal for cleanup purposes if needed,
		// but for a direct call, it's likely an error (trying to delete non-existent file with wrong SHA).
		// For now, treat all errors as fatal.
		return nil, err
	}
	return deleteResponse, nil
}

func handleGetCommit(ctx context.Context, client *goGithub.Client, params interface{}) (interface{}, error) {
	p := params.(*GetCommitParams)
	// Use Repositories.GetCommit for richer information including files changed
	commit, resp, err := client.Repositories.GetCommit(ctx, p.Owner, p.Repo, p.Ref, nil) // No specific options needed here for ListOptions
	if err := handleGitHubResponse(resp, err, "get commit"); err != nil {
		return nil, err
	}
	return commit, nil
}

func handleListPullRequests(ctx context.Context, client *goGithub.Client, params interface{}) (interface{}, error) {
	p := params.(*ListPullRequestsParams)
	opts := &goGithub.PullRequestListOptions{
		State:     p.State,
		Head:      p.Head,
		Base:      p.Base,
		Sort:      p.Sort,
		Direction: p.Direction,
		ListOptions: goGithub.ListOptions{
			Page:    p.Page,
			PerPage: p.PerPage,
		},
	}
	prs, resp, err := client.PullRequests.List(ctx, p.Owner, p.Repo, opts)
	if err := handleGitHubResponse(resp, err, "list pull requests"); err != nil {
		return nil, err
	}
	return prs, nil
}

func handleGetPullRequest(ctx context.Context, client *goGithub.Client, params interface{}) (interface{}, error) {
	p := params.(*GetPullRequestParams)
	pr, resp, err := client.PullRequests.Get(ctx, p.Owner, p.Repo, p.Number)
	if err := handleGitHubResponse(resp, err, "get pull request"); err != nil {
		return nil, err
	}
	return pr, nil
}

func handleCreatePullRequest(ctx context.Context, client *goGithub.Client, params interface{}) (interface{}, error) {
	p := params.(*CreatePullRequestParams)
	newPR := &goGithub.NewPullRequest{
		Title:               &p.Title,
		Head:                &p.Head,
		Base:                &p.Base,
		Body:                &p.Body, // Will be nil if empty string
		MaintainerCanModify: p.MaintainerCanModify,
		Draft:               p.Draft,
		Issue:               p.Issue,
	}
	// Note: Body, MaintainerCanModify, Draft, Issue are optional and handled correctly if nil/empty

	pr, resp, err := client.PullRequests.Create(ctx, p.Owner, p.Repo, newPR)
	if err := handleGitHubResponse(resp, err, "create pull request"); err != nil {
		return nil, err
	}
	return pr, nil
}

func handleListRepositoryIssues(ctx context.Context, client *goGithub.Client, params interface{}) (interface{}, error) {
	p := params.(*ListRepositoryIssuesParams)
	opts := &goGithub.IssueListByRepoOptions{
		Milestone: p.Milestone,
		State:     p.State,
		Assignee:  p.Assignee,
		Creator:   p.Creator,
		Mentioned: p.Mentioned,
		Labels:    strToSlice(p.Labels), // Helper function needed for comma-separated string
		Sort:      p.Sort,
		Direction: p.Direction,
		ListOptions: goGithub.ListOptions{
			Page:    p.Page,
			PerPage: p.PerPage,
		},
	}
	// Assign Since only if it's not nil, dereferencing the pointer
	if p.Since != nil {
		opts.Since = *p.Since
	}
	issues, resp, err := client.Issues.ListByRepo(ctx, p.Owner, p.Repo, opts)
	if err := handleGitHubResponse(resp, err, "list repository issues"); err != nil {
		return nil, err
	}
	return issues, nil
}

func handleUpdatePullRequest(ctx context.Context, client *goGithub.Client, params interface{}) (interface{}, error) {
	p := params.(*UpdatePullRequestParams)
	updateReq := &goGithub.PullRequest{
		// Only set fields if they are provided in the params
	}
	if p.Title != "" {
		updateReq.Title = &p.Title
	}
	if p.Body != "" {
		updateReq.Body = &p.Body
	}
	if p.State != "" {
		updateReq.State = &p.State
	}
	if p.Base != "" {
		updateReq.Base = &goGithub.PullRequestBranch{Ref: &p.Base}
	}
	if p.MaintainerCanModify != nil {
		updateReq.MaintainerCanModify = p.MaintainerCanModify
	}

	// Note: client.PullRequests.Edit takes number and the PullRequest object containing updates.
	// It seems counter-intuitive, but we pass the updates within the PullRequest struct itself.
	pr, resp, err := client.PullRequests.Edit(ctx, p.Owner, p.Repo, p.Number, updateReq)
	if err := handleGitHubResponse(resp, err, "update pull request"); err != nil {
		return nil, err
	}
	return pr, nil
}

func handleDeleteRef(ctx context.Context, client *goGithub.Client, params interface{}) (interface{}, error) {
	p := params.(*DeleteRefParams)
	resp, err := client.Git.DeleteRef(ctx, p.Owner, p.Repo, p.Ref)
	if err := handleGitHubResponse(resp, err, "delete ref"); err != nil {
		// Handle cases where the ref doesn't exist gracefully? For cleanup, maybe non-fatal.
		// However, for a direct call, non-existence is likely an error.
		// Let's treat it as an error for now.
		return nil, err
	}
	// DeleteRef returns only a response and error, no specific data structure on success.
	// Return a simple success indicator or nil.
	return map[string]interface{}{"success": true, "status": resp.Status}, nil
}

func handleGetRef(ctx context.Context, client *goGithub.Client, params interface{}) (interface{}, error) {
	p := params.(*GetRefParams)
	// The go-github library expects the ref path *without* the leading "refs/"
	// but the API docs often show it with it. Let's remove it if present.
	refPath := strings.TrimPrefix(p.Ref, "refs/")

	ref, resp, err := client.Git.GetRef(ctx, p.Owner, p.Repo, refPath)
	if err := handleGitHubResponse(resp, err, "get ref"); err != nil {
		return nil, err
	}
	return ref, nil
}

func handleCreateRef(ctx context.Context, client *goGithub.Client, params interface{}) (interface{}, error) {
	p := params.(*CreateRefParams)
	refInput := &goGithub.Reference{
		Ref: &p.Ref,
		Object: &goGithub.GitObject{
			SHA: &p.SHA,
		},
	}
	ref, resp, err := client.Git.CreateRef(ctx, p.Owner, p.Repo, refInput)
	if err := handleGitHubResponse(resp, err, "create ref"); err != nil {
		return nil, err
	}
	return ref, nil
}

func handleMergePullRequest(ctx context.Context, client *goGithub.Client, params interface{}) (interface{}, error) {
	p := params.(*MergePullRequestParams)
	opts := &goGithub.PullRequestOptions{
		SHA:         p.SHA,
		MergeMethod: p.MergeMethod,
		// Squash      p.Squash, // Deprecated
		// CommitTitle p.CommitTitle // These seem to be handled by the commitMessage param only now
		// MergeCommitSHA p.MergeCommitSHA // Deprecated
	}
	// Commit message parameter handles both title and message based on content
	mergeResult, resp, err := client.PullRequests.Merge(ctx, p.Owner, p.Repo, p.Number, p.CommitMessage, opts)
	if err := handleGitHubResponse(resp, err, "merge pull request"); err != nil {
		return nil, err
	}
	return mergeResult, nil
}

func handleCreatePullRequestReview(ctx context.Context, client *goGithub.Client, params interface{}) (interface{}, error) {
	p := params.(*CreatePullRequestReviewParams)
	reviewRequest := &goGithub.PullRequestReviewRequest{
		CommitID: &p.CommitID,
		Body:     &p.Body,
		Event:    &p.Event,
		Comments: p.Comments, // Pass the slice directly
	}
	// Set CommitID to nil if empty, as the API expects omitted or null for default head
	if p.CommitID == "" {
		reviewRequest.CommitID = nil
	}
	// Body is required for COMMENT or REQUEST_CHANGES
	if (p.Event == "COMMENT" || p.Event == "REQUEST_CHANGES") && p.Body == "" {
		return nil, fmt.Errorf("body is required for event type %s", p.Event)
	}
	// Event defaults to PENDING if body is set and event is empty, but API requires explicit event for APPROVE/REQUEST_CHANGES
	if p.Event == "" && p.Body == "" {
		// Cannot create an empty PENDING review via API? Library might default. Let's require an event or body.
		return nil, fmt.Errorf("either event or body must be provided to create a review")
	}

	review, resp, err := client.PullRequests.CreateReview(ctx, p.Owner, p.Repo, p.Number, reviewRequest)
	if err := handleGitHubResponse(resp, err, "create pull request review"); err != nil {
		return nil, err
	}
	return review, nil
}

func handleListPullRequestReviews(ctx context.Context, client *goGithub.Client, params interface{}) (interface{}, error) {
	p := params.(*ListPullRequestReviewsParams)
	opts := &goGithub.ListOptions{
		Page:    p.Page,
		PerPage: p.PerPage,
	}
	reviews, resp, err := client.PullRequests.ListReviews(ctx, p.Owner, p.Repo, p.Number, opts)
	if err := handleGitHubResponse(resp, err, "list pull request reviews"); err != nil {
		return nil, err
	}
	return reviews, nil
}

func handleGetGitTree(ctx context.Context, client *goGithub.Client, params interface{}) (interface{}, error) {
	p := params.(*GetGitTreeParams)
	recursive := false // Default value
	if p.Recursive != nil {
		recursive = *p.Recursive
	}
	tree, resp, err := client.Git.GetTree(ctx, p.Owner, p.Repo, p.TreeSHA, recursive)
	if err := handleGitHubResponse(resp, err, "get git tree"); err != nil {
		return nil, err
	}
	return tree, nil
}

func handleGetBlob(ctx context.Context, client *goGithub.Client, params interface{}) (interface{}, error) {
	p := params.(*GetBlobParams)
	blob, resp, err := client.Git.GetBlob(ctx, p.Owner, p.Repo, p.FileSHA)
	if err := handleGitHubResponse(resp, err, "get blob"); err != nil {
		return nil, err
	}
	return blob, nil
}
func handleCreateBlob(ctx context.Context, client *goGithub.Client, params interface{}) (interface{}, error) {
	p := params.(*CreateBlobParams)
	newBlob := &goGithub.Blob{
		Content:  &p.Content,
		Encoding: &p.Encoding,
	}
	blob, resp, err := client.Git.CreateBlob(ctx, p.Owner, p.Repo, newBlob)
	if err := handleGitHubResponse(resp, err, "create blob"); err != nil {
		return nil, err
	}
	return blob, nil
}

func handleCreateRepositoryFromTemplate(ctx context.Context, client *goGithub.Client, params interface{}) (interface{}, error) {
	p := params.(*CreateRepositoryFromTemplateParams)
	templateRepoReq := &goGithub.TemplateRepoRequest{
		Name:               &p.Name,
		Owner:              &p.Owner,             // Can be empty, defaults to authenticated user
		Description:        &p.Description,       // Can be empty
		Private:            p.Private,            // Pass pointer directly
		IncludeAllBranches: p.IncludeAllBranches, // Pass pointer directly
	}
	// Owner needs to be nil if empty string was passed, otherwise API fails
	if p.Owner == "" {
		templateRepoReq.Owner = nil
	}
	if p.Description == "" {
		templateRepoReq.Description = nil
	}

	repo, resp, err := client.Repositories.CreateFromTemplate(ctx, p.TemplateOwner, p.TemplateRepo, templateRepoReq)
	if err := handleGitHubResponse(resp, err, "create repository from template"); err != nil {
		return nil, err
	}
	return repo, nil
}
func handleStarRepository(ctx context.Context, client *goGithub.Client, params interface{}) (interface{}, error) {
	p := params.(*StarRepositoryParams)
	resp, err := client.Activity.Star(ctx, p.Owner, p.Repo)
	if err := handleGitHubResponse(resp, err, "star repository"); err != nil {
		return nil, err
	}
	return map[string]interface{}{"success": true, "status": resp.Status}, nil
}

func handleUnstarRepository(ctx context.Context, client *goGithub.Client, params interface{}) (interface{}, error) {
	p := params.(*UnstarRepositoryParams)
	resp, err := client.Activity.Unstar(ctx, p.Owner, p.Repo)
	// Handle 204 No Content as success
	if err != nil && (resp == nil || resp.StatusCode != http.StatusNoContent) {
		return nil, fmt.Errorf("failed to unstar repository: %w (status: %d)", err, resp.StatusCode)
	}
	// Special handling for 404 Not Found (already unstarred or never starred) - consider this success?
	// For now, treat 404 as an error unless specifically requested otherwise.
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("failed to unstar repository: %s (status code %d)", "Not Found - repository may not exist or was not starred", resp.StatusCode)
	}
	if resp.StatusCode >= http.StatusMultipleChoices && resp.StatusCode != http.StatusNoContent {
		return nil, fmt.Errorf("failed to unstar repository: %s (status code %d)", resp.Status, resp.StatusCode)
	}

	return map[string]interface{}{"success": true, "status": resp.Status}, nil
}

// handleGitHubResponse handles common GitHub API response processing
func handleGitHubResponse(resp *goGithub.Response, err error, operation string) error {
	if err != nil {
		return fmt.Errorf("failed to %s: %w", operation, err)
	}

	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("failed to %s: %s", operation, resp.Status)
	}

	return nil
}

// Helper function to convert comma-separated string to slice of strings
// Handles empty strings and trims whitespace.
func strToSlice(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	if len(result) == 0 {
		return nil // Return nil if all parts were empty/whitespace
	}
	return result
}
