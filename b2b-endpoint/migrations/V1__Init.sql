CREATE SCHEMA endpoint;

CREATE TABLE endpoint.clients
(
    id   UUID PRIMARY KEY,
    name varchar(64) NOT NULL
);

CREATE TABLE endpoint.secrets
(
    api_key    UUID PRIMARY KEY,
    api_secret TEXT                                  NOT NULL,
    client_id  UUID REFERENCES endpoint.clients (id) NOT NULL,
    is_active  BOOLEAN                               NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP                             NOT NULL DEFAULT NOW(),
    last_used  TIMESTAMP                             NOT NULL DEFAULT NOW()
);

CREATE TABLE endpoint.whitelist_ips
(
    api_key    UUID REFERENCES endpoint.secrets (api_key) ON DELETE CASCADE NOT NULL,
    ip_address VARCHAR(39)                                                  NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()                                      NOT NULL,

    PRIMARY KEY (api_key, ip_address)
);

CREATE TABLE endpoint.request_logs
(
    api_key      UUID REFERENCES endpoint.secrets (api_key) NOT NULL,
    request_hash VARCHAR(64)                                NOT NULL,
    created_at   TIMESTAMP                                  NOT NULL,

    PRIMARY KEY (api_key, request_hash)
);

CREATE INDEX idx_api_key_request_hash ON endpoint.request_logs (api_key, request_hash);
