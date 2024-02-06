from .tool.func import *

def give_delete_admin_group_2(name = 'test'):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if admin_check() != 1:
            return re_error('/error/3')

        if flask.request.method == 'POST':
            admin_check(None, 'auth list delete (' + name + ')')

            curs.execute(db_change("delete from alist where name = ?"), [name])
            curs.execute(db_change("update user_set set data = 'user' where name = 'acl' and data = ?"), [name])

            conn.commit()

            return redirect('/auth/list')
        else:
            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang("delete_admin_group"), wiki_set(), wiki_custom(), wiki_css(['(' + name + ')', 0])],
                data = '''
                    <form method=post>
                        <button type=submit>''' + load_lang('delete') + '''</button>
                    </form>
                ''',
                menu = [['auth/list', load_lang('return')]]
            ))