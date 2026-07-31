<!doctype html>
<html>
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>{{ .Campaign.Subject }}</title>
</head>
<body style="margin:0;padding:16px;font-family:Arial,Helvetica,sans-serif;font-size:14px;line-height:1.5;color:#000;background:#fff;">
    {{ template "content" . }}

    <p style="margin-top:24px;font-size:11px;color:#888;">
        <a href="{{ UnsubscribeURL }}" style="color:#888;">{{ L.T "email.unsub" }}</a>
    </p>
    {{ TrackView }}
</body>
</html>
