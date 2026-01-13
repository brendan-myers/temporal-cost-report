# temporal-cost-report

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

┌───────────────────┬─────────┬─────────────┬──────────┬──────────────┬─────────────┬─────────┬────────────────┬───────────────┬──────────┬────────────┐
│ NAMESPACE         │ ACTIONS │ ACTION COST │ ACTIONS% │ ACTIVE (GBH) │ ACTIVE COST │ ACTIVE% │ RETAINED (GBH) │ RETAINED COST │ RETAINED%│ TOTAL COST │
├───────────────────┼─────────┼─────────────┼──────────┼──────────────┼─────────────┼─────────┼────────────────┼───────────────┼──────────┼────────────┤
│ prod-workflows    │ 12.35M  │     $617.28 │   89.37% │      1234.56 │      $51.85 │  89.37% │        5678.90 │         $5.96 │   90.10% │    $675.10 │
│ staging           │  1.23M  │      $61.73 │    8.94% │       123.45 │       $5.18 │   8.94% │         567.89 │         $0.60 │    9.01% │     $67.51 │
│ dev-team-alpha    │ 234.57K │      $11.73 │    1.70% │        23.45 │       $0.98 │   1.70% │          56.78 │         $0.06 │    0.90% │     $12.77 │
├───────────────────┼─────────┼─────────────┼──────────┼──────────────┼─────────────┼─────────┼────────────────┼───────────────┼──────────┼────────────┤
│ TOTAL             │ 13.81M  │     $690.74 │  100.00% │      1381.46 │      $58.02 │ 100.00% │        6303.57 │         $6.62 │  100.00% │    $755.38 │
└───────────────────┴─────────┴─────────────┴──────────┴──────────────┴─────────────┴─────────┴────────────────┴───────────────┴──────────┴────────────┘
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
      "totalCost": 675.10
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
