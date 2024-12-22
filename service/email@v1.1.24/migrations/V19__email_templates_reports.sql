--- worker report templates
INSERT INTO email_templates (whitelabel_id, template, type, language, subject, footer)
SELECT '00000000-0000-0000-0000-000000000000',
       '<!DOCTYPE html>
      <html lang="ru" xmlns:v="urn:schemas-microsoft-com:vml">
      <head>
        <meta charset="utf-8">
        <meta name="x-apple-disable-message-reformatting">
        <meta name="viewport" content="width=device-width, initial-scale=1">
        <meta name="format-detection" content="telephone=no, date=no, address=no, email=no, url=no">
        <meta name="color-scheme" content="light dark">
        <meta name="supported-color-schemes" content="light dark">
        <!--[if mso]>
        <noscript>
          <xml>
            <o:OfficeDocumentSettings xmlns:o="urn:schemas-microsoft-com:office:office">
              <o:PixelsPerInch>96</o:PixelsPerInch>
            </o:OfficeDocumentSettings>
          </xml>
        </noscript>
        <style>
          td,th,div,p,a,h1,h2,h3,h4,h5,h6 {font-family: "RoobertPRO", Helvetica, Arial, "Segoe UI", sans-serif; mso-line-height-rule: exactly;}
        </style>
        <![endif]-->
        <title>Добро пожаловать в EMCD</title>
        <style>
          .hover-no-underline:hover {
            text-decoration-line: none !important
          }
          @media (max-width: 640px) {
            .sm-block {
              display: block !important
            }
            .sm-table-row {
              display: table-row !important
            }
            .sm-hidden {
              display: none !important
            }
            .sm-min-h-130px {
              min-height: 130px !important
            }
            .sm-w-full {
              width: 100% !important
            }
            .sm-p-0 {
              padding: 0 !important
            }
            .sm-text-2xl {
              font-size: 24px !important
            }
            .sm-text-lg {
              font-size: 18px !important
            }
          }
          @media (max-width: 450px) {
            .xs-mt-1 {
              margin-top: 4px !important
            }
            .xs-table-row {
              display: table-row !important
            }
            .xs-w-full {
              width: 100% !important
            }
            .xs-px-4 {
              padding-left: 16px !important;
              padding-right: 16px !important
            }
            .xs-pl-6 {
              padding-left: 24px !important
            }
            .xs-pr-6 {
              padding-right: 24px !important
            }
            .xs-text-center {
              text-align: center !important
            }
          }
          @media (max-width: 375px) {
            .xxs-hidden {
              display: none !important
            }
          }
          @media (prefers-color-scheme: dark) {
            .dark-mode-bg-light-100 {
              background-color: #f5f5f5 !important
            }
            .dark-mode-bg-white {
              background-color: #ffffff !important
            }
            .dark-mode-text-_1E1E1E {
              color: #1E1E1E !important
            }
            .dark-mode-text-black {
              color: #000000 !important
            }
            .dark-mode-text-light-600 {
              color: #7e7e7e !important
            }
          }
        </style>
      </head>
      <body style="margin: 0; width: 100%; background-color: #f5f5f5; padding: 0; font-size: 14px; line-height: 20px; -webkit-font-smoothing: antialiased; word-break: break-word">

          <div style="display: none">
            &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847;
          </div>

        <div role="article" aria-roledescription="email" aria-label="Добро пожаловать в EMCD" lang="ru" style="width: 100%;">
          <div role="separator" style="line-height: 20px; mso-line-height-alt: 20px">&zwj;</div>
          <table align="center" class="sm-w-full" cellpadding="0" cellspacing="0" style="margin-left: auto; margin-right: auto; width: 600px; font-family: ''RoobertPRO'', ''Arial'', ''Helvetica'', ui-sans-serif, system-ui, -apple-system, ''Segoe UI'', sans-serif" role="none">

              <tr>
                <td class="sm-p-0" style="padding-left: 8px; padding-right: 8px">
                  <!-- header: с иконками соц.сетей -->
                  <div style="width: 100%; border-top-left-radius: 10px; border-top-right-radius: 10px; background-color: #000000; color: #ffffff">
                    <table cellpadding="0" cellspacing="0" style="margin: 0; width: 100%; border-top-left-radius: 10px; border-top-right-radius: 10px; padding: 0; background: linear-gradient(#000000, #000000); background-color: #000000;" role="none">
                      <tr style="vertical-align: middle">
                        <td class="xs-pl-6" style="width: 100px; padding-top: 24px; padding-bottom: 24px; padding-left: 32px; vertical-align: middle">
                          <a href="https://emcd.io" style="display: inline-block; cursor: pointer">
                            <img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718368264082_logo-white%20(1)_01J0BBNSJXY1NWTMMXGEPDZYW4.png" alt="EMCD" width="100" height="27" border="0" style="max-width: 100%; vertical-align: middle; line-height: 1">
                          </a>
                        </td>
                        <td class="xs-pr-6" style="padding-top: 24px; padding-bottom: 24px; padding-right: 32px; text-align: right; vertical-align: middle">
                          <table cellpadding="0" cellspacing="0" style="float: right; clear: both; height: 24px; width: auto" role="none">
                            <tr style="vertical-align: middle;">
                              <td style="margin: 0; height: 24px; width: 24px; padding: 0 12px 0 0">
                                <a href="https://www.instagram.com/emcd_io/" class="untracked" style="display: inline-block; width: 24px; cursor: pointer;"><img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1714985175903_u_instagram-alt_01HX6HA39WZC86BCPFGXJTHPA7.png" alt="Instagram" border="0" width="24" height="24" style="max-width: 100%; vertical-align: middle; line-height: 1;"></a>
                              </td>
                              <td style="margin: 0; height: 24px; width: 24px; padding: 0 12px 0 0;">
                                <a href="https://t.me/emcd_international" style="display: inline-block; width: 24px; cursor: pointer;"><img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1714985174519_u_telegram_01HX6HA1YKQ9RWE6K7FQXR3BNE.png" alt="Telegram" border="0" width="24" height="24" style="max-width: 100%; vertical-align: middle; line-height: 1;"></a>
                              </td>
                              <td style="margin: 0; height: 24px; width: 24px; padding: 0 12px 0 0;"><a href="https://www.youtube.com/channel/UCEOI3erKte4PIKpBflZ00ZQ" style="display: inline-block; width: 24px; cursor: pointer;"><img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1714985171801_u_youtube_01HX6H9Z9Q7PG6N2P17NVYTHME.png" alt="YouTube" border="0" width="24" height="24" style="max-width: 100%; vertical-align: middle; line-height: 1;"></a></td>

                              <td style="margin: 0; height: 24px; width: 24px; padding: 0 12px 0 0;">
                                <a href="https://www.linkedin.com/company/emcdtech/" style="display: inline-block; width: 24px; cursor: pointer;"><img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1714985169108_u_linkedin_01HX6H9WNHT4423NVPJZEGMQS4.png" alt="LinkedIn" border="0" width="24" height="24" style="max-width: 100%; vertical-align: middle; line-height: 1;"></a>
                              </td>
                              <td style="margin: 0; height: 24px; width: 24px; padding: 0;">
                                <a href="https://twitter.com/emcd_io" style="display: inline-block; width: 24px; cursor: pointer;"><img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1714985167624_Group%201410093531_01HX6H9V8TJRR90TW90PF32D5X.png" alt="X (Twitter)" border="0" width="24" height="24" style="max-width: 100%; vertical-align: middle; line-height: 1;"></a>
                              </td>
                            </tr>
                          </table>
                        </td>
                      </tr>
                    </table>
                  </div>
                  <!-- /end header: с иконками соц.сетей -->
                </td>
              </tr>

            <tr>
              <td class="sm-p-0" style="padding-left: 8px; padding-right: 8px">
                <div style="background-color: #ffffff; font-family: ''RoobertPRO'', ''Arial'', ''Helvetica'', ui-sans-serif, system-ui, -apple-system, ''Segoe UI'', sans-serif; font-size: 16px; line-height: 20px"> <!-- black bg intro -->

                  <div class="dark-mode-bg-white xs-px-4" style="background-color: #ffffff; padding: 32px; font-family: ''RoobertPRO'', ''Arial'', ''Helvetica'', ui-sans-serif, system-ui, -apple-system, ''Segoe UI'', sans-serif; font-size: 16px; line-height: 20px">
                    <!--  spacer -->
                    <div role="separator" style="line-height: 10px; mso-line-height-alt: 56px">&zwj;</div>
                    <!-- /end spacer -->
                    <!-- title -->
                    <div class="dark-mode-text-black" style="text-align: left; font-size: 20px; font-weight: 700; line-height: 1.15; color: #000000; mso-line-height-alt: 36px">
                      Hi there 👋
                    </div>
                    <!-- /end title -->
                    <!-- spacer -->
                    <div role="separator" style="line-height: 24px; mso-line-height-alt: 24px">&zwj;</div>
                    <!-- /end spacer -->

                  <!-- text: слева -->
                  <p style="margin: 0; text-align: left; line-height: 1.4; color: #4d4d4d">

      Your request to receive a CSV file containing the status of your workers has been successfully processed. You can download the file using the link below. </p>
      <div role="separator" style="line-height: 24px; mso-line-height-alt: 24px">&zwj;</div>
      <div class="xs-text-center" style="text-align: center">
                      <a href="{{.ReportLink}}" class="dark-mode-text-white dark-mode-bg-violet-300" rel="noopener noreferrer" style="box-sizing: border-box; display: inline-block; height: 56px; min-width: 200px; cursor: pointer; border-radius: 10px; background-color: #8f42ff; padding: 16px 24px; text-align: center; vertical-align: middle; font-family: ''RoobertPRO'', ''Arial'', ''Helvetica'', ui-sans-serif, system-ui, -apple-system, ''Segoe UI'', sans-serif; font-size: 16px; font-weight: 700; line-height: 1.3; color: #ffffff; text-decoration-line: none">
                        <!--[if mso]><i style="mso-font-width: -100%; letter-spacing: 32px; mso-text-raise: 30px;" hidden>&nbsp;</i><![endif]-->
                        <span style="vertical-align: middle; color: inherit; mso-text-raise: 16px;">Download</span>
                        <!--[if mso]><i style="mso-font-width: -100%; letter-spacing: 32px;" hidden>&nbsp;</i><![endif]-->
                      </a>
                    </div>

                  <!-- /end text: слева -->
      <div role="separator" style="background-color: #EBEBEB; height: 1px; line-height: 1px; margin: 32px 0">&zwj;</div>
      <!-- text: слева -->

                  <!-- /end text: слева -->




                    <!-- spacer -->





                    <!-- /end spacer -->


                </div>
              </td>
            </tr>
            <tr>
              <td class="sm-p-0" style="padding-left: 8px; padding-right: 8px">
                <!-- footer: на светлом фоне, краткий -->
                <div class="dark-mode-bg-white dark-mode-text-light-600 xs-px-4" style="background-color: #ffffff; padding-left: 32px; padding-right: 32px; color: #7e7e7e">
                  <div class="dark-mode-bg-light-100" style="border-radius: 10px; background-color: #f5f5f5; padding: 24px 24px 8px; text-align: center">
      <img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718356995357_logo-black_01J0B0XWJZC50MPY0JKHMPDS8Z.png" alt="EMCD" width="57" height="18" border="0" style="max-width: 100%; vertical-align: middle; line-height: 1">
                    <div role="separator" style="line-height: 8px; mso-line-height-alt: 8px">&zwj;</div>
                    <div class="dark-mode-text-black sm-text-2xl" style="text-align: center; font-size: 30px; font-weight: 700; line-height: 1.15; color: #000000; mso-line-height-alt: 36px"> A crypto fintech platform where everything is simple
                    </div>
                    <div role="separator" style="line-height: 8px; mso-line-height-alt: 8px">&zwj;</div>
                    <!-- button: две кнопки -->
                    <div style="width: 100%;">
                      <table cellpadding="0" cellspacing="16" style="margin-left: auto; margin-right: auto; width: auto; vertical-align: middle;" role="none">
                        <tr style="vertical-align: middle;">
                          <td class="sm-table-row">
                            <!-- button: черная с картинкой -->
                            <div style="text-align: center;">
                              <a href="https://emcd.onelink.me/FCtc/x4ojb23m" rel="noopener noreferrer" style="box-sizing: border-box; display: inline-block; cursor: pointer; text-decoration-line: none;"> <img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718624674835_store-android-dark_01J0K06SZC7G7BFWYM0EN64E14.png" alt="Google Play" border="0" height="56" width="200" style="max-width: 100%; vertical-align: middle; line-height: 1; display: block; height: 56px; width: 200px; border-radius: 10px">
                              </a>
                            </div>
                            <!-- /end button: черная с картинкой -->
                          </td>
                          <td class="sm-table-row">
                            <!-- button: черная с картинкой -->
                            <div style="text-align: center;">
                              <a href="https://emcd.onelink.me/FCtc/x4ojb23m	" rel="noopener noreferrer" style="box-sizing: border-box; display: inline-block; cursor: pointer; text-decoration-line: none;"> <img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718624678280_store-ios-dark_01J0K06X34Q7DS84H8YMTS7VRB.png" alt="Apple Store" border="0" height="56" width="200" style="max-width: 100%; vertical-align: middle; line-height: 1; display: block; height: 56px; width: 200px; border-radius: 10px;">
                              </a>
                            </div>
                            <!-- /end button: черная с картинкой -->
                          </td>
                        </tr>
                      </table>
                    </div>
                    <!-- /end button: две кнопки -->
                  </div>
                  <div role="separator" style="line-height: 10px; mso-line-height-alt: 32px">&zwj;</div>
        <div> <a href="https://www.trustpilot.com/review/emcd.io" style="display: block; text-decoration-line: none;">
                      <div style="width: 100%;">
                        <img class="sm-hidden" src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1724838104440_%F0%9F%9A%A3%20Trustpilot%20web_01J6C5SMHWHKHNFF1YYPZSJMRR.png" alt style="max-width: 100%; vertical-align: middle; line-height: 1; height: auto; width: 100%">
                        <img class="sm-block" src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1724838093974_%F0%9F%9A%A3%20Trustpilot%20mob_01J6C5SJ6GX8X6HN5QQC5CTBGB.png" alt style="max-width: 100%; vertical-align: middle; line-height: 1; display: none; height: auto; width: 100%">
                      </div>
          </a>

                  </div>
      <div role="separator" style="line-height: 36px; mso-line-height-alt: 32px">&zwj;</div>
                  <div style="text-align: center;">
                    <p class="dark-mode-text-black" style="margin: 0 0 8px; font-size: 20px; font-weight: 700; color: #000000">Do you have any questions?</p>
                    <p style="margin: 0; font-size: 12px">Email us at   <a href="mailto:support@emcd.io" class="hover-no-underline dark-mode-text-light-600" style="display: inline-block; cursor: pointer; color: #7e7e7e; text-decoration-line: underline; text-underline-offset: 2px">support@emcd.io</a> </p>
                  </div>
                  <div role="separator" style="line-height: 24px; mso-line-height-alt: 24px">&zwj;</div>
                  <div role="separator" style="background-color: #EBEBEB; height: 1px; line-height: 1px; margin: 8px 0;">&zwj;</div>
                  <div role="separator" style="line-height: 24px; mso-line-height-alt: 24px">&zwj;</div>
                  <div style="text-align: center;">
                    <p style="margin: 0 0 8px; font-size: 12px;">Copyright © 2024 <i>EMCD</i> Tech ltd.</p>

                    <p style="margin: 0; font-size: 12px;">If the email is not displayed correctly, open it in your browser.</p>

                  </div>
                  <div class="xs-text-center">
                    <a href="{% view_in_browser_url %}" class="dark-mode-text-light-600 dark-mode-bg-light-100" rel="noopener noreferrer" style="margin-top: 24px; box-sizing: border-box; display: inline-block; height: 56px; width: 100%; min-width: 200px; cursor: pointer; border-radius: 10px; background-color: #f5f5f5; padding: 16px 24px; text-align: center; vertical-align: middle; font-family: ''RoobertPRO'', ''Arial'', ''Helvetica'', ui-sans-serif, system-ui, -apple-system, ''Segoe UI'', sans-serif; font-size: 13px; font-weight: 400; line-height: 1.7; color: #7e7e7e; text-decoration-line: none">
                      <!--[if mso]><i style="mso-font-width: -100%; letter-spacing: 32px; mso-text-raise: 30px;" hidden>&nbsp;</i><![endif]-->
                      <span style="vertical-align: middle; color: inherit; mso-text-raise: 16px;">
            Open web version of the email</span>
                      <!--[if mso]><i style="mso-font-width: -100%; letter-spacing: 32px;" hidden>&nbsp;</i><![endif]-->
                    </a>
                  </div>

                  <div role="separator" style="line-height: 24px; mso-line-height-alt: 24px">&zwj;</div>
                  <div role="separator" style="background-color: #EBEBEB; height: 1px; line-height: 1px; margin: 8px 0;">&zwj;</div>
                  <div role="separator" style="line-height: 10px; mso-line-height-alt: 10px">&zwj;</div>
                  <table cellpadding="0" cellspacing="0" style="margin-left: auto; margin-right: auto; width: auto; border-style: none;" role="none">
                    <tr style="vertical-align: middle;">
                      <td style="padding-left: 8px; padding-right: 8px; vertical-align: middle;">
                        <table cellpadding="0" cellspacing="0" style="width: auto; border-style: none;" role="none">
                          <tr style="vertical-align: middle;">

                              <td class="xxs-hidden">
                                <a href="https://emcd.io/pool/dashboard?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;"> <img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718356980754_footer-sr-thunder_01J0B0XEAKQ8QVPR1WB7G2P1TT.png" border="0" alt width="16" height="16" style="max-width: 100%; vertical-align: middle; line-height: 1; margin-right: 8px">
                                </a>
                              </td>

                            <td>
                              <a href="https://emcd.io/pool/dashboard?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;">
                                <span class="dark-mode-text-_1E1E1E" style="font-size: 12px; line-height: 12px; color: #1E1E1E">
                      Mining</span>
                              </a>
                            </td>
                          </tr>
                        </table>
                      </td>
                      <td style="padding-left: 8px; padding-right: 8px; vertical-align: middle;">
                        <table cellpadding="0" cellspacing="0" style="width: auto; border-style: none;" role="none">
                          <tr style="vertical-align: middle;">

                              <td class="xxs-hidden">
                                <a href="https://emcd.io/wallets?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;"> <img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718356985346_footer-sr-twoarrow_01J0B0XJTEMH8772NQB74BS2AY.png" border="0" alt width="16" height="16" style="max-width: 100%; vertical-align: middle; line-height: 1; margin-right: 8px;">
                                </a>
                              </td>

                            <td>
                              <a href="https://emcd.io/wallets?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;">
                                <span class="dark-mode-text-_1E1E1E" style="font-size: 12px; line-height: 12px; color: #1E1E1E">
                      Wallet</span>
                              </a>
                            </td>
                          </tr>
                        </table>
                      </td>
                      <td style="padding-left: 8px; padding-right: 8px; vertical-align: middle;">
                        <table cellpadding="0" cellspacing="0" style="width: auto; border-style: none;" role="none">
                          <tr style="vertical-align: middle;">

                              <td class="xxs-hidden">
                                <a href="https://emcd.io/deposits?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;"> <img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718356974038_footer-sr-arrow-down_01J0B0X7SJVMRNQ6QAZ37TXW0T.png" border="0" alt width="16" height="16" style="max-width: 100%; vertical-align: middle; line-height: 1; margin-right: 8px;">
                                </a>
                              </td>

                            <td>
                              <a href="https://emcd.io/deposits?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;">
                                <span class="dark-mode-text-_1E1E1E" style="font-size: 12px; line-height: 12px; color: #1E1E1E">Coinhold</span>
                              </a>
                            </td>
                          </tr>
                        </table>
                      </td>
                      <td style="padding-left: 8px; padding-right: 8px; vertical-align: middle;">
                        <table cellpadding="0" cellspacing="0" style="width: auto; border-style: none;" role="none">
                          <tr style="vertical-align: middle;">

                              <td class="xxs-hidden">
                                <a href="https://emcd.io/p2p?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;"> <img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718356977464_footer-sr-coin-dollar_01J0B0XB3WCNB24H8YB39BZFR5.png" border="0" alt width="16" height="16" style="max-width: 100%; vertical-align: middle; line-height: 1; margin-right: 8px;">
                                </a>
                              </td>

                            <td>
                              <a href="https://emcd.io/p2p?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;">
                                <span class="dark-mode-text-_1E1E1E" style="font-size: 12px; line-height: 12px; color: #1E1E1E">P2P</span>
                              </a>
                            </td>
                          </tr>
                        </table>
                      </td>
                    </tr>
                  </table>
                  <div role="separator" style="line-height: 24px; mso-line-height-alt: 24px">&zwj;</div>
                </div>
                <!-- /end footer: на светлом фоне, краткий -->
              </td>
            </tr>
          </table>
          <div role="separator" style="line-height: 40px; mso-line-height-alt: 40px">&zwj;</div>
        </div>
      </body>',
       'worker report',
       'en',
       'Your Worker Status Data Export Request',
       NULL
