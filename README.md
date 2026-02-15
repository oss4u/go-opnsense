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

Run only Unbound tests:

```bash
go test ./opnsense/core/unbound/...
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
make pre-commit-run
```

Run the local Dagger CI pipeline (tidy + build + tests):

```bash
make ci
```

### 3) Commit message policy

This repository enforces Conventional Commits via `commit-msg` hook.

Examples:

- `feat(unbound): add settings toggle wrapper`
- `fix(overrides): handle empty get response`
- `chore(ci): pin golangci-lint version`

### 4) Branch/CI behavior

- Push to `develop`:
    - `develop-lint` (`golangci-lint`)
    - `develop-build` (Dagger build)
    - `develop-test` (Dagger tests, after build)
- Pull request to `main`:
    - `build-and-test` (Dagger CI)
- Push to `main`:
    - `build-and-test` + `release`

## Usage

### Create client

```go
api := opnsense.GetOpnSenseClient("", "", "")
unboundApi := unbound.New(api)
```

Passing empty values uses `OPNSENSE_ADDRESS`, `OPNSENSE_KEY`, and `OPNSENSE_SECRET` from the environment.

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
