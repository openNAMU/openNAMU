<!-- 읽기 -->
% rebase('base.tpl')
<div id="tool">
    <nav class="menu">
        <a class="menu-item selected" href="#" onclick="return false">
            문서
        </a>
        % if(data_none == 1):
            <a class="menu-item" href="/edit/{{page}}">생성</a>
        % else:
            <a class="menu-item" href="/edit/{{page}}">수정</a>
        % end
        <a class="menu-item" id="{{topic}}" href="/topic/{{page}}">토론</a>
        % if(data_none == 0):
            <a class="menu-item" href="/delete/{{page}}">삭제</a>
        % end
        <a class="menu-item" href="/move/{{page}}">이동</a>
        % if(data_none == 0):
            <a class="menu-item" href="/raw/{{page}}">원본</a>
        % end
        <a class="menu-item" href="/history/{{page}}">역사</a>
        <a class="menu-item" href="/xref/{{page}}">역링크</a>
        % if(redirect):
            <a class="menu-item" href="/w/{{page}}">넘기기</a>
        % end
        % if(data_none == 0 and admin == "ACL"):
            <a class="menu-item" href="/acl/{{page}}">ACL</a>
        % end
        % if(not uppage == 0):
            <a class="menu-item" href="/w/{{uppage}}">상위</a>
        % end
        % if(down == 1):
            <a class="menu-item" style="{{down}}" href="/down/{{page}}">하위</a>
        % end
    </nav>
</div>
<h1 class="title">
    {{title}}
    % if(defined('acl')):
        <sub> {{acl}}</sub>
    % end
    % if(not sub == 0):
        <sub> ({{sub}})</sub>
    % end
</h1>
% if(redirect):
    <li style="margin-top: -20px;"><a href="/w/{{redirect}}/from/{{page}}">{{!redirect or 'None'}}</a>에서 넘어 왔습니다.</li>
    <br>
% end
<div>
    {{!data}}
</div>