from .tool.func import *

async def recent_record_reset(name = 'Test'):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if await acl_check('', 'owner_auth', '', '') == 1:
            return await re_error(conn, 3)

        if flask.request.method == 'POST':
            await acl_check(tool = 'owner_auth', memo = 'record reset ' + name)

            curs.execute(db_change("delete from history where ip = ?"), [name])

            return redirect(conn, '/record/' + url_pas(name))
        else:
            return easy_minify(flask.render_template(await skin_check(),
                imp = [name, await wiki_set(), await wiki_custom(conn), wiki_css(['(' + await get_lang('record_reset') + ')', 0])],
                data = '''
                    <form method="post">
                        <span>''' + await get_lang('delete_warning') + '''</span>
                        <hr class="main_hr">
                        <button type="submit">''' + await get_lang('reset') + '''</button>
                    </form>
                ''',
                menu = [['record/' + url_pas(name), await get_lang('return')]]
            ))