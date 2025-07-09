from .tool.func import *

async def w_watch_list(name, num = 1, do_type = 'watch_list'):
    return easy_minify(flask.render_template(await skin_check(),
        imp = [name, await wiki_set(), await wiki_custom(), wiki_css(['(' + await get_lang(do_type if do_type == 'star_doc' else 'watchlist') + ')', 0])],
        data = '' + \
            '<div id="opennamu_w_watch_list"></div>' + \
            '<script defer src="/views/main_css/js/route/w_watch_list.js' + cache_v() + '"></script>' + \
            '<script>window.addEventListener("DOMContentLoaded", function() { opennamu_w_watch_list(); });</script>' + \
        '',
        menu = [['w/' + url_pas(name), await get_lang('return')]]
    ))