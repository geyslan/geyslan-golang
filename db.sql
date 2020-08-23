CREATE USER oowlish PASSWORD 'oowlish';
CREATE DATABASE oowlish OWNER oowlish;
GRANT ALL PRIVILEGES ON DATABASE oowlish to oowlish;
\c oowlish

CREATE TABLE logs (
    id bigserial primary key,
    client_host varchar NOT NULL,
    rfc1413 varchar,
    remote_user varchar,
    date_time timestamptz NOT NULL,
    method varchar,
    resource varchar NOT NULL,
    protocol varchar,
    status smallint,
    size bigint NOT NULL,
    referer varchar,
    user_agent varchar,
    UNIQUE (client_host, date_time, resource, size)
);

ALTER TABLE logs OWNER TO oowlish;
GRANT ALL PRIVILEGES ON TABLE logs to oowlish;
GRANT USAGE, SELECT ON SEQUENCE logs_id_seq TO oowlish;


GRANT ALL ON ALL TABLES IN SCHEMA public TO oowlish;