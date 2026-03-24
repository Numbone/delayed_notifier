-- name: GetNotificationById :one
SELECT
    id,
    channel,
    recipient,
    subject,
    body,
    status,
    send_at,
    created_at,
    updated_at,
    sent_at
FROM notifications
WHERE id = $1;

-- name: CreateNotification :one
INSERT INTO notifications (
    channel,
    recipient,
    subject,
    body,
    status,
    send_at
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING
    id,
    channel,
    recipient,
    subject,
    body,
    status,
    send_at,
    created_at,
    updated_at,
    sent_at;

-- name: DeleteNotification :exec
DELETE FROM notifications
WHERE id = $1;
