INSERT INTO public.email_templates(whitelabel_id, template, type, language, subject)
VALUES ('00000000-0000-0000-0000-000000000000', '<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Email Notification</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f7f7f7;
            color: #333;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            background-color: #ffffff;
            border-radius: 5px;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
        }
        .header {
            font-size: 18px;
            font-weight: bold;
            color: #333;
        }
        .content {
            margin-top: 20px;
            font-size: 16px;
            line-height: 1.6;
        }
        .link {
            color: #1a73e8;
            text-decoration: none;
            font-weight: bold;
        }
        .footer {
            margin-top: 30px;
            font-size: 14px;
            color: #666;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">Уважаемый пользователь!</div>
        <div class="content">
            <p>Ваш обмен успешно инициализирован. Статус и вся информация о нем доступны по следующей ссылке:</p>
            <p><a href="{{.Link}}" class="link">{{.Link}}</a></p>
            <p>Если вы не инициировали этот обмен, пожалуйста, немедленно обратитесь в нашу службу поддержки.</p>
        </div>
        <div class="footer">
            Спасибо за использование нашего сервиса!
        </div>
    </div>
</body>
</html>', 'swap message', 'ru', 'Ваш обмен на EMCD'),
       ('00000000-0000-0000-0000-000000000000', '<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Email Notification</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f7f7f7;
            color: #333;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            background-color: #ffffff;
            border-radius: 5px;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
        }
        .header {
            font-size: 18px;
            font-weight: bold;
            color: #333;
        }
        .content {
            margin-top: 20px;
            font-size: 16px;
            line-height: 1.6;
        }
        .link {
            color: #1a73e8;
            text-decoration: none;
            font-weight: bold;
        }
        .footer {
            margin-top: 30px;
            font-size: 14px;
            color: #666;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">Dear User,</div>
        <div class="content">
            <p>Your exchange has been successfully initiated. The status and all information about it are available at the following link:</p>
            <p><a href="{{.Link}}" class="link">{{.Link}}</a></p>
            <p>If you did not initiate this exchange, please contact our support team immediately.</p>
        </div>
        <div class="footer">
            Thank you for using our service!
        </div>
    </div>
</body>
</html>', 'swap message', 'en', 'Your Swap in EMCD');