from .tool.func import *

async def login_find():
    return await render_template(
        await get_lang('password_search'),
        '''
            <ul>
                <li><a href="/login/find/email">''' + await get_lang('email') + '''</a></li>
                <li><a href="/login/find/key">''' + await get_lang('key') + '''</a></li>
            </ul>
        ''',
        0,
        [['user', await get_lang('return')]]
    )
