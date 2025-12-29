from .tool.func import *

async def recent_history_reset(name = 'Test'):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if await acl_check('', 'owner_auth', '', '') == 1:
            return await re_error(conn, 3)

        if flask.request.method == 'POST':
            await acl_check(tool = 'owner_auth', memo = 'history reset ' + name)

            curs.execute(db_change("delete from history where title = ?"), [name])

            return redirect(conn, '/history/' + url_pas(name))
        else:
            return await render_template(
                name,
                '''
                    <form method="post">
                        <span>''' + await get_lang('delete_warning') + '''</span>
                        <hr class="main_hr">
                        <button type="submit">''' + await get_lang('reset') + '''</button>
                    </form>
                ''',
                '(' + await get_lang('history_reset') + ')',
                [['history/' + url_pas(name), await get_lang('return')]]
            )
