-- +goose Up
CREATE TABLE users (
    id INT PRIMARY KEY IDENTITY(1,1),
    mobile_no NVARCHAR(50),
    password_hash NVARCHAR(255)
);

-- +goose Down
DROP TABLE users;
