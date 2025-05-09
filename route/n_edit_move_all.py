from .tool.func import *

async def edit_move_all():
    with get_db_connect() as conn:
        return easy_minify(conn, flask.render_template(skin_check(conn),
            imp = [get_lang(conn, 'multiple_move'), await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
            data = '' + \
                '<div id="opennamu_edit_move_all"></div>' + \
                '<script defer src="/views/main_css/js/route/edit_move_all.js' + cache_v() + '"></script>' + \
                '<script>window.addEventListener("DOMContentLoaded", function() { opennamu_edit_move_all(); });</script>' + \
            '',
            menu = [['other', get_lang(conn, 'return')]]
        ))