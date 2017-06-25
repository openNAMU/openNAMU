<!-- 관리자 -->
% rebase('base.tpl')
<div id="tool">
    <nav class="menu">
        <a class="menu-item selected" href="#" onclick="return false">
            {{sub}}
        </a>
        <a class="menu-item" href="javascript:history.back(-1);">뒤로</a>
    </nav>
</div>
<h1 class="title">
    {{title}}
    <sub> ({{sub}})</sub>
</h1>
<form id="usrform" method="POST" action="/admin/{{page}}">
    % if(now == '권한 부여'):
        <select name="select">
            {{!datalist}}
        </select>
        <br>
        <br>
    % end
    <button class="btn btn-primary" type="submit">{{now}}</button>
</form>