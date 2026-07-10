# Marketplace Directory
*Last synced: 2026-07-09 · Source: [`anthropics/claude-plugins-official`](https://github.com/anthropics/claude-plugins-official) marketplace manifest · Evaluation passes: 2026-07-03 (52 desk-checked; 15 exercised hands-on), 2026-07-08 (204 desk-checked)*

The lookup directory for official-marketplace plugins that have NOT been promoted to the solutions catalog — consulted by grep on stack or goal relevance ("user said SAPUI5 → grep"), never enumerated or read whole. All plugins here are installable via `/plugin install <name>@claude-plugins-official`; none are installed by default. Scope decision (2026-07-03): externally-maintained plugins listed in the official marketplace ARE in scope — "official" means Anthropic-curated, not Anthropic-authored.

**Promotion rule:** a plugin graduates out of this directory into its own `solutions/<id>.md` record when it is hands-on validated (author-agnostic — our own evaluation is the evidence that matters) and earns a ranked row in at least one `problems/<goal>.md` table, where it competes with every other approach on equal terms. Promoted plugins are ordinary solutions — teachable by default, indexed, profile-tracked, never listed here. This directory holds the rest: the long tail plus everything not (yet) meeting that bar, reachable only by stack-match grep or reactive mention. Together, this file's rows and the `kind: plugin` records in `solutions/` cover the full marketplace manifest — the drift check verifies the union.

Each row: backticked name = the plugin's `id` (a legal profile capability id — the profile doesn't care where a capability is cataloged), goal column = the `problems/<goal>.md` it maps to, verdict = trust tier. Directory plugins are NOT teachable-by-default: they enter a user's ignorance map only on stack or goal relevance (a user who never touches SAP is never taught `ui5-modernization` as their daily surprise).

Verdicts are produced by the repeatable protocol in `evals/plugin-evaluation.md` — same fixture, same per-plugin exercises, same criteria on every run, so evaluations stay comparable over time.

**Verdict legend** — every plugin carries one:

- ✅ **hands-on (date)** — installed, exercised against its mapped goal, and it worked; caveats noted verbatim
- ☑️ **desk-checked** — manifest, components, freshness, and provenance reviewed (2026-07-03); not exercised. For MCP integrations this usually means hands-on needs an external account or infrastructure we don't have — an honest label, not a defect
- ⚠️ **caution** — works, but overlaps a built-in feature or has a sharp edge; lead with the alternative named

The mentor recommends ✅ plugins freely, offers ☑️ ones with the "not hands-on evaluated" label, and never presents a ⚠️ without its caveat.

## Anthropic-built plugins

### Dev workflow

| Plugin | What it does | Relevant goal | Verdict |
|--------|-------------|--------------|---------|
| `code-review` | Multi-agent PR review with confidence scoring | `code-review` | ⚠️ duplicates the built-in `/code-review`, `/review`, and `/code-review ultra` — recommend the built-ins first |
| `code-simplifier` | Agent for clarity and maintainability refactors | `refactoring` | ⚠️ overlaps the built-in `/simplify` skill — recommend the built-in first |
| `ralph-loop` | Continuous while-true agent loops re-running the same prompt until completion | `migration` | ⚠️ overlaps the built-in `/loop` and `/goal` — recommend the built-ins first |
| `playground` | Interactive single-file HTML playgrounds with visual controls and live preview | `greenfield` | ☑️ desk-checked — partially overlaps the built-in Artifact tool for shareable pages |

### Hooks & output styles

| Plugin | What it does | Relevant goal | Verdict |
|--------|-------------|--------------|---------|
| `explanatory-output-style` | SessionStart hook injecting educational insights about implementation choices | `onboarding` | ☑️ desk-checked — mimics a deprecated output style; niche |
| `learning-output-style` | Prompts users to write meaningful code at decision points | `onboarding` | ☑️ desk-checked — mimics an unshipped output style; niche |

### Plugin & SDK development

| Plugin | What it does | Relevant goal | Verdict |
|--------|-------------|--------------|---------|
| `skill-creator` | Creates and improves skills, measures skill performance | `building-skills-plugins` | ☑️ desk-checked — in active daily use by this catalog's maintainer, which is stronger evidence than most desk checks |
| `mcp-tunnels` | Connects Claude to a private MCP server through an Anthropic MCP tunnel | `building-mcp-integrations` | ☑️ desk-checked — needs Docker Compose infrastructure to exercise |

### Project & session management

| Plugin | What it does | Relevant goal | Verdict |
|--------|-------------|--------------|---------|

### Language servers (LSPs)

Drop-in LSP integrations for code intelligence: `clangd-lsp` (C/C++), `csharp-lsp`, `gopls-lsp` (Go), `jdtls-lsp` (Java), `kotlin-lsp`, `lua-lsp`, `php-lsp`, `pyright-lsp` (Python), `ruby-lsp`, `rust-analyzer-lsp`, `swift-lsp`, `typescript-lsp`.

☑️ desk-checked as a family — uniform official wrappers around standard language servers, low risk; each requires its language-server binary on `$PATH` (the plugin errors visibly if missing). Recommend freely when the user's language matches and the binary exists or is easily installed.

### Specialty

Rarely relevant to everyday engineering, listed for completeness (☑️ desk-checked): `math-olympiad` (competition math solving with adversarial proof verification) and `cwc-makers` (onboarding for the Code-with-Claude Makers Cardputer hardware kit).

## External plugins (partner-maintained)

Hands-on evaluation of most integrations requires accounts or infrastructure (Slack workspaces, Figma files, cloud projects); those carry ☑️ with that caveat rather than a pretend verdict.

| Plugin | What it does | Relevant goal | Verdict |
|--------|-------------|--------------|---------|
| `github` (GitHub) | Official GitHub MCP: issues, PRs, code review, repo management | `code-review` | ☑️ desk-checked — first-party GitHub; needs repo auth to exercise |
| `gitlab` (GitLab) | GitLab MCP: merge requests, CI/CD, pipelines, issues | `ci-automation` | ☑️ desk-checked — first-party GitLab; needs instance auth |
| `playwright` (Microsoft) | Browser automation and E2E testing MCP server | `testing` | ☑️ desk-checked — first-party Microsoft; needs browser install; note the built-in Chrome integration covers some of this |
| `serena` (Oraios) | Semantic code analysis MCP for refactoring and code understanding | `code-understanding` | ☑️ desk-checked — note built-in LSP plugins cover much of the navigation ground |
| `greptile` (Greptile) | AI PR review agent for GitHub and GitLab | `code-review` | ☑️ desk-checked — needs a Greptile account; overlaps built-in review skills |
| `linear` (Linear) | Linear issue tracking: create issues, manage projects, search | `devops` | ☑️ desk-checked — needs workspace auth |
| `asana` (Asana) | Create and manage tasks, search projects, update assignments | `devops` | ☑️ desk-checked — needs workspace auth |
| `firebase` (Google) | Firestore, auth, cloud functions, and hosting via Firebase MCP | `devops` | ☑️ desk-checked — needs a Firebase project |
| `terraform` (HashiCorp) | Terraform MCP for IaC registry integration and module management | `devops` | ☑️ desk-checked — first-party HashiCorp, fresh (2026-03) |
| `laravel-boost` (Laravel) | Laravel development toolkit MCP server | `greenfield` | ☑️ desk-checked — first-party Laravel; needs a Laravel app |
| `telegram` | Telegram messaging bridge with access control (channels) | `devops` | ☑️ desk-checked — fresh (2026-04); covered by `solutions/channels.md`; needs a bot token |
| `discord` | Discord messaging bridge with access control (channels) | `devops` | ☑️ desk-checked — fresh (2026-04); covered by `solutions/channels.md`; needs a bot |
| `imessage` | iMessage bridge (reads `chat.db`, sends via AppleScript; channels) | `devops` | ☑️ desk-checked — macOS only, Full Disk Access required; covered by `solutions/channels.md` |
| `fakechat` | Localhost chat UI for testing channel flows — no tokens, no access control | `testing` | ☑️ desk-checked — the intended channels demo; requires Bun |

