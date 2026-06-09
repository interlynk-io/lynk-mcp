# Customer MCP Gaps: Ticketing and Vulnerability Triage

As of 2026-06-09, this document captures customer-reported blockers around AI-assisted vulnerability management, verifies what can be confirmed from the `lynk-mcp` codebase, and turns the issues into an implementation plan.

This review is local-code verification only across `lynk-mcp` and the sibling `../lynk-api` checkout. The reported `500 Internal server error` cases and timeout behavior still need real request/response captures or a reproducible staging dataset to confirm API-side root cause.

## Customer Reports

### Ticketing status blockers

1. `get_ticketing_status(product_id=...)` returns `500 Internal server error` after retries for every repo tested:
   - `platform-ai-agent`, 0 versions
   - `bngai-web-service`, 158 versions
   - `geosurvey-backend`, 228 versions
   - Customer conclusion: not size-related; blocks per-repo, on-demand fetching.

2. Org-wide `get_ticketing_status` cannot paginate:
   - Response reports `productsTotalCount: 178` and `productsHasMore: true`.
   - MCP input has `products_limit`, but no cursor or offset.
   - `products_limit: 200` times out.
   - Approximately 78 repos are unreachable from MCP.

3. No light source for ticketing fields:
   - `importStatus` and per-environment JIRA `issueTrackerSettings` are only exposed through `get_ticketing_status`.
   - `list_products`, `get_product`, and `get_environment` do not expose those fields.
   - `get_ticketing_status` also scans component vulnerability ticket links, reported at 270k+ links per org-wide call.

4. Customer asks for any one unblocking path:
   - Fix per-product `get_ticketing_status` 500.
   - Or add `importStatus` and JIRA config to `list_products`.
   - Or add a products cursor to `get_ticketing_status`.

5. Product-level JIRA defaults need clarification:
   - Customer asks how the UI's product-level "JIRA Defaults" map to the API's per-environment `issueTrackerSettings`.

### Vulnerability triage improvements

Must-haves:

1. Add pagination to `list_vulnerabilities`.
2. Filter `list_vulnerabilities` by component purl or component ID.
3. Add bulk VEX write instead of one `update_component_vex` call per finding.

High value:

4. Populate `fixedIn` consistently.
5. Add direct version lookup by product name and version string.
6. Return richer SBOM-matched component details when filtering vulnerabilities by component.

Nice to have:

7. Add per-component vulnerability summary to `get_version`.
8. Trigger CycloneDX VEX export programmatically and return a download URL.

## Local Verification

### Confirmed in `lynk-mcp`

1. `get_ticketing_status` exposes no cursor.
   - Tool schema has `product_id`, `products_limit`, `policies_limit`, and `ticket_links_limit`, but no `products_after`, `policies_after`, or `ticket_links_after`.
   - Code: `internal/mcp/server.go:405`.
   - API input also lacks cursor fields.
   - Code: `internal/api/policies.go:95`.

2. The GraphQL ticketing queries request cursors but the client does not return them.
   - `TicketingStatusQuery` and `ProductTicketingStatusQuery` request `pageInfo { hasNextPage endCursor }`.
   - `TicketingStatus` only exposes `ProductsHasNextPage`, `PoliciesHasNextPage`, and `TicketsHasNextPage`, not end cursors.
   - Code: `internal/graphql/queries.go:564`, `internal/api/policies.go:103`.

3. `get_ticketing_status` scans component vulnerability ticket links by default.
   - Org-wide query requests `componentVulns(first: $ticketsFirst)` and includes `externalIssueTrackerLinks`.
   - Default `ticket_links_limit` is 500, but `totalCount` can represent the full org scan surface.
   - Code: `internal/graphql/queries.go:667`, `internal/api/policies.go:465`.

4. `importStatus` and `issueTrackerSettings` are currently ticketing-status-only fields.
   - Ticketing query asks for `importedRepository.importStatus`.
   - Ticketing query asks for `projects.externalIssueTrackerSettings`.
   - Product queries used by `list_products`, `get_product`, and `get_environment` do not request those fields.
   - Code: `internal/graphql/queries.go:41`, `internal/graphql/queries.go:78`, `internal/graphql/queries.go:95`, `internal/graphql/queries.go:564`.

