from .tool.func import *

async def give_auth(user_name = ''):
    with get_db_connect() as conn:
        if user_name == '':
            user_name = get_lang(conn, 'authorize')
            sub = 0
        else:
            sub = '(' + get_lang(conn, 'authorize') + ')'

        return easy_minify(conn, flask.render_template(skin_check(conn),
            imp = [user_name, await wiki_set(), await wiki_custom(conn), wiki_css([sub, 0])],
            data = '' + \
                '<div id="opennamu_give_auth"></div>' + \
                '<script defer src="/views/main_css/js/route/give_auth.js' + cache_v() + '"></script>' + \
                '<script>window.addEventListener("DOMContentLoaded", function() { opennamu_give_auth(); });</script>' + \
            '',
            menu = [['manager', get_lang(conn, 'return')]]
        ))