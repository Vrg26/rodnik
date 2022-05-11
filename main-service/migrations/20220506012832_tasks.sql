-- +goose Up
CREATE TABLE statuses (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    createdOn timestamp DEFAULT now(),
    name text,
    description text
);
CREATE TABLE tasks(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    title  text not null,
    description text,
    status_id uuid not null REFERENCES statuses ("id"),
    creator_id uuid not null REFERENCES users ("id"),
    helper_id uuid REFERENCES users ("id"),
    cost decimal ,
    date_relevance timestamp ,
    created_on timestamp NOT NULL DEFAULT now(),
    overdue boolean
);

-- DATA --
INSERT INTO statuses ("id","name", "description")
VALUES ('11b812e7-f6c0-4007-b161-b28ca41e5d13','Открыто','Для задачи еще нет исполнителя'),
   ('3911ab7b-2bbf-4ec2-9de0-f651ec002692','В работе','Для задачи найден исполнитель'),
   ('9d3ca853-8e38-4e41-90b4-21c46314734d','Закрыто','Задача выпонена');
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE tasks;
DROP TABLE statuses;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
