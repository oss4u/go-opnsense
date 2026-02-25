# go-opnsense

Go client library for OPNsense APIs.

This documentation currently covers only the implemented **Unbound** resources.

## Scope

Implemented OPNsense module:

- `unbound/diagnostics`
- `unbound/overview`
- `unbound/service`
- `unbound/settings`
- `unbound/settings` host overrides and host aliases (CRUD wrappers in `overrides` package)
- generic `core` API resources (`/api/core/<controller>/{add,get,search,set,del,toggle}`)
- generic plugin API resources (`/api/<plugin>/<controller>/{add,get,search,set,del,toggle}`)

## Requirements

- Go `1.20+` (project currently uses `v1.26` via `mise`)
- Access to an OPNsense instance with API key and secret

## Setup

### 1) Install toolchain

If you use `mise`:

```bash
mise install
```

Or install Go manually and ensure `go version` is at least `1.20`.

### 2) Install dependencies

```bash
go mod tidy
```

### 3) Configure credentials

Set the following environment variables:

```bash
export OPNSENSE_ADDRESS="https://<your-opnsense-host>"
export OPNSENSE_KEY="<api-key>"
export OPNSENSE_SECRET="<api-secret>"
```

## Build and test

Run all tests:

```bash
go test ./...
```

Or use Taskfile commands:

```bash
task build
task test
```

Run only Unbound tests:

```bash
go test ./opnsense/core/unbound/...
```

Run manual integration tests against a real OPNsense instance:

```bash
go test -tags manual ./...
```

or via Taskfile:

```bash
task test:manual:local
```

Optional Dagger wrapper:

```bash
task test:manual
```

## Development

### 1) Install development hooks

Install `pre-commit` (for example via `pipx`):

```bash
pipx install pre-commit
```

Install repository hooks:

```bash
make pre-commit-install
```

### 2) Run checks locally before pushing

Run the same lint checks as configured in CI:

```bash
task pre-commit:run
```

Run the local CI task (tidy + lint + build + tests):

```bash
task ci
```

`task ci` runs in an isolated Dagger environment and is the same entry-point used by CI.

Run only the Pact pipeline (also via Dagger):

```bash
task pact:ci
```

### 3) Commit message policy

This repository enforces Conventional Commits via `commit-msg` hook.

Examples:

- `feat(unbound): add settings toggle wrapper`
- `fix(overrides): handle empty get response`
- `chore(ci): pin golangci-lint version`

### 4) Branch/CI behavior

- Push to `develop`:
    - `ci` (runs `task ci`)
- Pull request to `main`:
    - `ci`
    - `pact-tests` (runs `task pact:ci`)
- Push to `main`:
    - `ci`
    - `pact-tests`
    - `release`

## Taskfile commands

This repository now uses `Taskfile.yml` as the primary command runner.

Core tasks:

- `task init` — `go mod tidy`
- `task lint` — run `golangci-lint` with `.golangci.yml`
- `task build` — `go build -v ./...`
- `task test` — `go test ./...`
- `task test:ci` — CI-mode tests
- `task ci` — full pipeline (`init`, `lint`, `build`, `test:ci`)
- `task pact:ci` — Pact pipeline in Dagger
- `task pact:install` — install Pact FFI runtime to `.pact/lib`
- `task test:pact` — run `go test -tags pact ./...`
- `task test:manual:local` — run `go test -tags manual ./...` against a real OPNsense instance
- `task test:manual` — run manual tests in Dagger
- `task release` — build release artifacts via Dagger
- `task vagrant:start|stop|destroy` — local Vagrant helper tasks

### Manual integration test environment

Required variables:

```bash
export OPNSENSE_ADDRESS="https://<your-opnsense-host>"
export OPNSENSE_KEY="<api-key>"
export OPNSENSE_SECRET="<api-secret>"
```

Optional variables for plugin and override resource tests:

```bash
export OPNSENSE_MANUAL_PLUGIN="caddy"
export OPNSENSE_MANUAL_PLUGIN_CONTROLLER="service"
export OPNSENSE_MANUAL_HOST_OVERRIDE_UUID="<existing-host-override-uuid>"
export OPNSENSE_MANUAL_ALIAS_OVERRIDE_UUID="<existing-alias-override-uuid>"
```

## Pact runtime notes

Pact consumer tests require the native `libpact_ffi` library.

The Taskfile handles platform-specific runtime names automatically:

- macOS: `libpact_ffi.dylib`
- Linux: `libpact_ffi.so` (or versioned `libpact_ffi.so.*`)

`task pact:install` forces installation into `.pact/lib`, and `task test:pact`
resolves the correct library file before running `go test -tags pact ./...`.

Local setup:

```bash
task pact:install
task test:pact
```

In CI, the `pact-tests` job runs the same `Taskfile` task through Dagger.

## Usage

### Create client

```go
api := opnsense.GetOpnSenseClient("", "", "")
unboundApi := unbound.New(api)

coreResourceApi := coreresources.New(api, "service")
pluginResourceApi := pluginresources.New(api, "caddy", "service")

coreApi := core.New(api)
pluginApi := plugins.New(api, "caddy")

viaCoreController := coreApi.Controller("service")
viaPluginController := pluginApi.Controller("service")

typedCoreService := coreApi.Service()
typedCaddyService := pluginApi.Caddy().Service()

registry := plugins.NewRegistry(api)
samePlugin := registry.Plugin("caddy")
_ = samePlugin.Controller("service")
```

