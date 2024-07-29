from .tool.func import *

def recent_history_delete(name = 'Test', rev = 1):
    with get_db_connect() as conn:
        curs = conn.cursor()

        num = str(rev)

        if acl_check('', 'owner_auth', '', '') == 1:
            return re_error(conn, 3)

        if flask.request.method == 'POST':
            acl_check(tool = 'owner_auth', memo = 'history delete ' + name + ' r' + num)

            curs.execute(db_change("delete from history where id = ? and title = ?"), [num, name])

            return redirect(conn, '/history/' + url_pas(name))
        else:
            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [name, wiki_set(conn), wiki_custom(conn), wiki_css(['(' + get_lang(conn, 'history_delete') + ') (r' + num + ')', 0])],
                data = '''
                    <form method="post">
                        <span>''' + get_lang(conn, 'delete_warning') + '''</span>
                        <hr class="main_hr">
                        <button type="submit">''' + get_lang(conn, 'delete') + '''</button>
                    </form>
                ''',
                menu = [['history/' + url_pas(name), get_lang(conn, 'return')]]
            ))