## Externally-hosted plugins

Plugins listed in the marketplace manifest whose source lives in the author's own repository (SHA-pinned by the marketplace). All desk-checked per `evals/plugin-evaluation.md`; hands-on evaluation generally needs the author's product or an account, so verdicts stay ☑️ until a plugin lands on a recommendation path.

### Automation

| Plugin | What it does | Author | Relevant goal | Verdict |
|--------|-------------|--------|--------------|---------|
| `zyte-web-data` | Zyte API web scraping: spiders, extraction schemas, Scrapy Cloud deployment | Zyte | `research` | ☑️ desk-checked 2026-07-08 — 15 skills; active 2026-07; needs Zyte API account |

### Database

| Plugin | What it does | Author | Relevant goal | Verdict |
|--------|-------------|--------|--------------|---------|
| `alloydb` | Create, connect, and query AlloyDB for PostgreSQL databases | Google LLC | `devops` | ☑️ desk-checked 2026-07-08 — 7 skills (MCP Toolbox fetched via skills, none declared in manifest); active 2026-07; needs Google Cloud project with AlloyDB |
| `alloydb-omni` | Create, connect, and query AlloyDB Omni databases | Google LLC | `devops` | ☑️ desk-checked 2026-07-08 — 9 skills (incl. Kubernetes/container ops, none declared as MCP in manifest); active 2026-07; needs AlloyDB Omni installation |
| `azure-cosmos-db-assistant` | Azure Cosmos DB data modeling, query optimization, and performance tuning | unlabeled | `devops` | ☑️ desk-checked 2026-07-08 — MCP server (.mcp.example.json) + 1 skill + 1 agent + 3 commands; active 2026-04; needs Azure Cosmos DB account |
| `bigdata-com` | Financial research, analytics, and intelligence tools powered by Bigdata MCP | RavenPack | `research` | ☑️ desk-checked 2026-07-08 — MCP server (.mcp.json) + 1 skill + 27 commands; active 2026-06; needs Bigdata.com account |
| `bigquery-data-analytics` | Connect, query, and generate insights from BigQuery datasets | Google LLC | `devops` | ☑️ desk-checked 2026-07-08 — 3 skills (MCP toolbox-based); active 2026-07; needs GCP/BigQuery account |
| `clickhouse` | ClickHouse Cloud MCP for schema browsing, read-only SQL, backups, billing, ClickPipes | ClickHouse | `devops` | ☑️ desk-checked 2026-07-08 — MCP server + 2 skills; active 2026-07; needs ClickHouse Cloud account |
| `clickhouse-best-practices` | 28 prioritized ClickHouse rules for schema design, query optimization, and ingestion | ClickHouse Inc | `performance` | ☑️ desk-checked 2026-07-08 — 10 skills; active 2026-07 |
| `cloud-sql-mysql` | Connect and interact with Cloud SQL for MySQL databases | Google LLC | `devops` | ☑️ desk-checked 2026-07-08 — 4 skills; active 2026-07; needs GCP Cloud SQL instance |
| `cloud-sql-postgresql` | Create, connect, and query Cloud SQL for PostgreSQL databases | Google LLC | `devops` | ☑️ desk-checked 2026-07-08 — 8 skills; active 2026-07; needs GCP project + Cloud SQL PostgreSQL instance |
| `cloud-sql-sqlserver` | Connect to Cloud SQL for SQL Server databases | Google LLC | `devops` | ☑️ desk-checked 2026-07-08 — 4 skills; active 2026-07; needs GCP project + Cloud SQL SQL Server instance |
| `cockroachdb` | CockroachDB cluster integration: schema exploration, SQL optimization, cluster management | Cockroach Labs | `devops` | ☑️ desk-checked 2026-07-08 — MCP server (2 backends) + 3 agents + skills + safety hooks; active 2026-06; needs CockroachDB cluster |
| `convex` | Convex backend development: schema design, real-time features, auth, scheduled jobs | Convex | `greenfield` | ☑️ desk-checked 2026-07-08 — MCP server + 7 skills + 2 agents + hooks + runtime-error monitor; active 2026-07; needs Convex deployment |
| `databases-on-aws` | AWS database guidance for schema design, queries, migrations, and engine selection | unlabeled | `devops` | ☑️ desk-checked 2026-07-07 — 1 skill (dsql) + MCP server + hooks; active 2026-06; needs AWS account |
| `databricks` | Databricks skills for CLI, apps, model serving, pipelines, and serverless migration | Databricks | `devops` | ☑️ desk-checked 2026-07-07 — 28 skills + 2 commands + hooks; active 2026-07; needs Databricks workspace |
| `datahub-skills` | DataHub toolkit for catalog search, metadata enrichment, lineage, and data quality | DataHub | `devops` | ☑️ desk-checked 2026-07-07 — 12 skills + 8 commands + 4 agents; active 2026-05; needs DataHub instance |
| `dataproc` | Manage Google Cloud Dataproc clusters and jobs | Google LLC | `devops` | ☑️ desk-checked 2026-07-07 — 1 skill; active 2026-07; needs GCP account |
| `dataverse` | Microsoft Dataverse skills with Dataverse MCP, PAC CLI, and Python SDK | unlabeled | `devops` | ☑️ desk-checked 2026-07-07 — 8 skills; active 2026-06; needs Microsoft Dataverse environment |
| `duckdb-skills` | DuckDB skills to read data files, query databases, and search docs | DuckDB Foundation | `devops` | ☑️ desk-checked 2026-07-08 — 9 skills, no MCP; last commit 2026-04 |
| `firestore-native` | Query and manage Firestore databases, collections, and documents | Google LLC | `devops` | ☑️ desk-checked 2026-07-08 — 1 skill, no MCP in plugin.json; active 2026-07; needs Google Cloud project (required userConfig) |
| `knowledge-catalog` | Discover, manage, monitor, and govern data and AI artifacts | Google LLC | `devops` | ☑️ desk-checked 2026-07-08 — 1 skill (no .mcp.json at pinned commit); active 2026-07; needs Google Cloud Knowledge Catalog access |
| `looker` | Connect to Looker and query data using LookML | Google LLC | `devops` | ☑️ desk-checked 2026-07-08 — 2 skills + userConfig; active 2026-07; needs Looker instance API credentials |
| `mongodb` | MongoDB MCP and skills for data exploration, query optimization, schema design | unlabeled | `devops` | ☑️ desk-checked 2026-07-08 — MCP server + 7 skills (connection, querying, query-optimizer, schema-design, search-and-ai, stream-processing, mcp-setup); active 2026-07; needs a MongoDB/Atlas deployment |
| `neon` | Manage Neon Postgres projects and databases via agent skill and MCP server | unlabeled | `devops` | ☑️ desk-checked 2026-07-08 — MCP server + 3 skills (neon, neon-postgres, neon-postgres-branches); active 2026-06; needs Neon account |
| `oracledb` | Connect, query, and interact with Oracle databases and their data | Google LLC | `devops` | ☑️ desk-checked 2026-07-07 — 1 skill + connection userConfig (connection string, username, password, wallet); active 2026-07; needs Oracle DB credentials |
| `pinecone` | Pinecone MCP and skills for vector index management, querying, and RAG prototyping | unlabeled | `llm-features` | ☑️ desk-checked 2026-07-07 — MCP server + 1 command + 9 skills; active 2026-05; needs Pinecone account/API key |
| `planetscale` | Hosted PlanetScale MCP for databases, branches, schema, and slow-query Insights | unlabeled | `devops` | ☑️ desk-checked 2026-07-07 — MCP server + skills (git submodule); active 2026-05; needs PlanetScale account (authenticated hosted MCP) |
| `qdrant-skills` | Qdrant skills for scaling, performance, search quality, monitoring, upgrades, multi-language SDKs | Qdrant | `devops` | ☑️ desk-checked 2026-07-07 — 9 skills; active 2026-06; needs a Qdrant deployment |
| `redis-development` | Redis best practices for data structures, vector search, caching, and performance | Redis | `devops` | ☑️ desk-checked 2026-07-08 — 8 skills; active 2026-05 |
| `sap-hana-cli` | 150+ SAP HANA database tools for queries, data, backups, and monitoring | SAP SE | `devops` | ☑️ desk-checked 2026-07-07 — MCP server (.mcp.json only); active 2026-06; needs SAP HANA Cloud or on-premise database |
| `spanner` | Natural-language querying and interaction with Google Cloud Spanner data | Google LLC | `devops` | ☑️ desk-checked 2026-07-08 — 1 skill (spanner-data); active 2026-06; needs Google Cloud Spanner instance |
| `supabase` | Supabase MCP for database, auth, storage, and real-time backend operations | unlabeled | `devops` | ☑️ desk-checked 2026-07-08 — MCP server + 2 skills + agents; active 2026-06; needs Supabase project |
| `vsql-extension-builder` | Builds VillageSQL MySQL extensions via 7-phase workflow; ports PostgreSQL extensions | VillageSQL | `migration` | ☑️ desk-checked 2026-07-08 — 1 skill (7-phase persona workflow); active 2026-06; needs a MySQL/VillageSQL build environment |
| `zilliz` | Zilliz Cloud vector database management: clusters, search, indexing, backups, monitoring | Zilliz | `devops` | ☑️ desk-checked 2026-07-08 — 21 skills + 2 commands; active 2026-05; needs Zilliz Cloud account |

