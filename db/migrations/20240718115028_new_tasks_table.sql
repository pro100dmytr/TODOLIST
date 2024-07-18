-- +goose Up
-- +goose StatementBegin
DO $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'fk_category') THEN
            ALTER TABLE tasks
                ADD CONSTRAINT fk_category
                    FOREIGN KEY (category_id)
                        REFERENCES categories(id)
                        ON DELETE CASCADE;
        END IF;
    END $$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tasks
    DROP CONSTRAINT IF EXISTS fk_category;
-- +goose StatementEnd