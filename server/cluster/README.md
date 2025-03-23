# Tinode Clustering System

This directory contains the clustering implementation for the Tinode chat server, enabling horizontal scaling across multiple server instances.

## Overview

The Tinode clustering system allows distributing the workload across multiple server nodes by:

1. **Distributing topic handling**: Using ring hash distribution to assign topics to specific nodes
2. **Routing messages**: Efficiently sending messages between nodes that need to communicate
3. **Managing node status**: Tracking node health, handling failures, and leader election

## Architecture

### Core Components

- **ClusterNode**: Represents a connection to a remote node with RPC endpoints
- **Cluster**: Central structure that manages the entire node communication system
- **Topic Master/Proxy Pattern**: Topics are owned by a single "master" node, with other nodes acting as "proxies"
- **Leader Election**: Based on Raft protocol concepts for leader node selection and failover

### Message Routing

The cluster routes several types of messages:
- Client requests from proxy to topic master
- Server responses from topic master to proxy
- User cache updates and push notifications between nodes

### Session Handling

- **Multiplexing Sessions**: Efficiently handle communication between nodes
- **ClusterSess**: Maintains information about remote sessions
- **Proxy Sessions**: Forward requests from one node to the topic master node

## Fault Tolerance

The clustering implementation includes:
- **Node Failure Detection**: Through heartbeats and fingerprints
- **Leader Election**: For coordinating cluster-wide decisions
- **Automatic Reconnection**: When nodes become available again
- **Ring Hash Redistribution**: When nodes join or leave the cluster

## Data Flow

1. When a client connects to any node, that node becomes the session's entry point
2. If the client accesses a topic hosted on another node, the entry point acts as a proxy
3. The proxy forwards client requests to the topic master
4. The topic master processes requests and returns responses to the proxy
5. The proxy forwards responses back to the client

This architecture allows Tinode to scale horizontally while maintaining consistent topic ownership and efficient message routing.