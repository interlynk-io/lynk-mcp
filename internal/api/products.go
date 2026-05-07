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
	"time"

	"github.com/interlynk-io/lynk-mcp/internal/graphql"
)

const getProductEnvironmentsPageSize = 100

// Product represents a product (formerly project group)
type Product struct {
	ID             string
	Name           string
	Description    string
	Enabled        bool
	OrganizationID string
	UpdatedAt      time.Time
	VersionsCount  int
	Environments   []Environment
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
	Product       *Product
}

// ProductsResult represents the result of listing products
type ProductsResult struct {
	Products    []Product
	TotalCount  int
	HasNextPage bool
	EndCursor   string
}

// ListProductsInput contains parameters for listing products
type ListProductsInput struct {
	First  int
	After  string
	Search string
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

	products := make([]Product, len(result.Organization.ProjectGroups.Nodes))
	for i, n := range result.Organization.ProjectGroups.Nodes {
		products[i] = Product{
			ID:             n.ID,
			Name:           n.Name,
			Description:    n.Description,
			Enabled:        n.Enabled,
			OrganizationID: n.OrganizationID,
			UpdatedAt:      n.UpdatedAt,
			VersionsCount:  n.SbomsCount,
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
				ID             string    `json:"id"`
				Name           string    `json:"name"`
				Description    string    `json:"description"`
				Enabled        bool      `json:"enabled"`
				OrganizationID string    `json:"organizationId"`
				UpdatedAt      time.Time `json:"updatedAt"`
				SbomsCount     int       `json:"sbomsCount"`
				Projects       struct {
					Nodes []struct {
						ID          string    `json:"id"`
						Name        string    `json:"name"`
						Description string    `json:"description"`
						Enabled     bool      `json:"enabled"`
						UpdatedAt   time.Time `json:"updatedAt"`
						SbomsCount  int       `json:"sbomsCount"`
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

	var product *Product
	if result.Project.ProjectGroup.ID != "" {
		product = &Product{
			ID:   result.Project.ProjectGroup.ID,
			Name: result.Project.ProjectGroup.Name,
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
		Product:       product,
	}, nil
}
