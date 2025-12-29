from .tool.func import *

async def edit_move_all():
    return await render_template(
        await get_lang('multiple_move'),
        '' + \
            '<div id="opennamu_edit_move_all"></div>' + \
            '<script defer src="/views/main_css/js/route/edit_move_all.js' + cache_v() + '"></script>' + \
            '<script>window.addEventListener("DOMContentLoaded", function() { opennamu_edit_move_all(); });</script>' + \
        '',
        0,
        [['other', await get_lang('return')]]
    )
