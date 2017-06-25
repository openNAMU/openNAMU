<!-- 편집 -->
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