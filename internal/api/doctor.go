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

	"github.com/interlynk-io/lynk-mcp/internal/graphql"
)

// DoctorFinding represents a single SBOM Doctor finding for a version.
type DoctorFinding struct {
	CheckCode        string
	CheckName        string
	Severity         string
	Domain           string
	ComponentID      string
	ComponentName    string
	ComponentVersion string
	AutoFixable      bool
	Findings         interface{}
}

// DoctorPageInfo contains pagination metadata for Doctor findings.
type DoctorPageInfo struct {
	EndCursor       string
	HasNextPage     bool
	HasPreviousPage bool
	StartCursor     string
}

// DoctorResultsResult represents the result of listing Doctor findings.
type DoctorResultsResult struct {
	Findings   []DoctorFinding
	TotalCount int
	PageInfo   DoctorPageInfo
}

// ListDoctorResultsInput contains parameters for listing Doctor findings.
type ListDoctorResultsInput struct {
	VersionID     string
	Search        string
	ComponentID   string
	Severity      []string
	Domain        []string
	CheckCode     []string
	ComponentName []string
	ForceRefresh  *bool
	First         int
	Last          int
	After         string
	Before        string
}

// ListDoctorResults fetches SBOM Doctor findings for a version.
func (c *Client) ListDoctorResults(ctx context.Context, input ListDoctorResultsInput) (*DoctorResultsResult, error) {
	vars := map[string]interface{}{
		"sbomId": input.VersionID,
	}
	if input.Search != "" {
		vars["search"] = input.Search
	}
	if input.ComponentID != "" {
		vars["componentId"] = input.ComponentID
	}
	if len(input.Severity) > 0 {
		vars["severity"] = input.Severity
	}
	if len(input.Domain) > 0 {
		vars["domain"] = input.Domain
	}
	if len(input.CheckCode) > 0 {
		vars["checkCode"] = input.CheckCode
	}
	if len(input.ComponentName) > 0 {
		vars["componentName"] = input.ComponentName
	}
	if input.ForceRefresh != nil {
		vars["forceRefresh"] = *input.ForceRefresh
	}
	if input.First > 0 {
		vars["first"] = input.First
	}
	if input.Last > 0 {
		vars["last"] = input.Last
	}
	if input.After != "" {
		vars["after"] = input.After
	}
	if input.Before != "" {
		vars["before"] = input.Before
	}
	if _, ok := vars["first"]; !ok {
		if _, ok := vars["last"]; !ok {
			vars["first"] = 25
		}
	}

	var result struct {
		DoctorResults struct {
			Nodes []struct {
				CheckCode        string      `json:"checkCode"`
				CheckName        string      `json:"checkName"`
				Severity         string      `json:"severity"`
				Domain           string      `json:"domain"`
				ComponentID      string      `json:"componentId"`
				ComponentName    string      `json:"componentName"`
				ComponentVersion string      `json:"componentVersion"`
				AutoFixable      bool        `json:"autoFixable"`
				Findings         interface{} `json:"findings"`
			} `json:"nodes"`
			TotalCount int `json:"totalCount"`
			PageInfo   struct {
				EndCursor       string `json:"endCursor"`
				HasNextPage     bool   `json:"hasNextPage"`
				HasPreviousPage bool   `json:"hasPreviousPage"`
				StartCursor     string `json:"startCursor"`
			} `json:"pageInfo"`
		} `json:"doctorResults"`
	}

	if err := c.gql.Execute(ctx, graphql.DoctorResultsQuery, vars, &result); err != nil {
		return nil, err
	}

	findings := make([]DoctorFinding, len(result.DoctorResults.Nodes))
	for i, n := range result.DoctorResults.Nodes {
		findings[i] = DoctorFinding{
			CheckCode:        n.CheckCode,
			CheckName:        n.CheckName,
			Severity:         n.Severity,
			Domain:           n.Domain,
			ComponentID:      n.ComponentID,
			ComponentName:    n.ComponentName,
			ComponentVersion: n.ComponentVersion,
			AutoFixable:      n.AutoFixable,
			Findings:         n.Findings,
		}
	}

	return &DoctorResultsResult{
		Findings:   findings,
		TotalCount: result.DoctorResults.TotalCount,
		PageInfo: DoctorPageInfo{
			EndCursor:       result.DoctorResults.PageInfo.EndCursor,
			HasNextPage:     result.DoctorResults.PageInfo.HasNextPage,
			HasPreviousPage: result.DoctorResults.PageInfo.HasPreviousPage,
			StartCursor:     result.DoctorResults.PageInfo.StartCursor,
		},
	}, nil
}
