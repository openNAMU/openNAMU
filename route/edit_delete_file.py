from .tool.func import *

from .edit_delete import edit_delete

def edit_delete_file(name = 'test.jpg'):
    with get_db_connect() as conn:
        curs = conn.cursor()

        ip = ip_check()
        if admin_check(conn) == 0:
            return re_error(conn, '/ban')

        mime_type = re.search(r'([^.]+)$', name)
        mime_type_str = 'jpg'
        if mime_type:
            mime_type_str = mime_type.group(1)

        file_name = re.sub(r'\.([^.]+)$', '', name)
        file_name = re.sub(r'^file:', '', file_name)

        file_all_name = sha224_replace(file_name) + '.' + mime_type_str
        file_directory = os.path.join(load_image_url(conn), file_all_name)

        if not os.path.exists(file_directory):
            return redirect(conn, '/w/' + url_pas(name))

        if flask.request.method == 'POST':
            admin_check(conn, None, 'file del (' + name + ')')
            os.remove(file_directory)

            if flask.request.form.get('with_doc', '') != '':
                edit_delete(name)

            return redirect(conn, '/w/' + url_pas(name))
        else:
            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [name, wiki_set(conn), wiki_custom(conn), wiki_css(['(' + get_lang(conn, 'file_delete') + ')', 0])],
                data = '''
                    <form method="post">
                        <img src="/image/''' + url_pas(file_all_name) + '''">
                        <hr class="main_hr">
                        <a href="/image/''' + url_pas(file_all_name) + '''">/image/''' + url_pas(file_all_name) + '''</a>
                        <hr class="main_hr">
                        <input name="with_doc" type="checkbox" checked> ''' + get_lang(conn, 'file_delete_with_document') + '''
                        <hr class="main_hr">
                        <button type="submit">''' + get_lang(conn, 'file_delete') + '''</button>
                    </form>
                ''',
                menu = [['w/' + url_pas(name), get_lang(conn, 'return')]]
            ))