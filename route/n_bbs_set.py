from .tool.func import *

def bbs_set():
    with get_db_connect() as conn:
        return easy_minify(conn, flask.render_template(skin_check(conn),
            imp = [get_lang(conn, 'bbs_set'), wiki_set(conn), wiki_custom(conn), wiki_css([0, 0])],
            data = '' + \
                '<div id="opennamu_bbs_set"></div>' + \
                '<script defer src="/views/main_css/js/route/bbs_set.js' + cache_v() + '"></script>' + \
                '<script>window.addEventListener("DOMContentLoaded", function() { opennamu_bbs_set(); });</script>' + \
            '',
            menu = [['bbs/main', get_lang(conn, 'bbs_main')]]
        ))