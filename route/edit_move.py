from .tool.func import *

def edit_move(name):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if acl_check(name, 'document_move') == 1:
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
            
            has_error = 0

            move_option = flask.request.form.get('move_option', 'none')
            move_option_topic = flask.request.form.get('move_topic_option', 'none')
            document_set_option = flask.request.form.get('document_set_option', 'none')
            
            if do_edit_send_check(send) == 1:
                return re_error('/error/37')
            
            if do_edit_text_bottom_check_box_check(agree) == 1:
                return re_error('/error/29')

            # 문서 이동 파트 S
            curs.execute(db_change("select title from history where title = ?"), [move_title])
            if curs.fetchall():
                if (
                    move_option == 'merge' and 
                    admin_check(None, 'merge documents (' + name + ') (' + move_title + ')') == 1
                ):
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

                    curs.execute(db_change('select data from other where name = "count_all_title"'))
                    curs.execute(db_change("update other set data = ? where name = 'count_all_title'"), [str(int(curs.fetchall()[0][0]) - 1)])

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
                elif move_option == 'reverse':
                    i = 0
                    var_name = ''
                    while var_name == '':
                        curs.execute(db_change("select title from history where title = ?"), ['test ' + str(i)])
                        if not curs.fetchall():
                            var_name = 'test ' + str(i)
                        else:
                            i += 1

                    curs.execute(db_change("select data from data where title = ?"), [name])
                    data = curs.fetchall()
                    if data:
                        curs.execute(db_change("update data set title = ? where title = ?"), [var_name, name])
                        curs.execute(db_change("update back set link = ? where link = ?"), [var_name, name])

                    curs.execute(db_change("update history set title = ? where title = ?"), [var_name, name])
                    curs.execute(db_change("update rc set title = ? where title = ?"), [var_name, name])

                    for title_name in [[move_title, name], [var_name, move_title]]:
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
                            t_check = '<a>' + (title_name[0] if title_name[0] != var_name else name) + '</a> - <a>' + title_name[1] + '</a> move',
                            mode = 'move'
                        )

                        curs.execute(db_change("update history set title = ? where title = ?"), [title_name[1], title_name[0]])
                        curs.execute(db_change("update rc set title = ? where title = ?"), [title_name[1], title_name[0]])
                elif move_option != 'none':
                    has_error = 1
            elif move_option != 'none':                
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
                
            # 문서 이동 파트 E
            
            # 토론 이동 파트 S
            if (
                move_option_topic == 'merge' and
                admin_check(None, 'merge document\'s topics (' + name + ') (' + move_title + ')') == 1
            ):
                curs.execute(db_change("update rd set title = ? where title = ?"), [move_title, name])
            elif move_option_topic == 'reverse':
                i = 0
                var_name = ''
                while var_name == '':
                    curs.execute(db_change("select title from rd where title = ?"), ['test ' + str(i)])
                    if not curs.fetchall():
                        var_name = 'test ' + str(i)
                    else:
                        i += 1
                
                curs.execute(db_change("update rd set title = ? where title = ?"), [var_name, move_title])
                curs.execute(db_change("update rd set title = ? where title = ?"), [move_title, name])
                curs.execute(db_change("update rd set title = ? where title = ?"), [name, var_name])
            elif move_option_topic == 'normal':
                curs.execute(db_change("select title from rd where title = ?"), [move_title])
                if curs.fetchall():
                    has_error = 1
                else:
                    curs.execute(db_change("update rd set title = ? where title = ?"), [move_title, name])

            # 토론 이동 파트 E

            # data_set 이동 파트 S
            if document_set_option == 'reverse':
                i = 0
                var_name = ''
                while var_name == '':
                    curs.execute(db_change("select title from rd where title = ?"), ['test ' + str(i)])
                    if not curs.fetchall():
                        var_name = 'test ' + str(i)
                    else:
                        i += 1
                
                # create_data['data_set'] = ['doc_name', 'doc_rev', 'set_name', 'set_data']
                # create_data['acl'] = ['title', 'data', 'type']
                curs.execute(db_change("update data_set set doc_name = ? where doc_name = ?"), [var_name, move_title])
                curs.execute(db_change("update data_set set doc_name = ? where doc_name = ?"), [move_title, name])
                curs.execute(db_change("update data_set set doc_name = ? where doc_name = ?"), [name, var_name])

                curs.execute(db_change("update acl set title = ? where title = ?"), [var_name, move_title])
                curs.execute(db_change("update acl set title = ? where title = ?"), [move_title, name])
                curs.execute(db_change("update acl set title = ? where title = ?"), [name, var_name])
            elif document_set_option == 'normal':
                curs.execute(db_change("delete from data_set where doc_name = ?"), [move_title])
                curs.execute(db_change("update data_set set doc_name = ? where doc_name = ?"), [move_title, name])

                curs.execute(db_change("delete from acl where title = ?"), [move_title])
                curs.execute(db_change("update acl set title = ? where title = ?"), [move_title, name])

            if document_set_option != 'reverse':
                curs.execute(db_change("select data from data where title = ?"), [name])
                db_data = curs.fetchall()
                if db_data:
                    render_set(
                        doc_name = name,
                        doc_data = db_data[0][0],
                        data_type = 'backlink'
                    )

                curs.execute(db_change("select data from data where title = ?"), [move_title])
                db_data = curs.fetchall()
                if db_data:
                    render_set(
                        doc_name = move_title,
                        doc_data = db_data[0][0],
                        data_type = 'backlink'
                    )

            # data_set 이동 파트 E
                
            conn.commit()

            if has_error == 0:
                return redirect('/w/' + url_pas(move_title))
            else:
                return re_error('/error/19')
        else:
            owner_auth = admin_check()

            return easy_minify(flask.render_template(skin_check(),
                imp = [name, wiki_set(), wiki_custom(), wiki_css(['(' + load_lang('move') + ')', 0])],
                data = '''
                    <form method="post">
                        <span>''' + load_lang('document_name') + '''</span>
                        <hr class="main_hr">
                        <input value="''' + name + '''" name="title" type="text">
                        <hr class="main_hr">
                        
                        <input placeholder="''' + load_lang('why') + '''" name="send" type="text">
                        <hr class="main_hr">
                        
                        <h2>''' + load_lang('document') + '''</h2>
                        <select name="move_option">
                            <option value="none"> ''' + load_lang('dont_move') + '''</option>
                            <option value="normal"> ''' + load_lang('normal') + '''</option>
                            <option value="reverse"> ''' + load_lang('replace_move') + '''</option>
                            ''' + ('<option value="merge"> ' + load_lang('merge_move') + '</option>' if owner_auth == 1 else '') + '''
                        </select>
                        <hr class="main_hr">
                        <!-- <input type="checkbox" name="move_redirect_make"> ''' + load_lang('move_redirect_make') + '''
                        <hr class="main_hr"> -->
                        
                        <h2>''' + load_lang('discussion') + '''</h2>
                        <select name="move_topic_option">
                            <option value="none"> ''' + load_lang('dont_move') + '''</option>
                            <option value="normal"> ''' + load_lang('normal') + '''</option>
                            <option value="reverse"> ''' + load_lang('replace_move') + '''</option>
                            ''' + ('<option value="merge"> ' + load_lang('merge_move') + '</option>' if owner_auth == 1 else '') + '''
                        </select>
                        <hr class="main_hr">

                        ''' + ((
                            '''<h2>''' + load_lang('document_set') + '''</h2>
                            <select name="document_set_option">
                                <option value="none"> ''' + load_lang('dont_move') + '''</option>
                                <option value="normal"> ''' + load_lang('normal') + '''</option>
                                <option value="reverse"> ''' + load_lang('replace_move') + '''</option>
                            </select>
                            <hr class="main_hr">
                            '''
                        ) if owner_auth == 1 else '') + '''

                        ''' + captcha_get() + ip_warning() + get_edit_text_bottom_check_box() + get_edit_text_bottom() + '''
                        
                        <button type="submit">''' + load_lang('move') + '''</button>
                    </form>
                ''',
                menu = [['w/' + url_pas(name), load_lang('return')]]
            ))