5. `list_products` has API-level cursor support but MCP does not expose it.
   - Tool schema only exposes `search` and `limit`.
   - `api.ListProductsInput` already has `After`.
   - Code: `internal/mcp/server.go:66`, `internal/api/products.go:52`.

6. `list_vulnerabilities` has API-level cursor support but MCP does not expose it.
   - Tool schema only exposes `limit`, not `after`.
   - Handler constructs `ListVersionVulnsInput` with `VersionID` and `First`, but never fills `After`.
   - API input and GraphQL query already support `After` and return `EndCursor`.
   - Code: `internal/mcp/server.go:192`, `internal/mcp/tools.go:583`, `internal/api/vulnerabilities.go:90`, `internal/graphql/queries.go:330`.

7. `list_vulnerabilities` cannot filter by component ID or purl today.
   - Tool schema has free-text `search`, severity, VEX status, KEV, EPSS, CVSS, and match mode.
   - `ListVersionVulnsInput` has no `ComponentID` or `Purl`.
   - The GraphQL `sbom.vulns` query has no component filter variable.
   - Code: `internal/mcp/server.go:192`, `internal/api/vulnerabilities.go:90`, `internal/graphql/queries.go:330`.

8. Bulk VEX write is not implemented in MCP, but the API appears to already expose bulk mutations.
   - MCP exposes only `update_component_vex`.
   - `lynk-mcp` API client exposes only `UpdateComponentVex`.
   - `lynk-mcp` GraphQL query definitions include only `ComponentVexUpdateMutation`.
   - MCP code: `internal/mcp/server.go:221`, `internal/api/mutations.go:290`, `internal/graphql/queries.go:301`.
   - Sibling API exposes `component_vex_bulk_update` and `vuln_vex_bulk_update`.
   - API code: `../lynk-api/app/graphql/types/mutation_type.rb:71`, `../lynk-api/app/graphql/types/mutation_type.rb:103`, `../lynk-api/app/graphql/mutations/component_vex_bulk_update.rb:10`, `../lynk-api/app/graphql/mutations/vuln_vex_bulk_update.rb:12`.

9. `fixedIn` is already requested and surfaced by MCP, but data quality depends on API/data population.
   - `fixedIn` and `fixedVersions` are requested for vulnerability list queries.
   - `formatComponentVulns` returns `fixedIn`.
   - Code: `internal/graphql/queries.go:338`, `internal/api/vulnerabilities.go:138`, `internal/mcp/tools.go:1682`.

10. `get_version` stats are totals only.
    - MCP returns `version.Stats.VulnStats` without a per-component breakdown.
    - Code: `internal/mcp/tools.go:250`.

11. `sbom.vulns` in `../lynk-api` supports some filters that MCP does not expose, but not component ID or purl on that path.
    - `SbomVulnsResolver` supports `component_name`, `selected_sbom_id`, `vulnerability_source`, `direct`, `vex_complete`, and `assigned_before`.
    - It does not expose `component_id`, `component_ids`, or purl filtering.
    - API code: `../lynk-api/app/graphql/resolvers/sbom_vulns_resolver.rb:14`, `../lynk-api/app/graphql/resolvers/sbom_vulns_resolver.rb:20`, `../lynk-api/app/graphql/resolvers/sbom_vulns_resolver.rb:32`.

12. Org/global `componentVulns` in `../lynk-api` supports component ID filtering, but MCP's version-scoped `list_vulnerabilities` does not use that API path.
    - `ComponentVulnsResolver` has `component_ids`.
    - API code: `../lynk-api/app/graphql/resolvers/component_vulns_resolver.rb:36`, `../lynk-api/app/graphql/resolvers/component_vulns_resolver.rb:113`.

