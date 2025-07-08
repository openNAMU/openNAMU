from .tool.func import *

async def bbs_main():
    with get_db_connect() as conn:
        return easy_minify(flask.render_template(await skin_check(),
            imp = [await get_lang('bbs_main'), await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
            data = '' + \
                '<div id="opennamu_bbs_main"></div>' + \
                '<script defer src="/views/main_css/js/route/bbs_main.js' + cache_v() + '"></script>' + \
                '<script>window.addEventListener("DOMContentLoaded", function() { opennamu_bbs_main(); });</script>' + \
            '',
            menu = [['other', await get_lang('other_tool')], ['bbs/make', await get_lang('add')]]
        ))