WHERE NOT EXISTS(
        SELECT type, language FROM email_templates WHERE type = 'worker report' AND language = 'en'
    );

INSERT INTO email_templates (whitelabel_id, template, type, language, subject, footer)
SELECT '00000000-0000-0000-0000-000000000000',
       '<!DOCTYPE html>
            <html lang="ru" xmlns:v="urn:schemas-microsoft-com:vml">
            <head>
              <meta charset="utf-8">
              <meta name="x-apple-disable-message-reformatting">
              <meta name="viewport" content="width=device-width, initial-scale=1">
              <meta name="format-detection" content="telephone=no, date=no, address=no, email=no, url=no">
              <meta name="color-scheme" content="light dark">
              <meta name="supported-color-schemes" content="light dark">
              <!--[if mso]>
              <noscript>
                <xml>
                  <o:OfficeDocumentSettings xmlns:o="urn:schemas-microsoft-com:office:office">
                    <o:PixelsPerInch>96</o:PixelsPerInch>
                  </o:OfficeDocumentSettings>
                </xml>
              </noscript>
              <style>
                td,th,div,p,a,h1,h2,h3,h4,h5,h6 {font-family: "RoobertPRO", Helvetica, Arial, "Segoe UI", sans-serif; mso-line-height-rule: exactly;}
              </style>
              <![endif]-->
              <title>Добро пожаловать в EMCD</title>
              <style>
                .hover-no-underline:hover {
                  text-decoration-line: none !important
                }
                @media (max-width: 640px) {
                  .sm-block {
                    display: block !important
                  }
                  .sm-table-row {
                    display: table-row !important
                  }
                  .sm-hidden {
                    display: none !important
                  }
                  .sm-min-h-130px {
                    min-height: 130px !important
                  }
                  .sm-w-full {
                    width: 100% !important
                  }
                  .sm-p-0 {
                    padding: 0 !important
                  }
                  .sm-text-2xl {
                    font-size: 24px !important
                  }
                  .sm-text-lg {
                    font-size: 18px !important
                  }
                }
                @media (max-width: 450px) {
                  .xs-mt-1 {
                    margin-top: 4px !important
                  }
                  .xs-table-row {
                    display: table-row !important
                  }
                  .xs-w-full {
                    width: 100% !important
                  }
                  .xs-px-4 {
                    padding-left: 16px !important;
                    padding-right: 16px !important
                  }
                  .xs-pl-6 {
                    padding-left: 24px !important
                  }
                  .xs-pr-6 {
                    padding-right: 24px !important
                  }
                  .xs-text-center {
                    text-align: center !important
                  }
                }
                @media (max-width: 375px) {
                  .xxs-hidden {
                    display: none !important
                  }
                }
                @media (prefers-color-scheme: dark) {
                  .dark-mode-bg-light-100 {
                    background-color: #f5f5f5 !important
                  }
                  .dark-mode-bg-white {
                    background-color: #ffffff !important
                  }
                  .dark-mode-text-_1E1E1E {
                    color: #1E1E1E !important
                  }
                  .dark-mode-text-black {
                    color: #000000 !important
                  }
                  .dark-mode-text-light-600 {
                    color: #7e7e7e !important
                  }
                }
              </style>
            </head>
            <body style="margin: 0; width: 100%; background-color: #f5f5f5; padding: 0; font-size: 14px; line-height: 20px; -webkit-font-smoothing: antialiased; word-break: break-word">

                <div style="display: none">
                  &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847;
                </div>

              <div role="article" aria-roledescription="email" aria-label="Добро пожаловать в EMCD" lang="ru" style="width: 100%;">
                <div role="separator" style="line-height: 20px; mso-line-height-alt: 20px">&zwj;</div>
                <table align="center" class="sm-w-full" cellpadding="0" cellspacing="0" style="margin-left: auto; margin-right: auto; width: 600px; font-family: ''RoobertPRO'', ''Arial'', ''Helvetica'', ui-sans-serif, system-ui, -apple-system, ''Segoe UI'', sans-serif" role="none">

                    <tr>
                      <td class="sm-p-0" style="padding-left: 8px; padding-right: 8px">
                        <!-- header: с иконками соц.сетей -->
                        <div style="width: 100%; border-top-left-radius: 10px; border-top-right-radius: 10px; background-color: #000000; color: #ffffff">
                          <table cellpadding="0" cellspacing="0" style="margin: 0; width: 100%; border-top-left-radius: 10px; border-top-right-radius: 10px; padding: 0; background: linear-gradient(#000000, #000000); background-color: #000000;" role="none">
                            <tr style="vertical-align: middle">
                              <td class="xs-pl-6" style="width: 100px; padding-top: 24px; padding-bottom: 24px; padding-left: 32px; vertical-align: middle">
                                <a href="https://emcd.io" style="display: inline-block; cursor: pointer">
                                  <img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718368264082_logo-white%20(1)_01J0BBNSJXY1NWTMMXGEPDZYW4.png" alt="EMCD" width="100" height="27" border="0" style="max-width: 100%; vertical-align: middle; line-height: 1">
                                </a>
                              </td>
                              <td class="xs-pr-6" style="padding-top: 24px; padding-bottom: 24px; padding-right: 32px; text-align: right; vertical-align: middle">
                                <table cellpadding="0" cellspacing="0" style="float: right; clear: both; height: 24px; width: auto" role="none">
                                  <tr style="vertical-align: middle;">
                                    <td style="margin: 0; height: 24px; width: 24px; padding: 0 12px 0 0"><a href="https://instagram.com/emcd_io?link_id={% cio_link_id %}" class="untracked" style="display: inline-block; width: 24px; cursor: pointer;"><img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1714985175903_u_instagram-alt_01HX6HA39WZC86BCPFGXJTHPA7.png" alt="Instagram" border="0" width="24" height="24" style="max-width: 100%; vertical-align: middle; line-height: 1;"></a></td>
                                    <td style="margin: 0; height: 24px; width: 24px; padding: 0 12px 0 0;"><a href="https://t.me/Emcdnews" style="display: inline-block; width: 24px; cursor: pointer;"><img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1714985174519_u_telegram_01HX6HA1YKQ9RWE6K7FQXR3BNE.png" alt="Telegram" border="0" width="24" height="24" style="max-width: 100%; vertical-align: middle; line-height: 1;"></a></td>
                                    <td style="margin: 0; height: 24px; width: 24px; padding: 0 12px 0 0;"><a href="https://vk.com/emcd_io" style="display: inline-block; width: 24px; cursor: pointer;"><img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1714985173166_u_vk_01HX6HA0MA0RB8MPF11T6YXN4N.png" alt="VK" border="0" width="24" height="24" style="max-width: 100%; vertical-align: middle; line-height: 1;"></a></td>
                                    <td style="margin: 0; height: 24px; width: 24px; padding: 0 12px 0 0;"><a href="https://www.youtube.com/channel/UCOjxcJiwe-9F187C1MoS84w" style="display: inline-block; width: 24px; cursor: pointer;"><img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1714985171801_u_youtube_01HX6H9Z9Q7PG6N2P17NVYTHME.png" alt="YouTube" border="0" width="24" height="24" style="max-width: 100%; vertical-align: middle; line-height: 1;"></a></td>
                                    <td style="margin: 0; height: 24px; width: 24px; padding: 0;"><a href="https://vc.ru/u/2383972-emcd-mining-pool" style="display: inline-block; width: 24px; cursor: pointer;"><img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1714985170431_ic-vcru-x3_01HX6H9XYV03RS2KNFPE48ASSM.png" alt="VC.ru" border="0" width="24" height="24" style="max-width: 100%; vertical-align: middle; line-height: 1;"></a></td>
                                  </tr>
                                </table>
                              </td>
                            </tr>
                          </table>
                        </div>
                        <!-- /end header: с иконками соц.сетей -->
                      </td>
                    </tr>

                  <tr>
                    <td class="sm-p-0" style="padding-left: 8px; padding-right: 8px">
                      <div style="background-color: #ffffff; font-family: ''RoobertPRO'', ''Arial'', ''Helvetica'', ui-sans-serif, system-ui, -apple-system, ''Segoe UI'', sans-serif; font-size: 16px; line-height: 20px"> <!-- black bg intro -->

                        <div class="dark-mode-bg-white xs-px-4" style="background-color: #ffffff; padding: 32px; font-family: ''RoobertPRO'', ''Arial'', ''Helvetica'', ui-sans-serif, system-ui, -apple-system, ''Segoe UI'', sans-serif; font-size: 16px; line-height: 20px">
                          <!--  spacer -->
                          <div role="separator" style="line-height: 10px; mso-line-height-alt: 56px">&zwj;</div>
                          <!-- /end spacer -->
                          <!-- title -->
                          <div class="dark-mode-text-black" style="text-align: left; font-size: 20px; font-weight: 700; line-height: 1.15; color: #000000; mso-line-height-alt: 36px">
                            Привет 👋
                          </div>
                          <!-- /end title -->
                          <!-- spacer -->
                          <div role="separator" style="line-height: 24px; mso-line-height-alt: 24px">&zwj;</div>
                          <!-- /end spacer -->

                        <!-- text: слева -->
                        <p style="margin: 0; text-align: left; line-height: 1.4; color: #4d4d4d">

            Твой запрос на получение CSV-таблицы с данными о состоянии твоих воркеров успешно обработан. Ты можешь скачать таблицу, используя ссылку ниже. </p>
            <div role="separator" style="line-height: 24px; mso-line-height-alt: 24px">&zwj;</div>
            <div class="xs-text-center" style="text-align: center">
                            <a href="{{.ReportLink}}" class="dark-mode-text-white dark-mode-bg-violet-300" rel="noopener noreferrer" style="box-sizing: border-box; display: inline-block; height: 56px; min-width: 200px; cursor: pointer; border-radius: 10px; background-color: #8f42ff; padding: 16px 24px; text-align: center; vertical-align: middle; font-family: ''RoobertPRO'', ''Arial'', ''Helvetica'', ui-sans-serif, system-ui, -apple-system, ''Segoe UI'', sans-serif; font-size: 16px; font-weight: 700; line-height: 1.3; color: #ffffff; text-decoration-line: none">
                              <!--[if mso]><i style="mso-font-width: -100%; letter-spacing: 32px; mso-text-raise: 30px;" hidden>&nbsp;</i><![endif]-->
                              <span style="vertical-align: middle; color: inherit; mso-text-raise: 16px;">
            Скачать таблицу</span>
                              <!--[if mso]><i style="mso-font-width: -100%; letter-spacing: 32px;" hidden>&nbsp;</i><![endif]-->
                            </a>
                          </div>
                        <!-- /end text: слева -->
            <div role="separator" style="background-color: #EBEBEB; height: 1px; line-height: 1px; margin: 32px 0">&zwj;</div>
            <!-- text: слева -->

                        <!-- /end text: слева -->




                          <!-- spacer -->





                          <!-- /end spacer -->


                      </div>
                    </td>
                  </tr>
                  <tr>
                    <td class="sm-p-0" style="padding-left: 8px; padding-right: 8px">
                      <!-- footer: на светлом фоне, краткий -->
                      <div class="dark-mode-bg-white dark-mode-text-light-600 xs-px-4" style="background-color: #ffffff; padding-left: 32px; padding-right: 32px; color: #7e7e7e">
                        <div class="dark-mode-bg-light-100" style="border-radius: 10px; background-color: #f5f5f5; padding: 24px 24px 8px; text-align: center">
            <img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718356995357_logo-black_01J0B0XWJZC50MPY0JKHMPDS8Z.png" alt="EMCD" width="57" height="18" border="0" style="max-width: 100%; vertical-align: middle; line-height: 1">
                          <div role="separator" style="line-height: 8px; mso-line-height-alt: 8px">&zwj;</div>
                          <div class="dark-mode-text-black sm-text-2xl" style="text-align: center; font-size: 30px; font-weight: 700; line-height: 1.15; color: #000000; mso-line-height-alt: 36px"> Крипто финтех платформа, где все просто
                          </div>
                          <div role="separator" style="line-height: 8px; mso-line-height-alt: 8px">&zwj;</div>
                          <!-- button: две кнопки -->
                          <div style="width: 100%;">
                            <table cellpadding="0" cellspacing="16" style="margin-left: auto; margin-right: auto; width: auto; vertical-align: middle;" role="none">
                              <tr style="vertical-align: middle;">
                                <td class="sm-table-row">
                                  <!-- button: черная с картинкой -->
                                  <div style="text-align: center;">
                                    <a href="https://emcd.onelink.me/FCtc/x4ojb23m" rel="noopener noreferrer" style="box-sizing: border-box; display: inline-block; cursor: pointer; text-decoration-line: none;"> <img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718624674835_store-android-dark_01J0K06SZC7G7BFWYM0EN64E14.png" alt="Google Play" border="0" height="56" width="200" style="max-width: 100%; vertical-align: middle; line-height: 1; display: block; height: 56px; width: 200px; border-radius: 10px">
                                    </a>
                                  </div>
                                  <!-- /end button: черная с картинкой -->
                                </td>
                                <td class="sm-table-row">
                                  <!-- button: черная с картинкой -->
                                  <div style="text-align: center;">
                                    <a href="https://emcd.onelink.me/FCtc/x4ojb23m	" rel="noopener noreferrer" style="box-sizing: border-box; display: inline-block; cursor: pointer; text-decoration-line: none;"> <img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718624678280_store-ios-dark_01J0K06X34Q7DS84H8YMTS7VRB.png" alt="Apple Store" border="0" height="56" width="200" style="max-width: 100%; vertical-align: middle; line-height: 1; display: block; height: 56px; width: 200px; border-radius: 10px;">
                                    </a>
                                  </div>
                                  <!-- /end button: черная с картинкой -->
                                </td>
                              </tr>
                            </table>
                          </div>
                          <!-- /end button: две кнопки -->
                        </div>
                        <div role="separator" style="line-height: 10px; mso-line-height-alt: 32px">&zwj;</div>
              <div> <a href="https://www.trustpilot.com/review/emcd.io" style="display: block; text-decoration-line: none;">
                            <div style="width: 100%;">
                              <img class="sm-hidden" src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1724838104440_%F0%9F%9A%A3%20Trustpilot%20web_01J6C5SMHWHKHNFF1YYPZSJMRR.png" alt style="max-width: 100%; vertical-align: middle; line-height: 1; height: auto; width: 100%">
                              <img class="sm-block" src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1724838093974_%F0%9F%9A%A3%20Trustpilot%20mob_01J6C5SJ6GX8X6HN5QQC5CTBGB.png" alt style="max-width: 100%; vertical-align: middle; line-height: 1; display: none; height: auto; width: 100%">
                            </div>
                </a>

                        </div>
            <div role="separator" style="line-height: 36px; mso-line-height-alt: 32px">&zwj;</div>
                        <div style="text-align: center;">
                          <p class="dark-mode-text-black" style="margin: 0 0 8px; font-size: 20px; font-weight: 700; color: #000000">Остались вопросы?</p>
                          <p style="margin: 0; font-size: 12px">Напиши нам   <a href="mailto:support@emcd.io" class="hover-no-underline dark-mode-text-light-600" style="display: inline-block; cursor: pointer; color: #7e7e7e; text-decoration-line: underline; text-underline-offset: 2px">support@emcd.io</a> </p>
                        </div>
                        <div role="separator" style="line-height: 24px; mso-line-height-alt: 24px">&zwj;</div>
                        <div role="separator" style="background-color: #EBEBEB; height: 1px; line-height: 1px; margin: 8px 0;">&zwj;</div>
                        <div role="separator" style="line-height: 24px; mso-line-height-alt: 24px">&zwj;</div>
                        <div style="text-align: center;">
                          <p style="margin: 0 0 8px; font-size: 12px;">Copyright © 2024 <i>EMCD</i> Tech ltd.</p>

                      <p style="margin: 0; font-size: 12px;">Если письмо отображается некорректно, открой его в браузере.</p>

                        </div>
                        <div class="xs-text-center">
                          <a href="{% view_in_browser_url %}" class="dark-mode-text-light-600 dark-mode-bg-light-100" rel="noopener noreferrer" style="margin-top: 24px; box-sizing: border-box; display: inline-block; height: 56px; width: 100%; min-width: 200px; cursor: pointer; border-radius: 10px; background-color: #f5f5f5; padding: 16px 24px; text-align: center; vertical-align: middle; font-family: ''RoobertPRO'', ''Arial'', ''Helvetica'', ui-sans-serif, system-ui, -apple-system, ''Segoe UI'', sans-serif; font-size: 13px; font-weight: 400; line-height: 1.7; color: #7e7e7e; text-decoration-line: none">
                            <!--[if mso]><i style="mso-font-width: -100%; letter-spacing: 32px; mso-text-raise: 30px;" hidden>&nbsp;</i><![endif]-->
                            <span style="vertical-align: middle; color: inherit; mso-text-raise: 16px;">
                  Открыть письмо в браузере</span>
                            <!--[if mso]><i style="mso-font-width: -100%; letter-spacing: 32px;" hidden>&nbsp;</i><![endif]-->
                          </a>
                        </div>

                        <div role="separator" style="line-height: 24px; mso-line-height-alt: 24px">&zwj;</div>
                        <div role="separator" style="background-color: #EBEBEB; height: 1px; line-height: 1px; margin: 8px 0;">&zwj;</div>
                        <div role="separator" style="line-height: 10px; mso-line-height-alt: 10px">&zwj;</div>
                        <table cellpadding="0" cellspacing="0" style="margin-left: auto; margin-right: auto; width: auto; border-style: none;" role="none">
                          <tr style="vertical-align: middle;">
                            <td style="padding-left: 8px; padding-right: 8px; vertical-align: middle;">
                              <table cellpadding="0" cellspacing="0" style="width: auto; border-style: none;" role="none">
                                <tr style="vertical-align: middle;">

                                    <td class="xxs-hidden">
                                      <a href="https://emcd.io/pool/dashboard?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;"> <img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718356980754_footer-sr-thunder_01J0B0XEAKQ8QVPR1WB7G2P1TT.png" border="0" alt width="16" height="16" style="max-width: 100%; vertical-align: middle; line-height: 1; margin-right: 8px">
                                      </a>
                                    </td>

                                  <td>
                                    <a href="https://emcd.io/pool/dashboard?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;">
                                      <span class="dark-mode-text-_1E1E1E" style="font-size: 12px; line-height: 12px; color: #1E1E1E">
                            Майнинг</span>
                                    </a>
                                  </td>
                                </tr>
                              </table>
                            </td>
                            <td style="padding-left: 8px; padding-right: 8px; vertical-align: middle;">
                              <table cellpadding="0" cellspacing="0" style="width: auto; border-style: none;" role="none">
                                <tr style="vertical-align: middle;">

                                    <td class="xxs-hidden">
                                      <a href="https://emcd.io/wallets?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;"> <img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718356985346_footer-sr-twoarrow_01J0B0XJTEMH8772NQB74BS2AY.png" border="0" alt width="16" height="16" style="max-width: 100%; vertical-align: middle; line-height: 1; margin-right: 8px;">
                                      </a>
                                    </td>

                                  <td>
                                    <a href="https://emcd.io/wallets?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;">
                                      <span class="dark-mode-text-_1E1E1E" style="font-size: 12px; line-height: 12px; color: #1E1E1E">
                            Кошелек</span>
                                    </a>
                                  </td>
                                </tr>
                              </table>
                            </td>
                            <td style="padding-left: 8px; padding-right: 8px; vertical-align: middle;">
                              <table cellpadding="0" cellspacing="0" style="width: auto; border-style: none;" role="none">
                                <tr style="vertical-align: middle;">

                                    <td class="xxs-hidden">
                                      <a href="https://emcd.io/deposits?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;"> <img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718356974038_footer-sr-arrow-down_01J0B0X7SJVMRNQ6QAZ37TXW0T.png" border="0" alt width="16" height="16" style="max-width: 100%; vertical-align: middle; line-height: 1; margin-right: 8px;">
                                      </a>
                                    </td>

                                  <td>
                                    <a href="https://emcd.io/deposits?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;">
                                      <span class="dark-mode-text-_1E1E1E" style="font-size: 12px; line-height: 12px; color: #1E1E1E">Coinhold</span>
                                    </a>
                                  </td>
                                </tr>
                              </table>
                            </td>
                            <td style="padding-left: 8px; padding-right: 8px; vertical-align: middle;">
                              <table cellpadding="0" cellspacing="0" style="width: auto; border-style: none;" role="none">
                                <tr style="vertical-align: middle;">

                                    <td class="xxs-hidden">
                                      <a href="https://emcd.io/p2p?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;"> <img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718356977464_footer-sr-coin-dollar_01J0B0XB3WCNB24H8YB39BZFR5.png" border="0" alt width="16" height="16" style="max-width: 100%; vertical-align: middle; line-height: 1; margin-right: 8px;">
                                      </a>
                                    </td>

                                  <td>
                                    <a href="https://emcd.io/p2p?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;">
                                      <span class="dark-mode-text-_1E1E1E" style="font-size: 12px; line-height: 12px; color: #1E1E1E">P2P</span>
                                    </a>
                                  </td>
                                </tr>
                              </table>
                            </td>
                          </tr>
                        </table>
                        <div role="separator" style="line-height: 24px; mso-line-height-alt: 24px">&zwj;</div>
                      </div>
                      <!-- /end footer: на светлом фоне, краткий -->
                    </td>
                  </tr>
                </table>
                <div role="separator" style="line-height: 40px; mso-line-height-alt: 40px">&zwj;</div>
              </div>
            </body>
            </html>',
       'worker report',
       'ru',
       'Твой запрос на выгрузку данных о состоянии воркеров',
       NULL
