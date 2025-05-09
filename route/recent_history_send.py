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

                return easy_minify(conn, flask.render_template(skin_check(conn),
                    imp = [name, await wiki_set(), await wiki_custom(conn), wiki_css(['(' + get_lang(conn, 'send_edit') + ') (r' + num + ')', 0])],
                    data = '''
                        <form method="post">
                            <span>''' + get_lang(conn, 'delete_warning') + '''</span>
                            <hr class="main_hr">
                            <input value="''' + html.escape(send) + '''" name="send">
                            <hr class="main_hr">
                            <button type="submit">''' + get_lang(conn, 'edit') + '''</button>
                        </form>
                    ''',
                    menu = [['history/' + url_pas(name), get_lang(conn, 'return')]]
                ))
            else:
                return redirect(conn, '/history/' + url_pas(name))