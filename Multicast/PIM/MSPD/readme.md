# Anycast-RP Multicast Lab (PIM-SM + MSDP)

This repository documents a working **Anycast-RP multicast design** built using **PIM Sparse Mode** with **MSDP** to synchronise source state between RPs.

The lab demonstrates:

- Active/active RPs using a shared anycast RP address
- Receiver and source selection of the *closest* RP via IGP
- MSDP Source-Active (SA) exchange inside a **single PIM domain**
- End-to-end multicast forwarding verification

This design mirrors real-world service provider multicast deployments.

---

## Topology Overview

- **CS1 / CS2**  
  Core switches acting as **Anycast Rendezvous Points (RPs)**

- **TOR switches**  
  Act as first-hop / last-hop PIM routers for sources and receivers

- **Anycast RP Address**

  ```
  192.168.10.1/32
  ```

- **Unique Loopbacks (for MSDP)**

  ```
  CS1: 10.255.0.1/32
  CS2: 10.255.0.2/32
  ```

- **Multicast Group Range**

  ```
  239.0.0.0/8
  ```

IGP (OSPF) determines which RP is *closest* to each router.

---

## Design Goals

- Receivers join their **nearest RP**
- Sources register with their **nearest RP**
- Multiple RPs active simultaneously
- Seamless failover without RP re-election
- No BSR / CRP dependency
- MSDP used **only** to synchronise source knowledge between RPs

---

## Multicast Architecture

### RP Model

- **Anycast-RP**
- Static RP mapping everywhere:

  ```
  239.0.0.0/8 → 192.168.10.1
  ```

- No Candidate-RP (CRP)
- No Bootstrap Router (BSR)

### Why MSDP Is Required

With Anycast-RP, multiple RPs are active at the same time.  
Each RP only learns about sources that register *to itself*.

**MSDP shares Source-Active (SA) information between the RPs**, ensuring:

- Receivers using CS1 can learn about sources registered to CS2
- Receivers using CS2 can learn about sources registered to CS1

This is required **even in a single PIM domain** when using Anycast-RP.

---

## Protocol Summary

| Protocol | Purpose |
|--------|---------|
| OSPF | IGP for unicast reachability and RP selection |
| PIM-SM | Multicast tree construction |
| MSDP | RP-to-RP source discovery |
| IGMP | Receiver group membership |

---

## Configuration Highlights

### PIM

- Enabled on all routed VLANs
- Enabled on RP loopback VLAN
- Static RP mapping to anycast RP address
- **CRP / BSR explicitly not used**

### MSDP

- Peer between CS1 and CS2
- Peering sourced from unique loopbacks
- Mesh-group used to prevent loops

### OSPF

- Loopbacks advertised
- P2P links configured correctly
- Metric tuning can be used to influence RP preference

---

## Verification Checklist

### RP Selection

```
show iproute 192.168.10.1/32
```

### PIM State

```
show pim cache
```

### MSDP

```
show msdp peer
show msdp sa-cache
```

---

## Key Takeaways

- Anycast-RP ≠ Dynamic RP election
- Anycast-RP **requires MSDP**
- CRP/BSR must not be mixed with Anycast-RP
- MSDP is purely control-plane
- RP selection is driven entirely by unicast routing

---

## References

- RFC 4610 — Anycast-RP using MSDP
- RFC 3618 — MSDP
- RFC 4601 — PIM-SM

---