13. Programmatic SBOM download exists in the API, but MCP has no tool for it.
    - `Sbom.download` accepts `include_vulns`, `spec`, `spec_version`, and processing-stage readiness checks.
    - API code: `../lynk-api/app/graphql/types/sbom_type.rb:285`, `../lynk-api/app/graphql/resolvers/sbom_download_resolver.rb:14`, `../lynk-api/app/graphql/resolvers/sbom_download_resolver.rb:19`, `../lynk-api/app/graphql/resolvers/sbom_download_resolver.rb:29`.

14. Direct version search exists in the API, but MCP only exposes environment-scoped `list_versions`.
    - API has `project_versions` for org-wide version-string search.
    - MCP requires `environment_id` for `list_versions`.
    - API code: `../lynk-api/app/graphql/types/query_type.rb:115`, `../lynk-api/app/graphql/resolvers/project_versions_resolver.rb:11`.
    - MCP code: `internal/mcp/server.go:90`, `internal/api/versions.go:101`.

15. Product-level "JIRA Defaults" do not appear to be a field on `ProjectGroup`.
    - API `ExternalIssueTrackerSetting` is project/environment-scoped and contains the JIRA defaults used for ticket creation.
    - API `OrganizationProjectSettingDefault` contains project-setting defaults, but no JIRA issue tracker fields.
    - API code: `../lynk-api/app/graphql/types/external_issue_tracker_setting_type.rb:10`, `../lynk-api/app/graphql/types/external_issue_tracker_setting_type.rb:20`, `../lynk-api/app/graphql/types/external_issue_tracker_setting_type.rb:35`, `../lynk-api/app/graphql/types/organization_project_setting_default_type.rb:15`.

### Not yet verified

1. Per-product `get_ticketing_status` 500:
   - Needs real product IDs and request/response captures.
   - Likely root causes to test:
     - `projectGroup(id: $productId)` resolver issue.
     - Nested `componentVulns` under `projectGroup`.
     - `importedRepository` union fields.
     - `externalIssueTrackerSettings` field on `projects`.
     - Policy inclusion traversal.

2. Org-wide `products_limit: 200` timeout:
   - Needs staging/prod reproduction with timing and GraphQL tracing.
   - Likely aggravated by combining product config, policy inclusions, and ticket-link scan in one query.

3. UI "JIRA Defaults" mapping:
   - `../lynk-api` indicates ticketing defaults are per-environment/project `externalIssueTrackerSettings`.
   - `OrganizationProjectSettingDefault` is an org-level default for project settings, not JIRA issue tracker defaults.
   - The UI still needs review to explain why it labels this as product-level "JIRA Defaults".

## Implementation Plan

### Phase 0: Capture and reproduce

- [ ] Request full request/response captures for the three per-product 500s.
- [x] Recover real product IDs for the three reported products.
- [x] Attempt to reproduce `get_ticketing_status(product_id=...)` against the current MCP/API environment.
- [ ] Reproduce the reported per-product 500.
- [x] Attempt to reproduce org-wide `get_ticketing_status(products_limit: 200)` timeout against the current MCP/API environment.
- [ ] Reproduce the reported org-wide timeout with GraphQL timing.
- [ ] Identify which nested field causes the per-product 500 by running reduced GraphQL queries.
- [x] Confirm UI "JIRA Defaults" source of truth in the UI repository.

Phase 0 progress, 2026-06-09:

- Product IDs were recovered from `list_products(search=...)`:
  - `platform-ai-agent`: `e16cbb23-23a8-42a1-936c-9bb93c7d5219`, 0 versions.
  - `bngai-web-service`: `6107914e-1f34-4c95-b7aa-0de04dd1605d`, 158 versions.
  - `geosurvey-backend`: `c03a1346-d243-4593-9e59-9c849ce9913d`, 227 versions in the current API response.
- `get_ticketing_status(product_id=...)` did not reproduce the reported 500 for any of the three products from the current MCP/API environment:
  - `platform-ai-agent` returned successfully in about 3.9s, with `ticketsScannedCount: 0`.
  - `bngai-web-service` returned successfully in about 0.8s, with `ticketsScannedCount: 13361` and `ticketsHasMore: true`.
  - `geosurvey-backend` returned successfully in about 1.1s, with `ticketsScannedCount: 20083` and `ticketsHasMore: true`.
