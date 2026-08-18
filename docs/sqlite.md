# SQLite persistence

Pass `-db /path/to/perfcheck.db` to `perfcheckd` to persist audits across restarts.

```bash
mkdir -p data
cd api && go run ./cmd/perfcheckd -db ../data/perfcheck.db
```

An empty `-db` (the default) keeps the in-memory store, which is what tests and quick demos use.

## Docker

The database file must live on a volume, otherwise it disappears when the container stops:

```bash
docker run --rm -p 8080:8080 \
  -v perfcheck-data:/data \
  perfcheck
```

The image entrypoint will use `-db /data/perfcheck.db` once that flag is wired in the Dockerfile.
