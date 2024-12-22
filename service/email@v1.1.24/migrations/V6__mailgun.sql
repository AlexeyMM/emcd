CREATE TABLE mailgun_settings (
    domain text,
    api_key text,
    from_address text,
    whitelabel_id uuid,
    PRIMARY KEY (whitelabel_id)
);

CREATE OR REPLACE FUNCTION get_new_mailgun() RETURNS TRIGGER as $$
BEGIN
        PERFORM pg_notify('mailgun', row_to_json(NEW)::text);
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER notify_mailgun
    AFTER INSERT OR UPDATE ON mailgun_settings
                    FOR EACH ROW
                    EXECUTE PROCEDURE get_new_mailgun();