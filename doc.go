/*
Package pyrotest provides a Prometheus datamodel-specific DSL for writing
concise and easily readable assertions on metrics.

# Metric Families and Metrics

Package pyrotest conceals the slightly fussy hierarchical differentiation of the
Prometheus data model into metric families that only then contain individual
metrics (that is, the individual “timeseries”) as a nasty implementation detail.

# Motivation

This package isn't strictly necessary, as Gomega's matcher toolchest already
contains everything necessary. However, this requires extensive use of
`HaveField()` matchers in combination with protobuf-originating accessor
functions like `GetName()` in order to correctly deal with the level of pointer
indirection used in protobuf optional fields – which the Prometheus data model
likes to use basically everywhere. pyrotest brings back concise matcher design,
with build-time type-safety on top.
*/
package pyrotest