### Deployment

| Plugin | What it does | Author | Relevant goal | Verdict |
|--------|-------------|--------|--------------|---------|
| `azure` | Azure MCP plus skills for resources, deployments, diagnostics, cost optimization | unlabeled | `devops` | ☑️ desk-checked 2026-07-08 — MCP server (.mcp.json) + 28 skills + hooks; active 2026-07; needs Azure account |
| `cloudflare` | Cloudflare platform skills: Workers, Durable Objects, Agents SDK, Wrangler CLI | unlabeled | `devops` | ☑️ desk-checked 2026-07-08 — 11 skills + 2 commands + MCP server; active 2026-06; needs Cloudflare account for deploys |
| `deploy-on-aws` | Deploy applications to AWS with architecture recommendations, cost estimates, and IaC | unlabeled | `devops` | ☑️ desk-checked 2026-07-07 — 3 skills + MCP server + hooks; active 2026-06; needs AWS account |
| `hostinger` | Deploy and manage Hostinger websites, domains, email, VPS, and payments | Hostinger | `devops` | ☑️ desk-checked 2026-07-08 — MCP server only (no skills); active 2026-06; needs Hostinger account (OAuth or API token) |
| `railway` | Deploy and manage apps, databases, and infrastructure on Railway | unlabeled | `devops` | ☑️ desk-checked 2026-07-08 — MCP server + 1 skill + hooks; active 2026-07; needs Railway account |
| `render` | Deploy, debug, and monitor applications on Render | Render | `devops` | ☑️ desk-checked 2026-07-08 — 21 skills + 1 agent + 2 commands + render.yaml validation hook; active 2026-05; needs Render account |
| `valtown` | Val Town MCP and skills for HTTP vals, cron, SQLite, email, deployment | Val Town | `devops` | ☑️ desk-checked 2026-07-08 — MCP server + 10 skills; active 2026-07; needs Val Town account |
| `vercel` | Vercel deployment platform: deployments, build status, logs, domains, infrastructure | unlabeled | `devops` | ☑️ desk-checked 2026-07-08 — MCP server + 28 skills, 5 commands, 3 agents, hooks; active 2026-07; needs Vercel account |

### Design

| Plugin | What it does | Author | Relevant goal | Verdict |
|--------|-------------|--------|--------------|---------|
| `adobe-for-creativity` | Adobe creative AI tools for image editing, design automation, and retouching | Adobe | `greenfield` | ☑️ desk-checked 2026-07-08 — MCP server + 7 skills; active 2026-06; needs Adobe account |
| `canva` | Create, edit, review, resize, and brand-check Canva designs via Canva MCP server | Canva | `greenfield` | ☑️ desk-checked 2026-07-08 — MCP server + 6 skills; active 2026-06; needs Canva account |
| `figma` | Figma MCP to read design files and translate designs into code | unlabeled | `greenfield` | ☑️ desk-checked 2026-07-08 — MCP server + 13 skills (11 core, 2 workflow); active 2026-07; needs Figma account |
| `hyperframes` | HeyGen HyperFrames: write HTML, render video with animations, captions, voiceovers | HeyGen | `greenfield` | ☑️ desk-checked 2026-07-08 — 21 skills (no MCP); active 2026-07; needs HyperFrames CLI/runtime |
| `miro` | Miro board access to read context, create diagrams, and generate code | Miro | `greenfield` | ☑️ desk-checked 2026-07-08 — MCP server + 7 skills (diagram, doc, table, browse, code-spec, code-review, code-explain); active 2026-06; needs Miro account |
| `runway-api` | Generate videos, images, and audio at scale with Runway's API | Runway | `greenfield` | ☑️ desk-checked 2026-07-08 — 17 skills + scripts; active 2026-04; needs Runway API key |

### Development

