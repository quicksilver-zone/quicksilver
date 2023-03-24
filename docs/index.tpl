<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="utf-8" />
        <title>{{ .Title }}</title>
        <link rel="stylesheet" type="text/css" href="//unpkg.com/swagger-ui-dist@latest/swagger-ui.css" />
        <link rel="icon" type="image/png" href="//unpkg.com/swagger-ui-dist@latest/favicon-16x16.png" />
    </head>
    <body>
        <div id="swagger-ui"></div>

        <script src="//unpkg.com/swagger-ui-dist@latest/swagger-ui-bundle.js"></script>
        <script>
            window.onload = function() {
              window.ui = SwaggerUIBundle({
                url: {{ .URL }},
                dom_id: "#swagger-ui",
                deepLinking: true,
                layout: "BaseLayout",
              });
            }
        </script>
    </body>
</html>
