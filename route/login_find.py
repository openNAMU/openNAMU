from .tool.func import *

async def login_find():
    return easy_minify(flask.render_template(await skin_check(),
        imp = [await get_lang('password_search'), await wiki_set(), await wiki_custom(), wiki_css([0, 0])],
        data = '''
            <ul>
                <li><a href="/login/find/email">''' + await get_lang('email') + '''</a></li>
                <li><a href="/login/find/key">''' + await get_lang('key') + '''</a></li>
            </ul>
        ''',
        menu = [['user', await get_lang('return')]]
    ))