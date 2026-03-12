```mermaid 
sequenceDiagram
    participant Client
    participant UnixDomainSock
    participant Daemon
    participant DB
    
    Daemon-->>UnixDomainSock: Listen for requests
    
    Client->>UnixDomainSock: HTTP (GET /list)
    UnixDomainSock->>Daemon: Forward
    Daemon->>+DB: ListAllDocs

    DB-->>-Daemon
    Daemon-->>UnixDomainSock: HTTP Response
    UnixDomainSock-->>Client: All documentation available
```
