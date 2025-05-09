from .tool.func import *

async def topic_list(page = 1, name = 'Test'):
    with get_db_connect() as conn:
        return easy_minify(conn, flask.render_template(skin_check(conn),
            imp = [name, await wiki_set(), await wiki_custom(conn), wiki_css(['(' + get_lang(conn, 'discussion_list') + ')', 0])],
            data = '' + \
                '<div id="opennamu_topic_list"></div>' + \
                '<script defer src="/views/main_css/js/route/topic_list.js' + cache_v() + '"></script>' + \
                '<script>window.addEventListener("DOMContentLoaded", function() { opennamu_topic_list(); });</script>' + \
            '',
            menu = [['w/' + url_pas(name), get_lang(conn, 'document')]]
        ))