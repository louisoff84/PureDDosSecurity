# PureAntiDDoS API — attaques et télémétrie

Cette documentation décrit l'API HTTP locale de PureAntiDDoS et la façon de l'utiliser pour afficher l'état d'une VM et les attaques détectées.

## Base URL

Par défaut :

```
http://IP_DE_LA_VM:2456
```

Exemple local :

```
http://127.0.0.1:2456
```

Pour un dashboard public, il est préférable de ne pas exposer directement le port 2456 sur Internet. Utilisez une API centrale, un reverse proxy authentifié ou un réseau privé.

---

## 1. Vérifier que PureAntiDDoS fonctionne

### GET /health

Retourne l'état du service.

Requête :

```bash
curl http://127.0.0.1:2456/health
```

Réponse :

```json
{
  "ok": true,
  "service": "PureDDosSecurity",
  "time": "2026-10-03T09:00:00Z"
}
```

---

## 2. Obtenir l'état actuel de la VM

### GET /api/v1/status

C'est l'endpoint principal pour un dashboard.

```bash
curl http://127.0.0.1:2456/api/v1/status
```

La réponse contient :

- `attack` : `true` si PureAntiDDoS considère actuellement la VM sous attaque ;
- `host_id` : identifiant de la VM ;
- `time` : date de la mesure ;
- `network` : statistiques réseau ;
- `interfaces` : statistiques de chaque carte réseau ;
- `system` : CPU, RAM, load, uptime, disque et processus ;
- `events` : historique des événements conservés en mémoire.

Exemple :

```json
{
  "time": "2026-10-03T09:00:12Z",
  "host_id": "craftpick-vm-01",
  "attack": true,
  "network": {
    "packets_per_second": 84231,
    "bits_per_second": 583000000,
    "syn_per_second": 60412,
    "udp_per_second": 21000,
    "icmp_per_second": 2811,
    "unique_sources": 1932,
    "tcp_packets": 604120,
    "udp_packets": 210000,
    "icmp_packets": 28110
  },
  "system": {
    "cpu_percent": 91.2,
    "memory_percent": 68.4,
    "memory_used": 146028888064,
    "memory_total": 213070643200,
    "load1": 7.2,
    "load5": 5.9,
    "load15": 4.1,
    "goroutines": 17,
    "uptime_seconds": 91234,
    "disk_used_percent": 54.2,
    "processes": 184
  },
  "interfaces": [
    {
      "name": "eth0",
      "up": true,
      "rx_bytes": 92837718234,
      "tx_bytes": 22371718234,
      "rx_packets": 183771823,
      "tx_packets": 112717182,
      "rx_pps": 82412,
      "tx_pps": 1819,
      "rx_mbps": 571.4,
      "tx_mbps": 12.3,
      "addresses": [
        "192.0.2.10"
      ]
    }
  ],
  "events": [
    {
      "id": "1727946012000000000",
      "time": "2026-10-03T09:00:10Z",
      "host_id": "craftpick-vm-01",
      "severity": "high",
      "active": true,
      "reason": "pps 42150 >= 25000; unique-sources 1200 >= 500",
      "packets_per_second": 42150,
      "bits_per_second": 312000000,
      "unique_sources": 1200,
      "syn_ratio": 0.81,
      "interfaces": [
        "eth0"
      ]
    }
  ]
}
```

### Pour savoir si une attaque est active

Le champ à utiliser sur le site est :

```text
status.attack
```

Logique :

```js
if (status.attack) {
  // Afficher "ATTAQUE EN COURS"
} else {
  // Afficher "Aucune attaque détectée"
}
```

---

## 3. Consulter uniquement les attaques / événements

### GET /api/v1/events

Cet endpoint retourne l'historique des événements détectés, avec un maximum de 100 événements conservés en mémoire par agent.

```bash
curl http://127.0.0.1:2456/api/v1/events
```

Exemple :

```json
[
  {
    "id": "1727946012000000000",
    "time": "2026-10-03T09:00:10Z",
    "host_id": "craftpick-vm-01",
    "severity": "high",
    "active": true,
    "reason": "pps 42150 >= 25000; unique-sources 1200 >= 500",
    "packets_per_second": 42150,
    "bits_per_second": 312000000,
    "unique_sources": 1200,
    "syn_ratio": 0.81,
    "interfaces": [
      "eth0"
    ]
  }
]
```

### Champs d'un événement

| Champ | Description |
|---|---|
| `id` | Identifiant unique de l'événement |
| `time` | Date/heure UTC |
| `host_id` | VM concernée |
| `severity` | `medium`, `high` ou `critical` |
| `active` | État de l'événement |
| `reason` | Raisons ayant déclenché la détection |
| `packets_per_second` | PPS observé |
| `bits_per_second` | BPS observé |
| `unique_sources` | Nombre de sources IP uniques |
| `syn_ratio` | Ratio SYN mesuré |
| `interfaces` | Interfaces concernées |

---

## 4. API Prometheus

### GET /api/v1/metrics

Permet de récupérer quelques métriques simples.

```bash
curl http://127.0.0.1:2456/api/v1/metrics
```

Réponse :

```text
puredos_attack 1
puredos_pps 84231.000000
puredos_bps 583000000.000000
puredos_unique_sources 1932
```

Métriques :

