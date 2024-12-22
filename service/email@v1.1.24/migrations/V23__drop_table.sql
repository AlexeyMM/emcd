drop table mailgun_settings;
drop table smtp_settings;

update email_templates
set template = ''
where template is null;

alter table email_templates
    alter column template set not null;

update email_templates
set subject = ''
where subject is null;

alter table email_templates
    alter column subject set not null;


update email_templates
set footer = ''
where footer is null;

alter table email_templates
    alter column footer set not null;
