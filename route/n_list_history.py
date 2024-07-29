from .tool.func import *

def list_history(num = 1, set_type = 'normal', doc_name = 'Test'):
    with get_db_connect() as conn:
        return easy_minify(conn, flask.render_template(skin_check(conn),
            imp = [doc_name, wiki_set(conn), wiki_custom(conn), wiki_css(['(' + get_lang(conn, 'history') + ')', 0])],
            data = '' + \
                '<div id="opennamu_list_history"></div>' + \
                '<script defer src="/views/main_css/js/route/list_history.js' + cache_v() + '"></script>' + \
                '<script>window.addEventListener("DOMContentLoaded", function() { opennamu_list_history(); });</script>' + \
            '',
            menu = [['other', get_lang(conn, 'return')], ['history_add/' + url_pas(doc_name), get_lang(conn, 'history_add')], ['history_reset/' + url_pas(doc_name), get_lang(conn, 'history_reset')]]
        ))