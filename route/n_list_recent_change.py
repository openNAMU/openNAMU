from .tool.func import *

def list_recent_change(num = 1, set_type = 'normal'):
    with get_db_connect() as conn:
        return easy_minify(conn, flask.render_template(skin_check(conn),
            imp = [get_lang(conn, 'recent_change'), wiki_set(conn), wiki_custom(conn), wiki_css(['(' + get_lang(conn, set_type) + ')', 0])],
            data = '' + \
                '<div id="opennamu_list_recent_change"></div>' + \
                '<script defer src="/views/main_css/js/route/list_recent_change.js' + cache_v() + '"></script>' + \
                '<script>window.addEventListener("DOMContentLoaded", function() { opennamu_list_recent_change(); });</script>' + \
            '',
            menu = [['other', get_lang(conn, 'return')], ['recent_edit_request', get_lang(conn, 'edit_request')]]
        ))