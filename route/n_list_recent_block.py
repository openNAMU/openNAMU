from .tool.func import *

async def list_recent_block(user_name = 'Test', tool = 'all', num = 1, why = ''):
    with get_db_connect() as conn:
        sub = 0
        if tool == 'ongoing':
            sub = '(' + get_lang(conn, 'in_progress') + ')'
        elif tool == 'regex':
            sub = '(' + get_lang(conn, 'regex') + ')'
        elif tool == 'user':
            sub = '(' + get_lang(conn, 'blocked') + ')'
        elif tool == 'cidr':
            sub = '(' + get_lang(conn, 'cidr') + ')'
        elif tool == 'private':
            sub = '(' + get_lang(conn, 'private') + ')'
        elif tool == 'admin':
            sub = '(' + get_lang(conn, 'admin') + ')'

        return easy_minify(conn, flask.render_template(skin_check(conn),
            imp = [get_lang(conn, 'recent_ban'), await wiki_set(), await wiki_custom(conn), wiki_css([sub, 0])],
            data = '' + \
                '<div id="opennamu_list_recent_block"></div>' + \
                '<script defer src="/views/main_css/js/route/list_recent_block.js' + cache_v() + '"></script>' + \
                '<script>window.addEventListener("DOMContentLoaded", function() { opennamu_list_recent_block(); });</script>' + \
            '',
            menu = [['other', get_lang(conn, 'return')]]
        ))