| Plugin | What it does | Author | Relevant goal | Verdict |
|--------|-------------|--------|--------------|---------|
| `agentforce-adlc` | Author, scaffold, deploy, test, and optimize Salesforce Agentforce .agent files | unlabeled | `building-agents` | ☑️ desk-checked 2026-07-08 — 4 skills + 4 agents + hooks; active 2026-06; needs Salesforce Agentforce org |
| `apollo-skills` | Apollo GraphQL skills for Client, Server, Federation, Router, Rover, and MCP server | Apollo GraphQL | `api-design` | ☑️ desk-checked 2026-07-08 — MCP server + LSP config + 14 skills; active 2026-07 |
| `appwrite` | Appwrite SDK skills, MCP servers, and deployment commands | Appwrite | `devops` | ☑️ desk-checked 2026-07-08 — MCP server + 11 SDK skills + 2 deploy commands; active 2026-04; needs Appwrite account |
| `astronomer-data-agents` | Airflow and Astronomer data engineering: author DAGs, debug pipelines, lineage, migrations, deployments | unlabeled | `devops` | ☑️ desk-checked 2026-07-08 — 26 skills + astro-airflow-mcp server dir; active 2026-07; needs Airflow/Astronomer environment |
| `atomic-agents` | Workflows and agents for building with the Atomic Agents framework | unlabeled | `building-agents` | ☑️ desk-checked 2026-07-08 — 6 skills + 2 agents; active 2026-06 |
| `aws-agents` | Build, deploy, and operate AI agents with Amazon Bedrock AgentCore | Amazon Web Services | `building-agents` | ☑️ desk-checked 2026-07-08 — MCP server + 7 skills; active 2026-07; needs AWS account |
| `aws-agents-for-devsecops` | Incident investigation, code review, vulnerability scans, pentests via AWS DevOps/Security Agents | Amazon Web Services | `security` | ⚠️ code-review and diff-scanning skills overlap the built-in /code-review and /security-review — recommend the built-ins first; desk-checked 2026-07-08 — MCP server + 13 skills + 9 commands; active 2026-06; needs AWS account with DevOps and Security Agents |
| `aws-amplify` | Guided AWS Amplify Gen 2 workflows for auth, data, storage, and functions | unlabeled | `greenfield` | ☑️ desk-checked 2026-07-08 — MCP server + 1 skill; active 2026-05; needs AWS account |
| `aws-core` | Skills for building, deploying, and operating applications on AWS with IaC | Amazon Web Services | `devops` | ☑️ desk-checked 2026-07-08 — MCP server + 15 skills + hooks (secret-safety); active 2026-06; needs AWS account |
| `aws-data-analytics` | Data lake, analytics, and ETL workflows with S3 Tables, Glue, Athena | Amazon Web Services | `devops` | ☑️ desk-checked 2026-07-08 — MCP server + 8 skills; active 2026-06; needs AWS account |
| `aws-dev-toolkit` | AWS toolkit for building, migrating, and reviewing cloud architectures | aws-samples | `devops` | ☑️ desk-checked 2026-07-08 — 3 MCP servers (.mcp.json) + 35 skills + 11 agents; active 2026-05; needs AWS account |
| `aws-serverless` | Design, build, deploy, test, and debug AWS serverless applications | unlabeled | `devops` | ☑️ desk-checked 2026-07-08 — MCP server (.mcp.json) + 7 skills + hooks; active 2026-06; needs AWS account |
| `aws-startup-advisor` | Startup-focused AWS architecture, cost, security, and migration guidance | Amazon Web Services | `devops` | ☑️ desk-checked 2026-07-08 — MCP server (.mcp.json) + 5 skills; active 2026-06; needs AWS account |
| `base44` | Build and deploy Base44 full-stack apps with CLI and SDK | unlabeled | `greenfield` | ☑️ desk-checked 2026-07-08 — 5 skills (CLI-driven); active 2026-07; needs Base44 account |
| `boltz` | Predict structures, screen molecules, and design binders with Boltz | Boltz | `research` | ☑️ desk-checked 2026-07-08 — 8 skills; active 2026-06; needs Boltz API access |
| `buildkite` | Official Buildkite skills for pipelines, migration, preflight, agent runtime, CLI, API | Buildkite | `ci-automation` | ☑️ desk-checked 2026-07-08 — 6 skills; active 2026-07; needs Buildkite account |
| `cds-mcp` | SAP CAP development assistant; searches CDS models and CAP documentation | SAP SE | `code-understanding` | ☑️ desk-checked 2026-07-08 — MCP server (Node, .mcp.json); active 2026-06; needs a CAP project |
| `chrome-devtools-mcp` | Chrome DevTools MCP for browser automation, performance traces, network and console inspection | unlabeled | `debugging` | ⚠️ overlaps the built-in Chrome integration — recommend the built-in first; desk-checked 2026-07-08 — MCP server + 6 skills; active 2026-07; needs Chrome |
| `circle-skills` | Circle skills and MCP for USDC payments, cross-chain transfers, wallets, smart contracts | Circle | `greenfield` | ☑️ desk-checked 2026-07-08 — MCP server + 17 skills; active 2026-06; needs Circle account |
| `codspeed` | CodSpeed performance benchmarking, flamegraphs, and profiling via MCP | CodSpeed | `performance` | ☑️ desk-checked 2026-07-08 — MCP server + 2 skills; active 2026-06; needs CodSpeed account |
| `confidence` | Access Confidence feature flags, experiments, and migration tools | Spotify Confidence | `release-management` | ☑️ desk-checked 2026-07-08 — MCP server + 11 skills + 5 commands; active 2026-06; needs Confidence account |
| `data` | Apache Airflow and Astronomer data engineering: DAG authoring, debugging, lineage, migration | unlabeled | `devops` | ☑️ desk-checked 2026-07-08 — 26 skills + bundled Airflow MCP server; active 2026-07; needs Airflow/Astronomer environment |
| `data-agent-kit-starter-pack` | GCP data engineering skills for pipelines, dbt, BigQuery SQL, and workflow orchestration | Google LLC | `devops` | ☑️ desk-checked 2026-07-07 — 20 skills + MCP server; active 2026-07; needs GCP account |
| `datarobot-agent-skills` | DataRobot skills for model training, deployment, predictions, monitoring, and explainability | DataRobot | `devops` | ☑️ desk-checked 2026-07-07 — 13 skills; active 2026-07; needs DataRobot account |
| `dominodatalab` | Domino Data Lab platform support for workspaces, jobs, model deployment, experiment tracking | Domino Data Lab | `devops` | ☑️ desk-checked 2026-07-07 — 23 skills + 4 commands + 3 agents + hooks + MCP servers; active 2026-06; needs Domino platform |
| `expo` | Expo skills for building, deploying, upgrading, and debugging React Native apps | unlabeled | `greenfield` | ☑️ desk-checked 2026-07-08 — MCP server + 19 skills; active 2026-07 |
| `firecrawl` | Firecrawl web scraping, crawling, and structured data extraction skills | unlabeled | `research` | ☑️ desk-checked 2026-07-08 — 10 skills + 1 command, no MCP; active 2026-06; needs Firecrawl API key |
| `forge-skills` | Atlassian Forge skills to scaffold, deploy, review, and debug Forge apps | Atlassian | `greenfield` | ☑️ desk-checked 2026-07-08 — MCP server + 6 skills; active 2026-06; needs Atlassian Forge developer account |
| `huggingface-skills` | Build, train, evaluate, and use Hugging Face models, datasets, and spaces | unlabeled | `llm-features` | ☑️ desk-checked 2026-07-08 — MCP server + 19 skills; active 2026-07; needs Hugging Face account for hub operations |
| `idmp-plugin` | TDengine IDMP skills for discovery, schema inspection, and safe operational workflows | TaosData | `devops` | ☑️ desk-checked 2026-07-08 — 23 skills (no MCP); active 2026-05; needs TDengine IDMP instance |
| `liquid-lsp` | LSP for Shopify Liquid templates via Shopify CLI theme language server | Shopify | `code-understanding` | ☑️ desk-checked 2026-07-08 — LSP config (.lsp.json); last commit 2026-03; needs Shopify CLI |
| `liquid-skills` | Liquid fundamentals, coding standards, and WCAG accessibility patterns for Shopify themes | Shopify | `greenfield` | ☑️ desk-checked 2026-07-08 — 3 skills; last commit 2026-03 |
| `lovable` | Build, iterate, deploy, and manage Lovable apps via official MCP | Lovable | `greenfield` | ☑️ desk-checked 2026-07-08 — MCP server (remote, OAuth 2.1) + 3 commands; active 2026-06; needs Lovable account |
| `lumen` | Local semantic code search MCP with Go AST indexing and Ollama embeddings | Ory Corp | `code-understanding` | ☑️ desk-checked 2026-07-08 — MCP server (local Go binary) + 2 skills + hooks; last commit 2026-05; needs Ollama or LM Studio |
| `mcp-apps` | Skills for creating MCP Apps with the MCP Apps SDK | Anthropic / Model Context Protocol | `building-mcp-integrations` | ☑️ desk-checked 2026-07-08 — 4 skills (create-mcp-app, add-app-to-server, convert-web-app, migrate-oai-app); last active 2026-03 |
| `mercadopago` | Mercado Pago payment integration toolkit driven by the official MCP server | Mercado Pago Developer Experience | `greenfield` | ☑️ desk-checked 2026-07-08 — MCP server + 4 skills + 1 agent + commands + hooks; active 2026-06; needs Mercado Pago account and a live MCP connection (no offline mode) |
| `mergify` | Mergify CLI skills for merge queues, stacked PRs, and flaky-test insights | Mergify | `ci-automation` | ☑️ desk-checked 2026-07-08 — 5 skills (merge-queue, stack, ci, merge-protections, config); active 2026-07; needs Mergify CLI and account |
| `microsoft-docs` | Official Microsoft documentation, API references, and samples for Azure, .NET, Windows | unlabeled | `research` | ☑️ desk-checked 2026-07-08 — MCP server + 3 skills (docs, code-reference, skill-creator); active 2026-05 |
| `migration-to-aws` | Plans GCP and AI-workload migrations to AWS from IaC, code, billing data | Amazon Web Services | `migration` | ☑️ desk-checked 2026-07-08 — MCP server + skills (gcp-to-aws, heroku-to-aws) + tools + rules; active 2026-07; needs AWS account as migration target, processing stays local |
| `mintlify` | Build Mintlify documentation sites and convert files to formatted MDX pages | unlabeled | `documentation` | ☑️ desk-checked 2026-07-08 — MCP server + 1 skill; active 2026-06; needs a Mintlify docs project |
| `netlify-skills` | Netlify platform skills covering functions, blobs, forms, caching, and deployment | unlabeled | `devops` | ☑️ desk-checked 2026-07-08 — MCP server + 15 skills (functions, edge-functions, blobs, database, image-cdn, forms, config, deploy, caching, ai-gateway, more); active 2026-07; needs Netlify account |
| `netsuite-suitecloud` | NetSuite skills for SDF objects, UIF SPA components, and AI Service Connector | Oracle NetSuite | `greenfield` | ☑️ desk-checked 2026-07-07 — 10 skills; active 2026-06; needs NetSuite account for runtime connector guidance |
| `nvidia-skills` | NVIDIA accelerated-computing skills: cuOpt optimization, Omniverse, Dynamo, physical-AI infrastructure | NVIDIA | `performance` | ☑️ desk-checked 2026-07-07 — 12 skills; active 2026-07; needs NVIDIA GPU/SDK access |
| `oracle-ai-data-platform-workbench-databricks-migrator` | Automated Databricks to Oracle AIDP migration of notebooks, jobs, catalog DDL | Oracle | `migration` | ☑️ desk-checked 2026-07-07 — 10 skills + 4 commands + 2 agents + engine/references; active 2026-06; needs Oracle AIDP cluster and Databricks access |
| `oracle-ai-data-platform-workbench-engineer-agent` | 37-skill agent operating Oracle AIDP Spark/Delta lakehouse in natural language | Oracle | `devops` | ☑️ desk-checked 2026-07-07 — 37 skills + hooks + MCP template + scripts; active 2026-06; needs Oracle AIDP account and aidp CLI |
| `oracle-ai-data-platform-workbench-spark-connectors` | Spark connector skills for Oracle AIDP: databases, cloud storage, SaaS sources | Oracle | `devops` | ☑️ desk-checked 2026-07-07 — 25 skills + tests/examples/tools; active 2026-06; needs Oracle AIDP workbench cluster plus target data-source credentials |
| `outputai` | Output.ai workflow toolkit: agents, scaffolding, debugging, evaluation, credential skills | Output.ai | `llm-features` | ☑️ desk-checked 2026-07-07 — 49 skills + 5 agents + SessionStart hook; active 2026-07; needs Output.ai SDK/account |
| `pixeltable` | Pixeltable skills for multimodal AI apps: tables, embedding search, UDFs, agents | Pixeltable | `llm-features` | ☑️ desk-checked 2026-07-07 — 1 skill + 2 commands + 2 agents + hooks; active 2026-07; needs Pixeltable Python library |
| `postman` | Postman MCP for collections, client code, API tests, mocks, docs, security audits | unlabeled | `api-design` | ☑️ desk-checked 2026-07-07 — MCP server + 10 commands + 7 skills + 1 agent; active 2026-07; needs Postman account |
| `preset-cli-skills` | Skills for Preset/Superset sup CLI shell, scripting, and CI/CD workflows | Preset | `ci-automation` | ☑️ desk-checked 2026-07-07 — 2 skills; active 2026-06; needs Preset account + superset-sup CLI |
| `pydantic-ai` | Pydantic AI patterns for agents, tools, structured output, streaming, multi-agent apps | unlabeled | `building-agents` | ☑️ desk-checked 2026-07-07 — 1 skill; active 2026-06 |
| `qodo-skills` | Qodo skills for code quality checks, testing, security scanning, compliance validation | unlabeled | `code-review` | ☑️ desk-checked 2026-07-07 — 2 skills (qodo-get-rules, qodo-pr-resolver); active 2026-06; needs Qodo platform account |
| `qt-development-skills` | Qt C++/QML skills for code review, QML coding, and documentation | Qt Group | `greenfield` | ☑️ desk-checked 2026-07-07 — 12 skills + MCP server; active 2026-06; needs Qt toolchain |
| `quarkus-agent` | Quarkus MCP for project scaffolding, dev mode lifecycle, and documentation search | Quarkus | `greenfield` | ☑️ desk-checked 2026-07-08 — MCP server (.mcp.json, Java/Maven source) + plugin dir; active 2026-07; needs Java toolchain |
| `rc` | Configure RevenueCat projects, products, entitlements, and offerings from Claude Code | unlabeled | `devops` | ⚠️ identical source to the revenuecat plugin (same repo, path, and sha) — recommend revenuecat and list once; desk-checked 2026-07-08 — MCP server + 15 skills; active 2026-07; needs RevenueCat account |
| `resend` | Resend skills and MCP for email API integration and deliverability | Resend | `greenfield` | ☑️ desk-checked 2026-07-08 — MCP server + 5 skills; active 2026-06; needs Resend account |
| `revenuecat` | Configure RevenueCat projects, products, entitlements, and offerings from Claude Code | unlabeled | `devops` | ☑️ desk-checked 2026-07-08 — MCP server + 15 skills; active 2026-07; needs RevenueCat account (note: also listed as "rc", same source) |
| `rill` | Skills for developing and querying Rill business intelligence projects | Rill Data | `devops` | ☑️ desk-checked 2026-07-08 — MCP server + 8 skills; active 2026-07; needs Rill project |
| `sagemaker-ai` | AWS SageMaker expertise for building, training, and deploying AI models | unlabeled | `llm-features` | ☑️ desk-checked 2026-07-08 — MCP server + 19 skills; active 2026-06; needs AWS account |
| `sanity` | Sanity CMS MCP, skills, and commands for GROQ queries, schemas, content authoring | Sanity | `greenfield` | ☑️ desk-checked 2026-07-07 — MCP server + 7 skills + 4 commands; active 2026-06; needs Sanity account |
| `sap-cds-mcp` | SAP CAP development assistant searching CDS models and CAP documentation | SAP SE | `greenfield` | ☑️ desk-checked 2026-07-07 — MCP server; active 2026-06 |
| `sap-fiori-mcp-server` | MCP server for building and modifying SAP Fiori applications | SAP SE | `greenfield` | ☑️ desk-checked 2026-07-07 — MCP server + 5 skills; active 2026-07 |
| `sap-mdk-server` | MCP server for SAP Mobile Development Kit app building and scaffolding | SAP SE | `greenfield` | ☑️ desk-checked 2026-07-07 — MCP server; active 2026-06 |
| `servicenow-sdk` | Create, edit, and deploy ServiceNow applications with the Fluent SDK | ServiceNow | `greenfield` | ☑️ desk-checked 2026-07-07 — 1 skill; active 2026-06; needs ServiceNow instance |
| `shopify` | Shopify dev MCP for docs search and GraphQL, Liquid, extension code validation | Shopify | `greenfield` | ⚠️ overlaps shopify-ai-toolkit (same author, superset, newer) — recommend shopify-ai-toolkit first; desk-checked 2026-07-07 — MCP server only; last commit 2026-04; needs Shopify partner/store account |
| `shopify-ai-toolkit` | 18 Shopify skills covering docs, GraphQL, Liquid, Hydrogen, Polaris, CLI workflows | Shopify | `greenfield` | ☑️ desk-checked 2026-07-07 — MCP server + 20 skills + hooks; active 2026-06; needs Shopify partner/store account for CLI workflows |
| `snowflake-cortex-code` | Routes Snowflake prompts to Cortex Code with routing, run, and setup skills | Snowflake | `devops` | ☑️ desk-checked 2026-07-08 — 3 skills + hooks; active 2026-06; needs Snowflake account with Cortex Code |
| `sourcegraph` | Sourcegraph MCP for cross-repository code search, reference tracing, and impact analysis | unlabeled | `code-understanding` | ☑️ desk-checked 2026-07-08 — MCP server + 1 skill; last commit 2026-03; needs Sourcegraph access |
| `stripe` | Stripe development toolkit with MCP, best-practice skills, and upgrade commands | unlabeled | `greenfield` | ☑️ desk-checked 2026-07-08 — MCP server + 5 skills + 2 commands + 1 agent; active 2026-07; needs Stripe account |
| `sumup` | SumUp payment integration skills for POS apps, online checkout, and card readers | unlabeled | `greenfield` | ☑️ desk-checked 2026-07-08 — 6 skills; active 2026-06 |
| `superpowers` | Workflow skills: brainstorming, subagent development, TDD, systematic debugging, skill authoring | unlabeled | `greenfield` | ⚠️ overlaps the built-in /code-review, worktrees, and skill-creator — recommend the built-ins first; desk-checked 2026-07-08 — 14 skills + hooks; active 2026-07 |
| `tavily` | Tavily skills for search, extract, crawl, and research APIs in AI apps | Tavily Team | `llm-features` | ☑️ desk-checked 2026-07-08 — 8 skills; active 2026-06; needs Tavily API key |
| `teamcity-cli` | TeamCity CI/CD skills via teamcity CLI: builds, logs, queues, agents | JetBrains | `ci-automation` | ☑️ desk-checked 2026-07-08 — 2 skills (repo also bundles the Go CLI source); active 2026-07; needs TeamCity server + teamcity CLI installed |
| `togetherai-skills` | Together AI skills: inference, fine-tuning, embeddings, image/video generation, GPU clusters | Together AI | `llm-features` | ☑️ desk-checked 2026-07-08 — 12 skills; active 2026-06; needs Together AI API key |
| `twilio-developer-kit` | Twilio API skills for SMS, Voice, WhatsApp, Verify, SendGrid, 30+ products | Twilio | `greenfield` | ☑️ desk-checked 2026-07-08 — hosted docs MCP server + 2 skill trees (twilio, sendgrid); active 2026-06; needs Twilio account |
| `ui5` | SAPUI5/OpenUI5 project creation, validation, API docs, linter, best practices | SAP SE | `greenfield` | ☑️ desk-checked 2026-07-08 — MCP server + 7 skills; active 2026-07 |
| `ui5-modernization` | Workflow and fix patterns for modernizing SAPUI5/OpenUI5 applications | SAP SE | `migration` | ☑️ desk-checked 2026-07-08 — MCP server + 19 skills; active 2026-06 |
| `ui5-typescript-conversion` | Converts JavaScript-based UI5 projects to TypeScript | SAP SE | `migration` | ☑️ desk-checked 2026-07-08 — MCP server + 1 skill; active 2026-06 |
| `unreal-engine-skills-for-claude-code` | Control Unreal Editor via MCP: actors, blueprints, materials, Sequencer, testing | Epic Games | `greenfield` | ☑️ desk-checked 2026-07-08 — 3 skills + hooks (MCP server hosted inside Unreal Editor, no bundled .mcp.json); active 2026-06; needs Unreal Editor with Unreal MCP enabled |
| `wix` | Build, manage, and deploy Wix sites and apps with CLI skills and MCP | unlabeled | `greenfield` | ☑️ desk-checked 2026-07-08 — MCP server + 7 skills; active 2026-07; needs Wix account |
| `zoom-plugin` | Plan, build, and debug Zoom integrations across APIs, SDKs, webhooks, bots | unlabeled | `api-design` | ☑️ desk-checked 2026-07-08 — MCP server + 32 skills; active 2026-05; needs Zoom developer account |

