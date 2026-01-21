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

// ProjectGroup represents a product/project group
type ProjectGroup struct {
	ID             string
	Name           string
	Description    string
	Enabled        bool
	OrganizationID string
	UpdatedAt      time.Time
	SbomsCount     int
	Projects       []Project
}

// Project represents a project/stream
type Project struct {
	ID             string
	Name           string
	Description    string
	Enabled        bool
	ProjectGroupID string
	UpdatedAt      time.Time
	SbomsCount     int
	ProjectGroup   *ProjectGroup
}

// ProjectGroupsResult represents the result of listing project groups
type ProjectGroupsResult struct {
	ProjectGroups []ProjectGroup
	TotalCount    int
	HasNextPage   bool
	EndCursor     string
}

// ListProjectGroupsInput contains parameters for listing project groups
type ListProjectGroupsInput struct {
	First  int
	After  string
	Search string
}

// ListProjectGroups fetches project groups with pagination
func (c *Client) ListProjectGroups(ctx context.Context, input ListProjectGroupsInput) (*ProjectGroupsResult, error) {
	vars := make(map[string]interface{})
	if input.First > 0 {
		vars["first"] = input.First
	} else {
		vars["first"] = 20 // default
	}
	if input.After != "" {
		vars["after"] = input.After
	}
	if input.Search != "" {
		vars["search"] = input.Search
	}

	var result struct {
		Organization struct {
			ProjectGroups struct {
				Nodes []struct {
					ID             string    `json:"id"`
					Name           string    `json:"name"`
					Description    string    `json:"description"`
					Enabled        bool      `json:"enabled"`
					OrganizationID string    `json:"organizationId"`
					UpdatedAt      time.Time `json:"updatedAt"`
					SbomsCount     int       `json:"sbomsCount"`
				} `json:"nodes"`
				TotalCount int `json:"totalCount"`
				PageInfo   struct {
					HasNextPage bool   `json:"hasNextPage"`
					EndCursor   string `json:"endCursor"`
				} `json:"pageInfo"`
			} `json:"projectGroups"`
		} `json:"organization"`
	}

	if err := c.gql.Execute(ctx, graphql.ProjectGroupsQuery, vars, &result); err != nil {
		return nil, err
	}

	groups := make([]ProjectGroup, len(result.Organization.ProjectGroups.Nodes))
	for i, n := range result.Organization.ProjectGroups.Nodes {
		groups[i] = ProjectGroup{
			ID:             n.ID,
			Name:           n.Name,
			Description:    n.Description,
			Enabled:        n.Enabled,
			OrganizationID: n.OrganizationID,
			UpdatedAt:      n.UpdatedAt,
			SbomsCount:     n.SbomsCount,
		}
	}

	return &ProjectGroupsResult{
		ProjectGroups: groups,
		TotalCount:    result.Organization.ProjectGroups.TotalCount,
		HasNextPage:   result.Organization.ProjectGroups.PageInfo.HasNextPage,
		EndCursor:     result.Organization.ProjectGroups.PageInfo.EndCursor,
	}, nil
}

// GetProjectGroup fetches a single project group by ID
func (c *Client) GetProjectGroup(ctx context.Context, id string) (*ProjectGroup, error) {
	vars := map[string]interface{}{
		"id": id,
	}

	var result struct {
		ProjectGroup struct {
			ID             string    `json:"id"`
			Name           string    `json:"name"`
			Description    string    `json:"description"`
			Enabled        bool      `json:"enabled"`
			OrganizationID string    `json:"organizationId"`
			UpdatedAt      time.Time `json:"updatedAt"`
			SbomsCount     int       `json:"sbomsCount"`
			Projects       []struct {
				ID          string    `json:"id"`
				Name        string    `json:"name"`
				Description string    `json:"description"`
				Enabled     bool      `json:"enabled"`
				UpdatedAt   time.Time `json:"updatedAt"`
				SbomsCount  int       `json:"sbomsCount"`
			} `json:"projects"`
		} `json:"projectGroup"`
	}

	if err := c.gql.Execute(ctx, graphql.ProjectGroupQuery, vars, &result); err != nil {
		return nil, err
	}

	projects := make([]Project, len(result.ProjectGroup.Projects))
	for i, p := range result.ProjectGroup.Projects {
		projects[i] = Project{
			ID:             p.ID,
			Name:           p.Name,
			Description:    p.Description,
			Enabled:        p.Enabled,
			UpdatedAt:      p.UpdatedAt,
			SbomsCount:     p.SbomsCount,
			ProjectGroupID: result.ProjectGroup.ID,
		}
	}

	return &ProjectGroup{
		ID:             result.ProjectGroup.ID,
		Name:           result.ProjectGroup.Name,
		Description:    result.ProjectGroup.Description,
		Enabled:        result.ProjectGroup.Enabled,
		OrganizationID: result.ProjectGroup.OrganizationID,
		UpdatedAt:      result.ProjectGroup.UpdatedAt,
		SbomsCount:     result.ProjectGroup.SbomsCount,
		Projects:       projects,
	}, nil
}

// GetProject fetches a single project by ID
func (c *Client) GetProject(ctx context.Context, id string) (*Project, error) {
	vars := map[string]interface{}{
		"id": id,
	}

	var result struct {
		Project struct {
			ID             string    `json:"id"`
			Name           string    `json:"name"`
			Description    string    `json:"description"`
			Enabled        bool      `json:"enabled"`
			ProjectGroupID string    `json:"projectGroupId"`
			UpdatedAt      time.Time `json:"updatedAt"`
			SbomsCount     int       `json:"sbomsCount"`
			ProjectGroup   struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"projectGroup"`
		} `json:"project"`
	}

	if err := c.gql.Execute(ctx, graphql.ProjectQuery, vars, &result); err != nil {
		return nil, err
	}

	var pg *ProjectGroup
	if result.Project.ProjectGroup.ID != "" {
		pg = &ProjectGroup{
			ID:   result.Project.ProjectGroup.ID,
			Name: result.Project.ProjectGroup.Name,
		}
	}

	return &Project{
		ID:             result.Project.ID,
		Name:           result.Project.Name,
		Description:    result.Project.Description,
		Enabled:        result.Project.Enabled,
		ProjectGroupID: result.Project.ProjectGroupID,
		UpdatedAt:      result.Project.UpdatedAt,
		SbomsCount:     result.Project.SbomsCount,
		ProjectGroup:   pg,
	}, nil
}
