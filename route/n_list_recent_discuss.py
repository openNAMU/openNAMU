from .tool.func import *

def list_recent_discuss():
    with get_db_connect() as conn:
        return easy_minify(conn, flask.render_template(skin_check(conn),
            imp = [get_lang(conn, 'recent_discuss'), wiki_set(conn), wiki_custom(conn), wiki_css([0, 0])],
            data = '' + \
                '<div id="opennamu_list_recent_discuss"></div>' + \
                '<script src="/views/main_css/js/route/list_recent_discuss.js' + cache_v() + '"></script>' + \
                '<script>opennamu_list_recent_discuss();</script>' + \
            '',
            menu = [['other', get_lang(conn, 'return')]]
        ))