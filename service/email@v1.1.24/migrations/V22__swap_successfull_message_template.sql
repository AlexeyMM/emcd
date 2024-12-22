INSERT INTO public.email_templates(
    whitelabel_id, template, type, language, subject)
VALUES ('00000000-0000-0000-0000-000000000000', '<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Успешный обмен криптовалюты</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #F7F7F7;
            color: #333;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            background-color: #FFFFFF;
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
            color: #1A73E8;
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
            <p>Сообщаем вам, что операция по обмену криптовалюты была успешно завершена.</p>
            <div class="details">
                <p><strong>Детали обмена:</strong></p>
                <p>ID Транзакции: <strong>{{.SwapID}}</strong></p>
                <p>Дата и время транзакции: <strong>{{.ExecutionTime}}</strong></p>
                <p>Сумма отправки: <strong>{{.From}}</strong></p>
                <p>Сумма получения: <strong>{{.To}}</strong></p>
                <p>Кошелёк получения: <strong>{{.Address}}</strong></p>
            </div>
            <p>Средства были успешно зачислены на указанный вами кошелек, и вы можете проверить баланс.</p>
            <p>Спасибо за использование нашего сервиса!</p>
        </div>
    </div>
</body>
</html>', 'swap successful message', 'ru', 'Обмен успешно выполнен'),
       ('00000000-0000-0000-0000-000000000000', '<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Successful Cryptocurrency Exchange</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #F7F7F7;
            color: #333;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            background-color: #FFFFFF;
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
            color: #1A73E8;
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
            <p>We are pleased to inform you that your cryptocurrency exchange transaction has been successfully completed.</p>
            <div class="details">
                <p><strong>Exchange Details:</strong></p>
                <p>Transaction ID: <strong>{{.SwapID}}</strong></p>
                <p>Date and Time of Transaction: <strong>{{.ExecutionTime}}</strong></p>
                <p>Amount Sent: <strong>{{.From}}</strong></p>
                <p>Amount Received: <strong>{{.To}}</strong></p>
                <p>Receiving Wallet: <strong>{{.Address}}</strong></p>
            </div>
            <p>The funds have been successfully credited to your designated wallet, and you can verify your balance.</p>
            <p>Thank you for using our service!</p>
        </div>
    </div>
</body>
</html>', 'swap successful message', 'en', 'Successful swap execution');