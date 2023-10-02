-- migrate:up
CREATE TABLE brands(
    id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
);

-- migrate:down
DROP TABLE brands;