WHERE NOT EXISTS(
        SELECT type, language FROM email_templates WHERE type = 'worker report' AND language = 'ru'
    );

--- income report templates
INSERT INTO email_templates (whitelabel_id, template, type, language, subject, footer)
SELECT '00000000-0000-0000-0000-000000000000',
       '<!DOCTYPE html>
<html lang="ru" xmlns:v="urn:schemas-microsoft-com:vml">
<head>
  <meta charset="utf-8">
  <meta name="x-apple-disable-message-reformatting">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <meta name="format-detection" content="telephone=no, date=no, address=no, email=no, url=no">
  <meta name="color-scheme" content="light dark">
  <meta name="supported-color-schemes" content="light dark">
  <!--[if mso]>
  <noscript>
    <xml>
      <o:OfficeDocumentSettings xmlns:o="urn:schemas-microsoft-com:office:office">
        <o:PixelsPerInch>96</o:PixelsPerInch>
      </o:OfficeDocumentSettings>
    </xml>
  </noscript>
  <style>
    td,th,div,p,a,h1,h2,h3,h4,h5,h6 {font-family: "RoobertPRO", Helvetica, Arial, "Segoe UI", sans-serif; mso-line-height-rule: exactly;}
  </style>
  <![endif]-->
  <title>Добро пожаловать в EMCD</title>
  <style>
    .hover-no-underline:hover {
      text-decoration-line: none !important
    }
    @media (max-width: 640px) {
      .sm-block {
        display: block !important
      }
      .sm-table-row {
        display: table-row !important
      }
      .sm-hidden {
        display: none !important
      }
      .sm-min-h-130px {
        min-height: 130px !important
      }
      .sm-w-full {
        width: 100% !important
      }
      .sm-p-0 {
        padding: 0 !important
      }
      .sm-text-2xl {
        font-size: 24px !important
      }
      .sm-text-lg {
        font-size: 18px !important
      }
    }
    @media (max-width: 450px) {
      .xs-mt-1 {
        margin-top: 4px !important
      }
      .xs-table-row {
        display: table-row !important
      }
      .xs-w-full {
        width: 100% !important
      }
      .xs-px-4 {
        padding-left: 16px !important;
        padding-right: 16px !important
      }
      .xs-pl-6 {
        padding-left: 24px !important
      }
      .xs-pr-6 {
        padding-right: 24px !important
      }
      .xs-text-center {
        text-align: center !important
      }
    }
    @media (max-width: 375px) {
      .xxs-hidden {
        display: none !important
      }
    }
    @media (prefers-color-scheme: dark) {
      .dark-mode-bg-light-100 {
        background-color: #f5f5f5 !important
      }
      .dark-mode-bg-white {
        background-color: #ffffff !important
      }
      .dark-mode-text-_1E1E1E {
        color: #1E1E1E !important
      }
      .dark-mode-text-black {
        color: #000000 !important
      }
      .dark-mode-text-light-600 {
        color: #7e7e7e !important
      }
    }
  </style>
