/*
Package pyrotest provides a Prometheus datamodel-specific DSL for writing
concise and easily readable assertions on metrics.

# Metric Families and Metrics

As a timeseries database user you don't care about the gory internal details of
the Prometheus data model, you just work with “metrics” that have “labels”.

Consequently, package pyrotest conceals the confusing and fussy hierarchical
differentiation of the Prometheus data model into metric families that only then
contain individual metrics (that is, the individual “timeseries”).

# Usage

A typical usage might look like this:

	// import . "github.com/thediveo/pyrotest"

	metfams := CollectAndLint(mycollector)
	Expect(metfams).To(ContainMetrics(
		Counter(HaveName("foo_sprockets_total"),
			HaveUnit("sprockets"),
			HaveLabel("type=barz"), HaveLabel("anyway")),
		Gauge(HaveName(ContainsSubstring("rumpelpumpel")),
			HaveLabel("region=elsewhere")),
	))

# Motivation

This package isn't strictly necessary, as [Gomega's matcher toolchest] already
contains everything necessary. However, this requires extensive use of
[gomega.HaveField] matchers in combination with protobuf-originating accessor
functions like [gomega.GetName] in order to correctly deal with the level of
pointer indirection used in protobuf optional fields – which the Prometheus data
model likes to use basically everywhere. pyrotest brings back concise matcher
design, with build-time type-safety on top.

[Gomega's matcher toolchest]: https://onsi.github.io/gomega/
*/
package pyrotest
