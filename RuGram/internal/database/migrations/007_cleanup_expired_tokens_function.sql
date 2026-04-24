-- Create function to cleanup expired tokens
CREATE OR REPLACE FUNCTION cleanup_expired_tokens()
RETURNS void AS $$
BEGIN
    DELETE FROM user_tokens 
    WHERE expires_at < NOW() 
    AND revoked = true;
END;
$$ LANGUAGE plpgsql;

-- Create scheduled job for token cleanup (PostgreSQL 14+)
-- Note: This requires the pg_cron extension
-- CREATE EXTENSION IF NOT EXISTS pg_cron;
-- 
-- SELECT cron.schedule(
--     'cleanup-expired-tokens',  -- job name
--     '0 0 * * *',              -- daily at midnight
--     'SELECT cleanup_expired_tokens();'
-- );