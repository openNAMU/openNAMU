from .tool.func import *

async def w_watch_list(name, num = 1, do_type = 'watch_list'):
    with get_db_connect() as conn:
        return easy_minify(conn, flask.render_template(skin_check(conn),
            imp = [name, await wiki_set(), await wiki_custom(conn), wiki_css(['(' + get_lang(conn, do_type if do_type == 'star_doc' else 'watchlist') + ')', 0])],
            data = '' + \
                '<div id="opennamu_w_watch_list"></div>' + \
                '<script defer src="/views/main_css/js/route/w_watch_list.js' + cache_v() + '"></script>' + \
                '<script>window.addEventListener("DOMContentLoaded", function() { opennamu_w_watch_list(); });</script>' + \
            '',
            menu = [['w/' + url_pas(name), get_lang(conn, 'return')]]
        ))