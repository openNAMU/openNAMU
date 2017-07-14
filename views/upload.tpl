<!-- 업로드 -->
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
<div>
    <form action="" method=post enctype=multipart/form-data>
        <input class="btn" type="file" name="file">
        <input class="btn btn-primary" type="submit" value="Upload">
    </form>
</div>
<span>{{number}}MB 이하 파일만 업로드 가능하고 확장자는 jpg, png, gif, jpeg만 가능합니다.</span>