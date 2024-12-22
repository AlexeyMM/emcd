-- stage и prod
UPDATE email_templates
SET "template" = REPLACE("template", 'https://emcd.io/auth/restorePasswd/', 'https://{{.Domain}}/auth/restorePasswd/')
WHERE "template" LIKE '%https://emcd.io/auth/restorePasswd/%';

UPDATE email_templates
SET "template" = REPLACE("template", 'https://emcd.io/auth/verifyEmail/', 'https://{{.Domain}}/auth/verifyEmail/')
WHERE "template" LIKE '%https://emcd.io/auth/verifyEmail/%';

UPDATE email_templates
SET "template" = REPLACE("template", 'https://emcd.io/pool/dashboard/verifyEmail/', 'https://{{.Domain}}/pool/dashboard/verifyEmail/{{.Token}}')
WHERE "template" LIKE '%https://emcd.io/pool/dashboard/verifyEmail/%';

UPDATE email_templates
SET "template" = REPLACE("template", 'https://emcd.io/profile/confirm/google?token=', 'https://{{.Domain}}/profile/confirm/google?token=')
WHERE "template" LIKE '%https://emcd.io/profile/confirm/google?token=%';

UPDATE email_templates
SET "template" = REPLACE("template", 'https://emcd.io/profile/confirm/phone?token=', 'https://{{.Domain}}/profile/confirm/phone?token=')
WHERE "template" LIKE '%https://emcd.io/profile/confirm/phone?token=%';

UPDATE email_templates
SET "template" = REPLACE("template", 'https://emcd.io/profile/confirm/phone/delete?token=', 'https://{{.Domain}}/profile/confirm/phone/delete?token=')
WHERE "template" LIKE '%https://emcd.io/profile/confirm/phone/delete?token=%';

-- dev
UPDATE email_templates
SET "template" = REPLACE("template", 'https://eco.dev.emcd.io/auth/restorePasswd/', 'https://{{.Domain}}/auth/restorePasswd/')
WHERE "template" LIKE '%https://eco.dev.emcd.io/auth/restorePasswd/%';

UPDATE email_templates
SET "template" = REPLACE("template", 'https://eco.dev.emcd.io/auth/verifyEmail/', 'https://{{.Domain}}/auth/verifyEmail/')
WHERE "template" LIKE '%https://eco.dev.emcd.io/auth/verifyEmail/%';

UPDATE email_templates
SET "template" = REPLACE("template", 'https://eco.dev.emcd.io/pool/dashboard/verifyEmail/', 'https://{{.Domain}}/pool/dashboard/verifyEmail/{{.Token}}')
WHERE "template" LIKE '%https://eco.dev.emcd.io/pool/dashboard/verifyEmail/%';

UPDATE email_templates
SET "template" = REPLACE("template", 'https://eco.dev.emcd.io/profile/confirm/google?token=', 'https://{{.Domain}}/profile/confirm/google?token=')
WHERE "template" LIKE '%https://eco.dev.emcd.io/profile/confirm/google?token=%';

UPDATE email_templates
SET "template" = REPLACE("template", 'https://eco.dev.emcd.io/profile/confirm/phone?token=', 'https://{{.Domain}}/profile/confirm/phone?token=')
WHERE "template" LIKE '%https://eco.dev.emcd.io/profile/confirm/phone?token=%';

UPDATE email_templates
SET "template" = REPLACE("template", 'https://eco.dev.emcd.io/profile/confirm/phone/delete?token=', 'https://{{.Domain}}/profile/confirm/phone/delete?token=')
WHERE "template" LIKE '%https://eco.dev.emcd.io/profile/confirm/phone/delete?token=%';
