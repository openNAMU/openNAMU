from .tool.func import *

async def main_sys_shutdown():
    with get_db_connect() as conn:
        if await acl_check('', 'owner_auth', '', '') == 1:
            return await re_error(conn, 3)

        if flask.request.method == 'POST':
            await acl_check(tool = 'owner_auth', memo = 'shutdown')

            print('Shutdown')

            sys.exit()
        else:
            return await render_template(
                await get_lang('wiki_shutdown'),
                '''
                    <form method="post">
                        <button type="submit">''' + await get_lang('shutdown') + '''</button>
                    </form>
                ''',
                0,
                [['manager', await get_lang('return')]]
            )
