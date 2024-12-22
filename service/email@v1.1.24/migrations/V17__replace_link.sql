UPDATE email_templates
SET "template" = REPLACE("template", 'https://{{.Domain}}/pool/dashboard/verifyEmail/{{.Token}}', 'https://{{.Domain}}/auth/verifyEmail/{{.Token}}')
WHERE "template" LIKE '%https://{{.Domain}}/pool/dashboard/verifyEmail/{{.Token}}%';