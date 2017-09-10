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
</a>
</h1>
<h2 style="margin-top: -15px;">{{toron}}</h2>
<br>
{{!rows}}
<a id="reload" href="javascript:void(0);" onclick="location.href.endsWith('#reload') ?  location.reload(true) : location.href = '#reload'"><i aria-hidden="true" class="fa fa-refresh"></i></a>
% if(not ban == 1):
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