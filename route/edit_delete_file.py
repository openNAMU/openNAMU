from .tool.func import *
from .edit_delete import edit_delete

# 처음으로 차세대 코드 방법론 적용
# 앞으로 다 이렇게 작성할 예정
def edit_delete_file(name : str = 'test.jpg') -> str:
    with get_db_connect() as conn:
        curs : typing.Union[sqlite3.dbapi2.Cursor, pymysql.cursors.Cursor, None] = conn.cursor()

        ip : str = ip_check()
        if admin_check() == 0:
            return re_error('/ban')

        mime_type : typing.Union[re.Match, None] = re.search(r'([^.]+)$', name)
        if mime_type:
            mime_type = mime_type.group(1).lower()
        else:
            mime_type = 'jpg'

        file_name : str = re.sub(r'\.([^.]+)$', '', name)
        file_name = re.sub(r'^file:', '', file_name)

        file_all_name : str = sha224_replace(file_name) + '.' + mime_type
        file_directory : str = os.path.join(load_image_url(), file_all_name)

        if not os.path.exists(file_directory):
            return redirect('/w/' + url_pas(name))

        if flask.request.method == 'POST':
            admin_check(None, 'file del (' + name + ')')
            os.remove(file_directory)

            if flask.request.form.get('with_doc', '') != '':
                edit_delete(name)

            return redirect('/w/' + url_pas(name))
        else:
            return easy_minify(flask.render_template(skin_check(),
                imp = [name, wiki_set(), wiki_custom(), wiki_css(['(' + load_lang('file_delete') + ')', 0])],
                data = '''
                    <form method="post">
                        <img src="/image/''' + url_pas(file_all_name) + '''">
                        <hr class="main_hr">
                        <a href="/image/''' + url_pas(file_all_name) + '''">/image/''' + url_pas(file_all_name) + '''</a>
                        <hr class="main_hr">
                        <input name="with_doc" type="checkbox" checked> ''' + load_lang('file_delete_with_document') + '''
                        <hr class="main_hr">
                        <button type="submit">''' + load_lang('file_delete') + '''</button>
                    </form>
                ''',
                menu = [['w/' + url_pas(name), load_lang('return')]]
            ))