### Learning

| Plugin | What it does | Author | Relevant goal | Verdict |
|--------|-------------|--------|--------------|---------|
| `learn-with-coursera` | Personalized Coursera course, project, and learning-path recommendations | Coursera | `onboarding` | ☑️ desk-checked 2026-07-08 — 1 skill (3 reference workflows); last commit 2026-05; needs Coursera connector |

### Location

| Plugin | What it does | Author | Relevant goal | Verdict |
|--------|-------------|--------|--------------|---------|
| `amazon-location-service` | Guides adding maps, geocoding, and routing with Amazon Location Service | unlabeled | `greenfield` | ☑️ desk-checked 2026-07-08 — 1 skill; last touched 2026-04; needs AWS account |
| `mapbox` | Mapbox MCP and skills for building location-aware apps and geospatial tools | Mapbox | `greenfield` | ☑️ desk-checked 2026-07-08 — MCP server + 19 skills (web/iOS/Android/Flutter patterns, cartography, migrations); active 2026-06; needs Mapbox account/token |

### Migration

| Plugin | What it does | Author | Relevant goal | Verdict |
|--------|-------------|--------|--------------|---------|
| `aws-transform` | Migrate and modernize codebases to AWS: .NET, COBOL, VMware, databases | Amazon Web Services | `migration` | ☑️ desk-checked 2026-07-08 — MCP server (.mcp.json) + 1 skill; active 2026-07; needs AWS account |

