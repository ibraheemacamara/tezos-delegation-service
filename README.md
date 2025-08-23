# Tezos Delegation Service

This service is responsible for watching for new blocks and delegations and storing them in a database.

## Run

```bash
docker compose up
```

```bash
make run
```

## Run Tests

```bash
make test
```

## API Usage

```bash
curl http://localhost:3000/delegations
```
