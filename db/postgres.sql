BEGIN;

DROP TABLE IF EXISTS user_role;
DROP TABLE IF EXISTS domain_user;
DROP TABLE IF EXISTS role_permission;
DROP TABLE IF EXISTS domain CASCADE;
DROP TABLE IF EXISTS "permission" CASCADE;
DROP TABLE IF EXISTS "role" CASCADE;
DROP TABLE IF EXISTS "user" CASCADE;
DROP TABLE IF EXISTS user_session;

CREATE TABLE domain (
    domain_id serial NOT NULL,
    object_id character varying(255) NOT NULL,
    name character varying(255) NOT NULL,
    description text,
    is_enabled boolean,
    created_on timestamp NOT NULL,
    updated_on timestamp NOT NULL,
    CONSTRAINT domain_pkey PRIMARY KEY (domain_id),
    UNIQUE(object_id),
    UNIQUE(name)
);

CREATE TABLE "user" (
    user_id serial NOT NULL,
    object_id character varying(255) NOT NULL,
    name character varying(255) NOT NULL,
    passwd character varying(255) NOT NULL,
    is_enabled boolean,
    created_on timestamp NOT NULL,
    updated_on timestamp NOT NULL,
    CONSTRAINT user_pkey PRIMARY KEY (user_id),
    UNIQUE(object_id),
    UNIQUE(name)
);

CREATE TABLE domain_user (
    domain_id integer NOT NULL REFERENCES domain ON DELETE CASCADE,
    user_id integer NOT NULL REFERENCES "user" ON DELETE CASCADE,
    CONSTRAINT domain_user_pkey PRIMARY KEY (domain_id, user_id)
);

CREATE TABLE "role" (
    role_id serial NOT NULL,
    name character varying(255) NOT NULL,
    description text,
    is_enabled boolean,
    CONSTRAINT role_pkey PRIMARY KEY (role_id),
    UNIQUE(name)
);

CREATE TABLE "permission" (
    permission_id serial NOT NULL,
    name character varying(255) NOT NULL,
    description text,
    evaluation_rule text,
    is_enabled boolean,
    CONSTRAINT permission_pkey PRIMARY KEY (permission_id),
    UNIQUE(name)
);

CREATE TABLE role_permission (
    role_id integer NOT NULL REFERENCES "role" ON DELETE CASCADE,
    permission_id integer NOT NULL  REFERENCES "permission" ON DELETE CASCADE,
    CONSTRAINT role_permission_pkey PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE user_role (
    user_id integer NOT NULL REFERENCES "user" ON DELETE CASCADE,
    role_id integer NOT NULL REFERENCES "role" ON DELETE CASCADE,
    CONSTRAINT user_role_pkey PRIMARY KEY (user_id, role_id)
);

CREATE TABLE user_session (
    user_session_id character varying(255) NOT NULL,
    domain_id integer NOT NULL REFERENCES domain ON DELETE CASCADE,
    user_id integer NOT NULL REFERENCES "user" ON DELETE CASCADE,
    user_agent character varying(300),
    remote_addr character varying(20),
    created_on timestamp NOT NULL,
    updated_on timestamp NOT NULL,
    expires_on timestamp NOT NULL,
    CONSTRAINT user_session_pkey PRIMARY KEY (user_session_id)
);

COMMIT;