### Monitoring

| Plugin | What it does | Author | Relevant goal | Verdict |
|--------|-------------|--------|--------------|---------|
| `amplitude` | Amplitude analytics for instrumentation, charts, dashboards, experiments, and user insights | unlabeled | `devops` | ☑️ desk-checked 2026-07-08 — MCP server + 27 skills; active 2026-07; needs Amplitude account |
| `dash0` | OpenTelemetry tracing of Claude Code sessions to Dash0 or OTel backends | Dash0 | `devops` | ☑️ desk-checked 2026-07-08 — hooks + Go OTel collector binary + commands/skills; active 2026-07; needs Dash0 or OTel-compatible backend |
| `datadog` | Preconfigured Datadog MCP for querying logs, metrics, traces, and dashboards | Datadog | `incident-response` | ☑️ desk-checked 2026-07-07 — MCP server + 3 skills; active 2026-06; needs Datadog account |
| `fullstory` | Fullstory MCP for behavioral analytics, session replays, and CX insights | Fullstory | `incident-response` | ☑️ desk-checked 2026-07-08 — MCP server + 3 skills + 1 agent; active 2026-06; needs Fullstory account |
| `grafana-assistant` | Skills and rules for developing and using Grafana Assistant app and CLI | Grafana | `devops` | ☑️ desk-checked 2026-07-08 — 1 skill + rules/steering docs (no MCP); active 2026-05; needs Grafana Assistant app/CLI |
| `grafana-cloud-mcp` | Hosted MCP server for Grafana Cloud observability without local installation | Grafana | `devops` | ☑️ desk-checked 2026-07-08 — hosted MCP server + 1 skill; active 2026-05; needs Grafana Cloud account |
| `grafana-mcp` | MCP server for Grafana dashboards, datasources, alerting, and incident management | Grafana | `incident-response` | ☑️ desk-checked 2026-07-08 — MCP server + 1 skill; active 2026-06; needs a Grafana instance |
| `honeycomb` | Honeycomb observability skills: query patterns, production investigations, SLOs, OpenTelemetry instrumentation | Honeycomb | `incident-response` | ☑️ desk-checked 2026-07-08 — MCP server + 11 skills + 2 agents + 1 command + hooks; active 2026-06; needs Honeycomb account |
| `langfuse` | Skills for Langfuse LLM tracing, prompt management, and evaluation workflows | Langfuse | `llm-features` | ☑️ desk-checked 2026-07-08 — 1 skill (with references); active 2026-07; needs Langfuse account |
| `langfuse-observability` | Hooks that trace Claude Code sessions into Langfuse observability | Langfuse | `devops` | ☑️ desk-checked 2026-07-08 — hooks (hooks.json + Python hook); active 2026-06; needs Langfuse account |
| `logfire` | Adds Logfire observability and auto-instrumentation to Python applications | Pydantic | `devops` | ☑️ desk-checked 2026-07-08 — MCP server + 4 commands + 3 skills; active 2026-06; needs Logfire account |
| `logrocket` | Query LogRocket session replays, metrics, issues, and user behavior | LogRocket | `debugging` | ☑️ desk-checked 2026-07-08 — MCP server + 1 skill; active 2026-06; needs LogRocket account |
| `pagerduty` | PagerDuty risk scoring of pre-commit diffs against historical incident data | unlabeled | `incident-response` | ☑️ desk-checked 2026-07-07 — MCP server + 2 commands (pre-commit-risk-scoring, create-pagerduty-skill); active 2026-05; needs PagerDuty account |
| `posthog` | PostHog MCP for analytics, feature flags, experiments, error tracking, and insights | unlabeled | `devops` | ☑️ desk-checked 2026-07-07 — MCP server + 106 skills + 6 commands + 1 agent + hooks; active 2026-07; needs PostHog account |
| `rootly` | Incident management: deploy safety, on-call, incident response, and retrospectives | Rootly | `incident-response` | ☑️ desk-checked 2026-07-08 — MCP server + 18 skills + 3 agents + hook; active 2026-05; needs Rootly account |
| `sentry` | Sentry error monitoring integration for stack traces, issue search, production debugging | unlabeled | `incident-response` | ☑️ desk-checked 2026-07-07 — 35 skills + 1 command; active 2026-07; needs Sentry account |
| `sentry-cli` | Skills for driving Sentry from the command line via sentry-cli | Sentry | `devops` | ☑️ desk-checked 2026-07-07 — 1 skill; active 2026-07; needs sentry-cli and Sentry account |

