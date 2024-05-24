from .tool.func import *

def list_user_check_submit(name = 'Test'):
    with get_db_connect() as conn:
        return easy_minify(conn, flask.render_template(skin_check(conn),
            imp = [name, wiki_set(conn), wiki_custom(conn), wiki_css(['(' + get_lang(conn, 'check') + ')', 0])],
            data = '' + \
                '<div id="opennamu_list_user_check_submit"></div>' + \
                '<script defer src="/views/main_css/js/route/list_user_check_submit.js' + cache_v() + '"></script>' + \
                '<script>window.addEventListener("DOMContentLoaded", function() { opennamu_list_user_check_submit(); });</script>' + \
            '',
            menu = [['setting', get_lang(conn, 'return')]]
        ))