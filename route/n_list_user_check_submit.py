from .tool.func import *

async def list_user_check_submit(name = 'Test'):
    return await render_template(
        name,
        '' + \
            '<div id="opennamu_list_user_check_submit"></div>' + \
            '<script defer src="/views/main_css/js/route/list_user_check_submit.js' + cache_v() + '"></script>' + \
            '<script>window.addEventListener("DOMContentLoaded", function() { opennamu_list_user_check_submit(); });</script>' + \
        '',
        '(' + await get_lang('check') + ')',
        [['setting', await get_lang('return')]]
    )
