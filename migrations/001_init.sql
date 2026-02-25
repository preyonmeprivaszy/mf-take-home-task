CREATE TABLE IF NOT EXISTS products (
    sku VARCHAR(50) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    stock INT NOT NULL DEFAULT 0 CHECK (stock >= 0)
);


CREATE TABLE IF NOT EXISTS orders (
    id VARCHAR(100) PRIMARY KEY,
    sku VARCHAR(50) NOT NULL REFERENCES products(sku),
    quantity INT NOT NULL CHECK (quantity >=0),
    note VARCHAR(255),
    is_increment BOOL NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


-- seed
INSERT INTO products (sku, name, stock) VALUES
('M-FSH-IP16PR-CL-011', 'M - FullShock iPhone 16 Pro Clear Transparent MagSafe', 22),
('I-NAN-MI10L5-000', 'I - NanoGlass Mi 10 Lite 5G', 7),
('I-NAN-MI10T5-000', 'I - NanoGlass Mi 10T 5G', 12),
('M-PSH-IP16PR-BK-000', 'M - ProShock iPhone 16 Pro Black', 52),
('I-NAN-0GXA71-000', 'I - NanoGlass Galaxy A71', 41),
('I-FSH-00IPXS-010', 'I - FullShock iPhone X/XS Babyblue Transparent', 52),
('M-PSH-IP15PL-FG-000', 'I - ProShock iPhone 15 Plus Forest Green', 29),
('M-FSH-IP17PR-LV-000', 'M - FullShock iPhone 17 Pro Lavender', 65),
('I-NAN-GXS23P-000', 'I - NanoGlass Galaxy S23 Plus', 59),
('I-FSH-IP14PL-CL-010', 'I - FullShock iPhone 14 Plus Clear Transparent', 13),
('M-FSH-IP17PM-BK-001', 'M - FullShock iPhone 17 Pro Max Black MagSafe', 36),
('I-FSH-IP14PL-SM-010', 'I - FullShock iPhone 14 Plus Smokey Transparent', 10),
('M-FSH-IP17PM-FG-000', 'M - FullShock iPhone 17 Pro Max Forest Green', 41),
('I-NAN-IP12MN-000', 'I - NanoGlass iPhone 12 Mini', 43)
ON CONFLICT (sku) DO NOTHING;
