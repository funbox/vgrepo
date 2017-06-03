<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>Storage index</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="https://unpkg.com/basscss@8.0.2/css/basscss.min.css" rel="stylesheet">
    <link href="https://fonts.googleapis.com/css?family=Open+Sans:400,400i,700&amp;subset=cyrillic,cyrillic-ext,latin-ext" rel="stylesheet">
    <style>
      *,*:before,*:after{box-sizing:inherit}
      html{box-sizing:border-box;-ms-text-size-adjust:100%;-webkit-text-size-adjust:100%}
      body{margin:0;font-family:"Open Sans","sans-serif",system,-apple-system,BlinkMacSystemFont,"Helvetica Neue","Lucida Grande";line-height:1.6}
      h1,h2,h3,h4,h5,h6{margin-top:0;line-height:1.4;font-weight:800}
      .mb0 a,.mb0 a:visited{color:#000;text-decoration:none}
      h1 a{color:#000}p{margin-top:0}
      a:not(h1 a){color:#00f}
      a:not(h1 a):visited{color:purple}
      div.shadowed{color: #999;}
    </style>
</head>

<body class="p3">
<main class="max-width-3" role="main">
<header class="mb4">
<h1 class="mb0"><a href="/" class="text-decoration-none">Storage index</a></h1>
</header>

<section class="mb4">
{{ range .Repositories }}
<h2 class="mb2 h2">{{ .Name }} ¬</h2>
{{ range $i, $v := .Versions }}
<h3 class="mb2 h4">Release v{{ $v.Version }}</h3>
<ul class="mt0 list-reset">
  {{ range .Providers }}
  <li class="mb1">
    <div class="clearfix">
      <div class="col col-3"><a href="{{ .URL }}">{{ .Name }} {{ $v.Version }}</a></div>
      <div class="col col-4 shadowed">{{ .Checksum }}</div>
    </div>
  </li>
  {{ end }}
</ul>
{{ end }}
</section>
{{ end }}

</main>

</body>
</html>
