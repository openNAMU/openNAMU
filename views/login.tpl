<!-- 로그인 -->
% rebase('base.tpl')
<div id="tool">
    <nav class="menu">
        <a class="menu-item selected" href="#" onclick="return false">
            {{title}}
        </a>
        <a class="menu-item" href="/user">사용자</a>
    </nav>
</div>
<h1 class="title">
    {{title}}
</h1>
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