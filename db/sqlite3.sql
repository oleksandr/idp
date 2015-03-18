BEGIN TRANSACTION;

CREATE TABLE domain (
    domain_id INTEGER PRIMARY KEY AUTOINCREMENT,
    object_id TEXT,
    name TEXT,
    description TEXT,
    is_enabled BOOL,
    created_on DATETIME,
    updated_on DATETIME
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
    FOREIGN KEY(user_id) REFERENCES user(domain_id)
);

/*
CREATE TABLE audit_log (
    event_id INTEGER PRIMARY KEY AUTOINCREMENT,
    subject_type TEXT,
    subject_id TEXT,
    action TEXT,
    object_type TEXT,
    object_id TEXT,
    created_on DATETIME
);
*/

COMMIT;