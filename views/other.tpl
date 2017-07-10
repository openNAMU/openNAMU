<!-- 나머지 -->
% rebase('base.tpl')
<div id="tool">
    <nav class="menu">
        <a class="menu-item selected" href="#" onclick="return false">
            % if defined('sub'):
                {{sub}}
            % else:
                {{title}}
            % end
        </a>
        <a class="menu-item" href="javascript:history.back(-1);">뒤로</a>
    </nav>
</div>
<h1 class="title">
    {{title}}
    % if defined('sub'):
         ({{sub}})
    % end
</h1>
<div>
    % if(defined('select')):
        <form class="usrform" method='POST' action='/history/{{page}}'>
            <select name="a">
                {{!select}}
            </select>
            <select name="b">
                {{!select}}
            </select>
            <button class="btn btn-primary" type='submit'>리비전 비교</button>
        </form>
    % end
    {{!data}}
</div>