<!-- ACL -->
% rebase('base.tpl')
<div id="tool">
    <nav class="menu">
        <a class="menu-item selected" href="#" onclick="return false">
            {{sub}}
        </a>
        <a class="menu-item" href="/w/{{get('page', '')}}">문서</a>
    </nav>
</div>
<h1 class="title">
    {{title}}
    <sub> ({{sub}})</sub>
</h1>
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