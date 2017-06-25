<!-- 차단 -->
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
<form id="usrform" method="POST" action="/ban/{{page}}">
    % if(now == '차단' or now == '대역 차단'):
        <input class="form-control" name="end" style="width: 100%">
        <br>
        <br>
        <span>아무것도 안 적으면 무기한 차단입니다.</span>
        <br>
        <br>
        <span>차단 일 지정시 형식은 YYYY-MM-DD로 기록해야 합니다. (예시: 2017-01-20, 2017-10-15)</span>
        <br>
        <br>
        <span>지금 시각은 {{today}} 입니다.</span>
        <br>
        <br>
        <input class="form-control" name="why" style="width: 100%">
        <br>
        <br>
        <span>사유를 쓰는 곳입니다.</span>
        % if(defined('allif')):
            <br>
            <br>
            <input type="checkbox" name="band">
                <span>대역 차단</span>
            </input>
        % end
        <br>
        <br>
    % end
    <button class="btn btn-primary" type="submit">{{now}}</button>
</form>