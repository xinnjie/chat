# Topic Component in Tinode

The Topic is a fundamental component in Tinode's architecture, representing an isolated communication channel for users to exchange messages. It serves as the container for conversations and manages all aspects of user interactions within those conversations.

## Overview

A Topic in Tinode is an isolated communication channel that can represent:
- A group chat
- A 1:1 conversation
- A personal "me" topic for managing user's own data and subscriptions
- A discovery "fnd" topic for finding other users and groups

Topics handle message routing, access control, presence notifications, and subscription management, forming the core of Tinode's messaging capabilities.

## Topic Structure

Each Topic instance manages:

1. **Basic Metadata**:
   - Name (expanded/unique identifier)
   - Original name (user-facing name)
   - Topic category (me, p2p, grp, fnd)
   - Creation and update timestamps
   - Owner information

2. **Access Control**:
   - Default access modes for authenticated and anonymous users
   - Per-user access permissions (want and given)
   - Union access modes for proxy topics

3. **Data Management**:
   - Topic's public data (visible to all subscribers)
   - Topic's trusted data (server-managed)
   - Per-user private data
   - Message sequencing (lastID, delID)

4. **Sessions**:
   - Map of attached sessions with their metadata
   - Subscription and presence tracking

5. **State Management**:
   - Status flags (loaded, paused, read-only, deleted)
   - Call management for p2p topics

## Topic Lifecycle

1. **Creation**:
   - Topics are created when a user subscribes to a non-existent topic
   - The Hub routes subscription requests to appropriate handlers
   - New topics are initialized with default parameters

2. **Active State**:
   - Topics receive and process messages
   - Handle subscription updates
   - Manage access control changes
   - Send presence notifications

3. **Termination**:
   - Topics can be paused (temporarily reject messages)
   - Marked as deleted (for cleanup)
   - Unregistered from the Hub
   - Resources released when no longer needed

## Key Functionalities

### Message Handling

Topics process several message types:
- `pub` messages (user content)
- `meta` messages (topic metadata operations)
- `sub`/`leave` messages (subscription management)
- `pres` messages (presence notifications)
- `del` messages (content deletion)

For each message type, the Topic:
1. Validates permissions
2. Processes the request
3. Updates local state
4. Persists changes to the database
5. Sends appropriate notifications

### Subscription Management

Topics handle:
- New subscription requests
- Subscription updates (changing permissions)
- Unsubscribe requests
- User eviction

Each operation involves:
- Access control checks
- Database updates
- Notification dispatch
- Session management

### Access Control

Topics enforce a sophisticated access control system:
- Each user has "want" permissions (what they request)
- Each user has "given" permissions (what the topic grants)
- Effective permissions are the intersection of want and given
- Permissions control abilities like reading, writing, and managing the topic

### Presence Notifications

Topics generate and route presence notifications:
- User joining/leaving topics
- Subscription changes
- User coming online/going offline
- Message read/received receipts

## Integration with Other Components

### Hub Integration

The Topic works closely with the Hub:
- Hub creates and registers topics
- Hub routes messages to appropriate topics
- Topics notify Hub when they should be unregistered
- Hub manages topic shutdown during system termination

### Session Integration

Topics maintain connections to user sessions:
- Track which sessions are subscribed
- Route messages to appropriate sessions
- Handle session attachment/detachment
- Manage background vs. foreground session states

### Store Integration

Topics interact with the persistence layer:
- Load subscriptions and topic data
- Save messages and metadata changes
- Manage message deletion (soft and hard)
- Handle topic deletion

## Special Topic Types

### Me Topic

The "me" topic is special:
- One per user
- Manages user's contact list
- Handles credential management
- Routes notifications to the user
- Tracks subscriptions to other topics

### P2P Topics

P2P (person-to-person) topics:
- Connect exactly two users
- Have special naming conventions (composite of two user IDs)
- Manage specialized presence information
- Support video call functionality

### Group Topics

Group topics:
- Support multiple subscribers
- Have owner-based permissions
- Can be channels (broadcast-only)
- Support discovery via tags

### Fnd Topic

The "fnd" (find) topic:
- Handles user and group discovery
- Processes search queries
- Returns results as ephemeral subscriptions

## Clustering Support

In clustered environments:
- Topics can be "proxy" topics (local representation of remote topics)
- Master topics handle the actual data and logic
- Proxy topics forward messages to master topics
- System handles topic rehashing during cluster changes

## Performance Considerations

Topics are designed for performance:
- Non-blocking message processing
- Efficient permission checking
- Optimized presence notification routing
- Background session support to reduce notification overhead
