ALTER TABLE email_templates RENAME COLUMN restore_password_template TO template;

ALTER TABLE email_templates ADD COLUMN type text;

ALTER TABLE email_templates ADD COLUMN language text;

ALTER TABLE email_templates DROP CONSTRAINT email_templates_pkey;

ALTER TABLE email_templates ADD PRIMARY KEY (whitelabel_id,language,type);


DROP TRIGGER notify_templates ON email_templates;


CREATE OR REPLACE FUNCTION get_new_template()  RETURNS TRIGGER as $$
BEGIN
        PERFORM pg_notify('templates', json_build_object('whitelabel_id',NEW.whitelabel_id,
                'language', NEW.language, 'type', NEW.type)::text);
RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER notify_templates
    AFTER INSERT OR UPDATE ON email_templates
                        FOR EACH ROW
                        EXECUTE PROCEDURE get_new_template();


