-- +goose Up
-- +goose StatementBegin

-- Insert sample exchange rates
INSERT INTO exchange_rates (from_currency, to_currency, rate) VALUES
    -- Crypto to USD
    ('BTC', 'USD', 43500.00),
    ('ETH', 'USD', 2280.50),
    ('USDT', 'USD', 1.00),
    ('BNB', 'USD', 315.75),
    ('SOL', 'USD', 98.30),
    
    -- USD to Crypto
    ('USD', 'BTC', 0.000023),
    ('USD', 'ETH', 0.000438),
    ('USD', 'USDT', 1.00),
    ('USD', 'BNB', 0.003167),
    ('USD', 'SOL', 0.010173),
    
    -- Fiat conversions
    ('USD', 'EUR', 0.92),
    ('USD', 'RUB', 92.50),
    ('USD', 'GBP', 0.79),
    ('EUR', 'USD', 1.09),
    ('RUB', 'USD', 0.0108),
    ('GBP', 'USD', 1.27);

-- Insert sample users
INSERT INTO users (id, email, first_name, last_name, phone) VALUES
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'john.doe@example.com', 'John', 'Doe', '+1234567890'),
    ('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'jane.smith@example.com', 'Jane', 'Smith', '+1234567891'),
    ('c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'bob.johnson@example.com', 'Bob', 'Johnson', '+1234567892');

-- Insert sample accounts
INSERT INTO accounts (user_id, currency, balance) VALUES
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'USD', 10000.00),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'EUR', 5000.00),
    ('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'USD', 15000.00),
    ('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'RUB', 500000.00),
    ('c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'USD', 8000.00),
    ('c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'GBP', 3000.00);

-- Insert sample crypto wallets
INSERT INTO crypto_wallets (user_id, crypto_type, balance, address) VALUES
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'BTC', 0.5, '1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa'),
    ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'ETH', 10.0, '0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb1'),
    ('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'BTC', 1.2, '1BvBMSEYstWetqTFn5Au4m4GFg7xJaNVN2'),
    ('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'USDT', 5000.0, '0x8e23ee67d1332ad560396262c48ffbb01f93d052'),
    ('c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'ETH', 5.5, '0x5a0b54d5dc17e0aadc383d2db43b0a0d3e029c4c'),
    ('c2eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'SOL', 100.0, 'SoLaRiUm7tKK6wd78SXA9B5gvKGzNaZVjHGgCsYFDMv');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM crypto_wallets;
DELETE FROM accounts;
DELETE FROM users;
DELETE FROM exchange_rates;

-- +goose StatementEnd

