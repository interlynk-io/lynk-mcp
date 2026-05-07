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

// LicenseInput contains license data accepted by component mutations.
type LicenseInput struct {
	LicensesExp string `json:"licensesExp"`
}

// ChecksumInput contains checksum data accepted by component mutations.
type ChecksumInput struct {
	Alg     string `json:"alg"`
	Content string `json:"content"`
}

// ExternalURLInput contains external URL data accepted by component mutations.
type ExternalURLInput struct {
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

// ComponentVulnCustomFieldAttributeInput contains custom VEX field values.
type ComponentVulnCustomFieldAttributeInput struct {
	ID                                   string `json:"id,omitempty"`
	ComponentVulnCustomFieldDefinitionID string `json:"componentVulnCustomFieldDefinitionId,omitempty"`
	Value                                string `json:"value,omitempty"`
	Destroy                              *bool  `json:"_destroy,omitempty"`
}

// ComponentChecksum represents a checksum returned for a component.
type ComponentChecksum struct {
	Alg     string
	Content string
}

// ComponentExternalURL represents an external URL returned for a component.
type ComponentExternalURL struct {
	Name string
	URL  string
}

// UpdateComponentInput contains mutable component metadata.
type UpdateComponentInput struct {
	ID               string
	VersionID        string
	Kind             *string
	Name             *string
	Description      *string
	Copyright        *string
	Version          *string
	Group            *string
	Licenses         *LicenseInput
	Cpes             *[]string
	Purl             *string
	Primary          *bool
	Internal         *bool
	GenerateUniqueID *bool
	Scope            *string
	SupportLevel     *string
	EndOfSupport     *string
	Notice           *string
	Checksums        *[]ChecksumInput
	ExternalURLs     *[]ExternalURLInput
}

// UpdateComponentResult contains the component mutation result.
type UpdateComponentResult struct {
	Component *VersionComponent
	Errors    []string
}

// ComponentSupplier represents a component supplier.
type ComponentSupplier struct {
	ID           string
	Name         string
	URL          string
	ContactName  string
	ContactEmail string
}

// UpdateComponentSupplierInput contains mutable component supplier metadata.
type UpdateComponentSupplierInput struct {
	ID           string
	Name         *string
	URL          *string
	ContactName  *string
	ContactEmail *string
}

// UpdateComponentSupplierResult contains the component supplier mutation result.
type UpdateComponentSupplierResult struct {
	Supplier *ComponentSupplier
	Errors   []string
}

// UpdateComponentVexInput contains mutable VEX fields for a component vulnerability.
type UpdateComponentVexInput struct {
	ComponentVulnID                    string
	CurrentVersionID                   string
	VexStatusID                        *string
	VexJustificationID                 *string
	CDXResponseID                      *string
	Note                               *string
	Impact                             *string
	Detail                             *string
	Action                             *string
	FixedIn                            *string
	PropagateVex                       *bool
	ResolutionDate                     *string
	ComponentVulnCustomFieldAttributes *[]ComponentVulnCustomFieldAttributeInput
}

// UpdateComponentVexResult contains the VEX mutation result.
type UpdateComponentVexResult struct {
	ComponentVuln *ComponentVuln
	Errors        []string
}

// UpdateComponent updates mutable component metadata.
func (c *Client) UpdateComponent(ctx context.Context, input UpdateComponentInput) (*UpdateComponentResult, error) {
	vars := map[string]interface{}{
		"id":     input.ID,
		"sbomId": input.VersionID,
	}
	addStringVar(vars, "kind", input.Kind)
	addStringVar(vars, "name", input.Name)
	addStringVar(vars, "description", input.Description)
	addStringVar(vars, "copyright", input.Copyright)
	addStringVar(vars, "version", input.Version)
	addStringVar(vars, "group", input.Group)
	addStringVar(vars, "purl", input.Purl)
	addStringVar(vars, "scope", input.Scope)
	addStringVar(vars, "supportLevel", input.SupportLevel)
	addStringVar(vars, "endOfSupport", input.EndOfSupport)
	addStringVar(vars, "notice", input.Notice)
	addBoolVar(vars, "primary", input.Primary)
	addBoolVar(vars, "internal", input.Internal)
	addBoolVar(vars, "generateUniqueId", input.GenerateUniqueID)
	if input.Licenses != nil {
		vars["licenses"] = input.Licenses
	}
	if input.Cpes != nil {
		vars["cpes"] = *input.Cpes
	}
	if input.Checksums != nil {
		vars["checksums"] = *input.Checksums
	}
	if input.ExternalURLs != nil {
		vars["externalUrls"] = *input.ExternalURLs
	}

	var result struct {
		ComponentUpdate struct {
			Component *struct {
				ID           string   `json:"id"`
				Name         string   `json:"name"`
				Version      string   `json:"version"`
				Kind         string   `json:"kind"`
				Purl         string   `json:"purl"`
				Cpes         []string `json:"cpes"`
				LicensesExp  string   `json:"licensesExp"`
				Group        string   `json:"group"`
				Description  string   `json:"description"`
				Scope        string   `json:"scope"`
				Copyright    string   `json:"copyright"`
				Primary      bool     `json:"primary"`
				Internal     bool     `json:"internal"`
				UniqueID     string   `json:"uniqueId"`
				SbomID       string   `json:"sbomId"`
				Notice       string   `json:"notice"`
				SupportLevel string   `json:"supportLevel"`
				EndOfSupport string   `json:"endOfSupport"`
				Checksums    []struct {
					Alg     string `json:"alg"`
					Content string `json:"content"`
				} `json:"checksums"`
				ExternalURLs []struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"externalUrls"`
			} `json:"component"`
			Errors []string `json:"errors"`
		} `json:"componentUpdate"`
	}

	if err := c.gql.Execute(ctx, graphql.ComponentUpdateMutation, vars, &result); err != nil {
		return nil, err
	}

	mutation := result.ComponentUpdate
	updateResult := &UpdateComponentResult{Errors: mutation.Errors}
	if mutation.Component != nil {
		component := &VersionComponent{
			ID:           mutation.Component.ID,
			Name:         mutation.Component.Name,
			Version:      mutation.Component.Version,
			Kind:         mutation.Component.Kind,
			Purl:         mutation.Component.Purl,
			Cpes:         mutation.Component.Cpes,
			LicensesExp:  mutation.Component.LicensesExp,
			Group:        mutation.Component.Group,
			Description:  mutation.Component.Description,
			Scope:        mutation.Component.Scope,
			Copyright:    mutation.Component.Copyright,
			Primary:      mutation.Component.Primary,
			Internal:     mutation.Component.Internal,
			UniqueID:     mutation.Component.UniqueID,
			VersionID:    mutation.Component.SbomID,
			Notice:       mutation.Component.Notice,
			SupportLevel: mutation.Component.SupportLevel,
			EndOfSupport: mutation.Component.EndOfSupport,
		}
		component.Checksums = make([]ComponentChecksum, len(mutation.Component.Checksums))
		for i, checksum := range mutation.Component.Checksums {
			component.Checksums[i] = ComponentChecksum{Alg: checksum.Alg, Content: checksum.Content}
		}
		component.ExternalURLs = make([]ComponentExternalURL, len(mutation.Component.ExternalURLs))
		for i, externalURL := range mutation.Component.ExternalURLs {
			component.ExternalURLs[i] = ComponentExternalURL{Name: externalURL.Name, URL: externalURL.URL}
		}
		updateResult.Component = component
	}

	return updateResult, nil
}

// UpdateComponentSupplier updates component supplier metadata.
func (c *Client) UpdateComponentSupplier(ctx context.Context, input UpdateComponentSupplierInput) (*UpdateComponentSupplierResult, error) {
	vars := map[string]interface{}{"id": input.ID}
	addStringVar(vars, "name", input.Name)
	addStringVar(vars, "url", input.URL)
	addStringVar(vars, "contactName", input.ContactName)
	addStringVar(vars, "contactEmail", input.ContactEmail)

	var result struct {
		CompSupplierUpdate struct {
			CompSupplier *struct {
				ID           string `json:"id"`
				Name         string `json:"name"`
				URL          string `json:"url"`
				ContactName  string `json:"contactName"`
				ContactEmail string `json:"contactEmail"`
			} `json:"compSupplier"`
			Errors []string `json:"errors"`
		} `json:"compSupplierUpdate"`
	}

	if err := c.gql.Execute(ctx, graphql.ComponentSupplierUpdateMutation, vars, &result); err != nil {
		return nil, err
	}

	mutation := result.CompSupplierUpdate
	updateResult := &UpdateComponentSupplierResult{Errors: mutation.Errors}
	if mutation.CompSupplier != nil {
		updateResult.Supplier = &ComponentSupplier{
			ID:           mutation.CompSupplier.ID,
			Name:         mutation.CompSupplier.Name,
			URL:          mutation.CompSupplier.URL,
			ContactName:  mutation.CompSupplier.ContactName,
			ContactEmail: mutation.CompSupplier.ContactEmail,
		}
	}

	return updateResult, nil
}

// UpdateComponentVex updates VEX data for a component vulnerability.
func (c *Client) UpdateComponentVex(ctx context.Context, input UpdateComponentVexInput) (*UpdateComponentVexResult, error) {
	vars := map[string]interface{}{
		"componentVulnId": input.ComponentVulnID,
		"currentSbomId":   input.CurrentVersionID,
	}
	addStringVar(vars, "vexStatusId", input.VexStatusID)
	addStringVar(vars, "vexJustificationId", input.VexJustificationID)
	addStringVar(vars, "cdxResponseId", input.CDXResponseID)
	addStringVar(vars, "note", input.Note)
	addStringVar(vars, "impact", input.Impact)
	addStringVar(vars, "detail", input.Detail)
	addStringVar(vars, "action", input.Action)
	addStringVar(vars, "fixedIn", input.FixedIn)
	addBoolVar(vars, "propagateVex", input.PropagateVex)
	addStringVar(vars, "resolutionDate", input.ResolutionDate)
	if input.ComponentVulnCustomFieldAttributes != nil {
		vars["componentVulnCustomFieldAttributes"] = *input.ComponentVulnCustomFieldAttributes
	}

	var result struct {
		ComponentVexUpdate struct {
			ComponentVuln *struct {
				ID          string `json:"id"`
				ComponentID string `json:"componentId"`
				VulnID      string `json:"vulnId"`
				SbomID      string `json:"sbomId"`
				FixedIn     string `json:"fixedIn"`
				Detail      string `json:"detail"`
				Impact      string `json:"impact"`
				ActionStmt  string `json:"actionStmt"`
				VexStatus   *struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				} `json:"vexStatus"`
				VexJustification *struct {
					ID   string `json:"id"`
					Name string `json:"name"`
				} `json:"vexJustification"`
			} `json:"componentVuln"`
			Errors []string `json:"errors"`
		} `json:"componentVexUpdate"`
	}

	if err := c.gql.Execute(ctx, graphql.ComponentVexUpdateMutation, vars, &result); err != nil {
		return nil, err
	}

	mutation := result.ComponentVexUpdate
	updateResult := &UpdateComponentVexResult{Errors: mutation.Errors}
	if mutation.ComponentVuln != nil {
		componentVuln := &ComponentVuln{
			ID:          mutation.ComponentVuln.ID,
			ComponentID: mutation.ComponentVuln.ComponentID,
			VulnID:      mutation.ComponentVuln.VulnID,
			VersionID:   mutation.ComponentVuln.SbomID,
			FixedIn:     mutation.ComponentVuln.FixedIn,
			Detail:      mutation.ComponentVuln.Detail,
			Impact:      mutation.ComponentVuln.Impact,
			ActionStmt:  mutation.ComponentVuln.ActionStmt,
		}
		if mutation.ComponentVuln.VexStatus != nil {
			componentVuln.VexStatus = &VexStatus{
				ID:   mutation.ComponentVuln.VexStatus.ID,
				Name: mutation.ComponentVuln.VexStatus.Name,
			}
		}
		if mutation.ComponentVuln.VexJustification != nil {
			componentVuln.VexJustification = &VexJustification{
				ID:   mutation.ComponentVuln.VexJustification.ID,
				Name: mutation.ComponentVuln.VexJustification.Name,
			}
		}
		updateResult.ComponentVuln = componentVuln
	}

	return updateResult, nil
}

func addStringVar(vars map[string]interface{}, key string, value *string) {
	if value != nil {
		vars[key] = *value
	}
}

func addBoolVar(vars map[string]interface{}, key string, value *bool) {
	if value != nil {
		vars[key] = *value
	}
}
