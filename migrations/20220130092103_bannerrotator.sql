-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS slot (
    slot_id SERIAL PRIMARY KEY,
    description text NOT NULL
);

CREATE TABLE IF NOT EXISTS banner_to_slot (
    banner_to_slot_id PRIMARY KEY
    banner_id integer REFERENCES banner (banner_id),
    slot_id integer REFERENCES slot (slot_id),
    UNIQUE (banner_id, slot_id)
);

CREATE TABLE IF NOT EXISTS banner (
    banner_id SERIAL PRIMARY KEY,
    description text NOT NULL
);

CREATE TABLE IF NOT EXISTS social_group (
    social_group_id SERIAL PRIMARY KEY,
    description text NOT NULL
);


CREATE TABLE IF NOT EXISTS stat (
    stat_id SERIAL PRIMARY KEY,
    banner_id integer REFERENCES banner (banner_id) NOT NULL,
    slot_id integer REFERENCES slot (slot_id) NOT NULL,
    social_group_id integer REFERENCES social_group (social_group_id) NOT NULL,
    display integer DEFAULT 1,
    click integer DEFAULT 0,
    UNIQUE (banner_id, slot_id, social_group_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS stat
DROP TABLE IF EXISTS banner_to_slot
DROP TABLE IF EXISTS social_group
DROP TABLE IF EXISTS banner
DROP TABLE IF EXISTS slot
-- +goose StatementEnd