-- +goose Up
CREATE TABLE staff (
    id INT PRIMARY KEY IDENTITY(1,1),
    name NVARCHAR(100),
    mobile_no NVARCHAR(50) UNIQUE NOT NULL,
    image NVARCHAR(MAX)
);

-- +goose Down
DROP TABLE staff;