</head>
<body style="margin: 0; width: 100%; background-color: #f5f5f5; padding: 0; font-size: 14px; line-height: 20px; -webkit-font-smoothing: antialiased; word-break: break-word">

    <div style="display: none">
      &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847;
    </div>

  <div role="article" aria-roledescription="email" aria-label="Добро пожаловать в EMCD" lang="ru" style="width: 100%;">
    <div role="separator" style="line-height: 20px; mso-line-height-alt: 20px">&zwj;</div>
    <table align="center" class="sm-w-full" cellpadding="0" cellspacing="0" style="margin-left: auto; margin-right: auto; width: 600px; font-family: ''RoobertPRO'', ''Arial'', ''Helvetica'', ui-sans-serif, system-ui, -apple-system, ''Segoe UI'', sans-serif" role="none">

        <tr>
          <td class="sm-p-0" style="padding-left: 8px; padding-right: 8px">
            <!-- header: с иконками соц.сетей -->
            <div style="width: 100%; border-top-left-radius: 10px; border-top-right-radius: 10px; background-color: #000000; color: #ffffff">
              <table cellpadding="0" cellspacing="0" style="margin: 0; width: 100%; border-top-left-radius: 10px; border-top-right-radius: 10px; padding: 0; background: linear-gradient(#000000, #000000); background-color: #000000;" role="none">
                <tr style="vertical-align: middle">
                  <td class="xs-pl-6" style="width: 100px; padding-top: 24px; padding-bottom: 24px; padding-left: 32px; vertical-align: middle">
                    <a href="https://emcd.io" style="display: inline-block; cursor: pointer">
                      <img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718368264082_logo-white%20(1)_01J0BBNSJXY1NWTMMXGEPDZYW4.png" alt="EMCD" width="100" height="27" border="0" style="max-width: 100%; vertical-align: middle; line-height: 1">
                    </a>
                  </td>
                  <td class="xs-pr-6" style="padding-top: 24px; padding-bottom: 24px; padding-right: 32px; text-align: right; vertical-align: middle">
                    <table cellpadding="0" cellspacing="0" style="float: right; clear: both; height: 24px; width: auto" role="none">
                      <tr style="vertical-align: middle;">
                        <td style="margin: 0; height: 24px; width: 24px; padding: 0 12px 0 0">
                          <a href="https://www.instagram.com/emcd_io/" class="untracked" style="display: inline-block; width: 24px; cursor: pointer;"><img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1714985175903_u_instagram-alt_01HX6HA39WZC86BCPFGXJTHPA7.png" alt="Instagram" border="0" width="24" height="24" style="max-width: 100%; vertical-align: middle; line-height: 1;"></a>
                        </td>
                        <td style="margin: 0; height: 24px; width: 24px; padding: 0 12px 0 0;">
                          <a href="https://t.me/emcd_international" style="display: inline-block; width: 24px; cursor: pointer;"><img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1714985174519_u_telegram_01HX6HA1YKQ9RWE6K7FQXR3BNE.png" alt="Telegram" border="0" width="24" height="24" style="max-width: 100%; vertical-align: middle; line-height: 1;"></a>
                        </td>
                        <td style="margin: 0; height: 24px; width: 24px; padding: 0 12px 0 0;"><a href="https://www.youtube.com/channel/UCEOI3erKte4PIKpBflZ00ZQ" style="display: inline-block; width: 24px; cursor: pointer;"><img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1714985171801_u_youtube_01HX6H9Z9Q7PG6N2P17NVYTHME.png" alt="YouTube" border="0" width="24" height="24" style="max-width: 100%; vertical-align: middle; line-height: 1;"></a></td>

                        <td style="margin: 0; height: 24px; width: 24px; padding: 0 12px 0 0;">
                          <a href="https://www.linkedin.com/company/emcdtech/" style="display: inline-block; width: 24px; cursor: pointer;"><img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1714985169108_u_linkedin_01HX6H9WNHT4423NVPJZEGMQS4.png" alt="LinkedIn" border="0" width="24" height="24" style="max-width: 100%; vertical-align: middle; line-height: 1;"></a>
                        </td>
                        <td style="margin: 0; height: 24px; width: 24px; padding: 0;">
                          <a href="https://twitter.com/emcd_io" style="display: inline-block; width: 24px; cursor: pointer;"><img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1714985167624_Group%201410093531_01HX6H9V8TJRR90TW90PF32D5X.png" alt="X (Twitter)" border="0" width="24" height="24" style="max-width: 100%; vertical-align: middle; line-height: 1;"></a>
                        </td>
                      </tr>
                    </table>
                  </td>
                </tr>
              </table>
            </div>
            <!-- /end header: с иконками соц.сетей -->
          </td>
        </tr>

      <tr>
        <td class="sm-p-0" style="padding-left: 8px; padding-right: 8px">
          <div style="background-color: #ffffff; font-family: ''RoobertPRO'', ''Arial'', ''Helvetica'', ui-sans-serif, system-ui, -apple-system, ''Segoe UI'', sans-serif; font-size: 16px; line-height: 20px"> <!-- black bg intro -->

            <div class="dark-mode-bg-white xs-px-4" style="background-color: #ffffff; padding: 32px; font-family: ''RoobertPRO'', ''Arial'', ''Helvetica'', ui-sans-serif, system-ui, -apple-system, ''Segoe UI'', sans-serif; font-size: 16px; line-height: 20px">
              <!--  spacer -->
              <div role="separator" style="line-height: 10px; mso-line-height-alt: 56px">&zwj;</div>
              <!-- /end spacer -->
              <!-- title -->
              <div class="dark-mode-text-black" style="text-align: left; font-size: 20px; font-weight: 700; line-height: 1.15; color: #000000; mso-line-height-alt: 36px">
                Hi there 👋
              </div>
              <!-- /end title -->
              <!-- spacer -->
              <div role="separator" style="line-height: 24px; mso-line-height-alt: 24px">&zwj;</div>
              <!-- /end spacer -->

            <!-- text: слева -->
            <p style="margin: 0; text-align: left; line-height: 1.4; color: #4d4d4d">

Your request to receive a CSV file containing your earnings data has been successfully processed. You can download the file using the link below.</p>
<div role="separator" style="line-height: 24px; mso-line-height-alt: 24px">&zwj;</div>
<div class="xs-text-center" style="text-align: center">
                <a href="{{.ReportLink}}" class="dark-mode-text-white dark-mode-bg-violet-300" rel="noopener noreferrer" style="box-sizing: border-box; display: inline-block; height: 56px; min-width: 200px; cursor: pointer; border-radius: 10px; background-color: #8f42ff; padding: 16px 24px; text-align: center; vertical-align: middle; font-family: ''RoobertPRO'', ''Arial'', ''Helvetica'', ui-sans-serif, system-ui, -apple-system, ''Segoe UI'', sans-serif; font-size: 16px; font-weight: 700; line-height: 1.3; color: #ffffff; text-decoration-line: none">
                  <!--[if mso]><i style="mso-font-width: -100%; letter-spacing: 32px; mso-text-raise: 30px;" hidden>&nbsp;</i><![endif]-->
                  <span style="vertical-align: middle; color: inherit; mso-text-raise: 16px;">Download</span>
                  <!--[if mso]><i style="mso-font-width: -100%; letter-spacing: 32px;" hidden>&nbsp;</i><![endif]-->
                </a>
              </div>

            <!-- /end text: слева -->
<div role="separator" style="background-color: #EBEBEB; height: 1px; line-height: 1px; margin: 32px 0">&zwj;</div>
<!-- text: слева -->

            <!-- /end text: слева -->




              <!-- spacer -->





              <!-- /end spacer -->


          </div>
        </td>
      </tr>
      <tr>
        <td class="sm-p-0" style="padding-left: 8px; padding-right: 8px">
          <!-- footer: на светлом фоне, краткий -->
          <div class="dark-mode-bg-white dark-mode-text-light-600 xs-px-4" style="background-color: #ffffff; padding-left: 32px; padding-right: 32px; color: #7e7e7e">
            <div class="dark-mode-bg-light-100" style="border-radius: 10px; background-color: #f5f5f5; padding: 24px 24px 8px; text-align: center">
<img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718356995357_logo-black_01J0B0XWJZC50MPY0JKHMPDS8Z.png" alt="EMCD" width="57" height="18" border="0" style="max-width: 100%; vertical-align: middle; line-height: 1">
              <div role="separator" style="line-height: 8px; mso-line-height-alt: 8px">&zwj;</div>
              <div class="dark-mode-text-black sm-text-2xl" style="text-align: center; font-size: 30px; font-weight: 700; line-height: 1.15; color: #000000; mso-line-height-alt: 36px"> A crypto fintech platform where everything is simple
              </div>
              <div role="separator" style="line-height: 8px; mso-line-height-alt: 8px">&zwj;</div>
              <!-- button: две кнопки -->
              <div style="width: 100%;">
                <table cellpadding="0" cellspacing="16" style="margin-left: auto; margin-right: auto; width: auto; vertical-align: middle;" role="none">
                  <tr style="vertical-align: middle;">
                    <td class="sm-table-row">
                      <!-- button: черная с картинкой -->
                      <div style="text-align: center;">
                        <a href="https://emcd.onelink.me/FCtc/x4ojb23m" rel="noopener noreferrer" style="box-sizing: border-box; display: inline-block; cursor: pointer; text-decoration-line: none;"> <img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718624674835_store-android-dark_01J0K06SZC7G7BFWYM0EN64E14.png" alt="Google Play" border="0" height="56" width="200" style="max-width: 100%; vertical-align: middle; line-height: 1; display: block; height: 56px; width: 200px; border-radius: 10px">
                        </a>
                      </div>
                      <!-- /end button: черная с картинкой -->
                    </td>
                    <td class="sm-table-row">
                      <!-- button: черная с картинкой -->
                      <div style="text-align: center;">
                        <a href="https://emcd.onelink.me/FCtc/x4ojb23m	" rel="noopener noreferrer" style="box-sizing: border-box; display: inline-block; cursor: pointer; text-decoration-line: none;"> <img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718624678280_store-ios-dark_01J0K06X34Q7DS84H8YMTS7VRB.png" alt="Apple Store" border="0" height="56" width="200" style="max-width: 100%; vertical-align: middle; line-height: 1; display: block; height: 56px; width: 200px; border-radius: 10px;">
                        </a>
                      </div>
                      <!-- /end button: черная с картинкой -->
                    </td>
                  </tr>
                </table>
              </div>
              <!-- /end button: две кнопки -->
            </div>
            <div role="separator" style="line-height: 10px; mso-line-height-alt: 32px">&zwj;</div>
  <div> <a href="https://www.trustpilot.com/review/emcd.io" style="display: block; text-decoration-line: none;">
                <div style="width: 100%;">
                  <img class="sm-hidden" src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1724838104440_%F0%9F%9A%A3%20Trustpilot%20web_01J6C5SMHWHKHNFF1YYPZSJMRR.png" alt style="max-width: 100%; vertical-align: middle; line-height: 1; height: auto; width: 100%">
                  <img class="sm-block" src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1724838093974_%F0%9F%9A%A3%20Trustpilot%20mob_01J6C5SJ6GX8X6HN5QQC5CTBGB.png" alt style="max-width: 100%; vertical-align: middle; line-height: 1; display: none; height: auto; width: 100%">
                </div>
    </a>

            </div>
<div role="separator" style="line-height: 36px; mso-line-height-alt: 32px">&zwj;</div>
            <div style="text-align: center;">
              <p class="dark-mode-text-black" style="margin: 0 0 8px; font-size: 20px; font-weight: 700; color: #000000">Do you have any questions?</p>
              <p style="margin: 0; font-size: 12px">Email us at   <a href="mailto:support@emcd.io" class="hover-no-underline dark-mode-text-light-600" style="display: inline-block; cursor: pointer; color: #7e7e7e; text-decoration-line: underline; text-underline-offset: 2px">support@emcd.io</a> </p>
            </div>
            <div role="separator" style="line-height: 24px; mso-line-height-alt: 24px">&zwj;</div>
            <div role="separator" style="background-color: #EBEBEB; height: 1px; line-height: 1px; margin: 8px 0;">&zwj;</div>
            <div role="separator" style="line-height: 24px; mso-line-height-alt: 24px">&zwj;</div>
            <div style="text-align: center;">
              <p style="margin: 0 0 8px; font-size: 12px;">Copyright © 2024 <i>EMCD</i> Tech ltd.</p>

              <p style="margin: 0; font-size: 12px;">If the email is not displayed correctly, open it in your browser.</p>

            </div>
            <div class="xs-text-center">
              <a href="{% view_in_browser_url %}" class="dark-mode-text-light-600 dark-mode-bg-light-100" rel="noopener noreferrer" style="margin-top: 24px; box-sizing: border-box; display: inline-block; height: 56px; width: 100%; min-width: 200px; cursor: pointer; border-radius: 10px; background-color: #f5f5f5; padding: 16px 24px; text-align: center; vertical-align: middle; font-family: ''RoobertPRO'', ''Arial'', ''Helvetica'', ui-sans-serif, system-ui, -apple-system, ''Segoe UI'', sans-serif; font-size: 13px; font-weight: 400; line-height: 1.7; color: #7e7e7e; text-decoration-line: none">
                <!--[if mso]><i style="mso-font-width: -100%; letter-spacing: 32px; mso-text-raise: 30px;" hidden>&nbsp;</i><![endif]-->
                <span style="vertical-align: middle; color: inherit; mso-text-raise: 16px;">
      Open web version of the email</span>
                <!--[if mso]><i style="mso-font-width: -100%; letter-spacing: 32px;" hidden>&nbsp;</i><![endif]-->
              </a>
            </div>

            <div role="separator" style="line-height: 24px; mso-line-height-alt: 24px">&zwj;</div>
            <div role="separator" style="background-color: #EBEBEB; height: 1px; line-height: 1px; margin: 8px 0;">&zwj;</div>
            <div role="separator" style="line-height: 10px; mso-line-height-alt: 10px">&zwj;</div>
            <table cellpadding="0" cellspacing="0" style="margin-left: auto; margin-right: auto; width: auto; border-style: none;" role="none">
              <tr style="vertical-align: middle;">
                <td style="padding-left: 8px; padding-right: 8px; vertical-align: middle;">
                  <table cellpadding="0" cellspacing="0" style="width: auto; border-style: none;" role="none">
                    <tr style="vertical-align: middle;">

                        <td class="xxs-hidden">
                          <a href="https://emcd.io/pool/dashboard?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;"> <img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718356980754_footer-sr-thunder_01J0B0XEAKQ8QVPR1WB7G2P1TT.png" border="0" alt width="16" height="16" style="max-width: 100%; vertical-align: middle; line-height: 1; margin-right: 8px">
                          </a>
                        </td>

                      <td>
                        <a href="https://emcd.io/pool/dashboard?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;">
                          <span class="dark-mode-text-_1E1E1E" style="font-size: 12px; line-height: 12px; color: #1E1E1E">
                Mining</span>
                        </a>
                      </td>
                    </tr>
                  </table>
                </td>
                <td style="padding-left: 8px; padding-right: 8px; vertical-align: middle;">
                  <table cellpadding="0" cellspacing="0" style="width: auto; border-style: none;" role="none">
                    <tr style="vertical-align: middle;">

                        <td class="xxs-hidden">
                          <a href="https://emcd.io/wallets?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;"> <img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718356985346_footer-sr-twoarrow_01J0B0XJTEMH8772NQB74BS2AY.png" border="0" alt width="16" height="16" style="max-width: 100%; vertical-align: middle; line-height: 1; margin-right: 8px;">
                          </a>
                        </td>

                      <td>
                        <a href="https://emcd.io/wallets?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;">
                          <span class="dark-mode-text-_1E1E1E" style="font-size: 12px; line-height: 12px; color: #1E1E1E">
                Wallet</span>
                        </a>
                      </td>
                    </tr>
                  </table>
                </td>
                <td style="padding-left: 8px; padding-right: 8px; vertical-align: middle;">
                  <table cellpadding="0" cellspacing="0" style="width: auto; border-style: none;" role="none">
                    <tr style="vertical-align: middle;">

                        <td class="xxs-hidden">
                          <a href="https://emcd.io/deposits?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;"> <img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718356974038_footer-sr-arrow-down_01J0B0X7SJVMRNQ6QAZ37TXW0T.png" border="0" alt width="16" height="16" style="max-width: 100%; vertical-align: middle; line-height: 1; margin-right: 8px;">
                          </a>
                        </td>

                      <td>
                        <a href="https://emcd.io/deposits?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;">
                          <span class="dark-mode-text-_1E1E1E" style="font-size: 12px; line-height: 12px; color: #1E1E1E">Coinhold</span>
                        </a>
                      </td>
                    </tr>
                  </table>
                </td>
                <td style="padding-left: 8px; padding-right: 8px; vertical-align: middle;">
                  <table cellpadding="0" cellspacing="0" style="width: auto; border-style: none;" role="none">
                    <tr style="vertical-align: middle;">

                        <td class="xxs-hidden">
                          <a href="https://emcd.io/p2p?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;"> <img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718356977464_footer-sr-coin-dollar_01J0B0XB3WCNB24H8YB39BZFR5.png" border="0" alt width="16" height="16" style="max-width: 100%; vertical-align: middle; line-height: 1; margin-right: 8px;">
                          </a>
                        </td>

                      <td>
                        <a href="https://emcd.io/p2p?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;">
                          <span class="dark-mode-text-_1E1E1E" style="font-size: 12px; line-height: 12px; color: #1E1E1E">P2P</span>
                        </a>
                      </td>
                    </tr>
                  </table>
                </td>
              </tr>
            </table>
            <div role="separator" style="line-height: 24px; mso-line-height-alt: 24px">&zwj;</div>
          </div>
          <!-- /end footer: на светлом фоне, краткий -->
        </td>
      </tr>
    </table>
    <div role="separator" style="line-height: 40px; mso-line-height-alt: 40px">&zwj;</div>
  </div>
