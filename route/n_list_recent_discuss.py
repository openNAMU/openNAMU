from .tool.func import *

async def list_recent_discuss(num = 1, tool = 'normal'):
    with get_db_connect() as conn:
        m_sub = 0
        if tool == 'close':
            m_sub = '(' + get_lang(conn, 'close_discussion') + ')'
        elif tool == 'open':
            m_sub = '(' + get_lang(conn, 'open_discussion') + ')'

        return easy_minify(conn, flask.render_template(skin_check(conn),
            imp = [get_lang(conn, 'recent_discussion'), await wiki_set(), await wiki_custom(conn), wiki_css([m_sub, 0])],
            data = '' + \
                '<div id="opennamu_list_recent_discuss"></div>' + \
                '<script defer src="/views/main_css/js/route/list_recent_discuss.js' + cache_v() + '"></script>' + \
                '<script>window.addEventListener("DOMContentLoaded", function() { opennamu_list_recent_discuss(); });</script>' + \
            '',
            menu = [['other', get_lang(conn, 'return')]]
        ))