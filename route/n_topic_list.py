from .tool.func import *

async def topic_list(page = 1, name = 'Test'):
    return await render_template(
        name,
        '' + \
            '<div id="opennamu_topic_list"></div>' + \
            '<script defer src="/views/main_css/js/route/topic_list.js' + cache_v() + '"></script>' + \
            '<script>window.addEventListener("DOMContentLoaded", function() { opennamu_topic_list(); });</script>' + \
        '',
        '(' + await get_lang('discussion_list') + ')',
        [['w/' + url_pas(name), await get_lang('document')]]
    )
