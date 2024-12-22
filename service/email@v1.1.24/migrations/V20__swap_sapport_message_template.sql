INSERT INTO public.email_templates(whitelabel_id, template, type, language, subject)
VALUES ('00000000-0000-0000-0000-000000000000', '<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Support Request</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #F4F4F4;
            padding: 20px;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            background-color: #fff;
            padding: 20px;
            border-radius: 8px;
            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
        }
        h2 {
            color: #333;
        }
        p {
            line-height: 1.6;
            color: #555;
        }
        .footer {
            margin-top: 20px;
            font-size: 0.9em;
            color: #888;
        }
    </style>
</head>
<body>
    <div class="container">
        <h2>Support Request</h2>
        <p><strong>Name:</strong> {{.Name}}</p>
        <p><strong>Email:</strong> {{.Email}}</p>
        <p><strong>Message:</strong></p>
        <p>{{.Text}}</p>
    </div>
</body>
</html>
', 'swap support message', 'ru', 'Support Request');