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
	"fmt"
	"sort"
	"time"

	"github.com/interlynk-io/lynk-mcp/internal/graphql"
)

const getProductEnvironmentsPageSize = 100

// Product represents a product (formerly project group)
type Product struct {
	ID                       string
	Name                     string
	Description              string
	Enabled                  bool
	OrganizationID           string
	UpdatedAt                time.Time
	VersionsCount            int
	Labels                   []Label
	Repository               *ImportedRepository
	TicketingDefaultsSummary *TicketingDefaultsSummary
	Environments             []Environment
}

// Label represents an organization label applied to products.
type Label struct {
	ID             string
	Name           string
	Color          string
	OrganizationID string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// Environment represents an environment (formerly project)
type Environment struct {
	ID            string
	Name          string
	Description   string
	Enabled       bool
	ProductID     string
	UpdatedAt     time.Time
	VersionsCount int
	JiraDefaults  *JiraDefaults
	Product       *Product
}

// TicketingDefaultsSummary contains compact per-product ticketing defaults.
type TicketingDefaultsSummary struct {
	JiraConfigured       bool
	JiraProjectKeys      []string
	EnvironmentsWithJira int
}

// JiraDefaults contains compact per-environment Jira defaults used for ticket creation.
type JiraDefaults struct {
	ID         string
	ProjectKey string
	IssueType  string
	Assignee   string
	Reporter   string
	Epic       string
	Components interface{}
	EnableSync bool
	UpdatedAt  time.Time
}

// ProductsResult represents the result of listing products
type ProductsResult struct {
	Products    []Product
	TotalCount  int
	HasNextPage bool
	EndCursor   string
}

// LabelsResult represents the result of listing labels.
type LabelsResult struct {
	Labels      []Label
	TotalCount  int
	HasNextPage bool
	EndCursor   string
}

// ListProductsInput contains parameters for listing products
type ListProductsInput struct {
	First    int
	After    string
	Search   string
	LabelIDs []string
}

// ListLabelsInput contains parameters for listing labels.
type ListLabelsInput struct {
	First  int
	After  string
	Search string
}

// ListLabels fetches labels with pagination.
func (c *Client) ListLabels(ctx context.Context, input ListLabelsInput) (*LabelsResult, error) {
	vars := make(map[string]interface{})
	if input.First > 0 {
		vars["first"] = input.First
	} else {
		vars["first"] = 50
	}
	if input.After != "" {
		vars["after"] = input.After
	}
	if input.Search != "" {
		vars["search"] = input.Search
	}

	var result struct {
		Labels struct {
			Nodes      []labelNode `json:"nodes"`
			TotalCount int         `json:"totalCount"`
			PageInfo   struct {
				HasNextPage bool   `json:"hasNextPage"`
				EndCursor   string `json:"endCursor"`
			} `json:"pageInfo"`
		} `json:"labels"`
	}

	if err := c.gql.Execute(ctx, graphql.LabelsQuery, vars, &result); err != nil {
		return nil, err
	}

	return &LabelsResult{
		Labels:      mapLabelNodes(result.Labels.Nodes),
		TotalCount:  result.Labels.TotalCount,
		HasNextPage: result.Labels.PageInfo.HasNextPage,
		EndCursor:   result.Labels.PageInfo.EndCursor,
	}, nil
}

// ListProducts fetches products with pagination
func (c *Client) ListProducts(ctx context.Context, input ListProductsInput) (*ProductsResult, error) {
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
	if len(input.LabelIDs) > 0 {
		vars["labelIds"] = input.LabelIDs
	}

	var result struct {
		Organization struct {
			ProjectGroups struct {
				Nodes []struct {
					ID                 string                 `json:"id"`
					Name               string                 `json:"name"`
					Description        string                 `json:"description"`
					Enabled            bool                   `json:"enabled"`
					OrganizationID     string                 `json:"organizationId"`
					UpdatedAt          time.Time              `json:"updatedAt"`
					SbomsCount         int                    `json:"sbomsCount"`
					Labels             []labelNode            `json:"labels"`
					ImportedRepository *productRepositoryNode `json:"importedRepository"`
					Projects           struct {
						Nodes []struct {
							ExternalIssueTrackerSettings []productIssueTrackerSettingNode `json:"externalIssueTrackerSettings"`
						} `json:"nodes"`
					} `json:"projects"`
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

	products := make([]Product, len(result.Organization.ProjectGroups.Nodes))
	for i, n := range result.Organization.ProjectGroups.Nodes {
		products[i] = Product{
			ID:                       n.ID,
			Name:                     n.Name,
			Description:              n.Description,
			Enabled:                  n.Enabled,
			OrganizationID:           n.OrganizationID,
			UpdatedAt:                n.UpdatedAt,
			VersionsCount:            n.SbomsCount,
			Labels:                   mapLabelNodes(n.Labels),
			Repository:               mapProductRepository(n.ImportedRepository),
			TicketingDefaultsSummary: summarizeTicketingDefaults(n.Projects.Nodes),
		}
	}

	return &ProductsResult{
		Products:    products,
		TotalCount:  result.Organization.ProjectGroups.TotalCount,
		HasNextPage: result.Organization.ProjectGroups.PageInfo.HasNextPage,
		EndCursor:   result.Organization.ProjectGroups.PageInfo.EndCursor,
	}, nil
}

// GetProduct fetches a single product by ID
func (c *Client) GetProduct(ctx context.Context, id string) (*Product, error) {
	var product *Product
	var environments []Environment
	var after string

	for {
		vars := map[string]interface{}{
			"id":            id,
			"projectsFirst": getProductEnvironmentsPageSize,
		}
		if after != "" {
			vars["projectsAfter"] = after
		}

		var result struct {
			ProjectGroup struct {
				ID                 string                 `json:"id"`
				Name               string                 `json:"name"`
				Description        string                 `json:"description"`
				Enabled            bool                   `json:"enabled"`
				OrganizationID     string                 `json:"organizationId"`
				UpdatedAt          time.Time              `json:"updatedAt"`
				SbomsCount         int                    `json:"sbomsCount"`
				Labels             []labelNode            `json:"labels"`
				ImportedRepository *productRepositoryNode `json:"importedRepository"`
				Projects           struct {
					Nodes []struct {
						ID                           string                           `json:"id"`
						Name                         string                           `json:"name"`
						Description                  string                           `json:"description"`
						Enabled                      bool                             `json:"enabled"`
						UpdatedAt                    time.Time                        `json:"updatedAt"`
						SbomsCount                   int                              `json:"sbomsCount"`
						ExternalIssueTrackerSettings []productIssueTrackerSettingNode `json:"externalIssueTrackerSettings"`
					} `json:"nodes"`
					PageInfo struct {
						HasNextPage bool   `json:"hasNextPage"`
						EndCursor   string `json:"endCursor"`
					} `json:"pageInfo"`
				} `json:"projects"`
			} `json:"projectGroup"`
		}

		if err := c.gql.Execute(ctx, graphql.ProjectGroupQuery, vars, &result); err != nil {
			return nil, err
		}

		if product == nil {
			product = &Product{
				ID:             result.ProjectGroup.ID,
				Name:           result.ProjectGroup.Name,
				Description:    result.ProjectGroup.Description,
				Enabled:        result.ProjectGroup.Enabled,
				OrganizationID: result.ProjectGroup.OrganizationID,
				UpdatedAt:      result.ProjectGroup.UpdatedAt,
				VersionsCount:  result.ProjectGroup.SbomsCount,
				Labels:         mapLabelNodes(result.ProjectGroup.Labels),
				Repository:     mapProductRepository(result.ProjectGroup.ImportedRepository),
			}
		}

		for _, p := range result.ProjectGroup.Projects.Nodes {
			environments = append(environments, Environment{
				ID:            p.ID,
				Name:          p.Name,
				Description:   p.Description,
				Enabled:       p.Enabled,
				UpdatedAt:     p.UpdatedAt,
				VersionsCount: p.SbomsCount,
				ProductID:     result.ProjectGroup.ID,
				JiraDefaults:  mapJiraDefaults(p.ExternalIssueTrackerSettings),
			})
		}

		if !result.ProjectGroup.Projects.PageInfo.HasNextPage {
			break
		}

		after = result.ProjectGroup.Projects.PageInfo.EndCursor
		if after == "" {
			return nil, fmt.Errorf("project group projects pageInfo has next page without end cursor")
		}
	}

	sort.SliceStable(environments, func(i, j int) bool {
		if environments[i].VersionsCount != environments[j].VersionsCount {
			return environments[i].VersionsCount > environments[j].VersionsCount
		}
		return environments[i].Name < environments[j].Name
	})

	product.Environments = environments
	return product, nil
}

// GetEnvironment fetches a single environment by ID
func (c *Client) GetEnvironment(ctx context.Context, id string) (*Environment, error) {
	vars := map[string]interface{}{
		"id": id,
	}

	var result struct {
		Project struct {
			ID                           string                           `json:"id"`
			Name                         string                           `json:"name"`
			Description                  string                           `json:"description"`
			Enabled                      bool                             `json:"enabled"`
			ProjectGroupID               string                           `json:"projectGroupId"`
			UpdatedAt                    time.Time                        `json:"updatedAt"`
			SbomsCount                   int                              `json:"sbomsCount"`
			ExternalIssueTrackerSettings []productIssueTrackerSettingNode `json:"externalIssueTrackerSettings"`
			ProjectGroup                 struct {
				ID                 string                 `json:"id"`
				Name               string                 `json:"name"`
				ImportedRepository *productRepositoryNode `json:"importedRepository"`
			} `json:"projectGroup"`
		} `json:"project"`
	}

	if err := c.gql.Execute(ctx, graphql.ProjectQuery, vars, &result); err != nil {
		return nil, err
	}

	var product *Product
	if result.Project.ProjectGroup.ID != "" {
		product = &Product{
			ID:         result.Project.ProjectGroup.ID,
			Name:       result.Project.ProjectGroup.Name,
			Repository: mapProductRepository(result.Project.ProjectGroup.ImportedRepository),
		}
	}

	return &Environment{
		ID:            result.Project.ID,
		Name:          result.Project.Name,
		Description:   result.Project.Description,
		Enabled:       result.Project.Enabled,
		ProductID:     result.Project.ProjectGroupID,
		UpdatedAt:     result.Project.UpdatedAt,
		VersionsCount: result.Project.SbomsCount,
		JiraDefaults:  mapJiraDefaults(result.Project.ExternalIssueTrackerSettings),
		Product:       product,
	}, nil
}

type productRepositoryNode struct {
	Type           string `json:"__typename"`
	ID             string `json:"id"`
	Name           string `json:"name"`
	FullName       string `json:"fullName"`
	Owner          string `json:"owner"`
	DefaultBranch  string `json:"defaultBranch"`
	Slug           string `json:"slug"`
	Workspace      string `json:"workspace"`
	FullPath       string `json:"fullPath"`
	GitlabID       string `json:"gitlabId"`
	ImportStatus   string `json:"importStatus"`
	WebhookEnabled bool   `json:"webhookEnabled"`
}

type labelNode struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Color          string    `json:"color"`
	OrganizationID string    `json:"organizationId"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type productIssueTrackerSettingNode struct {
	ID             string      `json:"id"`
	Provider       string      `json:"provider"`
	ProjectKey     string      `json:"projectKey"`
	IssueType      string      `json:"issueType"`
	Assignee       string      `json:"assignee"`
	Reporter       string      `json:"reporter"`
	Epic           string      `json:"epic"`
	Components     interface{} `json:"components"`
	EnableJiraSync bool        `json:"enableJiraSync"`
	UpdatedAt      time.Time   `json:"updatedAt"`
}

func mapLabelNodes(nodes []labelNode) []Label {
	if len(nodes) == 0 {
		return nil
	}
	labels := make([]Label, len(nodes))
	for i, node := range nodes {
		labels[i] = Label{
			ID:             node.ID,
			Name:           node.Name,
			Color:          node.Color,
			OrganizationID: node.OrganizationID,
			CreatedAt:      node.CreatedAt,
			UpdatedAt:      node.UpdatedAt,
		}
	}
	return labels
}

func mapProductRepository(node *productRepositoryNode) *ImportedRepository {
	if node == nil {
		return nil
	}
	return &ImportedRepository{
		Type:           node.Type,
		ID:             node.ID,
		Name:           node.Name,
		FullName:       node.FullName,
		Owner:          node.Owner,
		DefaultBranch:  node.DefaultBranch,
		Slug:           node.Slug,
		Workspace:      node.Workspace,
		FullPath:       node.FullPath,
		GitlabID:       node.GitlabID,
		ImportStatus:   node.ImportStatus,
		WebhookEnabled: node.WebhookEnabled,
	}
}

func mapJiraDefaults(settings []productIssueTrackerSettingNode) *JiraDefaults {
	for _, setting := range settings {
		if setting.Provider != "jira" {
			continue
		}
		return &JiraDefaults{
			ID:         setting.ID,
			ProjectKey: setting.ProjectKey,
			IssueType:  setting.IssueType,
			Assignee:   setting.Assignee,
			Reporter:   setting.Reporter,
			Epic:       setting.Epic,
			Components: setting.Components,
			EnableSync: setting.EnableJiraSync,
			UpdatedAt:  setting.UpdatedAt,
		}
	}
	return nil
}

func summarizeTicketingDefaults(environments []struct {
	ExternalIssueTrackerSettings []productIssueTrackerSettingNode `json:"externalIssueTrackerSettings"`
}) *TicketingDefaultsSummary {
	keys := map[string]bool{}
	environmentsWithJira := 0
	for _, env := range environments {
		hasJira := false
		for _, setting := range env.ExternalIssueTrackerSettings {
			if setting.Provider != "jira" {
				continue
			}
			hasJira = true
			if setting.ProjectKey != "" {
				keys[setting.ProjectKey] = true
			}
		}
		if hasJira {
			environmentsWithJira++
		}
	}

	summary := &TicketingDefaultsSummary{
		JiraConfigured:       environmentsWithJira > 0,
		JiraProjectKeys:      []string{},
		EnvironmentsWithJira: environmentsWithJira,
	}
	for key := range keys {
		summary.JiraProjectKeys = append(summary.JiraProjectKeys, key)
	}
	sort.Strings(summary.JiraProjectKeys)
	return summary
}
