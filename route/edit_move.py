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

            # 역링크 관련 패치 해야할 듯

            # 문서 이동 파트 S
            curs.execute(db_change("select title from history where title = ?"), [move_title])
            if curs.fetchall():
                if move_option == 'merge' and admin_check(None, 'merge documents (' + name + ') (' + move_title + ')') == 1:
                    curs.execute(db_change("select data from data where title = ?"), [move_title])
                    data = curs.fetchall()
                    if data:
                        curs.execute(db_change("delete from data where title = ?"), [move_title])
                        curs.execute(db_change("delete from back where link = ?"), [move_title])

                    curs.execute(db_change("select data from data where title = ?"), [name])
                    data = curs.fetchall()
                    data_in = data[0][0] if data else ''

                    curs.execute(db_change("update data set title = ? where title = ?"), [move_title, name])
                    curs.execute(db_change("update back set link = ? where link = ?"), [move_title, name])

                    # 역링크 S
                    # 문서 합치기이므로 기존 문서 쪽은 no 역링크 생성, 이동하는 곳에는 no 역링크 제거
                    curs.execute(db_change("select distinct link from back where title = ?"), [name])
                    backlink = [[for_a[0], name, 'no', ''] for for_a in curs.fetchall()]
                    curs.executemany(db_change("insert into back (link, title, type, data) values (?, ?, ?, ?)"), backlink)
                    curs.execute(db_change("delete from back where title = ? and type = 'no'"), [move_title])
                    # 역링크 E

                    curs.execute(db_change("select id from history where title = ? order by id + 0 desc limit 1"), [move_title])
                    num = curs.fetchall()[0][0]

                    curs.execute(db_change("select id from history where title = ? order by id + 0 asc"), [name])
                    data = curs.fetchall()
                    for move in data:
                        curs.execute(db_change("update rc set title = ?, id = ? where title = ? and id = ?"), [move_title, str(int(num) + int(move[0])), name, move[0]])
                        curs.execute(db_change("update history set title = ?, id = ? where title = ? and id = ?"), [move_title, str(int(num) + int(move[0])), name, move[0]])

                    history_plus(move_title, data_in, time, ip, send, '0',
                        t_check = 'merge <a>' + name + '</a> - <a>' + move_title + '</a> move',
                        mode = 'move'
                    )
                elif move_option == 'reverse':
                    # 전체적인 구조 변경 필요
                    # 중간 문서 거치지 않고 불러와서 바로 변경하도록
                    # 문서 이동 말고 나머지도 그렇게 변경 필요함
                    i = 0
                    var_name = ''
                    while var_name == '':
                        temp_title = 'test ' + load_random_key() + ' ' + str(i)
                        curs.execute(db_change("select title from history where title = ? limit 1"), [temp_title])
                        if not curs.fetchall():
                            var_name = temp_title
                        else:
                            i += 1

                    for title_name in [[name, var_name], [move_title, name], [var_name, move_title]]:
                        curs.execute(db_change("update data set title = ? where title = ?"), [title_name[1], title_name[0]])
                        curs.execute(db_change("update back set link = ? where link = ?"), [title_name[1], title_name[0]])

                        curs.execute(db_change("update history set title = ? where title = ?"), [title_name[1], title_name[0]])
                        curs.execute(db_change("update rc set title = ? where title = ?"), [title_name[1], title_name[0]])

                    for title_name in [[name, move_title], [move_title, name]]:
                        curs.execute(db_change("select data from data where title = ?"), [name])
                        data = curs.fetchall()
                        data_in = data[0][0] if data else ''

                        history_plus(title_name[0], data_in, time, ip, send, '0',
                            t_check = '<a>' + title_name[0] + '</a> - <a>' + title_name[1] + '</a> move',
                            mode = 'move'
                        )
                elif move_option != 'none':
                    has_error = 1
            elif move_option != 'none':
                curs.execute(db_change("select data from data where title = ?"), [name])
                data = curs.fetchall()
                data_in = data[0][0] if data else ''

                curs.execute(db_change("update data set title = ? where title = ?"), [move_title, name])
                curs.execute(db_change("update back set link = ? where link = ?"), [move_title, name])

                # 역링크 S
                # 문서 합치기 쪽 역링크와 동일하게
                curs.execute(db_change("select distinct link from back where title = ?"), [name])
                backlink = [[for_a[0], name, 'no', ''] for for_a in curs.fetchall()]
                curs.executemany(db_change("insert into back (link, title, type, data) values (?, ?, ?, ?)"), backlink)
                curs.execute(db_change("delete from back where title = ? and type = 'no'"), [move_title])
                # 역링크 E

                # 역사와 최근 변경 이동 S
                curs.execute(db_change("update history set title = ? where title = ?"), [move_title, name])
                curs.execute(db_change("update rc set title = ? where title = ?"), [move_title, name])
                # 역사와 최근 변경 이동 E

                history_plus(move_title, data_in, time, ip, send, '0',
                    t_check = '<a>' + name + '</a> - <a>' + move_title + '</a> move',
                    mode = 'move'
                )

            # 문서 이동 파트 E
            
            # 토론 이동 파트 S
            curs.execute(db_change("select title from rd where title = ?"), [move_title])
            if curs.fetchall():
                if move_option_topic == 'merge' and admin_check(None, 'merge document\'s topics (' + name + ') (' + move_title + ')') == 1:
                    curs.execute(db_change("update rd set title = ? where title = ?"), [move_title, name])
                elif move_option_topic == 'reverse':
                    i = 0
                    var_name = ''
                    while var_name == '':
                        temp_title = 'test ' + load_random_key() + ' ' + str(i)
                        curs.execute(db_change("select title from rd where title = ? limit 1"), [temp_title])
                        if not curs.fetchall():
                            var_name = temp_title
                        else:
                            i += 1
                    
                    for title_name in [[name, var_name], [move_title, name], [var_name, move_title]]:
                        curs.execute(db_change("update rd set title = ? where title = ?"), [title_name[1], title_name[0]])
                else:
                    has_error = 1
            elif move_option_topic != 'none':
                curs.execute(db_change("update rd set title = ? where title = ?"), [move_title, name])

            # 토론 이동 파트 E

            # data_set 이동 파트 S
            if document_set_option == 'reverse':
                i = 0
                var_name = ''
                while var_name == '':
                    temp_title = 'test ' + load_random_key() + ' ' + str(i)
                    curs.execute(db_change("select title from history where title = ? limit 1"), [temp_title])
                    if not curs.fetchall():
                        var_name = temp_title
                    else:
                        i += 1
                
                for title_name in [[name, var_name], [move_title, name], [var_name, move_title]]:
                    curs.execute(db_change("update data_set set doc_name = ? where doc_name = ?"), [title_name[1], title_name[0]])
            elif document_set_option == 'normal':
                curs.execute(db_change("delete from data_set where doc_name = ?"), [move_title])
                curs.execute(db_change("delete from acl where title = ?"), [move_title])

                curs.execute(db_change("update data_set set doc_name = ? where doc_name = ?"), [move_title, name])
                curs.execute(db_change("update acl set title = ? where title = ?"), [move_title, name])

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