| Métrique | Signification |
|---|---|
| `puredos_attack` | 1 = attaque détectée, 0 = aucune attaque |
| `puredos_pps` | paquets/seconde |
| `puredos_bps` | bits/seconde |
| `puredos_unique_sources` | sources IP uniques |

---

## 5. Authentification de l'API

Pour protéger l'API locale, définir :

```env
PUREDDOS_API_TOKEN=CHANGE_ME
```

Puis redémarrer :

```bash
sudo systemctl restart puredos
```

Requête authentifiée :

```bash
curl \
  -H "Authorization: Bearer CHANGE_ME" \
  http://127.0.0.1:2456/api/v1/status
```

Une requête sans token valide reçoit :

```text
401 Unauthorized
```

---

## 6. Utilisation depuis un site web

Exemple JavaScript :

```js
async function loadAttackStatus() {
  const response = await fetch("https://anti.example.com/api/v1/status", {
    headers: {
      "Authorization": "Bearer CHANGE_ME"
    }
  });

  if (!response.ok) {
    throw new Error("PureAntiDDoS API inaccessible");
  }

  const status = await response.json();

  document.querySelector("#attack-status").textContent =
    status.attack ? "ATTAQUE EN COURS" : "Aucune attaque détectée";

  document.querySelector("#pps").textContent =
    Math.round(status.network.packets_per_second).toLocaleString("fr-FR");

  document.querySelector("#bps").textContent =
    Math.round(status.network.bits_per_second / 1000000) + " Mbit/s";
}

loadAttackStatus();
setInterval(loadAttackStatus, 1000);
```

Configuration CORS :

```env
PUREDDOS_CORS_ORIGIN=https://status.craftpick.fr
```

---

## 7. Architecture recommandée pour plusieurs VM

Pour Craftpick, le dashboard ne devrait pas interroger directement toutes les VM sur Internet.

Architecture recommandée :

```text
VM 01 ─┐
VM 02 ─┼──> PureAntiDDoS API centrale ──> status.craftpick.fr
VM 03 ─┤
VM 04 ─┘
```

Chaque agent envoie sa télémétrie vers l'API centrale :

```env
PUREDDOS_COLLECTOR=https://api.example.com
PUREDDOS_TOKEN=TOKEN_DE_L_AGENT
```

L'agent envoie :

```text
POST /api/v1/telemetry
POST /api/v1/events
```

Le dashboard peut alors demander à l'API centrale :

```text
GET /api/v1/hosts
GET /api/v1/hosts/{host_id}/status
GET /api/v1/attacks
```

Ces endpoints supplémentaires ne font pas partie de l'agent local actuel : ils doivent être implémentés par l'API centrale.

---

## 8. Afficher une attaque sur le dashboard

Un dashboard peut afficher au minimum :

```text
┌─────────────────────────────────────────┐
│ craftpick-vm-01                         │
│                                         │
│ 🔴 ATTAQUE EN COURS                     │
│                                         │
│ 84 231 PPS                              │
│ 583 Mbit/s                              │
│ 1 932 sources uniques                   │
│ SYN ratio : 81 %                        │
│ Interface : eth0                        │
│ Sévérité : HIGH                         │
│                                         │
│ Détectée : 03/10/2026 09:00:10 UTC      │
└─────────────────────────────────────────┘
```

Lorsque `attack === false` :

```text
🟢 Aucune attaque détectée
```

Pour une page multi-VM, utiliser `host_id` comme identifiant de chaque serveur.

---

## 9. Consulter les journaux du service

```bash
sudo systemctl status puredos
```

Logs en temps réel :

```bash
sudo journalctl -u puredos -f
```

Lorsqu'une attaque est détectée, l'agent écrit notamment :

```text
ATTACK DETECTED: pps 42150 >= 25000; unique-sources 1200 >= 500
```

---

## 10. Dépannage API

### Port 2456 inaccessible

Vérifier :

```bash
sudo systemctl status puredos
sudo ss -lntp | grep 2456
```

Tester localement :

```bash
curl http://127.0.0.1:2456/health
```

### Vérifier les interfaces détectées

```bash
ip -br link
```

Puis consulter :

```bash
curl http://127.0.0.1:2456/api/v1/status
```

La liste est disponible dans :

```json
{
  "interfaces": []
}
```

### Vérifier les attaques

```bash
curl http://127.0.0.1:2456/api/v1/events
```

### Vérifier les métriques

```bash
curl http://127.0.0.1:2456/api/v1/metrics
```

---

## 11. Sécurité

Ne publiez pas une API PureAntiDDoS sans authentification directement sur Internet.

Recommandations :

- utiliser `PUREDDOS_API_TOKEN` ;
- limiter le port 2456 avec le firewall ;
- utiliser HTTPS via un reverse proxy si l'API doit être accessible à distance ;
- pour plusieurs VM, privilégier une API centrale ;
- utiliser un token différent par agent ;
- ne pas utiliser l'API locale comme mécanisme de mitigation volumétrique.

## 12. Résumé des endpoints

| Méthode | Endpoint | Utilisation |
|---|---|---|
| GET | `/health` | Santé du service |
| GET | `/api/v1/status` | État temps réel d'une VM |
| GET | `/api/v1/events` | Historique des attaques détectées |
| GET | `/api/v1/metrics` | Métriques simples |

L'information la plus importante pour l'interface web est `/api/v1/status` et son champ `attack`.
