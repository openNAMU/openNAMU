from .tool.func import *

async def list_user_check_submit(name = 'Test'):
    with get_db_connect() as conn:
        return easy_minify(flask.render_template(await skin_check(conn),
            imp = [name, await wiki_set(), await wiki_custom(conn), wiki_css(['(' + await get_lang('check') + ')', 0])],
            data = '' + \
                '<div id="opennamu_list_user_check_submit"></div>' + \
                '<script defer src="/views/main_css/js/route/list_user_check_submit.js' + cache_v() + '"></script>' + \
                '<script>window.addEventListener("DOMContentLoaded", function() { opennamu_list_user_check_submit(); });</script>' + \
            '',
            menu = [['setting', await get_lang('return')]]
        ))