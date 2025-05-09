from .tool.func import *

async def login_find():
    with get_db_connect() as conn:
        return easy_minify(conn, flask.render_template(skin_check(conn),
            imp = [get_lang(conn, 'password_search'), await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
            data = '''
                <ul>
                    <li><a href="/login/find/email">''' + get_lang(conn, 'email') + '''</a></li>
                    <li><a href="/login/find/key">''' + get_lang(conn, 'key') + '''</a></li>
                </ul>
            ''',
            menu = [['user', get_lang(conn, 'return')]]
        ))