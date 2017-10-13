<html>
    <head>
        <meta charset="utf-8">
        <title>{{imp[0]}} - {{imp[1]}}</title>
        <link rel="stylesheet" href="/views/yousoro/css/primer.css">
        <link rel="stylesheet" href="/views/yousoro/css/style.css">
        <!-- 필수 CSS, JS -->
        <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/8.7/styles/default.min.css">
        <script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/8.7/highlight.min.js"></script>
        <script>hljs.initHighlightingOnLoad();</script>
        <link rel="stylesheet" href="/views/yousoro/css/awesome/font-awesome.min.css">
        <script type="text/x-mathjax-config">MathJax.Hub.Config({tex2jax: {inlineMath: [['$','$'], ['\\(','\\)']]}});</script>
        <script type="text/javascript" async src="https://cdnjs.cloudflare.com/ajax/libs/mathjax/2.7.1/MathJax.js?config=TeX-AMS_CHTML"></script>
        <script src='https://www.google.com/recaptcha/api.js'></script>
        <style>{{!imp[4]}}</style>
        <script>{{!imp[5]}}</script>
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
        <meta name="twitter:creator" content="@{{imp[1]}}">
        <meta name="twitter:title" content="{{imp[0]}}">
        <meta name="twitter:site" content="@{{imp[1]}}">
        <meta name="twitter:card" content="summary">
        <link rel="shortcut icon" href="/views/yousoro/img/on.ico">
        <meta name="viewport" content="width=device-width, initial-scale=1">
    </head>
    <body>
        <br>
        <div class="one-fifth column">
            <div id="top">
                <a href="/" id="logo">{{imp[1]}}</a>
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
                        % if(imp[3] == 1):
                            <i class="fa fa-user" aria-hidden="true"></i>
                        % elif(imp[3] == 0):
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
                        <input placeholder="검색" class="form-control" name="search" type="text">
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
                % if(not menu == 0):
                    <div id="tool">
                        <nav class="menu">
                            % for sub_d in menu:
                                % if(sub_d[1] == 1):
                                    <a class="menu-item" href="/{{sub_d[0]}}" id="open">토론</a>
                                % elif(sub_d[1] == 0):
                                    <a class="menu-item" href="/{{sub_d[0]}}">토론</a>
                                % else:
                                    <a class="menu-item" href="/{{sub_d[0]}}">{{sub_d[1]}}</a>
                                % end
                            % end
                        </nav>
                        <br>
                        <span id="edit_time">
                            % if(not imp[7] == 0):
                                최근 수정 : {{imp[7]}}
                            % end
                        </span>
                    </div>
                % end
                <h1 class="title">
                    {{imp[0]}}
                    % if(not imp[6] == 0):
                        <sub>{{imp[6]}}</sub>
                    % end
                </h1>
                <br>
                {{!data}}
                <hr id="last">
                <p>
                    {{!imp[2]}}
                </p>
                <div id="powered">
                    <a href="https://github.com/2DU/openNAMU"><img src="/views/yousoro/img/on2.png" width="100px"></a>
                </div>
            </div>
        </div>
    </body>
</html>
