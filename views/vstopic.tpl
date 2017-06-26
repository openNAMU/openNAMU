<!-- 토론장 -->
% rebase('base.tpl')
<div id="tool">
    <nav class="menu">
        <a class="menu-item selected" href="#" onclick="return false">
            {{sub}}
        </a>
        <a class="menu-item" href="/topic/{{page}}">토론 목록</a>
    </nav>
</div>
<h1 class="title">
    {{title}}
    <sub> ({{sub}})</sub>
</h1>
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
        <span>비 로그인 상태입니다. 비 로그인으로 작업 시 아이피가 토론에 기록됩니다.</span>
    % end
% end