CREATE OR REPLACE FUNCTION set_commands_created_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.created_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_set_commands_created_at BEFORE
INSERT
    ON commands FOR EACH ROW
EXECUTE FUNCTION set_commands_created_at ();