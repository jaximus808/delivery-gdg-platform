# WashU GDG Delivery Robot Monorepo

<img src="https://github.com/jaximus808/delivery-gdg-platform/blob/main/assets/gdg_logo.jpg" width="48">

---

This monorepo will hold our main components being our fullstack, authserver, and command server platforms

More work to come :)

## Running locally

Prereqs: Docker Desktop (with `docker compose` v2). Go 1.25+ and Node 22+ only if you want to run apps outside Docker.

```bash
# 1. Configure secrets (Supabase URL/keys + a JWT secret)
cp deployments/.env.example deployments/.env
#    ...edit deployments/.env

# 2. Build and start everything (Kafka, authoritative, command, web)
./scripts/rebuild.sh            # foreground; Ctrl-C to stop
./scripts/rebuild.sh -d         # or detached

# 3. Useful commands
./scripts/rebuild.sh ps         # status
./scripts/rebuild.sh logs -f    # tail logs (add a service name to filter)
./scripts/rebuild.sh down       # stop + remove containers
```

| What            | URL / port                      |
|-----------------|---------------------------------|
| Web client      | http://localhost:3000           |
| Kafka UI        | http://localhost:8085           |
| gRPC (authoritative) | localhost:50051            |
| Robot WebSocket | ws://localhost:8080/ws          |
| Command TCP / UDP | localhost:8082 / localhost:8081 |
| Kafka (host)    | localhost:9092                  |

Run only Kafka in Docker and the apps natively:

```bash
cd deployments && docker compose up -d kafka kafka-ui
cd apps/authoritative && go run ./cmd/authoritative   # reads .env, uses localhost:9092
cd apps/command && go run . -mode=server
cd apps/client/web && npm install && npm run dev
```

See `deployments/README.md` for the full port table and the production/GCP deploy flow.
