CREATE TABLE provider_settings
(
    whitelabel_id uuid      NOT NULL,
    providers     json      NOT NULL,
    created_at    timestamp NOT NULL,
    updated_at    timestamp NOT NULL,
    PRIMARY KEY (whitelabel_id)
);


with wl as (
    select whitelabel_id
    from mailgun_settings
    union
    select whitelabel_id
    from smtp_settings
),

     providers as (select wl.whitelabel_id,
                          (case
                               when smtp.whitelabel_id is not null
                                   then '{
            "name":"smtp",
            "setting":{
                "username":"' || smtp.username || '",
                "password":"' || smtp.password || '",
                "server_address":"' || smtp.server_address || '",
                "server_port": ' || smtp.server_port || ',
                "from_address":"' || smtp.from_address || '",
                "from_address_displayed_as":"' || smtp.from_address_displayed_as || '"
            }}'
                               else null end) ::json smtp_provider,
                           (case
                                when mailgun.whitelabel_id is not null
                                    then '{
            "name":"mailgun",
            "setting":{
                "domain":"' || mailgun.domain || '",
                "api_key":"' || mailgun.api_key || '",
                "from_address":"' || mailgun.from_address || '",
                "api_base":"' || mailgun.api_base || '",
                "from_address_displayed_as":"' || mailgun.from_address_displayed_as || '"
            }}'
                                else null end) ::json mailgun_provider
                   from wl
                            left join mailgun_settings mailgun on wl.whitelabel_id = mailgun.whitelabel_id
                            left join smtp_settings smtp on wl.whitelabel_id = smtp.whitelabel_id)
insert into provider_settings(whitelabel_id, providers, created_at, updated_at)
select providers.whitelabel_id,
       ('[' || case when providers.smtp_provider is not null
                        then providers.smtp_provider::text  else '' end
           || case when providers.smtp_provider is not null and providers.mailgun_provider is not null
                       then ',' else '' end
           || case when providers.mailgun_provider is not null
                       then providers.mailgun_provider::text else '' end || ']')::json as provides,
        CURRENT_TIMESTAMP,
       CURRENT_TIMESTAMP
from providers;