</body>
</html>',
       'income report',
       'en',
       'Your Earnings Data Export Request',
       NULL
WHERE NOT EXISTS(
        SELECT type, language FROM email_templates WHERE type = 'income report' AND language = 'en'
    );

INSERT INTO email_templates (whitelabel_id, template, type, language, subject, footer)
SELECT '00000000-0000-0000-0000-000000000000',
       '<!DOCTYPE html>
<html lang="ru" xmlns:v="urn:schemas-microsoft-com:vml">
<head>
  <meta charset="utf-8">
  <meta name="x-apple-disable-message-reformatting">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <meta name="format-detection" content="telephone=no, date=no, address=no, email=no, url=no">
  <meta name="color-scheme" content="light dark">
  <meta name="supported-color-schemes" content="light dark">
  <!--[if mso]>
  <noscript>
    <xml>
      <o:OfficeDocumentSettings xmlns:o="urn:schemas-microsoft-com:office:office">
        <o:PixelsPerInch>96</o:PixelsPerInch>
      </o:OfficeDocumentSettings>
    </xml>
  </noscript>
  <style>
    td,th,div,p,a,h1,h2,h3,h4,h5,h6 {font-family: "RoobertPRO", Helvetica, Arial, "Segoe UI", sans-serif; mso-line-height-rule: exactly;}
  </style>
  <![endif]-->
  <title>Добро пожаловать в EMCD</title>
  <style>
    .hover-no-underline:hover {
      text-decoration-line: none !important
    }
    @media (max-width: 640px) {
      .sm-block {
        display: block !important
      }
      .sm-table-row {
        display: table-row !important
      }
      .sm-hidden {
        display: none !important
      }
      .sm-min-h-130px {
        min-height: 130px !important
      }
      .sm-w-full {
        width: 100% !important
      }
      .sm-p-0 {
        padding: 0 !important
      }
      .sm-text-2xl {
        font-size: 24px !important
      }
      .sm-text-lg {
        font-size: 18px !important
      }
    }
    @media (max-width: 450px) {
      .xs-mt-1 {
        margin-top: 4px !important
      }
      .xs-table-row {
        display: table-row !important
      }
      .xs-w-full {
        width: 100% !important
      }
      .xs-px-4 {
        padding-left: 16px !important;
        padding-right: 16px !important
      }
      .xs-pl-6 {
        padding-left: 24px !important
      }
      .xs-pr-6 {
        padding-right: 24px !important
      }
      .xs-text-center {
        text-align: center !important
      }
    }
    @media (max-width: 375px) {
      .xxs-hidden {
        display: none !important
      }
    }
    @media (prefers-color-scheme: dark) {
      .dark-mode-bg-light-100 {
        background-color: #f5f5f5 !important
      }
      .dark-mode-bg-white {
        background-color: #ffffff !important
      }
      .dark-mode-text-_1E1E1E {
        color: #1E1E1E !important
      }
      .dark-mode-text-black {
        color: #000000 !important
      }
      .dark-mode-text-light-600 {
        color: #7e7e7e !important
      }
    }
  </style>
</head>
<body style="margin: 0; width: 100%; background-color: #f5f5f5; padding: 0; font-size: 14px; line-height: 20px; -webkit-font-smoothing: antialiased; word-break: break-word">

    <div style="display: none">
      &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847;
    </div>

  <div role="article" aria-roledescription="email" aria-label="Добро пожаловать в EMCD" lang="ru" style="width: 100%;">
    <div role="separator" style="line-height: 20px; mso-line-height-alt: 20px">&zwj;</div>
    <table align="center" class="sm-w-full" cellpadding="0" cellspacing="0" style="margin-left: auto; margin-right: auto; width: 600px; font-family: ''RoobertPRO'', ''Arial'', ''Helvetica'', ui-sans-serif, system-ui, -apple-system, ''Segoe UI'', sans-serif" role="none">

        <tr>
          <td class="sm-p-0" style="padding-left: 8px; padding-right: 8px">
            <!-- header: с иконками соц.сетей -->
            <div style="width: 100%; border-top-left-radius: 10px; border-top-right-radius: 10px; background-color: #000000; color: #ffffff">
              <table cellpadding="0" cellspacing="0" style="margin: 0; width: 100%; border-top-left-radius: 10px; border-top-right-radius: 10px; padding: 0; background: linear-gradient(#000000, #000000); background-color: #000000;" role="none">
                <tr style="vertical-align: middle">
                  <td class="xs-pl-6" style="width: 100px; padding-top: 24px; padding-bottom: 24px; padding-left: 32px; vertical-align: middle">
                    <a href="https://emcd.io" style="display: inline-block; cursor: pointer">
                      <img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718368264082_logo-white%20(1)_01J0BBNSJXY1NWTMMXGEPDZYW4.png" alt="EMCD" width="100" height="27" border="0" style="max-width: 100%; vertical-align: middle; line-height: 1">
                    </a>
                  </td>
                  <td class="xs-pr-6" style="padding-top: 24px; padding-bottom: 24px; padding-right: 32px; text-align: right; vertical-align: middle">
                    <table cellpadding="0" cellspacing="0" style="float: right; clear: both; height: 24px; width: auto" role="none">
                      <tr style="vertical-align: middle;">
                        <td style="margin: 0; height: 24px; width: 24px; padding: 0 12px 0 0"><a href="https://instagram.com/emcd_io?link_id={% cio_link_id %}" class="untracked" style="display: inline-block; width: 24px; cursor: pointer;"><img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1714985175903_u_instagram-alt_01HX6HA39WZC86BCPFGXJTHPA7.png" alt="Instagram" border="0" width="24" height="24" style="max-width: 100%; vertical-align: middle; line-height: 1;"></a></td>
                        <td style="margin: 0; height: 24px; width: 24px; padding: 0 12px 0 0;"><a href="https://t.me/Emcdnews" style="display: inline-block; width: 24px; cursor: pointer;"><img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1714985174519_u_telegram_01HX6HA1YKQ9RWE6K7FQXR3BNE.png" alt="Telegram" border="0" width="24" height="24" style="max-width: 100%; vertical-align: middle; line-height: 1;"></a></td>
                        <td style="margin: 0; height: 24px; width: 24px; padding: 0 12px 0 0;"><a href="https://vk.com/emcd_io" style="display: inline-block; width: 24px; cursor: pointer;"><img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1714985173166_u_vk_01HX6HA0MA0RB8MPF11T6YXN4N.png" alt="VK" border="0" width="24" height="24" style="max-width: 100%; vertical-align: middle; line-height: 1;"></a></td>
                        <td style="margin: 0; height: 24px; width: 24px; padding: 0 12px 0 0;"><a href="https://www.youtube.com/channel/UCOjxcJiwe-9F187C1MoS84w" style="display: inline-block; width: 24px; cursor: pointer;"><img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1714985171801_u_youtube_01HX6H9Z9Q7PG6N2P17NVYTHME.png" alt="YouTube" border="0" width="24" height="24" style="max-width: 100%; vertical-align: middle; line-height: 1;"></a></td>
                        <td style="margin: 0; height: 24px; width: 24px; padding: 0;"><a href="https://vc.ru/u/2383972-emcd-mining-pool" style="display: inline-block; width: 24px; cursor: pointer;"><img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1714985170431_ic-vcru-x3_01HX6H9XYV03RS2KNFPE48ASSM.png" alt="VC.ru" border="0" width="24" height="24" style="max-width: 100%; vertical-align: middle; line-height: 1;"></a></td>
                      </tr>
                    </table>
                  </td>
                </tr>
              </table>
            </div>
            <!-- /end header: с иконками соц.сетей -->
          </td>
        </tr>

      <tr>
        <td class="sm-p-0" style="padding-left: 8px; padding-right: 8px">
          <div style="background-color: #ffffff; font-family: ''RoobertPRO'', ''Arial'', ''Helvetica'', ui-sans-serif, system-ui, -apple-system, ''Segoe UI'', sans-serif; font-size: 16px; line-height: 20px"> <!-- black bg intro -->

            <div class="dark-mode-bg-white xs-px-4" style="background-color: #ffffff; padding: 32px; font-family: ''RoobertPRO'', ''Arial'', ''Helvetica'', ui-sans-serif, system-ui, -apple-system, ''Segoe UI'', sans-serif; font-size: 16px; line-height: 20px">
              <!--  spacer -->
              <div role="separator" style="line-height: 10px; mso-line-height-alt: 56px">&zwj;</div>
              <!-- /end spacer -->
              <!-- title -->
              <div class="dark-mode-text-black" style="text-align: left; font-size: 20px; font-weight: 700; line-height: 1.15; color: #000000; mso-line-height-alt: 36px">
                Привет 👋
              </div>
              <!-- /end title -->
              <!-- spacer -->
              <div role="separator" style="line-height: 24px; mso-line-height-alt: 24px">&zwj;</div>
              <!-- /end spacer -->

            <!-- text: слева -->
            <p style="margin: 0; text-align: left; line-height: 1.4; color: #4d4d4d">

Твой запрос на получение CSV-таблицы с данными о твоих начислениях успешно обработан. Ты можешь скачать таблицу, используя ссылку ниже.
</p>
<div role="separator" style="line-height: 24px; mso-line-height-alt: 24px">&zwj;</div>
<div class="xs-text-center" style="text-align: center">
                <a href="{{.ReportLink}}" class="dark-mode-text-white dark-mode-bg-violet-300" rel="noopener noreferrer" style="box-sizing: border-box; display: inline-block; height: 56px; min-width: 200px; cursor: pointer; border-radius: 10px; background-color: #8f42ff; padding: 16px 24px; text-align: center; vertical-align: middle; font-family: ''RoobertPRO'', ''Arial'', ''Helvetica'', ui-sans-serif, system-ui, -apple-system, ''Segoe UI'', sans-serif; font-size: 16px; font-weight: 700; line-height: 1.3; color: #ffffff; text-decoration-line: none">
                  <!--[if mso]><i style="mso-font-width: -100%; letter-spacing: 32px; mso-text-raise: 30px;" hidden>&nbsp;</i><![endif]-->
                  <span style="vertical-align: middle; color: inherit; mso-text-raise: 16px;">
Скачать таблицу</span>
                  <!--[if mso]><i style="mso-font-width: -100%; letter-spacing: 32px;" hidden>&nbsp;</i><![endif]-->
                </a>
              </div>
            <!-- /end text: слева -->
<div role="separator" style="background-color: #EBEBEB; height: 1px; line-height: 1px; margin: 32px 0">&zwj;</div>
<!-- text: слева -->

            <!-- /end text: слева -->




              <!-- spacer -->





              <!-- /end spacer -->


          </div>
        </td>
      </tr>
      <tr>
        <td class="sm-p-0" style="padding-left: 8px; padding-right: 8px">
          <!-- footer: на светлом фоне, краткий -->
          <div class="dark-mode-bg-white dark-mode-text-light-600 xs-px-4" style="background-color: #ffffff; padding-left: 32px; padding-right: 32px; color: #7e7e7e">
            <div class="dark-mode-bg-light-100" style="border-radius: 10px; background-color: #f5f5f5; padding: 24px 24px 8px; text-align: center">
<img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718356995357_logo-black_01J0B0XWJZC50MPY0JKHMPDS8Z.png" alt="EMCD" width="57" height="18" border="0" style="max-width: 100%; vertical-align: middle; line-height: 1">
              <div role="separator" style="line-height: 8px; mso-line-height-alt: 8px">&zwj;</div>
              <div class="dark-mode-text-black sm-text-2xl" style="text-align: center; font-size: 30px; font-weight: 700; line-height: 1.15; color: #000000; mso-line-height-alt: 36px"> Крипто финтех платформа, где все просто
              </div>
              <div role="separator" style="line-height: 8px; mso-line-height-alt: 8px">&zwj;</div>
              <!-- button: две кнопки -->
              <div style="width: 100%;">
                <table cellpadding="0" cellspacing="16" style="margin-left: auto; margin-right: auto; width: auto; vertical-align: middle;" role="none">
                  <tr style="vertical-align: middle;">
                    <td class="sm-table-row">
                      <!-- button: черная с картинкой -->
                      <div style="text-align: center;">
                        <a href="https://emcd.onelink.me/FCtc/x4ojb23m" rel="noopener noreferrer" style="box-sizing: border-box; display: inline-block; cursor: pointer; text-decoration-line: none;"> <img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718624674835_store-android-dark_01J0K06SZC7G7BFWYM0EN64E14.png" alt="Google Play" border="0" height="56" width="200" style="max-width: 100%; vertical-align: middle; line-height: 1; display: block; height: 56px; width: 200px; border-radius: 10px">
                        </a>
                      </div>
                      <!-- /end button: черная с картинкой -->
                    </td>
                    <td class="sm-table-row">
                      <!-- button: черная с картинкой -->
                      <div style="text-align: center;">
                        <a href="https://emcd.onelink.me/FCtc/x4ojb23m	" rel="noopener noreferrer" style="box-sizing: border-box; display: inline-block; cursor: pointer; text-decoration-line: none;"> <img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718624678280_store-ios-dark_01J0K06X34Q7DS84H8YMTS7VRB.png" alt="Apple Store" border="0" height="56" width="200" style="max-width: 100%; vertical-align: middle; line-height: 1; display: block; height: 56px; width: 200px; border-radius: 10px;">
                        </a>
                      </div>
                      <!-- /end button: черная с картинкой -->
                    </td>
                  </tr>
                </table>
              </div>
              <!-- /end button: две кнопки -->
            </div>
            <div role="separator" style="line-height: 10px; mso-line-height-alt: 32px">&zwj;</div>
  <div> <a href="https://www.trustpilot.com/review/emcd.io" style="display: block; text-decoration-line: none;">
                <div style="width: 100%;">
                  <img class="sm-hidden" src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1724838104440_%F0%9F%9A%A3%20Trustpilot%20web_01J6C5SMHWHKHNFF1YYPZSJMRR.png" alt style="max-width: 100%; vertical-align: middle; line-height: 1; height: auto; width: 100%">
                  <img class="sm-block" src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1724838093974_%F0%9F%9A%A3%20Trustpilot%20mob_01J6C5SJ6GX8X6HN5QQC5CTBGB.png" alt style="max-width: 100%; vertical-align: middle; line-height: 1; display: none; height: auto; width: 100%">
                </div>
    </a>

            </div>
