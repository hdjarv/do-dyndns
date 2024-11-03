# Digital Ocean Dynamic DNS Updater

Tool to implement dynamic DNS using Digital Ocean's [DNS](https://docs.digitalocean.com/reference/api/api-reference/#tag/Domain-Records) API.

## Prerequisites

- You must have [Go](https://go.dev) programming installed to build the software locally
- You must have [Docker](https://www.docker.com) installed to build a Docker image of the software
- You must have an account at [Digital Ocean](https://www.digitalocean.com) and a domain setup there. The domain must contain the DNS (A) record you want this tool to update.

## Build

To download the repository and build the project using locally installed Go tools:

```shell
git clone <repo-url>
cd go-do-dyndns
go mod download
CGO_ENABLED=0 go build -ldflags="-w -s" -o build/do-dyndns
```

> Note: Adjust the `<repo-url>` to where you fetched this repository from. It is published to multiple code hosting services.

## Build Docker image

To download the repository and build a Docker image:

```shell
git clone <repo-url>
cd go-dyn-dns
docker build -t do-dyndns .
```

> Note: Adjust the `<repo-url>` to where you fetched this repository from. It is published to multiple code hosting services.

## Install

```shell
sudo cp build/do-dyndns /usr/local/bin
```

## Configure

This tool is configured using environment variables.

| Name                 | Required | Description                                                                                                                                           |
| -------------------- | -------- | ----------------------------------------------------------------------------------------------------------------------------------------------------- |
| `DO_DYN_EXT_IP_URL`  | Yes      | A URL to use for checking the external IP-address. For example: `https://ipinfo.io`.                                                                  |
| `DO_DYN_DO_DOMAIN`   | Yes      | The domain in Digital Ocean DNS to update the record for. Must exist before running this tool.                                                        |
| `DO_DYN_RECORD_NAME` | Yes      | The DNS record in the Digital Ocean DNS domain to update. Must be an exiting `A` DNS record before running this tool.                                 |
| `DO_DYN_API_TOKEN`   | Yes      | An API token for the Digital Ocean API to use for updating the DNS record.                                                                            |
| `DO_DYN_IP_REGEX`    | No       | A regular expression for finding the IP-address in the data from `DO_DYN_EXT_IP_URL`. Default is `\\b(\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3})\\b`. |
| `DO_DYN_DNS_TTL`     | No       | Set the TTL (time to live) for the DNS record to update. Default is `60`.                                                                               |
| `DO_DYN_DRY_RUN`     | No       | Set to `1` to enable dry-run mode. In this mode the DNS record will not be updated. Tool will output what it would do. Default is `0`.                |

## Run

```shell
DO_DYN_EXT_IP_URL=https://ipinfo.io DO_DYN_API_TOKEN=__insert_digitalocean_api_token_here__ DO_DYN_DO_DOMAIN=example.com DO_DYN_RECORD_NAME=dynamic do-dyndns
```

## Run with Docker

```shell
docker run --rm --name do-dyndns --env DO_DYN_EXT_IP_URL=https://ipinfo.io --env DO_DYN_API_TOKEN=__insert_digitalocean_api_token_here__ --env DO_DYN_DO_DOMAIN=example.com --env DO_DYN_RECORD_NAME=dynamic do-dyndns
```

### Run as systemd service

Install the systemd unit files:

```shell
sudo cp systemd/do-dyndns.{timer,service} /etc/systemd/system/
```

Configure the service by editing the `/etc/systemd/system/do-dyndns.service` file to set the environment variables:

```shell
sudo nano /etc/systemd/system/do-dyndns.service
```

After editing the file you must reload the systemd daemon with:

```shell
sudo systemctl daemon-reload
```

Finally enable the systemd units and start the timer unit:

```shell
sudo systemctl enable do-dyndns.service
sudo systemctl enable do-dyndns.timer
sudo systemctl start do-dyndns.timer
```

You can check the service's logs with:

```shell
sudo journalctl -u do-dyndns.service
```

## License

This software is licensed under the MIT License - see the `LICENSE` file for details.
