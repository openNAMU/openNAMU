from .tool.func import *

def main_sys_shutdown():
    with get_db_connect() as conn:
        if acl_check('', 'owner_auth', '', '') == 1:
            return re_error(conn, 3)

        if flask.request.method == 'POST':
            acl_check(tool = 'owner_auth', memo = 'shutdown')

            print('Shutdown')

            os._exit(0)
        else:
            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [get_lang(conn, 'wiki_shutdown'), wiki_set(conn), wiki_custom(conn), wiki_css([0, 0])],
                data = '''
                    <form method="post">
                        <button type="submit">''' + get_lang(conn, 'shutdown') + '''</button>
                    </form>
                ''',
                menu = [['manager', get_lang(conn, 'return')]]
            ))