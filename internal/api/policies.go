// Copyright 2025 Interlynk.io
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"context"
	"time"

	"github.com/interlynk-io/lynk-mcp/internal/graphql"
)

// Policy represents a security policy
type Policy struct {
	ID          string
	Name        string
	Description string
	Enabled     bool
	ResultType  string
	UpdatedAt   time.Time
	PolicyRules []PolicyRule
}

// PolicyRule represents a rule within a policy
type PolicyRule struct {
	ID          string
	Name        string
	Subject     string
	Operator    string
	Value       string
	Enabled     bool
	FailMessage string
}

// PolicyResult represents a policy evaluation result
type PolicyResult struct {
	ID         string
	PolicyID   string
	VersionID  string
	ResultType string
	Result     string
	CreatedAt  time.Time
	Policy     *Policy
	Version    *Version
}

// PoliciesResult represents the result of listing policies
type PoliciesResult struct {
	Policies    []Policy
	TotalCount  int
	HasNextPage bool
	EndCursor   string
}

// PolicyResultsResult represents the result of listing policy results
type PolicyResultsResult struct {
	PolicyResults []PolicyResult
	TotalCount    int
	HasNextPage   bool
	EndCursor     string
}

// OrganizationLicense represents a license in the organization
type OrganizationLicense struct {
	ShortID            string
	Name               string
	DerivedState       string
	CopyLeft           string
	OsiApproved        bool
	FsfLibre           bool
	Deprecated         bool
	Attribution        string
	SourceDistribution string
	Modifications      string
}

// LicensesResult represents the result of listing licenses
type LicensesResult struct {
	Licenses    []OrganizationLicense
	TotalCount  int
	HasNextPage bool
	EndCursor   string
}

// ListPoliciesInput contains parameters for listing policies
type ListPoliciesInput struct {
	First  int
	After  string
	Search string
}

// ListPolicies fetches policies with pagination
func (c *Client) ListPolicies(ctx context.Context, input ListPoliciesInput) (*PoliciesResult, error) {
	vars := make(map[string]interface{})
	if input.First > 0 {
		vars["first"] = input.First
	} else {
		vars["first"] = 20
	}
	if input.After != "" {
		vars["after"] = input.After
	}
	if input.Search != "" {
		vars["search"] = input.Search
	}

	var result struct {
		Policies struct {
			Nodes []struct {
				ID          string    `json:"id"`
				Name        string    `json:"name"`
				Description string    `json:"description"`
				Enabled     bool      `json:"isEnabled"`
				ResultType  string    `json:"resultType"`
				UpdatedAt   time.Time `json:"updatedAt"`
			} `json:"nodes"`
			TotalCount int `json:"totalCount"`
			PageInfo   struct {
				HasNextPage bool   `json:"hasNextPage"`
				EndCursor   string `json:"endCursor"`
			} `json:"pageInfo"`
		} `json:"policies"`
	}

	if err := c.gql.Execute(ctx, graphql.PoliciesQuery, vars, &result); err != nil {
		return nil, err
	}

	policies := make([]Policy, len(result.Policies.Nodes))
	for i, n := range result.Policies.Nodes {
		policies[i] = Policy{
			ID:          n.ID,
			Name:        n.Name,
			Description: n.Description,
			Enabled:     n.Enabled,
			ResultType:  n.ResultType,
			UpdatedAt:   n.UpdatedAt,
		}
	}

	return &PoliciesResult{
		Policies:    policies,
		TotalCount:  result.Policies.TotalCount,
		HasNextPage: result.Policies.PageInfo.HasNextPage,
		EndCursor:   result.Policies.PageInfo.EndCursor,
	}, nil
}

// GetPolicy fetches a single policy by ID with its rules
func (c *Client) GetPolicy(ctx context.Context, id string) (*Policy, error) {
	vars := map[string]interface{}{
		"id": id,
	}

	var result struct {
		Policy struct {
			ID          string    `json:"id"`
			Name        string    `json:"name"`
			Description string    `json:"description"`
			Enabled     bool      `json:"isEnabled"`
			ResultType  string    `json:"resultType"`
			UpdatedAt   time.Time `json:"updatedAt"`
			PolicyRules []struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				Subject     string `json:"subject"`
				Operator    string `json:"operator"`
				Value       string `json:"value"`
				Enabled     bool   `json:"enabled"`
				FailMessage string `json:"failMessage"`
			} `json:"policyRules"`
		} `json:"policy"`
	}

	if err := c.gql.Execute(ctx, graphql.PolicyQuery, vars, &result); err != nil {
		return nil, err
	}

	rules := make([]PolicyRule, len(result.Policy.PolicyRules))
	for i, r := range result.Policy.PolicyRules {
		rules[i] = PolicyRule{
			ID:          r.ID,
			Name:        r.Name,
			Subject:     r.Subject,
			Operator:    r.Operator,
			Value:       r.Value,
			Enabled:     r.Enabled,
			FailMessage: r.FailMessage,
		}
	}

	return &Policy{
		ID:          result.Policy.ID,
		Name:        result.Policy.Name,
		Description: result.Policy.Description,
		Enabled:     result.Policy.Enabled,
		ResultType:  result.Policy.ResultType,
		UpdatedAt:   result.Policy.UpdatedAt,
		PolicyRules: rules,
	}, nil
}

// ListPolicyResultsInput contains parameters for listing policy results
type ListPolicyResultsInput struct {
	First      int
	After      string
	PolicyID   string
	VersionID  string
	ResultType string
}

