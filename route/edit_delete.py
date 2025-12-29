from .tool.func import *

async def edit_delete(name):
    with get_db_connect() as conn:
        curs = conn.cursor()

        ip = ip_check()
        if await acl_check(name, 'document_delete') == 1:
            return await re_error(conn, 0)

        curs.execute(db_change("select title from data where title = ?"), [name])
        if not curs.fetchall():
            return redirect(conn, '/w/' + url_pas(name))

        if flask.request.method == 'POST':
            if await captcha_post(conn, flask.request.form.get('g-recaptcha-response', flask.request.form.get('g-recaptcha', ''))) == 1:
                return await re_error(conn, 13)

            if await do_edit_slow_check(conn) == 1:
                return await re_error(conn, 24)
            
            send = flask.request.form.get('send', '')
            agree = flask.request.form.get('copyright_agreement', '')
            
            if await do_edit_send_check(conn, send) == 1:
                return await re_error(conn, 37)
            
            if do_edit_text_bottom_check_box_check(conn, agree) == 1:
                return await re_error(conn, 29)

            curs.execute(db_change("select data from data where title = ?"), [name])
            data = curs.fetchall()
            if data:
                today = get_time()
                leng = '-' + str(len(data[0][0]))

                history_plus(conn, 
                    name,
                    '',
                    today,
                    ip,
                    send,
                    leng,
                    mode = 'delete'
                )

                curs.execute(db_change("select title, link from back where title = ? and not type = 'cat' and not type = 'no'"), [name])
                for data in curs.fetchall():
                    curs.execute(db_change("insert into back (title, link, type, data) values (?, ?, 'no', '')"), [data[0], data[1]])

                curs.execute(db_change("delete from back where link = ?"), [name])
                curs.execute(db_change("delete from data where title = ?"), [name])

            return redirect(conn, '/w/' + url_pas(name))
        else:            
            return await render_template(
                name,
                '''
                    <form method="post">
                        <input class="__ON_INPUT__" placeholder="''' + await get_lang('why') + '''" name="send">
                        <hr class="main_hr">
                        ''' + await captcha_get(conn) + await ip_warning(conn) + get_edit_text_bottom_check_box(conn) + get_edit_text_bottom(conn, 'delete')  + '''
                        <button class="__ON_BUTTON__" type="submit">''' + await get_lang('delete') + '''</button>
                    </form>
                ''',
                '(' + await get_lang('delete') + ')',
                [['w/' + url_pas(name), await get_lang('return')]]
            )
