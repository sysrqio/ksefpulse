# ksefpulse

KSeF 2.0 **reconciliation** and **FA(3) pre-audit** CLI for ERP integrators.

Local-only `audit` — no network calls, no telemetry.

## Quick start

```bash
go build -o ksefpulse ./cmd/ksefpulse
./ksefpulse audit --state fixtures/state-orphan.json --format table
# exit code 2 = policy/reconciliation failure
```

## Commands

- `ksefpulse audit --state submission-state.json [--fa3 invoice.xml]`
- `ksefpulse validate --fa3 invoice.xml`
- `ksefpulse version`

## License

MIT

## Incubator

Product dossier: [sysrq-io/Incubator](https://github.com/sysrq-io/Incubator) — `KsefPulse/docs/`.
