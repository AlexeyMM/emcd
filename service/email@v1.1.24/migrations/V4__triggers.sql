CREATE OR REPLACE FUNCTION get_new_templates()  RETURNS TRIGGER as $$
    BEGIN
        PERFORM pg_notify('templates', row_to_json(NEW)::text);
        RETURN NEW;
    END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER notify_templates
    AFTER INSERT OR UPDATE ON email_templates
    FOR EACH ROW
    EXECUTE PROCEDURE get_new_templates();


CREATE OR REPLACE FUNCTION get_new_smtp()  RETURNS TRIGGER as $$
BEGIN
        PERFORM pg_notify('smtp', row_to_json(NEW)::text);
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER notify_smtp
    AFTER INSERT OR UPDATE ON smtp_settings
                        FOR EACH ROW
                        EXECUTE PROCEDURE get_new_smtp();

