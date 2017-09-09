<html>
    <head>
        <title>{{title}} - {{logo}}</title>
        <link rel="stylesheet" href="/views/css/primer.css">
        <link rel="stylesheet" href="/views/css/style.css">
        <!-- 필수 CSS, JS -->
        <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/8.7/styles/default.min.css">
        <script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/8.7/highlight.min.js"></script>
        <script>hljs.initHighlightingOnLoad();</script>
        <link rel="stylesheet" href="/views/css/awesome/font-awesome.min.css">
        <script type="text/x-mathjax-config">MathJax.Hub.Config({tex2jax: {inlineMath: [['$','$'], ['\\(','\\)']]}});</script>
        <script type="text/javascript" async src="https://cdnjs.cloudflare.com/ajax/libs/mathjax/2.7.1/MathJax.js?config=TeX-AMS_CHTML"></script>
        <script>
            function folding(num) {
                var fol = document.getElementById('folding_' + num);
                if(fol.style.display == 'inline-block' || fol.style.display == 'block') {
                    fol.style.display = 'none';
                } else {
                    if(num % 3 == 0) {
                        fol.style.display = 'block';
                    } else {
                        fol.style.display = 'inline-block';
                    }
                }
            }
        </script>
        <!-- 필수 부분 끝 -->
        <meta name="twitter:creator" content="@{{logo}}">
        <meta name="twitter:title" content="{{title}}">
        <meta name="twitter:site" content="@{{logo}}">
        <meta name="twitter:card" content="summary">
        <link rel="shortcut icon" href="/views/img/on.ico">
        <meta name="viewport" content="width=device-width, initial-scale=1">
        <style>{{!get('custom_css', '')}}</style>
        <script>{{!get('custom_js', '')}}</script>
    </head>
    <body>
        <br>
        <div class="one-fifth column">
            <div id="top">
                <a href="/" id="logo">{{logo}}</a>
                <div>
                    <a href="/recent_changes" id="RecentChanges">
                        <i class="fa fa-refresh" aria-hidden="true"></i>
                        <span id="is_mobile">최근 변경</span>
                    </a>
                    <a href="/recent_discuss" id="RecentChanges">
                        <i class="fa fa-comment" aria-hidden="true"></i>
                        <span id="is_mobile">최근 토론</span>
                    </a>
                    <a href="/random" id="log">
                        <i class="fa fa-random" aria-hidden="true"></i>
                    </a>
                    <a href="/user" id="log">
                        % if(login == 1):
                            <i class="fa fa-user" aria-hidden="true"></i>
                        % elif(login == 0):
                            <i class="fa fa-user-times" aria-hidden="true"></i>
                        % else:
                            <i class="fa fa-user-secret" aria-hidden="true"></i>
                        % end
                    </a>
                    <a href="/other" id="log">
                        <i class="fa fa-cogs" aria-hidden="true"></i>
                    </a>
                </div>
                <form method="POST" action="/search" id="search">
                    <div class="input-group">
                        <input class="form-control" name="search" type="text">
                        <span class="input-group-button">
                            <button id="goto" class="btn" formaction="/goto"><i class="fa fa-share" aria-hidden="true"></i></button>
                        </span>
                        <span class="input-group-button">
                            <button class="btn"><i class="fa fa-search" aria-hidden="true"></i></button>
                        </span>
                    </div>
                </form>
            </div>
        </div>
        <div class="scroll-buttons">
            <a class="scroll-toc" href="#toc"><i class="fa fa-list-alt" aria-hidden="true"></i></a>
            <a class="scroll-button" href="#" id="left"><i class="fa fa-arrow-up" aria-hidden="true"></i></a>
            <a class="scroll-bottom" href="#powered" id="right"><i class="fa fa-arrow-down" aria-hidden="true"></i></a>
        </div>
        <div id="all_contect">
            <div class="four-fifths column">
                {{!base}}
                <hr id="last">
                <p>
                    {{!license}}
                </p>
                <div id="powered">
                    <a href="https://github.com/2DU/openNAMU"><img src="/views/img/on2.png" width="100px"></a>
                </div>
            </div>
        </div>
    </body>
</html>
