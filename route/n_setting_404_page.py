from .tool.func import *

def setting_404_page():
    with get_db_connect() as conn:
        return easy_minify(conn, flask.render_template(skin_check(conn),
            imp = [get_lang(conn, '404_page_setting'), wiki_set(conn), wiki_custom(conn), wiki_css([0, 0])],
            data = '' + \
                '<div id="opennamu_setting_404_page"></div>' + \
                '<script defer src="/views/main_css/js/route/setting_404_page.js' + cache_v() + '"></script>' + \
                '<script>window.addEventListener("DOMContentLoaded", function() { opennamu_setting_404_page(); });</script>' + \
            '',
            menu = [['setting', get_lang(conn, 'return')]]
        ))