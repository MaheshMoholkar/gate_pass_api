-- +goose Up
CREATE TABLE visitor_form (
    id INT PRIMARY KEY IDENTITY(1,1),
    name NVARCHAR(100) NOT NULL,
    purpose NVARCHAR(255) NOT NULL,
    date DATETIME NOT NULL,
    address NVARCHAR(255),
    vehicle_no INT,
    mobile_no NVARCHAR(50) NOT NULL,
    image NVARCHAR(MAX),
    appointment NVARCHAR(255),
    [in] DATETIME NOT NULL,
    [out] DATETIME
);

-- +goose Down
DROP TABLE visitor_form;