- Org-wide `get_ticketing_status(products_limit: 200, policies_limit: 50, ticket_links_limit: 500)` did not reproduce the reported timeout from the current MCP/API environment:
  - Returned successfully in about 1.5s.
  - Returned `productsTotalCount: 178` and `productsHasMore: false`, so all currently visible products were reachable with limit 200.
  - Confirmed the expensive scan surface: `ticketsScannedCount: 270473` and `ticketsHasMore: true`.
- Because the 500 did not reproduce, the failing nested field is still unknown. The successful product-level responses exercised the full current nested query, including:
  - `organization.connections`
  - `jiraVulnManagementConfig`
  - `projectGroup(id: $productId)`
  - `importedRepository`
  - `projects.externalIssueTrackerSettings`
  - nested `projectGroup.componentVulns.externalIssueTrackerLinks`
  - `policies.policyInclusions`
- UI source-of-truth check:
  - The `JIRA Defaults` card is rendered by `../lynk-dash-app/src/components/forms/fields/JiraFields.jsx`.
  - It is fed by `GetProjectWithExternalIssueTrackerSettings` / `GetProjectSettings` in `../lynk-dash-app/src/graphQL/products/queries.js`.
  - The data source is `project(id: ...).externalIssueTrackerSettings`, confirmed by `../lynk-api/app/graphql/types/project_type.rb`.
  - `ExternalIssueTrackerSettingType` defines the JIRA defaults fields: `projectKey`, `assignee`, `reporter`, `issueType`, `epic`, `components`, and `enableJiraSync`.
  - Conclusion: UI "JIRA Defaults" are environment/project-scoped `externalIssueTrackerSettings`; they are not a product-level `ProjectGroup` field.

### Phase 1: Quick MCP-only fixes

- [x] Add `after` to `list_products` tool schema and handler.
- [x] Return `endCursor` from `list_products` MCP output.
- [x] Add `after` to `list_vulnerabilities` tool schema and handler.
- [x] Return `endCursor` from `list_vulnerabilities` MCP output.
- [x] Add tests for both pagination paths.
- [x] Update README tool docs.

Expected impact: automated clients can page through products and vulnerability findings where the underlying API already supports cursors.

### Phase 2: Ticketing status decomposition

- [x] Add cursor fields to `api.TicketingStatusInput`:
  - `ProductsAfter`
  - `PoliciesAfter`
  - `TicketsAfter`
- [x] Add matching GraphQL variables to ticketing queries.
- [x] Capture and return `productsEndCursor`, `policiesEndCursor`, and `ticketsEndCursor`.
- [x] Add MCP tool params:
  - `products_after`
  - `policies_after`
  - `ticket_links_after`
- [x] Return those cursors in `get_ticketing_status`.
- [x] Consider defaulting `ticket_links_limit` to `0` or adding `include_created_tickets` so config lookups do not scan ticket links unless requested.
- [x] Add tests covering org-wide pagination metadata and per-product behavior.

Expected impact: fixes the unreachable products problem and separates config discovery from expensive ticket-link discovery.

### Phase 3: Lightweight ticketing fields on product tools

- [x] Extend product GraphQL queries to include `importedRepository.importStatus` where schema permits it.
- [x] Add a compact JIRA config projection to `get_product` and `get_environment` if `externalIssueTrackerSettings` is available on those API types.
- [x] If UI review finds a product-level abstraction over environment settings, add explicit fields with names that match the API model rather than overloading `issueTrackerSettings`.
- [x] Keep `list_products` response compact; include only repo import status and default ticketing summary, not full environment settings for every product.
- [x] Add tests for GitHub, GitLab, Bitbucket, and nil repository cases.

Expected impact: customers can inspect repo import status and JIRA defaults without invoking heavy ticketing status scans.

### Phase 4: Component-scoped vulnerability triage

