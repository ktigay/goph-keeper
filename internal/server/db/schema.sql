DO ' BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = ''user_data_type'') THEN
        CREATE TYPE user_data_type AS ENUM (''TEXT'', ''BINARY'', ''CARD'');
    END IF;
END ';

CREATE TABLE IF NOT EXISTS "user"
(
    "uuid" UUID DEFAULT gen_random_uuid(),
    "login"      VARCHAR(255) NOT NULL,
    "password"   VARCHAR(60) NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY ("uuid"),
    CONSTRAINT login_idx UNIQUE ("login")
);

CREATE TABLE IF NOT EXISTS "user_data"
(
    "uuid" UUID DEFAULT gen_random_uuid(),
    "user_uuid" UUID,
    "title" TEXT NOT NULL,
    "type" user_data_type NOT NULL,
    "data" BYTEA,
    "metadata" JSONB,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY ("uuid")
);

CREATE INDEX IF NOT EXISTS "user_uuid_idx" ON "user_data" ("user_uuid");