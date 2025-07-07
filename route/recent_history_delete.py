from .tool.func import *

async def recent_history_delete(name = 'Test', rev = 1):
    with get_db_connect() as conn:
        curs = conn.cursor()

        num = str(rev)

        if await acl_check('', 'owner_auth', '', '') == 1:
            return await re_error(conn, 3)

        if flask.request.method == 'POST':
            await acl_check(tool = 'owner_auth', memo = 'history delete ' + name + ' r' + num)

            curs.execute(db_change("delete from history where id = ? and title = ?"), [num, name])

            return redirect(conn, '/history/' + url_pas(name))
        else:
            return easy_minify(flask.render_template(await skin_check(conn),
                imp = [name, await wiki_set(), await wiki_custom(conn), wiki_css(['(' + await get_lang('history_delete') + ') (r' + num + ')', 0])],
                data = '''
                    <form method="post">
                        <span>''' + await get_lang('delete_warning') + '''</span>
                        <hr class="main_hr">
                        <button type="submit">''' + await get_lang('delete') + '''</button>
                    </form>
                ''',
                menu = [['history/' + url_pas(name), await get_lang('return')]]
            ))