<div role="separator" style="line-height: 36px; mso-line-height-alt: 32px">&zwj;</div>
            <div style="text-align: center;">
              <p class="dark-mode-text-black" style="margin: 0 0 8px; font-size: 20px; font-weight: 700; color: #000000">Остались вопросы?</p>
              <p style="margin: 0; font-size: 12px">Напиши нам   <a href="mailto:support@emcd.io" class="hover-no-underline dark-mode-text-light-600" style="display: inline-block; cursor: pointer; color: #7e7e7e; text-decoration-line: underline; text-underline-offset: 2px">support@emcd.io</a> </p>
            </div>
            <div role="separator" style="line-height: 24px; mso-line-height-alt: 24px">&zwj;</div>
            <div role="separator" style="background-color: #EBEBEB; height: 1px; line-height: 1px; margin: 8px 0;">&zwj;</div>
            <div role="separator" style="line-height: 24px; mso-line-height-alt: 24px">&zwj;</div>
            <div style="text-align: center;">
              <p style="margin: 0 0 8px; font-size: 12px;">Copyright © 2024 <i>EMCD</i> Tech ltd.</p>

          <p style="margin: 0; font-size: 12px;">Если письмо отображается некорректно, открой его в браузере.</p>

            </div>
            <div class="xs-text-center">
              <a href="{% view_in_browser_url %}" class="dark-mode-text-light-600 dark-mode-bg-light-100" rel="noopener noreferrer" style="margin-top: 24px; box-sizing: border-box; display: inline-block; height: 56px; width: 100%; min-width: 200px; cursor: pointer; border-radius: 10px; background-color: #f5f5f5; padding: 16px 24px; text-align: center; vertical-align: middle; font-family: ''RoobertPRO'', ''Arial'', ''Helvetica'', ui-sans-serif, system-ui, -apple-system, ''Segoe UI'', sans-serif; font-size: 13px; font-weight: 400; line-height: 1.7; color: #7e7e7e; text-decoration-line: none">
                <!--[if mso]><i style="mso-font-width: -100%; letter-spacing: 32px; mso-text-raise: 30px;" hidden>&nbsp;</i><![endif]-->
                <span style="vertical-align: middle; color: inherit; mso-text-raise: 16px;">
      Открыть письмо в браузере</span>
                <!--[if mso]><i style="mso-font-width: -100%; letter-spacing: 32px;" hidden>&nbsp;</i><![endif]-->
              </a>
            </div>

            <div role="separator" style="line-height: 24px; mso-line-height-alt: 24px">&zwj;</div>
            <div role="separator" style="background-color: #EBEBEB; height: 1px; line-height: 1px; margin: 8px 0;">&zwj;</div>
            <div role="separator" style="line-height: 10px; mso-line-height-alt: 10px">&zwj;</div>
            <table cellpadding="0" cellspacing="0" style="margin-left: auto; margin-right: auto; width: auto; border-style: none;" role="none">
              <tr style="vertical-align: middle;">
                <td style="padding-left: 8px; padding-right: 8px; vertical-align: middle;">
                  <table cellpadding="0" cellspacing="0" style="width: auto; border-style: none;" role="none">
                    <tr style="vertical-align: middle;">

                        <td class="xxs-hidden">
                          <a href="https://emcd.io/pool/dashboard?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;"> <img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718356980754_footer-sr-thunder_01J0B0XEAKQ8QVPR1WB7G2P1TT.png" border="0" alt width="16" height="16" style="max-width: 100%; vertical-align: middle; line-height: 1; margin-right: 8px">
                          </a>
                        </td>

                      <td>
                        <a href="https://emcd.io/pool/dashboard?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;">
                          <span class="dark-mode-text-_1E1E1E" style="font-size: 12px; line-height: 12px; color: #1E1E1E">
                Майнинг</span>
                        </a>
                      </td>
                    </tr>
                  </table>
                </td>
                <td style="padding-left: 8px; padding-right: 8px; vertical-align: middle;">
                  <table cellpadding="0" cellspacing="0" style="width: auto; border-style: none;" role="none">
                    <tr style="vertical-align: middle;">

                        <td class="xxs-hidden">
                          <a href="https://emcd.io/wallets?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;"> <img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718356985346_footer-sr-twoarrow_01J0B0XJTEMH8772NQB74BS2AY.png" border="0" alt width="16" height="16" style="max-width: 100%; vertical-align: middle; line-height: 1; margin-right: 8px;">
                          </a>
                        </td>

                      <td>
                        <a href="https://emcd.io/wallets?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;">
                          <span class="dark-mode-text-_1E1E1E" style="font-size: 12px; line-height: 12px; color: #1E1E1E">
                Кошелек</span>
                        </a>
                      </td>
                    </tr>
                  </table>
                </td>
                <td style="padding-left: 8px; padding-right: 8px; vertical-align: middle;">
                  <table cellpadding="0" cellspacing="0" style="width: auto; border-style: none;" role="none">
                    <tr style="vertical-align: middle;">

                        <td class="xxs-hidden">
                          <a href="https://emcd.io/deposits?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;"> <img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718356974038_footer-sr-arrow-down_01J0B0X7SJVMRNQ6QAZ37TXW0T.png" border="0" alt width="16" height="16" style="max-width: 100%; vertical-align: middle; line-height: 1; margin-right: 8px;">
                          </a>
                        </td>

                      <td>
                        <a href="https://emcd.io/deposits?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;">
                          <span class="dark-mode-text-_1E1E1E" style="font-size: 12px; line-height: 12px; color: #1E1E1E">Coinhold</span>
                        </a>
                      </td>
                    </tr>
                  </table>
                </td>
                <td style="padding-left: 8px; padding-right: 8px; vertical-align: middle;">
                  <table cellpadding="0" cellspacing="0" style="width: auto; border-style: none;" role="none">
                    <tr style="vertical-align: middle;">

                        <td class="xxs-hidden">
                          <a href="https://emcd.io/p2p?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;"> <img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718356977464_footer-sr-coin-dollar_01J0B0XB3WCNB24H8YB39BZFR5.png" border="0" alt width="16" height="16" style="max-width: 100%; vertical-align: middle; line-height: 1; margin-right: 8px;">
                          </a>
                        </td>

                      <td>
                        <a href="https://emcd.io/p2p?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;">
                          <span class="dark-mode-text-_1E1E1E" style="font-size: 12px; line-height: 12px; color: #1E1E1E">P2P</span>
                        </a>
                      </td>
                    </tr>
                  </table>
                </td>
              </tr>
            </table>
            <div role="separator" style="line-height: 24px; mso-line-height-alt: 24px">&zwj;</div>
          </div>
          <!-- /end footer: на светлом фоне, краткий -->
        </td>
      </tr>
    </table>
    <div role="separator" style="line-height: 40px; mso-line-height-alt: 40px">&zwj;</div>
  </div>
</body>
</html>',
       'income report',
       'ru',
       'Твой запрос на выгрузку данных о начислениях',
       NULL
WHERE NOT EXISTS(
        SELECT type, language FROM email_templates WHERE type = 'income report' AND language = 'ru'
    );

--- payout report templates
INSERT INTO email_templates (whitelabel_id, template, type, language, subject, footer)
SELECT '00000000-0000-0000-0000-000000000000',
       '<!DOCTYPE html>
<html lang="ru" xmlns:v="urn:schemas-microsoft-com:vml">
<head>
  <meta charset="utf-8">
  <meta name="x-apple-disable-message-reformatting">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <meta name="format-detection" content="telephone=no, date=no, address=no, email=no, url=no">
  <meta name="color-scheme" content="light dark">
  <meta name="supported-color-schemes" content="light dark">
  <!--[if mso]>
  <noscript>
    <xml>
      <o:OfficeDocumentSettings xmlns:o="urn:schemas-microsoft-com:office:office">
        <o:PixelsPerInch>96</o:PixelsPerInch>
      </o:OfficeDocumentSettings>
    </xml>
  </noscript>
  <style>
    td,th,div,p,a,h1,h2,h3,h4,h5,h6 {font-family: "RoobertPRO", Helvetica, Arial, "Segoe UI", sans-serif; mso-line-height-rule: exactly;}
  </style>
  <![endif]-->
  <title>Добро пожаловать в EMCD</title>
  <style>
    .hover-no-underline:hover {
      text-decoration-line: none !important
    }
    @media (max-width: 640px) {
      .sm-block {
        display: block !important
      }
      .sm-table-row {
        display: table-row !important
      }
      .sm-hidden {
        display: none !important
      }
      .sm-min-h-130px {
        min-height: 130px !important
      }
      .sm-w-full {
        width: 100% !important
      }
      .sm-p-0 {
        padding: 0 !important
      }
      .sm-text-2xl {
        font-size: 24px !important
      }
      .sm-text-lg {
        font-size: 18px !important
      }
    }
    @media (max-width: 450px) {
      .xs-mt-1 {
        margin-top: 4px !important
      }
      .xs-table-row {
        display: table-row !important
      }
      .xs-w-full {
        width: 100% !important
      }
      .xs-px-4 {
        padding-left: 16px !important;
        padding-right: 16px !important
      }
      .xs-pl-6 {
        padding-left: 24px !important
      }
      .xs-pr-6 {
        padding-right: 24px !important
      }
      .xs-text-center {
        text-align: center !important
      }
    }
    @media (max-width: 375px) {
      .xxs-hidden {
        display: none !important
      }
    }
    @media (prefers-color-scheme: dark) {
      .dark-mode-bg-light-100 {
        background-color: #f5f5f5 !important
      }
      .dark-mode-bg-white {
        background-color: #ffffff !important
      }
      .dark-mode-text-_1E1E1E {
        color: #1E1E1E !important
      }
      .dark-mode-text-black {
        color: #000000 !important
      }
      .dark-mode-text-light-600 {
        color: #7e7e7e !important
      }
    }
  </style>
</head>
<body style="margin: 0; width: 100%; background-color: #f5f5f5; padding: 0; font-size: 14px; line-height: 20px; -webkit-font-smoothing: antialiased; word-break: break-word">

    <div style="display: none">
      &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847;
    </div>

  <div role="article" aria-roledescription="email" aria-label="Добро пожаловать в EMCD" lang="ru" style="width: 100%;">
    <div role="separator" style="line-height: 20px; mso-line-height-alt: 20px">&zwj;</div>
    <table align="center" class="sm-w-full" cellpadding="0" cellspacing="0" style="margin-left: auto; margin-right: auto; width: 600px; font-family: ''RoobertPRO'', ''Arial'', ''Helvetica'', ui-sans-serif, system-ui, -apple-system, ''Segoe UI'', sans-serif" role="none">

        <tr>
          <td class="sm-p-0" style="padding-left: 8px; padding-right: 8px">
            <!-- header: с иконками соц.сетей -->
            <div style="width: 100%; border-top-left-radius: 10px; border-top-right-radius: 10px; background-color: #000000; color: #ffffff">
              <table cellpadding="0" cellspacing="0" style="margin: 0; width: 100%; border-top-left-radius: 10px; border-top-right-radius: 10px; padding: 0; background: linear-gradient(#000000, #000000); background-color: #000000;" role="none">
                <tr style="vertical-align: middle">
                  <td class="xs-pl-6" style="width: 100px; padding-top: 24px; padding-bottom: 24px; padding-left: 32px; vertical-align: middle">
                    <a href="https://emcd.io" style="display: inline-block; cursor: pointer">
                      <img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718368264082_logo-white%20(1)_01J0BBNSJXY1NWTMMXGEPDZYW4.png" alt="EMCD" width="100" height="27" border="0" style="max-width: 100%; vertical-align: middle; line-height: 1">
                    </a>
                  </td>
                  <td class="xs-pr-6" style="padding-top: 24px; padding-bottom: 24px; padding-right: 32px; text-align: right; vertical-align: middle">
                    <table cellpadding="0" cellspacing="0" style="float: right; clear: both; height: 24px; width: auto" role="none">
                      <tr style="vertical-align: middle;">
                        <td style="margin: 0; height: 24px; width: 24px; padding: 0 12px 0 0">
                          <a href="https://www.instagram.com/emcd_io/" class="untracked" style="display: inline-block; width: 24px; cursor: pointer;"><img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1714985175903_u_instagram-alt_01HX6HA39WZC86BCPFGXJTHPA7.png" alt="Instagram" border="0" width="24" height="24" style="max-width: 100%; vertical-align: middle; line-height: 1;"></a>
                        </td>
                        <td style="margin: 0; height: 24px; width: 24px; padding: 0 12px 0 0;">
                          <a href="https://t.me/emcd_international" style="display: inline-block; width: 24px; cursor: pointer;"><img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1714985174519_u_telegram_01HX6HA1YKQ9RWE6K7FQXR3BNE.png" alt="Telegram" border="0" width="24" height="24" style="max-width: 100%; vertical-align: middle; line-height: 1;"></a>
                        </td>
                        <td style="margin: 0; height: 24px; width: 24px; padding: 0 12px 0 0;"><a href="https://www.youtube.com/channel/UCEOI3erKte4PIKpBflZ00ZQ" style="display: inline-block; width: 24px; cursor: pointer;"><img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1714985171801_u_youtube_01HX6H9Z9Q7PG6N2P17NVYTHME.png" alt="YouTube" border="0" width="24" height="24" style="max-width: 100%; vertical-align: middle; line-height: 1;"></a></td>

                        <td style="margin: 0; height: 24px; width: 24px; padding: 0 12px 0 0;">
                          <a href="https://www.linkedin.com/company/emcdtech/" style="display: inline-block; width: 24px; cursor: pointer;"><img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1714985169108_u_linkedin_01HX6H9WNHT4423NVPJZEGMQS4.png" alt="LinkedIn" border="0" width="24" height="24" style="max-width: 100%; vertical-align: middle; line-height: 1;"></a>
                        </td>
                        <td style="margin: 0; height: 24px; width: 24px; padding: 0;">
                          <a href="https://twitter.com/emcd_io" style="display: inline-block; width: 24px; cursor: pointer;"><img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1714985167624_Group%201410093531_01HX6H9V8TJRR90TW90PF32D5X.png" alt="X (Twitter)" border="0" width="24" height="24" style="max-width: 100%; vertical-align: middle; line-height: 1;"></a>
                        </td>
                      </tr>
                    </table>
                  </td>
                </tr>
              </table>
            </div>
            <!-- /end header: с иконками соц.сетей -->
          </td>
        </tr>

      <tr>
        <td class="sm-p-0" style="padding-left: 8px; padding-right: 8px">
          <div style="background-color: #ffffff; font-family: ''RoobertPRO'', ''Arial'', ''Helvetica'', ui-sans-serif, system-ui, -apple-system, ''Segoe UI'', sans-serif; font-size: 16px; line-height: 20px"> <!-- black bg intro -->

            <div class="dark-mode-bg-white xs-px-4" style="background-color: #ffffff; padding: 32px; font-family: ''RoobertPRO'', ''Arial'', ''Helvetica'', ui-sans-serif, system-ui, -apple-system, ''Segoe UI'', sans-serif; font-size: 16px; line-height: 20px">
              <!--  spacer -->
              <div role="separator" style="line-height: 10px; mso-line-height-alt: 56px">&zwj;</div>
              <!-- /end spacer -->
              <!-- title -->
              <div class="dark-mode-text-black" style="text-align: left; font-size: 20px; font-weight: 700; line-height: 1.15; color: #000000; mso-line-height-alt: 36px">
                Hi there 👋
              </div>
              <!-- /end title -->
              <!-- spacer -->
              <div role="separator" style="line-height: 24px; mso-line-height-alt: 24px">&zwj;</div>
              <!-- /end spacer -->

            <!-- text: слева -->
            <p style="margin: 0; text-align: left; line-height: 1.4; color: #4d4d4d">

Your request to receive a CSV file containing your payout data has been successfully processed. You can download the file using the link below.
</p>
<div role="separator" style="line-height: 24px; mso-line-height-alt: 24px">&zwj;</div>
<div class="xs-text-center" style="text-align: center">
                <a href="{{.ReportLink}}" class="dark-mode-text-white dark-mode-bg-violet-300" rel="noopener noreferrer" style="box-sizing: border-box; display: inline-block; height: 56px; min-width: 200px; cursor: pointer; border-radius: 10px; background-color: #8f42ff; padding: 16px 24px; text-align: center; vertical-align: middle; font-family: ''RoobertPRO'', ''Arial'', ''Helvetica'', ui-sans-serif, system-ui, -apple-system, ''Segoe UI'', sans-serif; font-size: 16px; font-weight: 700; line-height: 1.3; color: #ffffff; text-decoration-line: none">
                  <!--[if mso]><i style="mso-font-width: -100%; letter-spacing: 32px; mso-text-raise: 30px;" hidden>&nbsp;</i><![endif]-->
                  <span style="vertical-align: middle; color: inherit; mso-text-raise: 16px;">Download</span>
                  <!--[if mso]><i style="mso-font-width: -100%; letter-spacing: 32px;" hidden>&nbsp;</i><![endif]-->
                </a>
              </div>

            <!-- /end text: слева -->
