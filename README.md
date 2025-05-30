# nebula-sync

[![Release version](https://img.shields.io/github/v/release/lovelaze/nebula-sync)](https://github.com/lovelaze/nebula-sync/releases/latest)
[![Tests](https://img.shields.io/github/actions/workflow/status/lovelaze/nebula-sync/test.yml?branch=main&label=tests)](https://github.com/lovelaze/nebula-sync/actions/workflows/test.yml?query=branch%3Amain)
![Go version](https://img.shields.io/github/go-mod/go-version/lovelaze/nebula-sync)
[![Docker image size](https://img.shields.io/docker/image-size/lovelaze/nebula-sync/latest)](https://hub.docker.com/r/lovelaze/nebula-sync)

Synchronize Pi-hole v6.x configuration to replicas.

This project is not a part of the [official Pi-hole project](https://github.com/pi-hole), but uses the api provided by Pi-hole instances to perform the synchronization actions.

## Features
- **Full sync**: Use Pi-hole Teleporter for full synchronization.
- **Selective sync**: Selective feature synchronization.
- **Cron schedule**: Run on cron schedule.

## Installation


### Linux/OSX binary
Download binary from the [latest release](https://github.com/lovelaze/nebula-sync/releases/latest) or build from source:
```
go install github.com/lovelaze/nebula-sync@latest
```

Run binary:
```bash
# run
nebula-sync run

# read envs from file
nebula-sync run --env-file .env
```

### Docker Compose (recommended)

```yaml
---
services:
  nebula-sync:
    image: ghcr.io/lovelaze/nebula-sync:latest
    container_name: nebula-sync
    environment:
    - PRIMARY=http://ph1.example.com|password
    - REPLICAS=http://ph2.example.com|password,http://ph3.example.com|password
    - FULL_SYNC=true
    - RUN_GRAVITY=true
    - CRON=0 * * * *
```

### Docker CLI

```bash
docker run --rm \
  --name nebula-sync \
  -e PRIMARY="http://ph1.example.com|password" \
  -e REPLICAS="http://ph2.example.com|password" \
  -e FULL_SYNC=true \
  -e RUN_GRAVITY=true \
  ghcr.io/lovelaze/nebula-sync:latest
```

## Examples
Env and docker-compose examples can be found [here](https://github.com/lovelaze/nebula-sync/tree/main/examples)

## Configuration

The following environment variables can be specified:

### Required Environment Variables

| Name      | Default | Example                                          | Description                                              |
|-----------|---------|--------------------------------------------------|----------------------------------------------------------|
| `PRIMARY` | n/a     | `http://ph1.example.com\|password`                       | Specifies the primary Pi-hole configuration              |
| `REPLICAS`| n/a     | `http://ph2.example.com\|password,http://ph3.example.com\|password` | Specifies the list of replica Pi-hole configurations     |
| `FULL_SYNC` | n/a   | `true`                                           | Specifies whether to perform a full synchronization      |

> **Note:** When `FULL_SYNC=true`, the system will perform a full Teleporter import/export from the primary Pi-hole to the replicas. This will synchronize all settings and configurations.

> **Docker secrets:** `PRIMARY` and `REPLICAS` environment variables support Docker secrets when defined as `PRIMARY_FILE` and `REPLICAS_FILE`. See note regarding default user and Docker secrets example below.

### Optional Environment Variables

| Name                               | Default | Example         | Description                                        |
|------------------------------------|---------|-----------------|----------------------------------------------------|
| `CRON`                             | n/a     | `0 * * * *`     | Specifies the cron schedule for synchronization    |
| `RUN_GRAVITY`                      | false   | true            | Specifies whether to run gravity after syncing     |
| `TZ`                               | n/a     | `Europe/London` | Specifies the timezone for logs and cron           |
| `CLIENT_SKIP_TLS_VERIFICATION`     | false   | true            | Skips TLS certificate verification                 |
| `CLIENT_RETRY_DELAY_SECONDS`       | 1       | 5               | Seconds to delay between connection attempts       |
| `CLIENT_TIMEOUT_SECONDS`           | 20      | 60              | Http client timeout in seconds                     |

> **Note:** The following optional settings apply only if `FULL_SYNC=false`. They allow for granular control of synchronization if a full sync is not wanted.

| Name                              | Default | Description                            |
|-----------------------------------|---------|----------------------------------------|
| `SYNC_CONFIG_DNS`                  | false   | Synchronize DNS settings               |
| `SYNC_CONFIG_DHCP`                 | false   | Synchronize DHCP settings              |
| `SYNC_CONFIG_NTP`                  | false   | Synchronize NTP settings               |
| `SYNC_CONFIG_RESOLVER`             | false   | Synchronize resolver settings          |
| `SYNC_CONFIG_DATABASE`             | false   | Synchronize database settings          |
| `SYNC_CONFIG_MISC`                 | false   | Synchronize miscellaneous settings     |
| `SYNC_CONFIG_DEBUG`                | false   | Synchronize debug settings             |
| `SYNC_GRAVITY_DHCP_LEASES`         | false   | Synchronize DHCP leases                |
| `SYNC_GRAVITY_GROUP`               | false   | Synchronize groups                     |
| `SYNC_GRAVITY_AD_LIST`             | false   | Synchronize ad lists                   |
| `SYNC_GRAVITY_AD_LIST_BY_GROUP`    | false   | Synchronize ad lists by group          |
| `SYNC_GRAVITY_DOMAIN_LIST`         | false   | Synchronize domain lists               |
| `SYNC_GRAVITY_DOMAIN_LIST_BY_GROUP`| false   | Synchronize domain lists by group      |
| `SYNC_GRAVITY_CLIENT`              | false   | Synchronize clients                    |
| `SYNC_GRAVITY_CLIENT_BY_GROUP`     | false   | Synchronize clients by group           |


#### Config filters
> Allows including or excluding specific config keys.\
**Note:** `The SYNC_CONFIG_*_INCLUDE` and `SYNC_CONFIG_*_EXCLUDE` settings are mutually exclusive within each section. Additionally, config filters are only applied if `FULL_SYNC=false`.\
Config keys are relative to the section and are **case sensitive**. For example, the key `dns.upstreams` should be referred to as `upstreams`, and `dns.cache.size` should be referred to as `cache.size`.

| Name                              | Example                    | Description                                     |
|-----------------------------------|----------------------------|-------------------------------------------------|
| `SYNC_CONFIG_DNS_INCLUDE`         | upstreams,interface        | DNS config keys to include                     |
| `SYNC_CONFIG_DNS_EXCLUDE`         | upstreams,interface        | DNS config keys to exclude                     |
| `SYNC_CONFIG_DHCP_INCLUDE`        | active,start               | DHCP config keys to include                    |
| `SYNC_CONFIG_DHCP_EXCLUDE`        | active,start               | DHCP config keys to exclude                    |
| `SYNC_CONFIG_NTP_INCLUDE`         | ipv4,sync                  | NTP config keys to include                     |
| `SYNC_CONFIG_NTP_EXCLUDE`         | ipv4,sync                  | NTP config keys to exclude                     |
| `SYNC_CONFIG_RESOLVER_INCLUDE`    | resolveIPv4,networkNames   | Resolver config keys to include                |
| `SYNC_CONFIG_RESOLVER_EXCLUDE`    | resolveIPv4,networkNames   | Resolver config keys to exclude                |
| `SYNC_CONFIG_DATABASE_INCLUDE`    | DBimport,maxDBdays         | Database config keys to include                |
| `SYNC_CONFIG_DATABASE_EXCLUDE`    | DBimport,maxDBdays         | Database config keys to exclude                |
| `SYNC_CONFIG_MISC_INCLUDE`        | nice,delay_startup         | Misc config keys to include                    |
| `SYNC_CONFIG_MISC_EXCLUDE`        | nice,delay_startup         | Misc config keys to exclude                    |
| `SYNC_CONFIG_DEBUG_INCLUDE`       | database,networking        | Debug config keys to include                   |
| `SYNC_CONFIG_DEBUG_EXCLUDE`       | database,networking        | Debug config keys to exclude                   |

## Notes

### Default user of Docker container / Docker secrets example
By default, the Docker container runs as user `1001`. If you are using Docker secrets, the user that is running the container will need read permissions to the files that the Docker secrets reference. If the user does not have the right permissions you will receive an error `Failed to initialize service error="open /run/secrets/primary: permission denied"`. To avoid this error, either make sure to `chown 1001 ./your/secretfiles && chmod 400 ./your/secretfiles` or use the [`user` directive in Docker Compose](https://docs.docker.com/reference/compose-file/services/#user) to change the user that the container runs as to a user of your choice - and then make sure to update your secret files' ownership to that user. In the example [docker-compose-with-secrets.yml](examples/docker-compose-with-secrets.yml), user `1234` owns `./secrets/primary.txt` and `./secrets/replicas.txt` and both have `-r--------` permissions.


## Disclaimer

This project is an unofficial, community-maintained project and is not affiliated with the [official Pi-hole project](https://github.com/pi-hole). It aims to add sync/replication features not available in the core Pi-hole product but operates independently of Pi-hole LLC. Although tested across various environments, using any software from the Internet involves inherent risks. See the [license](https://github.com/lovelaze/nebula-sync/blob/main/LICENSE) for more details.

Pi-hole and the Pi-hole logo are [registered trademarks](https://pi-hole.net/trademark-rules-and-brand-guidelines) of Pi-hole LLC.


