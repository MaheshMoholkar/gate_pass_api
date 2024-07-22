-- +goose Up
CREATE TABLE staff_entry_form (
    id INT PRIMARY KEY IDENTITY(1,1),
    name NVARCHAR(100) NOT NULL,
    purpose NVARCHAR(255),
    [in] DATETIME,
    [out] DATETIME,
    mobile_no NVARCHAR(50),
    image NVARCHAR(MAX)
);

-- +goose Down
DROP TABLE staff_entry_form;