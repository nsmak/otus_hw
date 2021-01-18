-- +goose Up
CREATE TABLE IF NOT EXISTS  event (
    id varchar(36) NOT NULL,
    title varchar(100) NOT NULL DEFAULT '',
    start_date integer NOT NULL,
    end_date integer NOT NULL,
    description text NOT NULL DEFAULT '',
    owner_id varchar(36) NOT NULL DEFAULT '',
    remind_in integer NOT NULL DEFAULT 0,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS  notification (
    id varchar(36) NOT NULL,
    title varchar(100) NOT NULL DEFAULT '',
    start_date integer NOT NULL,
    PRIMARY KEY (id)
);

-- +goose Down
drop table event;
drop table notification;