<div role="separator" style="background-color: #EBEBEB; height: 1px; line-height: 1px; margin: 32px 0">&zwj;</div>
<!-- text: слева -->

            <!-- /end text: слева -->




              <!-- spacer -->





              <!-- /end spacer -->


          </div>
        </td>
      </tr>
      <tr>
        <td class="sm-p-0" style="padding-left: 8px; padding-right: 8px">
          <!-- footer: на светлом фоне, краткий -->
          <div class="dark-mode-bg-white dark-mode-text-light-600 xs-px-4" style="background-color: #ffffff; padding-left: 32px; padding-right: 32px; color: #7e7e7e">
            <div class="dark-mode-bg-light-100" style="border-radius: 10px; background-color: #f5f5f5; padding: 24px 24px 8px; text-align: center">
<img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718356995357_logo-black_01J0B0XWJZC50MPY0JKHMPDS8Z.png" alt="EMCD" width="57" height="18" border="0" style="max-width: 100%; vertical-align: middle; line-height: 1">
              <div role="separator" style="line-height: 8px; mso-line-height-alt: 8px">&zwj;</div>
              <div class="dark-mode-text-black sm-text-2xl" style="text-align: center; font-size: 30px; font-weight: 700; line-height: 1.15; color: #000000; mso-line-height-alt: 36px"> A crypto fintech platform where everything is simple
              </div>
              <div role="separator" style="line-height: 8px; mso-line-height-alt: 8px">&zwj;</div>
              <!-- button: две кнопки -->
              <div style="width: 100%;">
                <table cellpadding="0" cellspacing="16" style="margin-left: auto; margin-right: auto; width: auto; vertical-align: middle;" role="none">
                  <tr style="vertical-align: middle;">
                    <td class="sm-table-row">
                      <!-- button: черная с картинкой -->
                      <div style="text-align: center;">
                        <a href="https://emcd.onelink.me/FCtc/x4ojb23m" rel="noopener noreferrer" style="box-sizing: border-box; display: inline-block; cursor: pointer; text-decoration-line: none;"> <img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718624674835_store-android-dark_01J0K06SZC7G7BFWYM0EN64E14.png" alt="Google Play" border="0" height="56" width="200" style="max-width: 100%; vertical-align: middle; line-height: 1; display: block; height: 56px; width: 200px; border-radius: 10px">
                        </a>
                      </div>
                      <!-- /end button: черная с картинкой -->
                    </td>
                    <td class="sm-table-row">
                      <!-- button: черная с картинкой -->
                      <div style="text-align: center;">
                        <a href="https://emcd.onelink.me/FCtc/x4ojb23m	" rel="noopener noreferrer" style="box-sizing: border-box; display: inline-block; cursor: pointer; text-decoration-line: none;"> <img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718624678280_store-ios-dark_01J0K06X34Q7DS84H8YMTS7VRB.png" alt="Apple Store" border="0" height="56" width="200" style="max-width: 100%; vertical-align: middle; line-height: 1; display: block; height: 56px; width: 200px; border-radius: 10px;">
                        </a>
                      </div>
                      <!-- /end button: черная с картинкой -->
                    </td>
                  </tr>
                </table>
              </div>
              <!-- /end button: две кнопки -->
            </div>
            <div role="separator" style="line-height: 10px; mso-line-height-alt: 32px">&zwj;</div>
  <div> <a href="https://www.trustpilot.com/review/emcd.io" style="display: block; text-decoration-line: none;">
                <div style="width: 100%;">
                  <img class="sm-hidden" src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1724838104440_%F0%9F%9A%A3%20Trustpilot%20web_01J6C5SMHWHKHNFF1YYPZSJMRR.png" alt style="max-width: 100%; vertical-align: middle; line-height: 1; height: auto; width: 100%">
                  <img class="sm-block" src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1724838093974_%F0%9F%9A%A3%20Trustpilot%20mob_01J6C5SJ6GX8X6HN5QQC5CTBGB.png" alt style="max-width: 100%; vertical-align: middle; line-height: 1; display: none; height: auto; width: 100%">
                </div>
    </a>

            </div>
<div role="separator" style="line-height: 36px; mso-line-height-alt: 32px">&zwj;</div>
            <div style="text-align: center;">
              <p class="dark-mode-text-black" style="margin: 0 0 8px; font-size: 20px; font-weight: 700; color: #000000">Do you have any questions?</p>
              <p style="margin: 0; font-size: 12px">Email us at   <a href="mailto:support@emcd.io" class="hover-no-underline dark-mode-text-light-600" style="display: inline-block; cursor: pointer; color: #7e7e7e; text-decoration-line: underline; text-underline-offset: 2px">support@emcd.io</a> </p>
            </div>
            <div role="separator" style="line-height: 24px; mso-line-height-alt: 24px">&zwj;</div>
            <div role="separator" style="background-color: #EBEBEB; height: 1px; line-height: 1px; margin: 8px 0;">&zwj;</div>
            <div role="separator" style="line-height: 24px; mso-line-height-alt: 24px">&zwj;</div>
            <div style="text-align: center;">
              <p style="margin: 0 0 8px; font-size: 12px;">Copyright © 2024 <i>EMCD</i> Tech ltd.</p>

              <p style="margin: 0; font-size: 12px;">If the email is not displayed correctly, open it in your browser.</p>

            </div>
            <div class="xs-text-center">
              <a href="{% view_in_browser_url %}" class="dark-mode-text-light-600 dark-mode-bg-light-100" rel="noopener noreferrer" style="margin-top: 24px; box-sizing: border-box; display: inline-block; height: 56px; width: 100%; min-width: 200px; cursor: pointer; border-radius: 10px; background-color: #f5f5f5; padding: 16px 24px; text-align: center; vertical-align: middle; font-family: ''RoobertPRO'', ''Arial'', ''Helvetica'', ui-sans-serif, system-ui, -apple-system, ''Segoe UI'', sans-serif; font-size: 13px; font-weight: 400; line-height: 1.7; color: #7e7e7e; text-decoration-line: none">
                <!--[if mso]><i style="mso-font-width: -100%; letter-spacing: 32px; mso-text-raise: 30px;" hidden>&nbsp;</i><![endif]-->
                <span style="vertical-align: middle; color: inherit; mso-text-raise: 16px;">
      Open web version of the email</span>
                <!--[if mso]><i style="mso-font-width: -100%; letter-spacing: 32px;" hidden>&nbsp;</i><![endif]-->
              </a>
            </div>

            <div role="separator" style="line-height: 24px; mso-line-height-alt: 24px">&zwj;</div>
            <div role="separator" style="background-color: #EBEBEB; height: 1px; line-height: 1px; margin: 8px 0;">&zwj;</div>
            <div role="separator" style="line-height: 10px; mso-line-height-alt: 10px">&zwj;</div>
            <table cellpadding="0" cellspacing="0" style="margin-left: auto; margin-right: auto; width: auto; border-style: none;" role="none">
              <tr style="vertical-align: middle;">
                <td style="padding-left: 8px; padding-right: 8px; vertical-align: middle;">
                  <table cellpadding="0" cellspacing="0" style="width: auto; border-style: none;" role="none">
                    <tr style="vertical-align: middle;">

                        <td class="xxs-hidden">
                          <a href="https://emcd.io/pool/dashboard?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;"> <img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718356980754_footer-sr-thunder_01J0B0XEAKQ8QVPR1WB7G2P1TT.png" border="0" alt width="16" height="16" style="max-width: 100%; vertical-align: middle; line-height: 1; margin-right: 8px">
                          </a>
                        </td>

                      <td>
                        <a href="https://emcd.io/pool/dashboard?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;">
                          <span class="dark-mode-text-_1E1E1E" style="font-size: 12px; line-height: 12px; color: #1E1E1E">
                Mining</span>
                        </a>
                      </td>
                    </tr>
                  </table>
                </td>
                <td style="padding-left: 8px; padding-right: 8px; vertical-align: middle;">
                  <table cellpadding="0" cellspacing="0" style="width: auto; border-style: none;" role="none">
                    <tr style="vertical-align: middle;">

                        <td class="xxs-hidden">
                          <a href="https://emcd.io/wallets?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;"> <img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718356985346_footer-sr-twoarrow_01J0B0XJTEMH8772NQB74BS2AY.png" border="0" alt width="16" height="16" style="max-width: 100%; vertical-align: middle; line-height: 1; margin-right: 8px;">
                          </a>
                        </td>

                      <td>
                        <a href="https://emcd.io/wallets?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;">
                          <span class="dark-mode-text-_1E1E1E" style="font-size: 12px; line-height: 12px; color: #1E1E1E">
                Wallet</span>
                        </a>
                      </td>
                    </tr>
                  </table>
                </td>
                <td style="padding-left: 8px; padding-right: 8px; vertical-align: middle;">
                  <table cellpadding="0" cellspacing="0" style="width: auto; border-style: none;" role="none">
                    <tr style="vertical-align: middle;">

                        <td class="xxs-hidden">
                          <a href="https://emcd.io/deposits?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;"> <img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718356974038_footer-sr-arrow-down_01J0B0X7SJVMRNQ6QAZ37TXW0T.png" border="0" alt width="16" height="16" style="max-width: 100%; vertical-align: middle; line-height: 1; margin-right: 8px;">
                          </a>
                        </td>

                      <td>
                        <a href="https://emcd.io/deposits?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;">
                          <span class="dark-mode-text-_1E1E1E" style="font-size: 12px; line-height: 12px; color: #1E1E1E">Coinhold</span>
                        </a>
                      </td>
                    </tr>
                  </table>
                </td>
                <td style="padding-left: 8px; padding-right: 8px; vertical-align: middle;">
                  <table cellpadding="0" cellspacing="0" style="width: auto; border-style: none;" role="none">
                    <tr style="vertical-align: middle;">

                        <td class="xxs-hidden">
                          <a href="https://emcd.io/p2p?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;"> <img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718356977464_footer-sr-coin-dollar_01J0B0XB3WCNB24H8YB39BZFR5.png" border="0" alt width="16" height="16" style="max-width: 100%; vertical-align: middle; line-height: 1; margin-right: 8px;">
                          </a>
                        </td>

                      <td>
                        <a href="https://emcd.io/p2p?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;">
                          <span class="dark-mode-text-_1E1E1E" style="font-size: 12px; line-height: 12px; color: #1E1E1E">P2P</span>
                        </a>
                      </td>
                    </tr>
                  </table>
                </td>
              </tr>
            </table>
            <div role="separator" style="line-height: 24px; mso-line-height-alt: 24px">&zwj;</div>
          </div>
          <!-- /end footer: на светлом фоне, краткий -->
        </td>
      </tr>
    </table>
    <div role="separator" style="line-height: 40px; mso-line-height-alt: 40px">&zwj;</div>
  </div>
</body>
</html>',
       'payout report',
       'en',
       'Your Payout Data Export Request',
       NULL
WHERE NOT EXISTS(
        SELECT type, language FROM email_templates WHERE type = 'payout report' AND language = 'en'
    );

INSERT INTO email_templates (whitelabel_id, template, type, language, subject, footer)
SELECT '00000000-0000-0000-0000-000000000000',
       '<!DOCTYPE html>
<html lang="ru" xmlns:v="urn:schemas-microsoft-com:vml">
<head>
  <meta charset="utf-8">
  <meta name="x-apple-disable-message-reformatting">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <meta name="format-detection" content="telephone=no, date=no, address=no, email=no, url=no">
  <meta name="color-scheme" content="light dark">
  <meta name="supported-color-schemes" content="light dark">
  <!--[if mso]>
  <noscript>
    <xml>
      <o:OfficeDocumentSettings xmlns:o="urn:schemas-microsoft-com:office:office">
        <o:PixelsPerInch>96</o:PixelsPerInch>
      </o:OfficeDocumentSettings>
    </xml>
  </noscript>
  <style>
    td,th,div,p,a,h1,h2,h3,h4,h5,h6 {font-family: "RoobertPRO", Helvetica, Arial, "Segoe UI", sans-serif; mso-line-height-rule: exactly;}
  </style>
  <![endif]-->
  <title>Добро пожаловать в EMCD</title>
  <style>
    .hover-no-underline:hover {
      text-decoration-line: none !important
    }
    @media (max-width: 640px) {
      .sm-block {
        display: block !important
      }
      .sm-table-row {
        display: table-row !important
      }
      .sm-hidden {
        display: none !important
      }
      .sm-min-h-130px {
        min-height: 130px !important
      }
      .sm-w-full {
        width: 100% !important
      }
      .sm-p-0 {
        padding: 0 !important
      }
      .sm-text-2xl {
        font-size: 24px !important
      }
      .sm-text-lg {
        font-size: 18px !important
      }
    }
    @media (max-width: 450px) {
      .xs-mt-1 {
        margin-top: 4px !important
      }
      .xs-table-row {
        display: table-row !important
      }
      .xs-w-full {
        width: 100% !important
      }
      .xs-px-4 {
        padding-left: 16px !important;
        padding-right: 16px !important
      }
      .xs-pl-6 {
        padding-left: 24px !important
      }
      .xs-pr-6 {
        padding-right: 24px !important
      }
      .xs-text-center {
        text-align: center !important
      }
    }
    @media (max-width: 375px) {
      .xxs-hidden {
        display: none !important
      }
    }
    @media (prefers-color-scheme: dark) {
      .dark-mode-bg-light-100 {
        background-color: #f5f5f5 !important
      }
      .dark-mode-bg-white {
        background-color: #ffffff !important
      }
      .dark-mode-text-_1E1E1E {
        color: #1E1E1E !important
      }
      .dark-mode-text-black {
        color: #000000 !important
      }
      .dark-mode-text-light-600 {
        color: #7e7e7e !important
      }
    }
  </style>
