CREATE TABLE collection
(
    collection_id STRING(36) DEFAULT (GENERATE_UUID()),
    user_id       STRING(36),
    name          STRING(128),
    description   STRING(1024),
    rank          INT64 NOT NULL,
) PRIMARY KEY(collection_id);

CREATE TABLE todo
(
    todo_id       STRING(36) DEFAULT (GENERATE_UUID()),
    name          STRING(128),
    description   STRING(1024),
    done          BOOL,
    rank          INT64     NOT NULL,
    collection_id STRING(36) NOT NULL,
    due_date      TIMESTAMP,
    created_at    TIMESTAMP NOT NULL OPTIONS (
        allow_commit_timestamp = true
        ),
    CONSTRAINT fk_collection_id FOREIGN KEY (collection_id) REFERENCES collection (collection_id),
) PRIMARY KEY(todo_id);