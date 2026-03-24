CREATE TABLE IF NOT EXISTS notifications(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    channel VARCHAR(20) NOT NULL,
    recipient VARCHAR(255) NOT NULL,
    subject VARCHAR(255) NOT NULL DEFAULT '',
    body TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    send_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    sent_at TIMESTAMPTZ
);