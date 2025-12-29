from .tool.func import *

async def recent_history_send(name = 'Test', rev = 1):
    with get_db_connect() as conn:
        curs = conn.cursor()

        num = str(rev)

        if await acl_check('', 'owner_auth', '', '') == 1:
            return await re_error(conn, 3)

        if flask.request.method == 'POST':
            await acl_check(tool = 'owner_auth', memo = 'send edit ' + name + ' r' + num)

            curs.execute(db_change("select send from history where title = ? and id = ?"), [name, num])
            if curs.fetchall():
                curs.execute(db_change("update history set send = ? where title = ? and id = ?"), [
                    flask.request.form.get('send', ''),
                    name, 
                    num
                ])

            return redirect(conn, '/history/' + url_pas(name))
        else:
            curs.execute(db_change("select send from history where title = ? and id = ?"), [name, num])
            send = curs.fetchall()
            if send:
                send = send[0][0]

                return await render_template(
                    name,
                    '''
                        <form method="post">
                            <span>''' + await get_lang('delete_warning') + '''</span>
                            <hr class="main_hr">
                            <input class="__ON_INPUT__" value="''' + html.escape(send) + '''" name="send">
                            <hr class="main_hr">
                            <button class="__ON_BUTTON__" type="submit">''' + await get_lang('edit') + '''</button>
                        </form>
                    ''',
                    '(' + await get_lang('send_edit') + ') (r' + num + ')',
                    [['history/' + url_pas(name), await get_lang('return')]]
                )
            else:
                return redirect(conn, '/history/' + url_pas(name))
