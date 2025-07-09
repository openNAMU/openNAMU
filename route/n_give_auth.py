from .tool.func import *

async def give_auth(user_name = ''):
    if user_name == '':
        user_name = await get_lang('authorize')
        sub = 0
    else:
        sub = '(' + await get_lang('authorize') + ')'

    return easy_minify(flask.render_template(await skin_check(),
        imp = [user_name, await wiki_set(), await wiki_custom(), wiki_css([sub, 0])],
        data = '' + \
            '<div id="opennamu_give_auth"></div>' + \
            '<script defer src="/views/main_css/js/route/give_auth.js' + cache_v() + '"></script>' + \
            '<script>window.addEventListener("DOMContentLoaded", function() { opennamu_give_auth(); });</script>' + \
        '',
        menu = [['manager', await get_lang('return')]]
    ))