-- +goose Up
-- +goose StatementBegin
CREATE TYPE taxi_type AS ENUM ('economy', 'comfort', 'business', 'electro');

CREATE TABLE IF NOT EXISTS orders
(
    id SERIAL PRIMARY KEY,
    user_uuid uuid NOT NULL,
    driver_uuid uuid NOT NULL,
    origin VARCHAR(255) NOT NULL,
    destination VARCHAR(255) NOT NULL,
    taxi_type taxi_type NOT NULL,
    is_finished BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE OR REPLACE FUNCTION trigger_set_timestamp()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_timestamp
    AFTER UPDATE ON orders
    FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
DROP TYPE taxi_type;
DROP FUNCTION IF EXISTS trigger_set_timestamp();
DROP TRIGGER IF EXISTS set_timestamp on orders;
-- +goose StatementEnd