</head>
<body style="margin: 0; width: 100%; background-color: #f5f5f5; padding: 0; font-size: 14px; line-height: 20px; -webkit-font-smoothing: antialiased; word-break: break-word">

    <div style="display: none">
      &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847; &#8199;&#65279;&#847;
    </div>

  <div role="article" aria-roledescription="email" aria-label="Добро пожаловать в EMCD" lang="ru" style="width: 100%;">
    <div role="separator" style="line-height: 20px; mso-line-height-alt: 20px">&zwj;</div>
    <table align="center" class="sm-w-full" cellpadding="0" cellspacing="0" style="margin-left: auto; margin-right: auto; width: 600px; font-family: ''RoobertPRO'', ''Arial'', ''Helvetica'', ui-sans-serif, system-ui, -apple-system, ''Segoe UI'', sans-serif" role="none">

        <tr>
          <td class="sm-p-0" style="padding-left: 8px; padding-right: 8px">
            <!-- header: с иконками соц.сетей -->
            <div style="width: 100%; border-top-left-radius: 10px; border-top-right-radius: 10px; background-color: #000000; color: #ffffff">
              <table cellpadding="0" cellspacing="0" style="margin: 0; width: 100%; border-top-left-radius: 10px; border-top-right-radius: 10px; padding: 0; background: linear-gradient(#000000, #000000); background-color: #000000;" role="none">
                <tr style="vertical-align: middle">
                  <td class="xs-pl-6" style="width: 100px; padding-top: 24px; padding-bottom: 24px; padding-left: 32px; vertical-align: middle">
                    <a href="https://emcd.io" style="display: inline-block; cursor: pointer">
                      <img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718368264082_logo-white%20(1)_01J0BBNSJXY1NWTMMXGEPDZYW4.png" alt="EMCD" width="100" height="27" border="0" style="max-width: 100%; vertical-align: middle; line-height: 1">
                    </a>
                  </td>
                  <td class="xs-pr-6" style="padding-top: 24px; padding-bottom: 24px; padding-right: 32px; text-align: right; vertical-align: middle">
                    <table cellpadding="0" cellspacing="0" style="float: right; clear: both; height: 24px; width: auto" role="none">
                      <tr style="vertical-align: middle;">
                        <td style="margin: 0; height: 24px; width: 24px; padding: 0 12px 0 0"><a href="https://instagram.com/emcd_io?link_id={% cio_link_id %}" class="untracked" style="display: inline-block; width: 24px; cursor: pointer;"><img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1714985175903_u_instagram-alt_01HX6HA39WZC86BCPFGXJTHPA7.png" alt="Instagram" border="0" width="24" height="24" style="max-width: 100%; vertical-align: middle; line-height: 1;"></a></td>
                        <td style="margin: 0; height: 24px; width: 24px; padding: 0 12px 0 0;"><a href="https://t.me/Emcdnews" style="display: inline-block; width: 24px; cursor: pointer;"><img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1714985174519_u_telegram_01HX6HA1YKQ9RWE6K7FQXR3BNE.png" alt="Telegram" border="0" width="24" height="24" style="max-width: 100%; vertical-align: middle; line-height: 1;"></a></td>
                        <td style="margin: 0; height: 24px; width: 24px; padding: 0 12px 0 0;"><a href="https://vk.com/emcd_io" style="display: inline-block; width: 24px; cursor: pointer;"><img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1714985173166_u_vk_01HX6HA0MA0RB8MPF11T6YXN4N.png" alt="VK" border="0" width="24" height="24" style="max-width: 100%; vertical-align: middle; line-height: 1;"></a></td>
                        <td style="margin: 0; height: 24px; width: 24px; padding: 0 12px 0 0;"><a href="https://www.youtube.com/channel/UCOjxcJiwe-9F187C1MoS84w" style="display: inline-block; width: 24px; cursor: pointer;"><img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1714985171801_u_youtube_01HX6H9Z9Q7PG6N2P17NVYTHME.png" alt="YouTube" border="0" width="24" height="24" style="max-width: 100%; vertical-align: middle; line-height: 1;"></a></td>
                        <td style="margin: 0; height: 24px; width: 24px; padding: 0;"><a href="https://vc.ru/u/2383972-emcd-mining-pool" style="display: inline-block; width: 24px; cursor: pointer;"><img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1714985170431_ic-vcru-x3_01HX6H9XYV03RS2KNFPE48ASSM.png" alt="VC.ru" border="0" width="24" height="24" style="max-width: 100%; vertical-align: middle; line-height: 1;"></a></td>
                      </tr>
                    </table>
                  </td>
                </tr>
              </table>
            </div>
            <!-- /end header: с иконками соц.сетей -->
          </td>
        </tr>

      <tr>
        <td class="sm-p-0" style="padding-left: 8px; padding-right: 8px">
          <div style="background-color: #ffffff; font-family: ''RoobertPRO'', ''Arial'', ''Helvetica'', ui-sans-serif, system-ui, -apple-system, ''Segoe UI'', sans-serif; font-size: 16px; line-height: 20px"> <!-- black bg intro -->

            <div class="dark-mode-bg-white xs-px-4" style="background-color: #ffffff; padding: 32px; font-family: ''RoobertPRO'', ''Arial'', ''Helvetica'', ui-sans-serif, system-ui, -apple-system, ''Segoe UI'', sans-serif; font-size: 16px; line-height: 20px">
              <!--  spacer -->
              <div role="separator" style="line-height: 10px; mso-line-height-alt: 56px">&zwj;</div>
              <!-- /end spacer -->
              <!-- title -->
              <div class="dark-mode-text-black" style="text-align: left; font-size: 20px; font-weight: 700; line-height: 1.15; color: #000000; mso-line-height-alt: 36px">
                Привет 👋
              </div>
              <!-- /end title -->
              <!-- spacer -->
              <div role="separator" style="line-height: 24px; mso-line-height-alt: 24px">&zwj;</div>
              <!-- /end spacer -->

            <!-- text: слева -->
            <p style="margin: 0; text-align: left; line-height: 1.4; color: #4d4d4d">

Твой запрос на получение CSV-таблицы с данными о твоих выплатах успешно обработан. Ты можешь скачать таблицу, используя ссылку ниже.
</p>
<div role="separator" style="line-height: 24px; mso-line-height-alt: 24px">&zwj;</div>
<div class="xs-text-center" style="text-align: center">
                <a href="{{.ReportLink}}" class="dark-mode-text-white dark-mode-bg-violet-300" rel="noopener noreferrer" style="box-sizing: border-box; display: inline-block; height: 56px; min-width: 200px; cursor: pointer; border-radius: 10px; background-color: #8f42ff; padding: 16px 24px; text-align: center; vertical-align: middle; font-family: ''RoobertPRO'', ''Arial'', ''Helvetica'', ui-sans-serif, system-ui, -apple-system, ''Segoe UI'', sans-serif; font-size: 16px; font-weight: 700; line-height: 1.3; color: #ffffff; text-decoration-line: none">
                  <!--[if mso]><i style="mso-font-width: -100%; letter-spacing: 32px; mso-text-raise: 30px;" hidden>&nbsp;</i><![endif]-->
                  <span style="vertical-align: middle; color: inherit; mso-text-raise: 16px;">
Скачать таблицу</span>
                  <!--[if mso]><i style="mso-font-width: -100%; letter-spacing: 32px;" hidden>&nbsp;</i><![endif]-->
                </a>
              </div>
            <!-- /end text: слева -->
<div role="separator" style="background-color: #EBEBEB; height: 1px; line-height: 1px; margin: 32px 0">&zwj;</div>
<!-- text: слева -->

            <!-- /end text: слева -->




              <!-- spacer -->





              <!-- /end spacer -->


          </div>
        </td>
      </tr>
      <tr>
        <td class="sm-p-0" style="padding-left: 8px; padding-right: 8px">
          <!-- footer: на светлом фоне, краткий -->
          <div class="dark-mode-bg-white dark-mode-text-light-600 xs-px-4" style="background-color: #ffffff; padding-left: 32px; padding-right: 32px; color: #7e7e7e">
            <div class="dark-mode-bg-light-100" style="border-radius: 10px; background-color: #f5f5f5; padding: 24px 24px 8px; text-align: center">
<img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718356995357_logo-black_01J0B0XWJZC50MPY0JKHMPDS8Z.png" alt="EMCD" width="57" height="18" border="0" style="max-width: 100%; vertical-align: middle; line-height: 1">
              <div role="separator" style="line-height: 8px; mso-line-height-alt: 8px">&zwj;</div>
              <div class="dark-mode-text-black sm-text-2xl" style="text-align: center; font-size: 30px; font-weight: 700; line-height: 1.15; color: #000000; mso-line-height-alt: 36px"> Крипто финтех платформа, где все просто
              </div>
              <div role="separator" style="line-height: 8px; mso-line-height-alt: 8px">&zwj;</div>
              <!-- button: две кнопки -->
              <div style="width: 100%;">
                <table cellpadding="0" cellspacing="16" style="margin-left: auto; margin-right: auto; width: auto; vertical-align: middle;" role="none">
                  <tr style="vertical-align: middle;">
                    <td class="sm-table-row">
                      <!-- button: черная с картинкой -->
                      <div style="text-align: center;">
                        <a href="https://emcd.onelink.me/FCtc/x4ojb23m" rel="noopener noreferrer" style="box-sizing: border-box; display: inline-block; cursor: pointer; text-decoration-line: none;"> <img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718624674835_store-android-dark_01J0K06SZC7G7BFWYM0EN64E14.png" alt="Google Play" border="0" height="56" width="200" style="max-width: 100%; vertical-align: middle; line-height: 1; display: block; height: 56px; width: 200px; border-radius: 10px">
                        </a>
                      </div>
                      <!-- /end button: черная с картинкой -->
                    </td>
                    <td class="sm-table-row">
                      <!-- button: черная с картинкой -->
                      <div style="text-align: center;">
                        <a href="https://emcd.onelink.me/FCtc/x4ojb23m  " rel="noopener noreferrer" style="box-sizing: border-box; display: inline-block; cursor: pointer; text-decoration-line: none;"> <img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718624678280_store-ios-dark_01J0K06X34Q7DS84H8YMTS7VRB.png" alt="Apple Store" border="0" height="56" width="200" style="max-width: 100%; vertical-align: middle; line-height: 1; display: block; height: 56px; width: 200px; border-radius: 10px;">
                        </a>
                      </div>
                      <!-- /end button: черная с картинкой -->
                    </td>
                  </tr>
                </table>
              </div>
              <!-- /end button: две кнопки -->
            </div>
            <div role="separator" style="line-height: 10px; mso-line-height-alt: 32px">&zwj;</div>
  <div> <a href="https://www.trustpilot.com/review/emcd.io" style="display: block; text-decoration-line: none;">
                <div style="width: 100%;">
                  <img class="sm-hidden" src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1724838104440_%F0%9F%9A%A3%20Trustpilot%20web_01J6C5SMHWHKHNFF1YYPZSJMRR.png" alt style="max-width: 100%; vertical-align: middle; line-height: 1; height: auto; width: 100%">
                  <img class="sm-block" src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1724838093974_%F0%9F%9A%A3%20Trustpilot%20mob_01J6C5SJ6GX8X6HN5QQC5CTBGB.png" alt style="max-width: 100%; vertical-align: middle; line-height: 1; display: none; height: auto; width: 100%">
                </div>
    </a>

            </div>
<div role="separator" style="line-height: 36px; mso-line-height-alt: 32px">&zwj;</div>
            <div style="text-align: center;">
              <p class="dark-mode-text-black" style="margin: 0 0 8px; font-size: 20px; font-weight: 700; color: #000000">Остались вопросы?</p>
              <p style="margin: 0; font-size: 12px">Напиши нам   <a href="mailto:support@emcd.io" class="hover-no-underline dark-mode-text-light-600" style="display: inline-block; cursor: pointer; color: #7e7e7e; text-decoration-line: underline; text-underline-offset: 2px">support@emcd.io</a> </p>
            </div>
            <div role="separator" style="line-height: 24px; mso-line-height-alt: 24px">&zwj;</div>
            <div role="separator" style="background-color: #EBEBEB; height: 1px; line-height: 1px; margin: 8px 0;">&zwj;</div>
            <div role="separator" style="line-height: 24px; mso-line-height-alt: 24px">&zwj;</div>
            <div style="text-align: center;">
              <p style="margin: 0 0 8px; font-size: 12px;">Copyright © 2024 <i>EMCD</i> Tech ltd.</p>

          <p style="margin: 0; font-size: 12px;">Если письмо отображается некорректно, открой его в браузере.</p>

            </div>
            <div class="xs-text-center">
              <a href="{% view_in_browser_url %}" class="dark-mode-text-light-600 dark-mode-bg-light-100" rel="noopener noreferrer" style="margin-top: 24px; box-sizing: border-box; display: inline-block; height: 56px; width: 100%; min-width: 200px; cursor: pointer; border-radius: 10px; background-color: #f5f5f5; padding: 16px 24px; text-align: center; vertical-align: middle; font-family: ''RoobertPRO'', ''Arial'', ''Helvetica'', ui-sans-serif, system-ui, -apple-system, ''Segoe UI'', sans-serif; font-size: 13px; font-weight: 400; line-height: 1.7; color: #7e7e7e; text-decoration-line: none">
                <!--[if mso]><i style="mso-font-width: -100%; letter-spacing: 32px; mso-text-raise: 30px;" hidden>&nbsp;</i><![endif]-->
                <span style="vertical-align: middle; color: inherit; mso-text-raise: 16px;">
      Открыть письмо в браузере</span>
                <!--[if mso]><i style="mso-font-width: -100%; letter-spacing: 32px;" hidden>&nbsp;</i><![endif]-->
              </a>
            </div>

            <div role="separator" style="line-height: 24px; mso-line-height-alt: 24px">&zwj;</div>
            <div role="separator" style="background-color: #EBEBEB; height: 1px; line-height: 1px; margin: 8px 0;">&zwj;</div>
            <div role="separator" style="line-height: 10px; mso-line-height-alt: 10px">&zwj;</div>
            <table cellpadding="0" cellspacing="0" style="margin-left: auto; margin-right: auto; width: auto; border-style: none;" role="none">
              <tr style="vertical-align: middle;">
                <td style="padding-left: 8px; padding-right: 8px; vertical-align: middle;">
                  <table cellpadding="0" cellspacing="0" style="width: auto; border-style: none;" role="none">
                    <tr style="vertical-align: middle;">

                        <td class="xxs-hidden">
                          <a href="https://emcd.io/pool/dashboard?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;"> <img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718356980754_footer-sr-thunder_01J0B0XEAKQ8QVPR1WB7G2P1TT.png" border="0" alt width="16" height="16" style="max-width: 100%; vertical-align: middle; line-height: 1; margin-right: 8px">
                          </a>
                        </td>

                      <td>
                        <a href="https://emcd.io/pool/dashboard?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;">
                          <span class="dark-mode-text-_1E1E1E" style="font-size: 12px; line-height: 12px; color: #1E1E1E">
                Майнинг</span>
                        </a>
                      </td>
                    </tr>
                  </table>
                </td>
                <td style="padding-left: 8px; padding-right: 8px; vertical-align: middle;">
                  <table cellpadding="0" cellspacing="0" style="width: auto; border-style: none;" role="none">
                    <tr style="vertical-align: middle;">

                        <td class="xxs-hidden">
                          <a href="https://emcd.io/wallets?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;"> <img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718356985346_footer-sr-twoarrow_01J0B0XJTEMH8772NQB74BS2AY.png" border="0" alt width="16" height="16" style="max-width: 100%; vertical-align: middle; line-height: 1; margin-right: 8px;">
                          </a>
                        </td>

                      <td>
                        <a href="https://emcd.io/wallets?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;">
                          <span class="dark-mode-text-_1E1E1E" style="font-size: 12px; line-height: 12px; color: #1E1E1E">
                Кошелек</span>
                        </a>
                      </td>
                    </tr>
                  </table>
                </td>
                <td style="padding-left: 8px; padding-right: 8px; vertical-align: middle;">
                  <table cellpadding="0" cellspacing="0" style="width: auto; border-style: none;" role="none">
                    <tr style="vertical-align: middle;">

                        <td class="xxs-hidden">
                          <a href="https://emcd.io/deposits?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;"> <img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718356974038_footer-sr-arrow-down_01J0B0X7SJVMRNQ6QAZ37TXW0T.png" border="0" alt width="16" height="16" style="max-width: 100%; vertical-align: middle; line-height: 1; margin-right: 8px;">
                          </a>
                        </td>

                      <td>
                        <a href="https://emcd.io/deposits?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;">
                          <span class="dark-mode-text-_1E1E1E" style="font-size: 12px; line-height: 12px; color: #1E1E1E">Coinhold</span>
                        </a>
                      </td>
                    </tr>
                  </table>
                </td>
                <td style="padding-left: 8px; padding-right: 8px; vertical-align: middle;">
                  <table cellpadding="0" cellspacing="0" style="width: auto; border-style: none;" role="none">
                    <tr style="vertical-align: middle;">

                        <td class="xxs-hidden">
                          <a href="https://emcd.io/p2p?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;"> <img src="https://userimg-assets-eu.customeriomail.com/images/client-env-107315/1718356977464_footer-sr-coin-dollar_01J0B0XB3WCNB24H8YB39BZFR5.png" border="0" alt width="16" height="16" style="max-width: 100%; vertical-align: middle; line-height: 1; margin-right: 8px;">
                          </a>
                        </td>

                      <td>
                        <a href="https://emcd.io/p2p?link_id={% cio_link_id %}" class="untracked" style="cursor: pointer; text-decoration-line: none;">
                          <span class="dark-mode-text-_1E1E1E" style="font-size: 12px; line-height: 12px; color: #1E1E1E">P2P</span>
                        </a>
                      </td>
                    </tr>
                  </table>
                </td>
              </tr>
            </table>
            <div role="separator" style="line-height: 24px; mso-line-height-alt: 24px">&zwj;</div>
          </div>
          <!-- /end footer: на светлом фоне, краткий -->
        </td>
      </tr>
    </table>
    <div role="separator" style="line-height: 40px; mso-line-height-alt: 40px">&zwj;</div>
  </div>
</body>
</html>',
       'payout report',
       'ru',
       'Твой запрос на выгрузку данных о выплатах',
       NULL
WHERE NOT EXISTS(
        SELECT type, language FROM email_templates WHERE type = 'payout report' AND language = 'ru'
    );