### Productivity

| Plugin | What it does | Author | Relevant goal | Verdict |
|--------|-------------|--------|--------------|---------|
| `airtable` | Airtable MCP for creating bases, schema, records, and shared collaboration views | Airtable | `devops` | ☑️ desk-checked 2026-07-08 — MCP server + 8 skills; active 2026-06; needs Airtable account |
| `airwallex-agentos` | Airwallex finance skills and MCP for invoices, suppliers, and cash positions | Airwallex | `devops` | ☑️ desk-checked 2026-07-08 — MCP server + 5 skills; active 2026-06; needs Airwallex account and public CLI |
| `apollo` | Apollo.io MCP for prospecting, lead enrichment, outreach sequences, and sales analytics | Apollo.io | `research` | ☑️ desk-checked 2026-07-08 — MCP server + 4 skills; last touched 2026-04; needs Apollo.io account |
| `atlassian` | Jira and Confluence integration for issues, docs, sprints, and status reports | unlabeled | `documentation` | ☑️ desk-checked 2026-07-08 — MCP server + 6 skills; active 2026-06; needs Atlassian account |
| `box` | Search, organize, and query Box content and documents via Box AI | unlabeled | `documentation` | ☑️ desk-checked 2026-07-08 — 5 skills + rules; active 2026-07; needs Box account |
| `carta-cap-table` | Query Carta cap tables, grants, SAFEs, 409A valuations, and waterfall scenarios | Carta Engineering | `research` | ☑️ desk-checked 2026-07-08 — hosted MCP server + 15 skills + hooks; active 2026-07; needs Carta account |
| `carta-crm` | Manage Carta CRM investors, companies, contacts, deals, notes, and fundraisings conversationally | Carta Engineering | `research` | ☑️ desk-checked 2026-07-08 — hosted MCP server + 21 skills + hooks; active 2026-06; needs Carta CRM account |
| `carta-investors` | Query Carta investor data, benchmarks, regulatory reporting, and AGM deck generation | Carta Engineering | `research` | ☑️ desk-checked 2026-07-08 — hosted MCP server + 15 skills + hooks; active 2026-07; needs Carta account |
| `circleback` | Circleback MCP for searching meetings, emails, and calendar events | unlabeled | `research` | ☑️ desk-checked 2026-07-08 — MCP server only; last commit 2026-01; needs Circleback account |
| `coderabbit` | CodeRabbit AI code review with 40+ static analyzers and suggested fixes | unlabeled | `code-review` | ⚠️ overlaps the built-in /code-review — recommend the built-in first; desk-checked 2026-07-08 — 2 skills + 1 agent + 1 command; active 2026-06; needs CodeRabbit (free per manifest) |
| `desktop-commander` | MCP for terminal commands, process management, and multi-format file operations | Desktop Commander | `devops` | ⚠️ overlaps Claude Code's built-in Bash and file tools (terminal, process, file ops) — recommend the built-ins first; desk-checked 2026-07-07 — MCP server + 6 skills; active 2026-07 |
| `dropbox` | Dropbox MCP to search, organize, save, and share files from Claude | Dropbox | `research` | ☑️ desk-checked 2026-07-08 — MCP server + 6 skills; active 2026-06; needs Dropbox account |
| `exa` | Exa MCP for web search, deep research, and content extraction | Exa | `research` | ⚠️ overlaps built-in web search and the deep-research skill — recommend the built-ins first; desk-checked 2026-07-08 — hosted MCP server + 2 skills; active 2026-06; needs Exa account |
| `hunter` | Find and verify professional emails, search domain contacts, enrich company data | Hunter.io | `research` | ☑️ desk-checked 2026-07-08 — MCP server + 9 skills; active 2026-07; needs Hunter.io account |
| `intercom` | Search Intercom conversations, analyze support patterns, look up contacts and companies | unlabeled | `research` | ☑️ desk-checked 2026-07-08 — MCP server + 4 skills; active 2026-04; needs Intercom workspace |
| `legalzoom` | AI legal document review with risk flagging and attorney routing | unlabeled | `research` | ☑️ desk-checked 2026-07-08 — MCP server + 1 command + 1 skill; last commit 2026-02; needs LegalZoom connector |
| `lusha` | Prospect and enrich B2B leads with Lusha verified contact data | Lusha | `research` | ☑️ desk-checked 2026-07-08 — remote MCP server + 4 skills; active 2026-06; needs Lusha account |
| `monday-crm` | Run monday CRM in plain language: pipelines, briefings, forecasts, bulk cleanup | monday.com | `greenfield` | ☑️ desk-checked 2026-07-08 — MCP server + 5 skills (workspace-builder, daily-briefing, forecast, data-cleanup, meeting-to-deal); active 2026-06; needs monday.com account |
| `notion` | Notion workspace MCP: search pages, manage databases, documentation workflows | unlabeled | `documentation` | ☑️ desk-checked 2026-07-07 — MCP server + 7 commands + 1 skill; last commit 2026-01; needs Notion account |
| `pigment` | Analyze business data and build Pigment models, metrics, and boards | Pigment | `greenfield` | ☑️ desk-checked 2026-07-07 — MCP server + 11 skills; active 2026-06; needs Pigment account |
| `save-to-spotify` | Creates TTS audio episodes with cover images and saves them to Spotify | Spotify | `documentation` | ☑️ desk-checked 2026-07-07 — 1 skill; active 2026-05; needs save-to-spotify CLI and Spotify account |
| `slack` | Slack MCP for searching messages, channels, and threads for team context | unlabeled | `research` | ☑️ desk-checked 2026-07-08 — MCP server + 6 skills + 5 commands; active 2026-07; needs Slack workspace app/token |
| `spotify-ads-api` | Manage Spotify ad campaigns, reports, and OAuth through conversation | unlabeled | `devops` | ☑️ desk-checked 2026-07-08 — 14 skills + 1 command + 1 agent + hooks; active 2026-07; needs Spotify Ads account/OAuth |
| `vibe-prospecting` | Live B2B company and contact data for prospecting, enrichment, CRM workflows | vibeprospecting.ai | `research` | ☑️ desk-checked 2026-07-08 — 1 skill + helper scripts; active 2026-06; needs vibeprospecting.ai account |
| `windsor-ai` | Query 325+ marketing, sales, CRM, and analytics data sources via Windsor.ai | Windsor.ai | `research` | ☑️ desk-checked 2026-07-08 — MCP server + 1 skill + 3 commands + 1 agent; active 2026-04; needs Windsor.ai account |
| `youdotcom-agent-skills` | You.com search and research skills with agent SDK integration guides | You.com | `building-agents` | ☑️ desk-checked 2026-07-08 — 8 skills (API, CLI, and 6 agent-framework integrations); active 2026-05; needs You.com API key |
| `zapier` | Discover, enable, and execute Zapier actions across 8,000+ connected apps | unlabeled | `devops` | ☑️ desk-checked 2026-07-08 — MCP server + 4 skills + 1 agent; active 2026-06; needs Zapier account |
| `zoominfo` | ZoomInfo B2B company and contact search, enrichment, and sales workflows | ZoomInfo | `research` | ☑️ desk-checked 2026-07-08 — MCP server + 14 skills; active 2026-07; needs ZoomInfo account |

