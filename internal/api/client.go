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

type graphQLExecutor interface {
	Execute(ctx context.Context, query string, variables map[string]interface{}, result interface{}) error
}

// Client provides high-level operations for the Lynk API
type Client struct {
	gql graphQLExecutor
}

// NewClient creates a new API client
func NewClient(endpoint, token string, timeout time.Duration) *Client {
	return &Client{
		gql: graphql.NewClient(endpoint, token, timeout),
	}
}

// Organization represents organization data with metrics
type Organization struct {
	ID        string
	Name      string
	Email     string
	URL       string
	Status    string
	Tier      string
	UpdatedAt time.Time
	Metrics   *OrganizationMetrics
}

// OrganizationMetrics contains organization metrics
type OrganizationMetrics struct {
	ProjectCount   int
	VersionCount   int
	ComponentCount int
	VulnsMetric    map[string]interface{}
}

// GetOrganization fetches the current organization information
func (c *Client) GetOrganization(ctx context.Context) (*Organization, error) {
	var result struct {
		Organization struct {
			ID        string    `json:"id"`
			Name      string    `json:"name"`
			Email     string    `json:"email"`
			URL       string    `json:"url"`
			Status    string    `json:"status"`
			Tier      string    `json:"tier"`
			UpdatedAt time.Time `json:"updatedAt"`
		} `json:"organization"`
		OrganizationMetric struct {
			ProjectCount   int                    `json:"projectCount"`
			VersionCount   int                    `json:"versionCount"`
			ComponentCount int                    `json:"componentCount"`
			VulnsMetric    map[string]interface{} `json:"vulnsMetric"`
		} `json:"organizationMetric"`
	}

	if err := c.gql.Execute(ctx, graphql.OrganizationQuery, nil, &result); err != nil {
		return nil, err
	}

	return &Organization{
		ID:        result.Organization.ID,
		Name:      result.Organization.Name,
		Email:     result.Organization.Email,
		URL:       result.Organization.URL,
		Status:    result.Organization.Status,
		Tier:      result.Organization.Tier,
		UpdatedAt: result.Organization.UpdatedAt,
		Metrics: &OrganizationMetrics{
			ProjectCount:   result.OrganizationMetric.ProjectCount,
			VersionCount:   result.OrganizationMetric.VersionCount,
			ComponentCount: result.OrganizationMetric.ComponentCount,
			VulnsMetric:    result.OrganizationMetric.VulnsMetric,
		},
	}, nil
}
