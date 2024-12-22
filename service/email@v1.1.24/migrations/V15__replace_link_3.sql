UPDATE email_templates
SET "template" = REPLACE("template", '{{.Token}}{{.Token}}', '{{.Token}}')
WHERE "template" LIKE '%{{.Token}}{{.Token}}%';
