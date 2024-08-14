-- migrate:up
CREATE TABLE brands(
    id INTEGER
);

-- migrate:down
DROP TABLE brands;