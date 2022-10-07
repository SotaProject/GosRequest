CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE subnets_tags (
    uuid uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    name text NOT NULL
);
CREATE TABLE subnets (
    uuid uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    name text NOT NULL,
    tag_uuid uuid NOT NULL REFERENCES subnets_tags (uuid)
);

CREATE TABLE subnet_ranges (
    cidr cidr NOT NULL UNIQUE PRIMARY KEY,
    subnet_uuid uuid NOT NULL REFERENCES subnets (uuid)
);
CREATE INDEX ON subnet_ranges USING gist (cidr inet_ops);

CREATE TABLE requests (
    uuid uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    tracker_uuid uuid NOT NULL REFERENCES trackers (uuid),
    url text NOT NULL,
    ip inet NOT NULL,
    subnet_uuid uuid NOT NULL REFERENCES subnets (uuid),
    user_agent text NOT NULL,
    created_at timestamp DEFAULT now() NOT NULL
);
CREATE INDEX ON requests (tracker_uuid);