<!-- 많은 문서 삭제 -->
% rebase('base.tpl')
<div id="tool">
    <nav class="menu">
        <a class="menu-item selected" href="#" onclick="return false">
            {{title}}
        </a>
        <a class="menu-item" href="javascript:history.back(-1);">뒤로</a>
    </nav>
</div>
<h1 class="title">
    {{title}}
</h1>
% if(login == 0):
    <br>
    <span>비 로그인 상태입니다. 비 로그인으로 작업 시 아이피가 역사에 기록됩니다.</span>
    <br>
    <br>
% end
<br>
<span>문서명 A
<br>
문서명 B
<br>
문서명 C</span>
<br>
<br>
<span>이런 식으로 기제 하시오</span>
<br>
<br>
<form id="usrform" method="POST" action="/mdel">
    <textarea rows="30" cols="100" name="content" form="usrform"></textarea>
    <input style="margin-top:10px;width:100%;" class="form-control input-sm" name="send" type="text">
    <br>
    <br>
    <div class="form-actions">
        <button class="btn btn-primary" type="submit">삭제</button>
    </div>
</form>