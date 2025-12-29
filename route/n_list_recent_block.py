from .tool.func import *

async def list_recent_block(user_name = 'Test', tool = 'all', num = 1, why = ''):
    sub = 0
    if tool == 'ongoing':
        sub = '(' + await get_lang('in_progress') + ')'
    elif tool == 'regex':
        sub = '(' + await get_lang('regex') + ')'
    elif tool == 'user':
        sub = '(' + await get_lang('blocked') + ')'
    elif tool == 'cidr':
        sub = '(' + await get_lang('cidr') + ')'
    elif tool == 'private':
        sub = '(' + await get_lang('private') + ')'
    elif tool == 'admin':
        sub = '(' + await get_lang('admin') + ')'

    return await render_template(
        await get_lang('recent_ban'),
        '' + \
            '<div id="opennamu_list_recent_block"></div>' + \
            '<script defer src="/views/main_css/js/route/list_recent_block.js' + cache_v() + '"></script>' + \
            '<script>window.addEventListener("DOMContentLoaded", function() { opennamu_list_recent_block(); });</script>' + \
        '',
        sub,
        [['other', await get_lang('return')]]
    )
