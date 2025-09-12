DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'programs' AND column_name = 'order_number'
    ) THEN
        ALTER TABLE programs
        RENAME COLUMN order_number TO idx;
    END IF;
END $$;