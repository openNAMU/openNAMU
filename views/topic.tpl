<!-- 토론 목록 -->
% rebase('base.tpl')
<div id="tool">
    <nav class="menu">
        <a class="menu-item selected" href="#" onclick="return false">
            {{sub}}
        </a>
        % if(defined('list')):
            <a class="menu-item" href="/w/{{get('page', '')}}">문서</a>
        % else:
            <a class="menu-item" href="/topic/{{page}}">토론 목록</a>
        % end
    </nav>
</div>
<h1 class="title">
    {{title}}
    <sub> ({{sub}})</sub>
</h1>
<form id="usrform" style="margin-top: -30px;" method="POST" action="/topic/{{page}}">
    {{!plus}}
    % if(defined('list')):
        <br>
        <a href="/topic/{{page}}/close">(닫힘)</a> <a href="/topic/{{page}}/agree">(합의)</a>
        <br>
        <br>
        <input class="form-control" name="topic" style="width: 100%">
        <br>
        <br>
        <button class="btn btn-primary" type="submit">만들기</button>
    % end
</form>