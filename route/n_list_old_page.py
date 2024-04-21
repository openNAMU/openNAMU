from .tool.func import *

def list_old_page(num = 1, set_type = 'old'):
    with get_db_connect() as conn:
        title = ''
        if set_type == 'old':
            title = get_lang(conn, 'old_page')
        else:
            title = get_lang(conn, 'new_page')

        return easy_minify(conn, flask.render_template(skin_check(conn),
            imp = [title, wiki_set(conn), wiki_custom(conn), wiki_css([0, 0])],
            data = '' + \
                '<div id="opennamu_list_old_page"></div>' + \
                '<script defer src="/views/main_css/js/route/list_old_page.js' + cache_v() + '"></script>' + \
                '<script>window.addEventListener("DOMContentLoaded", function() { opennamu_list_old_page(); });</script>' + \
            '',
            menu = [['other', get_lang(conn, 'return')]]
        ))