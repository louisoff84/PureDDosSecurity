# PureDDosSecurity

PureDDosSecurity is a Linux VM/network telemetry and DDoS detection agent for Craftpick infrastructure.

## Features

- Packet capture on every active non-loopback network interface.
- IPv4 and IPv6 source tracking.
- Packets/s and bits/s measurement.
- TCP/SYN, UDP and ICMP traffic counters.
- Unique source count.
- Kernel RX/TX counters for all interfaces.
- VM CPU, RAM, load, uptime and agent health.
- Live HTTP API on port 2456.
- Optional central collector with Bearer authentication.
- systemd service with automatic restart.
- Bounded in-memory incident history.
- Prometheus-style metrics endpoint.

## API

GET /health
GET /api/v1/status
GET /api/v1/events
GET /api/v1/metrics

If PUREDDOS_API_TOKEN is set, send Authorization: Bearer <token>. Set PUREDDOS_CORS_ORIGIN for browser dashboards hosted on another origin.

The status response contains the current attack flag, interface telemetry, system telemetry and network indicators.

## Central collector

Set PUREDDOS_COLLECTOR to your central API URL.

The agent sends:
- POST /api/v1/telemetry for live snapshots.
- POST /api/v1/events when a new attack is detected.

Example:

PUREDDOS_COLLECTOR=https://api.example.com
PUREDDOS_TOKEN=secret-token

The central API should authenticate the token, validate the HostID, rate-limit ingestion and store events in a database.

## Install

    git clone https://github.com/louisoff84/PureDDosSecurity.git
    cd PureDDosSecurity
    sudo bash scripts/install.sh

Then configure:

    sudo nano /etc/puredos/puredos.env
    sudo systemctl restart puredos

Check:

    systemctl status puredos
    curl http://127.0.0.1:2456/health
    curl http://127.0.0.1:2456/api/v1/status

## Detection

Detection is heuristic and uses a rolling window. Default thresholds are:

- 25,000 packets/s
- 100 Mbit/s
- 70% SYN ratio with meaningful packet volume
- 500 unique source IPs

Tune them for each VM. A threshold crossing is an alert signal, not proof of an attack.

## Architecture

VM
  -> PureDDosSecurity agent
  -> packet capture + interface counters + system metrics
  -> local HTTP API :2456
  -> optional central API
  -> dashboard/status website

PureDDosSecurity detects and reports. It does not magically absorb a volumetric attack that has already saturated the VM uplink. Provider-side filtering, Gcore, Cloudflare, firewall ACLs or upstream mitigation are still required for large attacks.