### Security

| Plugin | What it does | Author | Relevant goal | Verdict |
|--------|-------------|--------|--------------|---------|
| `42crunch-api-security-testing` | Audit OpenAPI specs, detect OWASP API vulnerabilities, and apply fixes with 42Crunch | 42Crunch | `security` | ☑️ desk-checked 2026-07-08 — 6 skills + references (no MCP); active 2026-07; needs 42Crunch account |
| `auth0` | Framework-aware skills for adding Auth0 login, SSO, MFA, and access control | Auth0 | `security` | ☑️ desk-checked 2026-07-08 — 45 framework-specific skills; active 2026-07; needs Auth0 tenant |
| `crowdstrike-falcon-foundry` | Build cybersecurity apps on CrowdStrike Falcon Foundry: UI, functions, workflows | CrowdStrike | `security` | ☑️ desk-checked 2026-07-08 — 10 skills + 5 hooks; active 2026-06; needs CrowdStrike Falcon Foundry access |
| `duende-skills` | OAuth/OIDC, IdentityServer, and ASP.NET Core identity security skills | Duende Software | `security` | ☑️ desk-checked 2026-07-08 — 24 skills + 2 agents; active 2026-06 |
| `jfrog` | JFrog Platform: Artifactory artifacts, security findings, package safety, platform administration | JFrog Ltd. | `security` | ☑️ desk-checked 2026-07-08 — MCP server + 3 skills + hooks; active 2026-07; needs JFrog Platform account |
| `semgrep` | Semgrep security scanning that flags vulnerabilities as Claude writes code | unlabeled | `security` | ☑️ desk-checked 2026-07-07 — MCP server + hooks; active 2026-06; needs Semgrep |
| `sonarqube` | SonarQube quality and security analysis enforced via hooks, MCP, and skills | SonarSource | `code-review` | ☑️ desk-checked 2026-07-08 — MCP server + 9 skills + hooks; active 2026-07; needs SonarQube server/token |
| `sonatype-guide` | Sonatype Guide MCP for dependency vulnerability analysis and secure version recommendations | unlabeled | `dependency-management` | ☑️ desk-checked 2026-07-08 — MCP server + 1 skill; last commit 2026-04 |
| `vanta` | Vanta MCP for fixing security compliance test failures with repo context | Vanta | `security` | ☑️ desk-checked 2026-07-08 — MCP server + 3 skills; active 2026-05; needs Vanta account |
| `vanta-mcp-plugin` | Vanta MCP for fixing security compliance test failures with repo context | Vanta | `security` | ⚠️ duplicate listing of the vanta plugin (identical repo and pinned SHA 345d86b) — recommend the vanta entry; desk-checked 2026-07-08 — MCP server + 3 skills; active 2026-05; needs Vanta account |
| `workos` | WorkOS skills for AuthKit, SSO, Directory Sync, RBAC, Vault, audit logs | WorkOS | `security` | ☑️ desk-checked 2026-07-08 — 2 skills; active 2026-06; needs WorkOS account |
| `zscaler` | Manage Zscaler security platform: policies, audits, connectivity, incident investigation | Zscaler | `security` | ☑️ desk-checked 2026-07-08 — MCP server + 8 skills + 20 commands; active 2026-06; needs Zscaler account |

### Uncategorized

| Plugin | What it does | Author | Relevant goal | Verdict |
|--------|-------------|--------|--------------|---------|
| `ai-plugins` | Endor Labs scanning to find and fix software supply chain security risks | unlabeled | `dependency-management` | ☑️ desk-checked 2026-07-08 — MCP server + 1 skill (endor-setup); active 2026-07; needs Endor Labs account and endorctl |
| `aikido` | Aikido SAST, secrets, and IaC vulnerability scanning via MCP server | unlabeled | `security` | ☑️ desk-checked 2026-07-08 — MCP server + 3 skills (setup, scan, issues); active 2026-06; needs Aikido account |
| `atlan` | Atlan data catalog MCP for asset search, lineage, glossary, and governance | unlabeled | `devops` | ☑️ desk-checked 2026-07-08 — MCP server only (no skills or commands); active 2026-06; needs Atlan account |
| `brightdata-plugin` | Web scraping, Google search, and structured data extraction via Bright Data | unlabeled | `research` | ☑️ desk-checked 2026-07-08 — MCP server (.mcp.json) + 21 skills; active 2026-06; needs Bright Data account |
| `cloudinary` | Manage Cloudinary assets, transformations, and media optimization from Claude | unlabeled | `devops` | ☑️ desk-checked 2026-07-08 — MCP server + 2 skills; last commit 2026-04; needs Cloudinary account |
| `data-engineering` | Astronomer plugin for warehouse exploration, pipeline authoring, and Airflow integration | unlabeled | `devops` | ☑️ desk-checked 2026-07-07 — 26 skills + bundled Airflow MCP; active 2026-07; needs Airflow/Astro environment |
| `fastly-agent-toolkit` | Fastly development tools and platform skills | unlabeled | `devops` | ☑️ desk-checked 2026-07-08 — 8 skills (CLI, compute, NGWAF, VCL), no MCP; active 2026-07; needs Fastly account |
| `fiftyone` | FiftyOne skills for dataset curation, model evaluation, and computer vision workflows | unlabeled | `llm-features` | ☑️ desk-checked 2026-07-08 — MCP server + 17 skills + 2 commands; active 2026-07 |
| `nightvision` | DAST and API discovery skills for finding exploitable web/API vulnerabilities | unlabeled | `security` | ☑️ desk-checked 2026-07-07 — 4 skills (api-discovery, scan-configuration, scan-triage, ci-cd-integration); active 2026-07; needs NightVision account/CLI |
| `nimble` | Nimble MCP and skills to search, extract, map, and crawl web data | unlabeled | `research` | ☑️ desk-checked 2026-07-07 — MCP server + 8 skills + 2 agents + 1 command; active 2026-06; needs Nimble account/API key |
| `postiz` | Postiz CLI for scheduling social posts, media, and analytics across 28+ platforms | unlabeled | `devops` | ☑️ desk-checked 2026-07-07 — 1 skill + bundled TypeScript CLI; active 2026-06; needs Postiz account/API key |
| `prisma` | Prisma MCP for Postgres provisioning, schema migrations, SQL queries, connection strings | unlabeled | `devops` | ☑️ desk-checked 2026-07-07 — MCP server only; last activity 2026-03; needs Prisma Postgres account |
| `remember` | Continuous memory compressing Claude Code conversations into tiered daily logs | unlabeled | `documentation` | ⚠️ overlaps the built-in auto-memory (persistent MEMORY.md across sessions) — recommend the built-in first; desk-checked 2026-07-08 — hooks + 1 skill + Python pipeline; active 2026-07; needs Python |
| `build-with-wordpress` | Craft production-grade WordPress sites and applications — themes, plugins, commerce, and deployment | unlabeled | `greenfield` | ☑️ desk-checked 2026-07-09 — renamed from `wordpress.com`; Automattic source (claude-code-wordpress.com); manifest and provenance re-verified, component inventory not re-counted since the rename |

