
CREATE TABLE domain (
    domain_id id SERIAL NOT NULL,
    object_id TEXT,
    name TEXT,
    description TEXT,
    is_enabled BOOL,
    created_on DATETIME,
    updated_on DATETIME,
    PRIMARY KEY (domain_id)
);
CREATE UNIQUE INDEX domain_object_id ON domain(object_id);
CREATE UNIQUE INDEX domain_name ON domain(name);

CREATE TABLE user (
    user_id INTEGER PRIMARY KEY AUTOINCREMENT,
    object_id TEXT,
    name TEXT,
    passwd TEXT,
    is_enabled BOOL,
    created_on DATETIME,
    updated_on DATETIME
);
CREATE UNIQUE INDEX user_object_id ON user(object_id);
CREATE UNIQUE INDEX user_name ON user(name);

CREATE TABLE domain_user (
    domain_id INTEGER,
    user_id INTEGER,
    PRIMARY KEY (domain_id, user_id),
    FOREIGN KEY(domain_id) REFERENCES domain(domain_id),
    FOREIGN KEY(user_id) REFERENCES user(user_id)
);

CREATE TABLE role (
    role_id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    description TEXT,
    is_enabled BOOL
);
CREATE UNIQUE INDEX role_name ON role(name);

CREATE TABLE permission (
    permission_id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT,
    description TEXT,
    evaluation_rule TEXT,
    is_enabled BOOL
);
CREATE UNIQUE INDEX permission_name ON permission(name);

CREATE TABLE role_permission (
    role_id INTEGER NOT NULL,
    permission_id INTEGER NO NULL,
    PRIMARY KEY (role_id, permission_id),
    FOREIGN KEY(role_id) REFERENCES role(role_id),
    FOREIGN KEY(permission_id) REFERENCES permission(permission_id)
);

CREATE TABLE user_role (
    user_id INTEGER NOT NULL,
    role_id INTEGER NO NULL,
    PRIMARY KEY (user_id, role_id)
);

CREATE TABLE user_session (
    user_session_id CHARACTER VARYING(255) NOT NULL,
    domain_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    user_agent CHARACTER VARYING(255),
    remote_addr CHARACTER VARYING(20),
    created_on TIMESTAMPTZ NOT NULL,
    updated_on TIMESTAMPTZ NOT NULL,
    expires_on TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (sid)
);

ALTER TABLE user_role ADD FOREIGN KEY (user_role_user) REFERENCES user (user_id);
ALTER TABLE user_role ADD FOREIGN KEY (user_role_role) REFERENCES role (role_id);

ALTER TABLE user_session ADD FOREIGN KEY (user_session_domain) REFERENCES domain (domain_id);
ALTER TABLE user_session ADD FOREIGN KEY (user_session_user) REFERENCES user (user_id);
