<html>
    <head>
        <title>{{title}} - {{logo}}</title>
        <link rel="stylesheet" href="/static/primer.css">
        <link rel="stylesheet" href="/static/style.css">
        <link rel="stylesheet" href="/static/font-awesome/css/font-awesome.min.css">
        <script type="text/x-mathjax-config">
          MathJax.Hub.Config({tex2jax: {inlineMath: [['$','$'], ['\\(','\\)']]}});
        </script>
        <script type="text/javascript" async
            src="https://cdnjs.cloudflare.com/ajax/libs/mathjax/2.7.1/MathJax.js?config=TeX-AMS_CHTML">
        </script>
        <meta name="twitter:creator" content="@{{logo}}">
        <meta name="twitter:title" content="{{title}}">
        <meta name="twitter:site" content="@{{logo}}">
        <meta name="twitter:card" content="summary">
        <meta name="twitter:description" content="{{get('data', 'None')}}" />
        <link rel="shortcut icon" href="/static/images/on.ico">
        <meta name="viewport" content="width=device-width, initial-scale=1">
    <style>
        {{get('custom', '')}}
    </style>
    </head>
    <body>
        <br>
        <div class="one-fifth column">
            <div id="top">
                <a href="/" id="logo">{{logo}}</a>
                <div>
                    <a href="/recentchanges" id="RecentChanges">
                        <i class="fa fa-refresh" aria-hidden="true"></i>
                        <span id="is_mobile">최근 변경</span>
                    </a>
                    <a href="/recentdiscuss" id="RecentChanges">
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
            <div id="left_bar">
                <a href="#">맨 위로</a>
                <br>
                <br>
                {{!get('left', '')}}
            </div>
            <div class="four-fifths column">
                <div id="tool">
                    <nav class="menu">
                        <a class="menu-item selected" href="#" onclick="return false">
                            % if(defined('sub')):
                                {{sub}}
                            % elif(defined('tn')):
                                % if(tn == 1):
                                    문서
                                % else:
                                    기타
                                % end
                            % else:
                                기타
                            % end
                        </a>
                        % if(defined('tn')):
                            % if(tn == 1):
                                % if(defined('data_none')):
                                    <a class="menu-item" href="/edit/{{page}}">생성</a>
                                % else:
                                    <a class="menu-item" href="/edit/{{page}}">수정</a>
                                % end
                                <a class="menu-item" id="{{topic}}" href="/topic/{{page}}">토론</a>
                                % if(not defined('data_none')):
                                    <a class="menu-item" href="/delete/{{page}}">삭제</a>
                                % end
                                <a class="menu-item" href="/move/{{page}}">이동</a>
                                % if(not defined('data_none')):
                                    <a class="menu-item" href="/raw/{{page}}">Raw</a>
                                % end
                                <a class="menu-item" href="/history/{{page}}/n/1">역사</a>
                                <a class="menu-item" href="/backlink/{{page}}/n/1">역링크</a>
                                % if(redirect):
                                    <a class="menu-item" href="/w/{{page}}">넘기기</a>
                                % end
                                % if(admin == "ACL"):
                                    <a class="menu-item" href="/acl/{{page}}">ACL</a>
                                % end
                                % if(uppage):
                                    <a class="menu-item" style="{{style}}" href="/w/{{uppage}}">상위</a>
                                % end
                            % elif(tn == 6 or tn == 13 or not tn and not sub == '역링크' or tn == 3):
                                <a class="menu-item" href="javascript:history.back(-1);">뒤로</a>
                            % elif(tn == 10 and not list or tn == 11):
                                <a class="menu-item" href="/topic/{{page}}">토론 목록</a>
                            % else:
                                <a class="menu-item" href="/w/{{get('page', '')}}">문서</a>
                            % end
                        % else:
                            <a class="menu-item" href="/w/{{get('page', '')}}">문서</a>
                        % end
                    </nav>
                </div>
                <h1 class="title">
                    {{title}}
                    % if(defined('acl')):
                        <sub> {{acl}}</sub>
                    % end
                    % if(defined('sub')):
                        <sub> ({{sub}})</sub>
                    % end
                </h1>
                % if(defined('tn')):
                    % if(tn == 1):
                        % if(redirect):
                            <li style="margin-top: -20px;"><a href="/w/{{redirect}}/from/{{page}}">{{!redirect or 'None'}}</a>에서 넘어 왔습니다.</li>
                            <br>
                        % end
                        <div>
                            {{!data}}
                        </div>
                    % elif(tn == 2):
                        % if(login == 0):
                            <br>
                            <span>비 로그인 상태입니다. 비 로그인으로 작업 시 아이피가 역사에 기록됩니다.</span>
                            <br>
                            <br>
                        % end
                        % if(defined('section')):
                            <form id="usrform" name="f1" method="POST" action="/edit/{{page}}/section/{{number}}">
                                <textarea rows="30" cols="100" name="content" form="usrform">{{data}}</textarea>
                                % if(defined('preview')):
                                    <textarea style="display:none;" rows="30" cols="100" name="otent" form="usrform">{{odata}}</textarea>
                                % else:
                                    <textarea style="display:none;" rows="30" cols="100" name="otent" form="usrform">{{data}}</textarea>
                                % end
                                <input name="send" style="margin-top:10px;width:100%" type="text">
                                <br>
                                <br>
                                <div class="form-actions">
                                    <button class="btn btn-primary" type="submit" onclick="f1.action='/edit/{{page}}/section/{{number}}';">저장</button>
                                    <button class="btn" type="submit" onclick="f1.action='/preview/{{page}}/section/{{number}}';">미리보기</button>
                                </div>
                            </form>
                        % else:
                            <form id="usrform" name="f1" method="POST" action="/edit/{{page}}">
                                <textarea rows="30" cols="100" name="content" form="usrform">{{data}}</textarea>
                                <input name="send" style="margin-top:10px;width:100%" type="text">
                                <br>
                                <br>
                                <div class="form-actions">
                                    <button class="btn btn-primary" type="submit" onclick="f1.action='/edit/{{page}}';">저장</button>
                                    <button class="btn" type="submit" onclick="f1.action='/preview/{{page}}';">미리보기</button>
                                </div>
                            </form>      
                        % end
                        % if(defined('preview')):
                            {{!enddata}}
                        % end
                    % elif(tn == 3):
                        {{!rows}}
                    % elif(tn == 5):
                        <form class="usrform" method='POST' action='/history/{{page}}/n/1'>
                        <select name="a">
                            {{!select}}
                        </select>
                        <select name="b">
                            {{!select}}
                        </select>
                        <button class="btn btn-primary" type='submit'>리비전 비교</button>
                        <br>
                        <br>
                        {{!rows}}
                    % elif(tn == 6):
                        <div>
                            {{!data}}
                        </div>
                    % elif(tn == 8):
                        <form id="usrform" method="POST" action="/delete/{{page}}">
                            {{plus}}
                            <br>
                            <br>
                            <button class="btn btn-primary" type="submit">삭제</button>
                        </form>
                        % if(login == 0):
                            <span>비 로그인 상태입니다. 비 로그인으로 작업 시 아이피가 역사에 기록됩니다.</span>
                        % end
                    % elif(tn == 9):
                        <form id="usrform" method="POST" action="/move/{{page}}">
                            {{plus}}
                            <br>
                            <br>
                            <input class="form-control input-sm" value="{{title}}" name="title" type="text">
                            <br>
                            <br>
                            <button class="btn btn-primary" type="submit">이동</button>
                        </form>
                        % if(login == 0):
                            <span>비 로그인 상태입니다. 비 로그인으로 작업 시 아이피가 역사에 기록됩니다.</span>
                        % end
                    % elif(tn == 10):
                        <form id="usrform" style="margin-top: -30px;" method="POST" action="/topic/{{page}}">
                            {{!plus}}
                            % if(list == 1):
                                <br>
                                <a href="/topic/{{page}}/close">(닫힌 토론)</a> <a href="/topic/{{page}}/agree">(합의된 토론)</a>
                                <br>
                                <br>
                                <input class="form-control" name="topic" style="width: 100%">
                                <br>
                                <br>
                                <button class="btn btn-primary" type="submit">새토론</button>
                            % end
                        </form>
                    % elif(tn == 11):
                        <h2 style="margin-top: -15px;">{{toron}}</h2>
                        <br>
                        {{!rows}}
                        % if(not ban == 1):
                            <a id="reload" href="javascript:window.location.reload(true);">(갱신)</a> <a href="#reload">(갱신 전에 누르시오)</a>
                            <form id="usrform" style="{{style}}" method="POST" action="/topic/{{page}}/sub/{{suburl}}">
                                <br>
                                <textarea rows="10" cols="100" name="content" form="usrform"></textarea>
                                <br>
                                <br>
                                <button class="btn btn-primary" type="submit">전송</button>
                            </form>
                            % if(login == 0 and style == ''):
                                <span>비 로그인 상태입니다. 비 로그인으로 작업 시 아이피가 역사에 기록됩니다.</span>
                            % end
                        % end
                    % elif(tn == 13):
                        <form id="usrform" method="POST" action="/revert/{{page}}/r/{{r}}">
                            {{plus}}
                            <br>
                            <br>
                            <button class="btn btn-primary" type="submit">되돌리기</button>
                        </form>
                        % if(login == 0):
                            <span>비 로그인 상태입니다. 비 로그인으로 작업 시 아이피가 역사에 기록됩니다.</span>
                        % end
                    % elif(tn == 15):
                        % if(title == '회원가입'):
                            <form id="usrform" method="POST" action="/register">
                        % elif(title == '비밀번호 변경'):
                            <form id="usrform" method="POST" action="/change">
                        % else:
                            <form id="usrform" method="POST" action="/login">
                        % end
                            <span>아이디</span>
                            <br>
                            <br>
                            <input name="id" type="text">
                            <br>
                            <br>
                            <span>
                                % if(title == '비밀번호 변경'):
                                    현재 
                                % end
                                비밀번호
                            </span>
                            <br>
                            <br>
                            <input name="pw" type="password">
                            <br>
                            <br>
                            % if(not title == '로그인'):
                                <span>
                                    % if(title == '비밀번호 변경'):
                                        변경할 비밀번호
                                    % else:
                                        재 확인
                                    % end
                                </span>
                                <br>
                                <br>
                                <input name="pw2" type="password">
                                <br>
                                <br>
                                % if(title == '비밀번호 변경'):
                                    <span>재 확인</span>
                                    <br>
                                    <br>
                                    <input name="pw3" type="password">
                                    <br>
                                    <br>
                                % end
                            % end
                            <button class="btn btn-primary" type="submit">{{enter}}</button>
                        </form>
                    % elif(tn == 16):
                        <form id="usrform" method="POST" action="/ban/{{page}}">
                            % if(now == '차단' or now == '대역 차단'):
                                <input class="form-control" name="end" style="width: 100%">
                                <br>
                                <br>
                                <span>아무것도 안 적으면 무기한 차단입니다.</span>
                                <br>
                                <br>
                                <span>차단 일 지정시 형식은 YYYY-MM-DD로 기록해야 합니다. (예시: 2017-01-20, 2017-10-15)</span>
                                <br>
                                <br>
                                <span>지금 시각은 {{today}} 입니다.</span>
                                <br>
                                <br>
                                <input class="form-control" name="why" style="width: 100%">
                                <br>
                                <br>
                                <span>사유를 쓰는 곳입니다.</span>
                                % if(defined('allif')):
                                    <br>
                                    <br>
                                    <input type="checkbox" name="band">
                                        <span>대역 차단</span>
                                    </input>
                                % end
                                <br>
                                <br>
                            % end
                            <button class="btn btn-primary" type="submit">{{now}}</button>
                        </form>
                    % elif(tn == 18):
                        <form id="usrform" method="POST" action="/admin/{{page}}">
                            % if(now == '권한 부여'):
                                <select name="select">
                                    <option value="admin" selected="selected">관리자</option>
                                    <option value="owner">소유자</option>
                                </select>
                                <br>
                                <br>
                            % end
                            <button class="btn btn-primary" type="submit">{{now}}</button>
                        </form>
                    % elif(tn == 19):
                        <br>
                        <span>{{now}}</span>
                        <br>
                        <br>
                        <form id="usrform" method="POST" action="/acl/{{page}}">
                            <select name="select">
                                <option value="admin" selected="selected">관리자만</option>
                                <option value="user">유저 이상</option>
                                <option value="normal">일반</option>
                            </select>
                            <br>
                            <br>
                            <button class="btn btn-primary" type="submit">ACL 변경</button>
                        </form>
                    % elif(tn == 21):
                        <div>
                            <form action="" method=post enctype=multipart/form-data>
                                <input type=file name=file>
                                <input type=submit value=Upload>
                            </form>
                        </div>
                        <span>{{number}}MB 이하 파일만 업로드 가능하고 확장자는 jpg, png, gif, jpeg만 가능합니다.</span>
                     % else:
                        <div>
                            {{!data}}
                        </div>
                    % end
                % else:
                    <div>
                        {{!data}}
                    </div>
                % end
                <hr id="last">
                <p>
                    {{!license}}
                </p>
                <div id="powered">
                    <a href="https://github.com/2DU/openNAMU"><img src="/static/images/on2.png" width="100px"></a>
                </div>
            </div>
        </div>
    </body>
</html>
