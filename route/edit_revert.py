from .tool.func import *

async def edit_revert(name, num):
    with get_db_connect() as conn:
        curs = conn.cursor()

        curs.execute(db_change("select title from history where title = ? and id = ? and hide = 'O'"), [name, str(num)])
        if curs.fetchall() and await acl_check(tool = 'hidel_auth') == 1:
            return await re_error(conn, 3)

        if await acl_check(name, 'document_edit') == 1:
            return await re_error(conn, 0)
        
        curs.execute(db_change("select data from history where title = ? and id = ?"), [name, str(num)])
        data = curs.fetchall()
        if not data:
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

            if await do_edit_filter(conn, data[0][0]) == 1:
                return await re_error(conn, 21)
            
            curs.execute(db_change("select data from other where name = 'document_content_max_length'"))
            db_data = curs.fetchall()
            if db_data and db_data[0][0] != '':
                if int(number_check(db_data[0][0])) < len(data[0][0]):
                    return await re_error(conn, 44)

            curs.execute(db_change("select data from data where title = ?"), [name])
            data_old = curs.fetchall()
            if data_old:
                leng = leng_check(len(data_old[0][0]), len(data[0][0]))
                curs.execute(db_change("update data set data = ? where title = ?"), [data[0][0], name])
            else:
                leng = '+' + str(len(data[0][0]))
                curs.execute(db_change("insert into data (title, data) values (?, ?)"), [name, data[0][0]])

            history_plus(conn, 
                name,
                data[0][0],
                get_time(),
                ip_check(),
                flask.request.form.get('send', ''),
                leng,
                t_check = 'r' + str(num),
                mode = 'revert'
            )

            render_set(conn, 
                doc_name = name,
                doc_data = data[0][0],
                data_type = 'backlink'
            )

            return redirect(conn, '/w/' + url_pas(name))
        else:
            if data:
                preview = '<hr class="main_hr"><pre>' + html.escape(data[0][0]) + '</pre>'
            else:
                preview = ''
            
            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [name, await wiki_set(), await wiki_custom(conn), wiki_css(['(r' + str(num) + ') (' + get_lang(conn, 'revert') + ')', 0])],
                data = '''
                    <form method="post">
                        <input placeholder="''' + get_lang(conn, 'why') + '''" name="send" type="text">
                        <hr class="main_hr">
                        ''' + await captcha_get(conn) + ip_warning(conn) + get_edit_text_bottom_check_box(conn) + get_edit_text_bottom(conn, 'revert')  + '''
                        <button type="submit">''' + get_lang(conn, 'revert') + '''</button>
                    </form>
                ''' + preview,
                menu = [['history/' + url_pas(name), get_lang(conn, 'history')], ['recent_changes', get_lang(conn, 'recent_change')]]
            ))