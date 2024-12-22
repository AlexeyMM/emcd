DROP TRIGGER notify_templates ON email_templates;

CREATE OR REPLACE FUNCTION get_new_templates()  RETURNS TRIGGER as $$
BEGIN
        PERFORM pg_notify('templates', NEW.whitelabel_id::text);
RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER notify_templates
    AFTER INSERT OR UPDATE ON email_templates
                        FOR EACH ROW
                        EXECUTE PROCEDURE get_new_templates();
