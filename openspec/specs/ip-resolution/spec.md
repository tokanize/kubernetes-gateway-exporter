# ip-resolution Specification

## Purpose
Define the portable mechanism for reading an IP address assigned to a Kubernetes Gateway.

> **Implementation Details:** See [ADR-002: Gateway Address Resolution](../../../docs/adrs/002-ip-resolution.md).

## Requirements

### Requirement: Retrieve IP from Gateway Status
The system SHALL read Gateway addresses exclusively from the Kubernetes Gateway API status.

#### Scenario: Gateway has an allocated IP
- GIVEN a `Gateway` contains an address with type `IPAddress`
- WHEN exposure metrics are collected
- THEN the first `IPAddress` value is exported in the `ip_address` label.

#### Scenario: Gateway address is not allocated yet
- GIVEN a valid exposure relationship exists
- AND the `Gateway` has no address with type `IPAddress`
- WHEN exposure metrics are collected
- THEN the relationship is still exported
- AND the `ip_address` label is an empty string.

#### Scenario: Gateway exposes only a hostname address
- GIVEN the `Gateway` contains only addresses with type `Hostname`
- WHEN exposure metrics are collected
- THEN the `ip_address` label is an empty string.

### Requirement: Provider Independence
The system SHALL NOT query cloud-provider APIs to resolve Gateway addresses.

#### Scenario: Address is missing from Gateway status
- GIVEN the `Gateway` has no `IPAddress` in `status.addresses`
- WHEN exposure metrics are collected
- THEN no external provider API is queried.
