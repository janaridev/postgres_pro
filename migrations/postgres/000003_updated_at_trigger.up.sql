CREATE OR REPLACE FUNCTION update_commands_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_commands_updated_at BEFORE
UPDATE ON commands FOR EACH ROW
EXECUTE FUNCTION update_commands_updated_at ();