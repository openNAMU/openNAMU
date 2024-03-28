from .tool.func import *

def give_delete_admin_group_2(name = 'test'):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if admin_check(conn) != 1:
            return re_error(conn, '/error/3')

        if flask.request.method == 'POST':
            admin_check(conn, None, 'auth list delete (' + name + ')')

            curs.execute(db_change("delete from alist where name = ?"), [name])
            curs.execute(db_change("update user_set set data = 'user' where name = 'acl' and data = ?"), [name])

            return redirect(conn, '/auth/list')
        else:
            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [get_lang(conn, "delete_admin_group"), wiki_set(conn), wiki_custom(conn), wiki_css(['(' + name + ')', 0])],
                data = '''
                    <form method=post>
                        <button type=submit>''' + get_lang(conn, 'delete') + '''</button>
                    </form>
                ''',
                menu = [['auth/list', get_lang(conn, 'return')]]
            ))