// ListPolicyResults fetches policy evaluation results
func (c *Client) ListPolicyResults(ctx context.Context, input ListPolicyResultsInput) (*PolicyResultsResult, error) {
	vars := make(map[string]interface{})
	if input.First > 0 {
		vars["first"] = input.First
	} else {
		vars["first"] = 50
	}
	if input.After != "" {
		vars["after"] = input.After
	}
	if input.PolicyID != "" {
		vars["policyId"] = []string{input.PolicyID}
	}
	if input.VersionID != "" {
		vars["sbomId"] = []string{input.VersionID}
	}
	if input.ResultType != "" {
		vars["resultType"] = []string{input.ResultType}
	}

	var result struct {
		PolicyResults struct {
			Nodes []struct {
				ID         string    `json:"id"`
				PolicyID   string    `json:"policyId"`
				SbomID     string    `json:"sbomId"`
				ResultType string    `json:"resultType"`
				Result     string    `json:"result"`
				CreatedAt  time.Time `json:"createdAt"`
				Policy     *struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				} `json:"policy"`
				Sbom *struct {
					ID             string `json:"id"`
					ProjectVersion string `json:"projectVersion"`
					Project        struct {
						ID   string `json:"id"`
						Name string `json:"name"`
					} `json:"project"`
				} `json:"sbom"`
			} `json:"nodes"`
			TotalCount int `json:"totalCount"`
			PageInfo   struct {
				HasNextPage bool   `json:"hasNextPage"`
				EndCursor   string `json:"endCursor"`
			} `json:"pageInfo"`
		} `json:"policyResults"`
	}

	if err := c.gql.Execute(ctx, graphql.PolicyResultsQuery, vars, &result); err != nil {
		return nil, err
	}

	results := make([]PolicyResult, len(result.PolicyResults.Nodes))
	for i, n := range result.PolicyResults.Nodes {
		pr := PolicyResult{
			ID:         n.ID,
			PolicyID:   n.PolicyID,
			VersionID:  n.SbomID,
			ResultType: n.ResultType,
			Result:     n.Result,
			CreatedAt:  n.CreatedAt,
		}
		if n.Policy != nil {
			pr.Policy = &Policy{
				ID:   n.Policy.ID,
				Name: n.Policy.Name,
			}
		}
		if n.Sbom != nil {
			pr.Version = &Version{
				ID:      n.Sbom.ID,
				Version: n.Sbom.ProjectVersion,
				Environment: &Environment{
					ID:   n.Sbom.Project.ID,
					Name: n.Sbom.Project.Name,
				},
			}
		}
		results[i] = pr
	}

	return &PolicyResultsResult{
		PolicyResults: results,
		TotalCount:    result.PolicyResults.TotalCount,
		HasNextPage:   result.PolicyResults.PageInfo.HasNextPage,
		EndCursor:     result.PolicyResults.PageInfo.EndCursor,
	}, nil
}

// ListLicensesInput contains parameters for listing licenses
type ListLicensesInput struct {
	First  int
	After  string
	Status string
	Search string
}

// ListLicenses fetches licenses with pagination
func (c *Client) ListLicenses(ctx context.Context, input ListLicensesInput) (*LicensesResult, error) {
	vars := make(map[string]interface{})
	if input.First > 0 {
		vars["first"] = input.First
	} else {
		vars["first"] = 50
	}
	if input.After != "" {
		vars["after"] = input.After
	}
	if input.Status != "" {
		vars["status"] = []string{input.Status}
	}
	if input.Search != "" {
		vars["search"] = input.Search
	}

	var result struct {
		Organization struct {
			Licenses struct {
				Nodes []struct {
					ID      string `json:"id"`
					Content struct {
						ShortID string `json:"shortId"` // For License type
						SpdxID  string `json:"spdxId"`  // For LicenseCustom type
						Name    string `json:"name"`
					} `json:"content"`
					State              string `json:"state"`
					CopyLeft           string `json:"copyLeft"`
					OsiApproved        bool   `json:"osiApproved"`
					FsfLibre           bool   `json:"fsfLibre"`
					Deprecated         bool   `json:"deprecated"`
					Attribution        string `json:"attribution"`
					SourceDistribution string `json:"sourceDistribution"`
					Modifications      string `json:"modifications"`
				} `json:"nodes"`
				TotalCount int `json:"totalCount"`
				PageInfo   struct {
					HasNextPage bool   `json:"hasNextPage"`
					EndCursor   string `json:"endCursor"`
				} `json:"pageInfo"`
			} `json:"licenses"`
		} `json:"organization"`
	}

	if err := c.gql.Execute(ctx, graphql.LicensesQuery, vars, &result); err != nil {
		return nil, err
	}

	licenses := make([]OrganizationLicense, len(result.Organization.Licenses.Nodes))
	for i, n := range result.Organization.Licenses.Nodes {
		// Use shortId if available, otherwise use spdxId
		shortID := n.Content.ShortID
		if shortID == "" {
			shortID = n.Content.SpdxID
		}
		licenses[i] = OrganizationLicense{
			ShortID:            shortID,
			Name:               n.Content.Name,
			DerivedState:       n.State,
			CopyLeft:           n.CopyLeft,
			OsiApproved:        n.OsiApproved,
			FsfLibre:           n.FsfLibre,
			Deprecated:         n.Deprecated,
			Attribution:        n.Attribution,
			SourceDistribution: n.SourceDistribution,
			Modifications:      n.Modifications,
		}
	}

	return &LicensesResult{
		Licenses:    licenses,
		TotalCount:  result.Organization.Licenses.TotalCount,
		HasNextPage: result.Organization.Licenses.PageInfo.HasNextPage,
		EndCursor:   result.Organization.Licenses.PageInfo.EndCursor,
	}, nil
}
