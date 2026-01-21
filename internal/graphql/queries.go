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

package graphql

// GraphQL query strings for the Lynk API

const (
	// OrganizationQuery fetches organization information with metrics
	OrganizationQuery = `
		query GetOrganization {
			organization {
				id
				name
				email
				url
				status
				tier
				updatedAt
			}
			organizationMetric {
				projectCount
				versionCount
				componentCount
				vulnsMetric
			}
		}
	`

	// ProjectGroupsQuery fetches project groups with pagination
	ProjectGroupsQuery = `
		query GetProjectGroups($first: Int, $after: String, $search: String) {
			organization {
				projectGroups(first: $first, after: $after, search: $search) {
					nodes {
						id
						name
						description
						enabled
						organizationId
						updatedAt
						sbomsCount
					}
					totalCount
					pageInfo {
						hasNextPage
						hasPreviousPage
						startCursor
						endCursor
					}
				}
			}
		}
	`

	// ProjectGroupQuery fetches a single project group by ID
	ProjectGroupQuery = `
		query GetProjectGroup($id: Uuid!) {
			projectGroup(id: $id) {
				id
				name
				description
				enabled
				organizationId
				updatedAt
				sbomsCount
				projects {
					id
					name
					description
					enabled
					updatedAt
					sbomsCount
				}
			}
		}
	`

	// ProjectQuery fetches a single project by ID
	ProjectQuery = `
		query GetProject($id: Uuid!) {
			project(id: $id) {
				id
				name
				description
				enabled
				projectGroupId
				updatedAt
				sbomsCount
				projectGroup {
					id
					name
				}
			}
		}
	`

	// ProjectSbomsQuery fetches SBOMs for a project
	ProjectSbomsQuery = `
		query GetProjectSboms($projectId: Uuid!, $first: Int, $after: String, $lifestage: [ProductLifecycleStageEnum!]) {
			project(id: $projectId) {
				sbomVersions(first: $first, after: $after, lifestage: $lifestage) {
					nodes {
						id
						projectVersion
						spec
						specVersion
						format
						lifecycle
						createdAt
						updatedAt
						projectId
						stats {
							compCount
							compPurlCount
							compCpeCount
							compLicenseCount
							compSupplierCount
							vulnStats
						}
					}
					totalCount
					pageInfo {
						hasNextPage
						hasPreviousPage
						startCursor
						endCursor
					}
				}
			}
		}
	`

	// SbomQuery fetches a single SBOM by ID
	SbomQuery = `
		query GetSbom($sbomId: Uuid!) {
			sbom(sbomId: $sbomId) {
				id
				projectVersion
				spec
				specVersion
				format
				lifecycle
				createdAt
				updatedAt
				projectId
				stats {
					compCount
					compPurlCount
					compCpeCount
					compLicenseCount
					compSupplierCount
					vulnStats
				}
				project {
					id
					name
					projectGroupId
				}
			}
		}
	`

	// SbomComponentsQuery fetches components for an SBOM
	SbomComponentsQuery = `
		query GetSbomComponents($sbomId: Uuid!, $first: Int, $after: String, $search: String, $kind: [String!], $direct: Boolean) {
			sbom(sbomId: $sbomId) {
				components(sbomId: $sbomId, first: $first, after: $after, search: $search, kind: $kind, direct: $direct) {
					nodes {
						id
						name
						version
						kind
						purl
						cpes
						licensesExp
						group
						description
						primary
						internal
						sbomId
						updatedAt
					}
					totalCount
					pageInfo {
						hasNextPage
						hasPreviousPage
						startCursor
						endCursor
					}
				}
			}
		}
	`

	// ComponentQuery fetches a single component by ID
	ComponentQuery = `
		query GetComponent($id: Uuid!, $sbomId: Uuid!) {
			component(id: $id, sbomId: $sbomId) {
				id
				name
				version
				kind
				purl
				cpes
				licensesExp
				group
				description
				primary
				internal
				sbomId
				updatedAt
				sbom {
					id
					projectVersion
					project {
						id
						name
					}
				}
			}
		}
	`

	// SbomVulnsQuery fetches vulnerabilities for an SBOM
	SbomVulnsQuery = `
		query GetSbomVulns($sbomId: Uuid!, $first: Int, $after: String, $severity: [String!], $status: [String!], $kev: Boolean, $search: String) {
			sbom(sbomId: $sbomId) {
				vulns(sbomId: $sbomId, first: $first, after: $after, severity: $severity, status: $status, kev: $kev, search: $search) {
					nodes {
						id
						componentId
						vulnId
						sbomId
						fixedIn
						fixedVersions
						detail
						impact
						actionStmt
						createdAt
						updatedAt
						component {
							id
							name
							version
							purl
						}
						vuln {
							id
							vulnId
							desc
							sev
							cvssScore
							cvssVector
							source
							publishedAt
							lastModifiedAt
							vulnInfo {
								cveId
								epssScore
								epssPercentile
								kev
								cwes
							}
						}
						vexStatus {
							id
							name
						}
						vexJustification {
							id
							name
						}
					}
					totalCount
					pageInfo {
						hasNextPage
						hasPreviousPage
						startCursor
						endCursor
					}
				}
			}
		}
	`

	// VulnQuery fetches a single vulnerability by internal UUID
	VulnQuery = `
		query GetVuln($id: Uuid!) {
			vuln(id: $id) {
				id
				vulnId
				desc
				sev
				cvssScore
				cvssVector
				source
				publishedAt
				lastModifiedAt
				updatedAt
				vulnInfo {
					id
					cveId
					epssScore
					epssPercentile
					kev
					cwes
					advisories
				}
			}
		}
	`

	// CveLookupQuery fetches vulnerability info by CVE ID
	CveLookupQuery = `
		query CveLookup($vulnId: String!) {
			cveLookup(vulnId: $vulnId) {
				vulnId
				description
				severity
				published
				lastModified
				cvssScore
				cvssVector
				cwes
				advisories
			}
		}
	`

	// ComponentVulnsQuery fetches component vulnerabilities with filters
	ComponentVulnsQuery = `
		query GetComponentVulns($first: Int, $after: String, $severity: [String!], $status: [String!], $kev: Boolean, $search: String, $projectIds: [Uuid!], $projectGroupIds: [Uuid!]) {
			componentVulns(first: $first, after: $after, severity: $severity, status: $status, kev: $kev, search: $search, projectIds: $projectIds, projectGroupIds: $projectGroupIds) {
				nodes {
					id
					componentId
					vulnId
					sbomId
					fixedIn
					fixedVersions
					createdAt
					updatedAt
					component {
						id
						name
						version
						purl
						sbomId
					}
					vuln {
						id
						vulnId
						desc
						sev
						cvssScore
						source
						vulnInfo {
							epssScore
							epssPercentile
							kev
						}
					}
					vexStatus {
						id
						name
					}
				}
				totalCount
				pageInfo {
					hasNextPage
					hasPreviousPage
					startCursor
					endCursor
				}
			}
		}
	`

	// PoliciesQuery fetches policies with pagination
	PoliciesQuery = `
		query GetPolicies($first: Int, $after: String, $search: String) {
			policies(first: $first, after: $after, search: $search) {
				nodes {
					id
					name
					description
					isEnabled
					resultType
					updatedAt
				}
				totalCount
				pageInfo {
					hasNextPage
					hasPreviousPage
					startCursor
					endCursor
				}
			}
		}
	`

	// PolicyQuery fetches a single policy by ID with rules
	PolicyQuery = `
		query GetPolicy($id: Uuid!) {
			policy(id: $id) {
				id
				name
				description
				isEnabled
				resultType
				updatedAt
				policyRules {
					id
					name
					subject
					operator
					value
				}
			}
		}
	`

	// PolicyResultsQuery fetches policy evaluation results
	PolicyResultsQuery = `
		query GetPolicyResults($first: Int, $after: String, $policyId: [Uuid!], $sbomId: [Uuid!], $resultType: [String!]) {
			policyResults(first: $first, after: $after, policyId: $policyId, sbomId: $sbomId, resultType: $resultType) {
				nodes {
					id
					policyId
					sbomId
					resultType
					result
					createdAt
					policy {
						id
						name
					}
					sbom {
						id
						projectVersion
						project {
							id
							name
						}
					}
				}
				totalCount
				pageInfo {
					hasNextPage
					hasPreviousPage
					startCursor
					endCursor
				}
			}
		}
	`

	// LicensesQuery fetches licenses with pagination
	LicensesQuery = `
		query GetLicenses($first: Int, $after: String, $status: [String!], $search: String) {
			organization {
				licenses(first: $first, after: $after, status: $status, search: $search) {
					nodes {
						id
						content {
							... on License {
								shortId
								name
							}
							... on LicenseCustom {
								spdxId
								name
							}
						}
						state
						copyLeft
						osiApproved
						fsfLibre
						deprecated
						attribution
						sourceDistribution
						modifications
					}
					totalCount
					pageInfo {
						hasNextPage
						hasPreviousPage
						startCursor
						endCursor
					}
				}
			}
		}
	`

	// SbomDriftQuery compares two SBOMs
	SbomDriftQuery = `
		query GetSbomDrift($sourceSbomId: Uuid!, $targetSbomId: Uuid!) {
			sbom(sbomId: $sourceSbomId) {
				sbomDrift(targetSbomId: $targetSbomId) {
					diffType
					diffTags
					subjectComponentId
					targetComponentId
					subjectComponent {
						id
						name
						version
						purl
					}
					targetComponent {
						id
						name
						version
						purl
					}
				}
			}
		}
	`

	// VexStatusesQuery fetches all VEX statuses
	VexStatusesQuery = `
		query GetVexStatuses {
			vexStatuses {
				id
				name
			}
		}
	`

	// VexJustificationsQuery fetches all VEX justifications
	VexJustificationsQuery = `
		query GetVexJustifications {
			vexJustifications {
				id
				name
			}
		}
	`
)