Passing empty values uses `OPNSENSE_ADDRESS`, `OPNSENSE_KEY`, and `OPNSENSE_SECRET` from the environment.

Imports for generic resource APIs:

```go
import coreresources "github.com/oss4u/go-opnsense/opnsense/core/resources"
import pluginresources "github.com/oss4u/go-opnsense/opnsense/plugins/resources"
import "github.com/oss4u/go-opnsense/opnsense/core"
import "github.com/oss4u/go-opnsense/opnsense/plugins"
import coreservice "github.com/oss4u/go-opnsense/opnsense/core/service"
import caddyservice "github.com/oss4u/go-opnsense/opnsense/plugins/caddy/service"
```

All requests are serialized by the client lock. Even when called concurrently,
only one HTTP request is in-flight at a time.

## Unbound API usage

### Diagnostics

Available methods:

- `Dumpcache()`
- `Dumpinfra()`
- `Listinsecure()`
- `Listlocaldata()`
- `Listlocalzones()`
- `Stats()`
- `TestBlocklist(payload any)`

Example:

```go
statsRaw, statusCode, err := unboundApi.Diagnostics.Stats()
_ = statsRaw
_ = statusCode
_ = err
```

### Overview

Available methods:

- `Rolling(timeperiod string, clients int)`
- `GetPolicies(uuid string)`
- `IsBlockListEnabled()`
- `IsEnabled()`
- `SearchQueries()`
- `Totals(maximum string)`

### Service

Available methods:

- `Dnsbl()`
- `Reconfigure()`
- `ReconfigureGeneral()`
- `Restart()`
- `Start()`
- `Status()`
- `Stop()`

### Settings (generic resources)

Supported resources:

- `unbound.SettingsResourceACL`
- `unbound.SettingsResourceDNSBL`
- `unbound.SettingsResourceForward`
- `unbound.SettingsResourceHostAlias`
- `unbound.SettingsResourceHostOverride`

Available generic methods:

- `Add(resource, payload)`
- `GetResource(resource, uuid)`
- `Search(resource, payload)`
- `SetResource(resource, uuid, payload)`
- `DeleteResource(resource, uuid)`
- `ToggleResource(resource, uuid, enabled)`

Additional settings methods:

- `Get()`
- `Set(payload)`
- `GetNameservers()`
- `UpdateBlocklist()`

Example:

```go
_, err := unboundApi.Settings.Add(
    unbound.SettingsResourceHostOverride,
    map[string]any{
        "host": map[string]any{
            "enabled": "1",
            "hostname": "srv01",
            "domain": "example.local",
            "rr": "A",
            "server": "10.0.0.10",
            "description": "example override",
        },
    },
)
_ = err
```

## Host overrides / aliases convenience API

Package: `opnsense/core/unbound/overrides`

### Hosts

- `Create(host *OverridesHost)`
- `Read(uuid string)`
- `Update(host *OverridesHost)`
- `Delete(host *OverridesHost)`
- `DeleteByID(uuid string)`

### Aliases

- `Create(alias *OverridesAlias)`
- `Read(uuid string)`
- `Update(alias *OverridesAlias)`
- `Delete(alias *OverridesAlias)`
- `DeleteByID(uuid string)`

Example:

```go
hostsApi := overrides.GetHostsOverrideApi(api)
created, err := hostsApi.Create(&overrides.OverridesHost{
    Host: overrides.OverridesHostDetails{
        Enabled: true,
        Hostname: "srv01",
        Domain: "example.local",
        Rr: "A",
        Server: "10.0.0.10",
        Description: "sample",
    },
})
_ = created
_ = err
```

## API call behavior

All OPNsense API calls are serialized in the client to avoid parallel request issues.

## Testing strategy

- Unit tests for types and wrappers
- Contract-style tests for Unbound endpoints in `opnsense/core/unbound/unbound_contract_test.go`
- Mocked endpoint CRUD tests for host overrides and aliases in `opnsense/core/unbound/overrides/api_mock_test.go`

## Linting policy

This project uses `golangci-lint` v2 with config version `2` from `.golangci.yml`.

Recommended local commands:

```bash
golangci-lint fmt --config=.golangci.yml
golangci-lint run --config=.golangci.yml
```

Notes about `revive` exclusions:

- Some `revive` style findings are intentionally excluded to preserve public API compatibility and avoid breaking downstream users.
- Specifically, we ignore:
    - package comment requirements on all packages
    - exported symbol comment requirements
    - naming/stuttering suggestions for existing exported API names
    - var-naming recommendations that would force broad API renames
- Functional and correctness-oriented checks (`govet`, `errcheck`, `staticcheck`, `ineffassign`, `unused`, `unconvert`, `misspell`) remain enabled.

## Notes

- This project currently documents and supports only the implemented **Unbound** API resources.
- Additional OPNsense modules can be added in the same pattern later.
