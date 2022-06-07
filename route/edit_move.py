from .tool.func import *

def edit_move(name):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if acl_check(name) == 1:
            return re_error('/ban')
        
        if do_title_length_check(name) == 1:
            return re_error('/error/38')

        if flask.request.method == 'POST':
            move_title = flask.request.form.get('title', 'test')
            if acl_check(move_title) == 1:
                return re_error('/ban')

            if captcha_post(flask.request.form.get('g-recaptcha-response', flask.request.form.get('g-recaptcha', ''))) == 1:
                return re_error('/error/13')
            else:
                captcha_post('', 0)

            if do_edit_slow_check() == 1:
                return re_error('/error/24')
            
            send = flask.request.form.get('send', '')
            agree = flask.request.form.get('copyright_agreement', '')
            time = get_time()
            ip = ip_check()
            
            if do_edit_send_check(send) == 1:
                return re_error('/error/37')
            
            if do_edit_text_bottom_check_box_check(agree) == 1:
                return re_error('/error/29')

            curs.execute(db_change("select title from history where title = ?"), [move_title])
            if curs.fetchall():
                if flask.request.form.get('move_option', 'normal') == 'merge' and admin_check(None, 'merge documents') == 1:
                    curs.execute(db_change("select data from data where title = ?"), [move_title])
                    data = curs.fetchall()
                    if data:
                        curs.execute(db_change("delete from data where title = ?"), [move_title])
                        curs.execute(db_change("delete from back where link = ?"), [move_title])

                    curs.execute(db_change("select data from data where title = ?"), [name])
                    data = curs.fetchall()
                    if data:
                        curs.execute(db_change("update data set title = ? where title = ?"), [move_title, name])
                        curs.execute(db_change("update back set link = ? where link = ?"), [move_title, name])

                        data_in = data[0][0]
                    else:
                        data_in = ''

                    history_plus(
                        name,
                        data_in,
                        time,
                        ip,
                        send,
                        '0',
                        t_check = 'merge <a>' + name + '</a> - <a>' + move_title + '</a> move',
                        mode = 'move'
                    )

                    curs.execute(db_change("update back set type = 'no' where title = ? and not type = 'cat' and not type = 'no'"), [name])
                    curs.execute(db_change("delete from back where title = ? and not type = 'cat' and type = 'no'"), [move_title])

                    curs.execute(db_change("select id from history where title = ? order by id + 0 desc limit 1"), [move_title])
                    data = curs.fetchall()

                    num = data[0][0]

                    curs.execute(db_change("select id from history where title = ? order by id + 0 asc"), [name])
                    data = curs.fetchall()
                    for move in data:
                        curs.execute(db_change("update rc set title = ?, id = ? where title = ? and id = ?"), [
                            move_title, 
                            str(int(num) + int(move[0])), 
                            name, 
                            move[0]
                        ])
                        curs.execute(db_change("update history set title = ?, id = ? where title = ? and id = ?"), [
                            move_title, 
                            str(int(num) + int(move[0])), 
                            name, 
                            move[0]
                        ])

                    conn.commit()

                    return redirect('/w/' + url_pas(move_title))
                elif flask.request.form.get('move_option', 'normal') == 'reverse':
                    var_name = ''
                    i = 0
                    while 1:
                        curs.execute(db_change("select title from history where title = ?"), ['test ' + str(i)])
                        if not curs.fetchall():
                            curs.execute(db_change("select data from data where title = ?"), [name])
                            data = curs.fetchall()
                            if data:
                                curs.execute(db_change("update data set title = ? where title = ?"), ['test ' + str(i), name])
                                curs.execute(db_change("update back set link = ? where link = ?"), ['test ' + str(i), name])

                            curs.execute(db_change("update history set title = ? where title = ?"), ['test ' + str(i), name])
                            curs.execute(db_change("update rc set title = ? where title = ?"), ['test ' + str(i), name])

                            break
                        else:
                            i += 1

                    for title_name in [[move_title, name], ['test ' + str(i), move_title]]:
                        curs.execute(db_change("select data from data where title = ?"), [title_name[0]])
                        data = curs.fetchall()
                        if data:
                            curs.execute(db_change("update data set title = ? where title = ?"), [title_name[1], title_name[0]])
                            curs.execute(db_change("update back set link = ? where link = ?"), [title_name[1], title_name[0]])

                            data_in = data[0][0]
                        else:
                            data_in = ''

                        history_plus(
                            title_name[0],
                            data_in,
                            time,
                            ip,
                            send,
                            '0',
                            t_check = '<a>' + (title_name[0] if title_name[0] != 'test ' + str(i) else name) + '</a> - <a>' + title_name[1] + '</a> move',
                            mode = 'move'
                        )

                        curs.execute(db_change("update history set title = ? where title = ?"), [title_name[1], title_name[0]])
                        curs.execute(db_change("update rc set title = ? where title = ?"), [title_name[1], title_name[0]])
                        conn.commit()

                    return redirect('/w/' + url_pas(move_title))
                else:
                    return re_error('/error/19')
            else:
                curs.execute(db_change("select data from data where title = ?"), [name])
                data = curs.fetchall()
                if data:
                    curs.execute(db_change("update data set title = ? where title = ?"), [move_title, name])
                    curs.execute(db_change("update back set link = ? where link = ?"), [move_title, name])

                    data_in = data[0][0]
                else:
                    data_in = ''

                history_plus(
                    name,
                    data_in,
                    time,
                    ip,
                    send,
                    '0',
                    t_check = '<a>' + name + '</a> - <a>' + move_title + '</a> move',
                    mode = 'move'
                )

                curs.execute(db_change("update back set type = 'no' where title = ? and not type = 'cat' and not type = 'no'"), [name])
                curs.execute(db_change("delete from back where title = ? and not type = 'cat' and type = 'no'"), [move_title])

                curs.execute(db_change("update history set title = ? where title = ?"), [move_title, name])
                curs.execute(db_change("update rc set title = ? where title = ?"), [move_title, name])
                conn.commit()

                return redirect('/w/' + url_pas(move_title))
        else:
            return easy_minify(flask.render_template(skin_check(),
                imp = [name, wiki_set(), wiki_custom(), wiki_css(['(' + load_lang('move') + ')', 0])],
                data = '''
                    <form method="post">
                        ''' + ip_warning() + '''
                        <input placeholder="''' + load_lang('document_name') + '" value="' + name + '''" name="title" type="text">
                        <hr class="main_hr">
                        <input placeholder="''' + load_lang('why') + '''" name="send" type="text">
                        <hr class="main_hr">
                        <select name="move_option">
                            <option value="normal"> ''' + load_lang('normal') + '''</option>
                            <option value="reverse"> ''' + load_lang('replace_move') + '''</option>
                            ''' + ('<option value="merge"> ' + load_lang('merge_move') + '</option>' if admin_check() == 1 else '') + '''
                        </select>
                        <hr class="main_hr">
                        ''' + captcha_get() + ip_warning() + get_edit_text_bottom_check_box() + get_edit_text_bottom() + '''
                        <button type="submit">''' + load_lang('move') + '''</button>
                    </form>
                ''',
                menu = [['w/' + url_pas(name), load_lang('return')]]
            ))