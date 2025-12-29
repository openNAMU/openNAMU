from .tool.func import *

async def list_recent_discuss(num = 1, tool = 'normal'):
    m_sub = 0
    if tool == 'close':
        m_sub = '(' + await get_lang('close_discussion') + ')'
    elif tool == 'open':
        m_sub = '(' + await get_lang('open_discussion') + ')'

    return await render_template(
        await get_lang('recent_discussion'),
        '' + \
            '<div id="opennamu_list_recent_discuss"></div>' + \
            '<script defer src="/views/main_css/js/route/list_recent_discuss.js' + cache_v() + '"></script>' + \
            '<script>window.addEventListener("DOMContentLoaded", function() { opennamu_list_recent_discuss(); });</script>' + \
        '',
        m_sub,
        [['other', await get_lang('return')]]
    )
