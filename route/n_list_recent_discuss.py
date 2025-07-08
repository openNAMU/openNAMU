from .tool.func import *

async def list_recent_discuss(num = 1, tool = 'normal'):
    with get_db_connect() as conn:
        m_sub = 0
        if tool == 'close':
            m_sub = '(' + await get_lang('close_discussion') + ')'
        elif tool == 'open':
            m_sub = '(' + await get_lang('open_discussion') + ')'

        return easy_minify(flask.render_template(await skin_check(),
            imp = [await get_lang('recent_discussion'), await wiki_set(), await wiki_custom(conn), wiki_css([m_sub, 0])],
            data = '' + \
                '<div id="opennamu_list_recent_discuss"></div>' + \
                '<script defer src="/views/main_css/js/route/list_recent_discuss.js' + cache_v() + '"></script>' + \
                '<script>window.addEventListener("DOMContentLoaded", function() { opennamu_list_recent_discuss(); });</script>' + \
            '',
            menu = [['other', await get_lang('return')]]
        ))