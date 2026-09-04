# Deployments

## Local

```bash
cp deployments/.env.example deployments/.env   # fill in Supabase + JWT values
./scripts/rebuild.sh                            # or: cd deployments && docker compose up --build
```

| Service        | Host port | Notes                                              |
|----------------|-----------|----------------------------------------------------|
| web            | 3000      | Next.js client + `/api/*`                          |
| authoritative  | 50051     | gRPC `OrderHandler` (web → authoritative)          |
| authoritative  | 8080      | robot WebSocket hub at `/ws`                       |
| command        | 8082/tcp  | TCP relay (container port 8080)                    |
| command        | 8081/udp  | UDP relay                                          |
| kafka          | 9092      | host listener for `go run` apps (`localhost:9092`) |
| kafka-ui       | 8085      | http://localhost:8085                              |

Inside the compose network apps use `KAFKA_BROKERS=kafka:9093` and
`GRPC_SERVER_URL=authoritative:50051`.

Just the broker, running the Go apps natively:

```bash
cd deployments && docker compose up kafka kafka-ui
cd apps/authoritative && go run ./cmd/authoritative      # uses localhost:9092
```

## Production (GCE VM)

`.github/workflows/deploy.yml` runs on every push to `main`: it writes
`deployments/.env` from the `prod` environment secrets, authenticates to GCP via
Workload Identity Federation, copies the repo to `/opt/delivery-gdg` on the VM,
and runs

```bash
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build
```

The prod overlay adds Caddy (auto-TLS for `DOMAIN_NAME` → `web:3000`) and
removes host port bindings for Kafka, kafka-ui, gRPC and the web server; only
80/443, 8080 (robot WS), 8082/tcp + 8081/udp (command) stay exposed.
Required secrets and VM setup are documented at the top of the workflow file.
