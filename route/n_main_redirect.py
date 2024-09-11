from .tool.func import *

def main_redirect(n = 1):
    with get_db_connect() as conn:
        return easy_minify(conn, flask.render_template(skin_check(conn),
            imp = [get_lang(conn, 'redirect'), wiki_set(conn), wiki_custom(conn), wiki_css([0, 0])],
            data = '' + \
                '<div id="opennamu_main_redirect"></div>' + \
                '<script defer src="/views/main_css/js/route/main_redirect.js' + cache_v() + '"></script>' + \
                '<script>window.addEventListener("DOMContentLoaded", function() { opennamu_main_redirect(); });</script>' + \
            '',
            menu = [['manager', get_lang(conn, 'return')]]
        ))