- [x] Check `lynk-api` schema for existing component filters on `sbom.vulns` or `componentVulns`.
- [x] Add MCP support for API-supported component ID filtering, likely via `componentVulns(componentIds: ...)` or a version-scoped API addition.
- [ ] Add API support for purl filtering if exact purl filtering is required server-side.
- [x] Treat purl filtering carefully:
  - Prefer exact purl match.
  - Define behavior for empty purl.
  - Decide whether purl is scoped by version ID, org, or component ID.
- [x] Return stable component identifiers in the vulnerability response:
  - `component.id`
  - `component.name`
  - `component.version`
  - `component.purl`
  - `component.bomRef` if available
  - `component.sbomId` or `versionId`
- [x] Add tests for component ID, purl, empty purl, and pagination combined with component filters.

Notes:
- `componentVulns(componentIds: ...)` is used for org-wide `search_vulnerabilities`; version-scoped `list_vulnerabilities` filters component IDs within `sbom.vulns` results.
- Exact `purl` filtering is implemented in MCP response handling. Server-side purl filtering remains open until the API exposes it.
- `component.bomRef` was not added because it is not exposed by the current API schema; responses include component id, name, version, purl, and `sbomId`/`versionId` where available.

Expected impact: customers can enumerate all findings for a component and perform class disposition with auditability.

### Phase 5: Bulk VEX write

- [x] Existing API bulk mutations found: `component_vex_bulk_update` and `vuln_vex_bulk_update`.
- [x] Add MCP client support for the existing API bulk mutation rather than designing a new API mutation first.
- [x] Decide whether the accepted payload should be native Interlynk update input, OpenVEX-like statements, CSAF-like statements, or a small normalized wrapper.
- [x] Add `bulk_update_component_vex` MCP tool with `confirm=true`.
- [x] Return per-item success/error results; avoid all-or-nothing unless the API guarantees transaction semantics.
- [x] Add partial failure tests.
- [x] Add rate-limit and retry behavior appropriate for bulk writes.

Notes:
- `bulk_update_component_vex` uses the existing `componentVexBulkUpdate` API mutation with a small MCP wrapper: `component_vuln_ids` plus the same shared update fields as `update_component_vex`.
- The tool returns requested, updated, and failed counts. Updated items are returned as component vulnerability records; requested IDs missing from the API response are reported as failed with API errors when present.
- Bulk writes are sent as one API mutation rather than a client-side loop. Existing GraphQL transient retry handling covers 5xx, HTTP 429, and network-level failures.

Expected impact: triage agents can write approved dispositions in one reviewable operation instead of issuing many single-finding mutations.

### Phase 6: Data-quality and discovery improvements

- [ ] Investigate why `fixedIn` is usually empty even though MCP requests it.
- [ ] Confirm whether `fixedVersions` has better population and should be emphasized in responses.
- [ ] Add direct version lookup by product name plus version string:
  - API has org-wide `project_versions` search by version string.
  - MCP still needs product/environment disambiguation for an exact lookup helper.
- [ ] Add per-component vulnerability summary to `get_version`:
  - Prefer API-side aggregation for performance.
  - Include severity counts and total count per component.
- [ ] Add programmatic CycloneDX VEX export:
  - API has `Sbom.download`; verify whether `include_vulns` plus CycloneDX spec/spec version satisfies CycloneDX VEX needs.
  - Return readiness/status plus content metadata or download content according to API behavior.
  - Include permission and expiration behavior in docs.

## Priority Recommendation

1. Ship MCP cursor exposure for `list_vulnerabilities` and `list_products` first. These are low-risk because the API layer already has cursor support.
2. Add pagination and an option to skip ticket-link scans in `get_ticketing_status`. This directly addresses the customer's unreachable repo issue and timeout risk.
3. Reproduce and fix the per-product 500. If the fault is API-side, reduce the MCP query as a temporary workaround while the resolver is fixed.
4. Add component-scoped vulnerability filters and bulk VEX write next. These are the core changes needed for reliable AI-assisted triage at customer scale.
