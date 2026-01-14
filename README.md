# temporal-cost-report

A command-line tool that fetches usage data from Temporal Cloud and generates cost reports per namespace.

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
go build -o temporal-cost-report .
```

## Prerequisites

You need a Temporal Cloud API key with permissions to read usage data. You can either:

1. Set it as an environment variable:
```bash
export TEMPORAL_API_KEY=your-api-key-here
```

2. Pass it as a command-line argument:
```bash
temporal-cost-report --api-key your-api-key-here
```

To create an API key, see the [Temporal Cloud documentation](https://docs.temporal.io/cloud/api-keys).

## Usage

```bash
# Generate report for current month (default)
temporal-cost-report

# Specify a custom date range
temporal-cost-report --start-date 2025-12-01 --end-date 2025-12-31

# Output as JSON
temporal-cost-report --format json

# Use custom pricing
temporal-cost-report --action-price 30 --active-storage-price 0.05 --retained-storage-price 0.001

# Pass API key directly
temporal-cost-report --api-key your-api-key-here
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

┌────────────────┬─────────────────────────┬────────────────────────┬────────────────────────┬─────────────────┐
│                │         ACTIONS         │     ACTIVE STORAGE     │    RETAINED STORAGE    │      TOTAL      │
├────────────────┼───────┬───────┬─────────┼──────┬───────┬─────────┼──────┬───────┬─────────┼───────┬─────────┤
│ NAMESPACE      │ COUNT │ COST  │    %    │ GBH  │ COST  │    %    │ GBH  │ COST  │    %    │ COST  │    %    │
├────────────────┼───────┼───────┼─────────┼──────┼───────┼─────────┼──────┼───────┼─────────┼───────┼─────────┤
│ prod-workflows │ 1.23M │ $61.73│   89.4% │12.34 │ $0.52 │   89.4% │56.79 │ $0.06 │   90.1% │ $62.31│   89.4% │
│ staging        │ 123.4K│  $6.17│    8.9% │ 1.23 │ $0.05 │    8.9% │ 5.68 │ $0.01 │    9.0% │  $6.23│    8.9% │
│ dev-team-alpha │ 23.46K│  $1.17│    1.7% │ 0.23 │ $0.01 │    1.7% │ 0.57 │ $0.00 │    0.9% │  $1.18│    1.7% │
├────────────────┼───────┼───────┼─────────┼──────┼───────┼─────────┼──────┼───────┼─────────┼───────┼─────────┤
│ TOTAL          │ 1.38M │ $69.07│  100.0% │13.81 │ $0.58 │  100.0% │63.04 │ $0.07 │  100.0% │ $69.72│  100.0% │
└────────────────┴───────┴───────┴─────────┴──────┴───────┴─────────┴──────┴───────┴─────────┴───────┴─────────┘

* Costs are estimates based on the provided pricing and may differ from actual invoiced amounts.
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
      "name": "prod-workflows",
      "actions": 12345678,
      "actionsPercent": 89.37,
      "activeStorageGBh": 1234.56,
      "activeStoragePercent": 89.37,
      "retainedStorageGBh": 5678.90,
      "retainedStoragePercent": 90.10,
      "actionCost": 617.28,
      "activeStorageCost": 51.85,
      "retainedStorageCost": 5.96,
      "totalCost": 675.10,
      "totalCostPercent": 89.38
    }
  ],
  "totals": {
    "actions": 13814812,
    "activeStorageGBh": 1381.46,
    "retainedStorageGBh": 6303.57,
    "actionCost": 690.74,
    "activeStorageCost": 58.02,
    "retainedStorageCost": 6.62,
    "totalCost": 755.38
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
