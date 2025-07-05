# `pyrotest`

[![PkgGoDev](https://img.shields.io/badge/-reference-blue?logo=go&logoColor=white&labelColor=505050)](https://pkg.go.dev/github.com/thediveo/pyrotest)
[![License](https://img.shields.io/github/license/thediveo/pyrotest)](https://img.shields.io/github/license/thediveo/pyrotest)
![build and test](https://github.com/thediveo/pyrotest/actions/workflows/buildandtest.yaml/badge.svg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/thediveo/pyrotest)](https://goreportcard.com/report/github.com/thediveo/pyrotestb)
![Coverage](https://img.shields.io/badge/Coverage-95.6%25-brightgreen)

`pyrotest` provides [Gomega matchers](https://onsi.github.io/gomega/) as well as
specially typed matchers for reasoning about Prometheus metrics.

In particular, it conceals the slightly fussy hierarchical differentiation of
the Prometheus data model into metric families that only then contain individual
metrics (that is, the individual “timeseries”) as a nasty implementation detail.
Not least, as a prometheus user you deal with the (ultimate) metrics, not
families.

## Example

```go
package some_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/thediveo/pyrotest"
)

var _ = Describe("example", func() {

    It("collects metrics, lints them, and reasons about them", func() {
        var coll prometheus.Collector = ...

        families := CollectAndLint(coll)
        Expect(families).To(ContainMetrics(
            Gauge(HaveName("bar_baz"),
                HaveHelp(Not(BeEmpty())),
                HaveLabel("foobar=baz")),
            Counter(HaveName(ContainSubstring("_total")),
                HaveHelp(ContainSubstring("no help")),
                HaveLabelWithValue("label", "scam")),
        ))
    })

})

```

## Contributing

Please see [CONTRIBUTING.md](CONTRIBUTING.md).

## DevContainer

> [!CAUTION]
>
> Do **not** use VSCode's "~~Dev Containers: Clone Repository in Container
> Volume~~" command, as it is utterly broken by design, ignoring
> `.devcontainer/devcontainer.json`.

1. `git clone https://github.com/thediveo/pyrotest`
2. in VSCode: Ctrl+Shift+P, "Dev Containers: Open Workspace in Container..."
3. select `pyrotest.code-workspace` and off you go...

## Go Version Support

`pyrotest` supports versions of Go that are noted by the Go release policy, that
is, major versions _N_ and _N_-1 (where _N_ is the current major version).

## Copyright and License

`pyrotest` is Copyright 2025 Harald Albrecht, and licensed under the Apache
License, Version 2.0.
