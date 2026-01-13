# tcost - Temporal Cloud Cost Report

A command-line tool that fetches usage data from Temporal Cloud and generates cost reports per namespace. Designed for platform teams that want to charge back Temporal Cloud costs to individual teams.

## Features

- Fetches usage data from the Temporal Cloud API
- Aggregates costs by namespace
- Configurable pricing for actions, active storage, and retained storage
- Supports table and JSON output formats
- Flexible date range selection

## Installation

```bash
go install github.com/brendan-myers/temporal-cost-report@latest
```

Or build from source:

```bash
git clone https://github.com/brendan-myers/temporal-cost-report.git
cd temporal-cost-report
go build -o tcost .
```

## Prerequisites

You need a Temporal Cloud API key with permissions to read usage data. You can either:

1. Set it as an environment variable:
```bash
export TEMPORAL_API_KEY=your-api-key-here
```

2. Pass it as a command-line argument:
```bash
tcost --api-key your-api-key-here
```

To create an API key, see the [Temporal Cloud documentation](https://docs.temporal.io/cloud/api-keys).

## Usage

```bash
# Generate report for current month (default)
tcost

# Specify a custom date range
tcost --start-date 2025-12-01 --end-date 2025-12-31

# Output as JSON
tcost --format json

# Use custom pricing
tcost --action-price 30 --active-storage-price 0.05 --retained-storage-price 0.001

# Pass API key directly
tcost --api-key your-api-key-here
```

## Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--api-key` | string | | Temporal Cloud API key (defaults to `TEMPORAL_API_KEY` env var) |
| `--start-date` | string | First day of current month | Start date (YYYY-MM-DD format) |
| `--end-date` | string | Today | End date (YYYY-MM-DD format) |
| `--action-price` | float | 50.0 | Price per million actions (USD) |
| `--active-storage-price` | float | 0.042 | Price per GBh of active storage (USD) |
| `--retained-storage-price` | float | 0.00105 | Price per GBh of retained storage (USD) |
| `--format` | string | table | Output format: `table` or `json` |

## Output Examples

### Table Format

```
Temporal Cloud Usage Report
Period: 2026-01-01 to 2026-01-14
Pricing: $50.00/M actions, $0.0420/GBh active, $0.00105/GBh retained

┌─────────────────────────┬─────────┬──────────────┬────────────────┬─────────────┬──────────────┬────────────┐
│ NAMESPACE               │ ACTIONS │ ACTIVE (GBh) │ RETAINED (GBh) │ ACTION COST │ STORAGE COST │ TOTAL COST │
├─────────────────────────┼─────────┼──────────────┼────────────────┼─────────────┼──────────────┼────────────┤
│ prod-workflows.acct     │ 12.35M  │ 1234.56      │ 5678.90        │ $308.64     │ $57.81       │ $366.45    │
│ staging.acct            │ 1.23M   │ 123.45       │ 567.89         │ $30.86      │ $5.78        │ $36.64     │
│ dev-team-alpha.acct     │ 234.57K │ 23.45        │ 56.78          │ $5.86       │ $1.04        │ $6.91      │
├─────────────────────────┼─────────┼──────────────┼────────────────┼─────────────┼──────────────┼────────────┤
│ TOTAL                   │ 13.81M  │ 1381.46      │ 6303.57        │ $345.37     │ $64.64       │ $410.01    │
└─────────────────────────┴─────────┴──────────────┴────────────────┴─────────────┴──────────────┴────────────┘
```

### JSON Format

```json
{
  "period": {
    "start": "2026-01-01",
    "end": "2026-01-14"
  },
  "pricing": {
    "actionPricePerMillion": 50,
    "activeStoragePricePerGBh": 0.042,
    "retainedStoragePricePerGBh": 0.00105
  },
  "namespaces": [
    {
      "name": "prod-workflows.acct",
      "actions": 12345678,
      "activeStorageGBh": 1234.56,
      "retainedStorageGBh": 5678.90,
      "actionCost": 308.64,
      "activeStorageCost": 51.85,
      "retainedStorageCost": 5.96,
      "totalCost": 366.45
    }
  ],
  "totals": {
    "actions": 13814812,
    "activeStorageGBh": 1381.46,
    "retainedStorageGBh": 6303.57,
    "actionCost": 345.37,
    "storageCost": 64.64,
    "totalCost": 410.01
  }
}
```

## Pricing

Default prices are based on Temporal Cloud's published rates:

- **Actions**: $50 per million (base rate; actual Temporal pricing is tiered with volume discounts)
- **Active Storage**: $0.042 per GBh
- **Retained Storage**: $0.00105 per GBh

See [Temporal Cloud Pricing](https://docs.temporal.io/cloud/pricing) for current rates.

## License

MIT
