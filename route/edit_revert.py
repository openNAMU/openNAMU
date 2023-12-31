from .tool.func import *

def edit_revert(name, num):
    with get_db_connect() as conn:
        curs = conn.cursor()

        curs.execute(db_change("select title from history where title = ? and id = ? and hide = 'O'"), [name, str(num)])
        if curs.fetchall() and admin_check(6) != 1:
            return re_error('/error/3')

        if acl_check(name, 'document_edit') == 1:
            return re_error('/ban')
        
        curs.execute(db_change("select data from history where title = ? and id = ?"), [name, str(num)])
        data = curs.fetchall()
        if not data:
            return redirect('/w/' + url_pas(name))

        if flask.request.method == 'POST':
            if captcha_post(flask.request.form.get('g-recaptcha-response', flask.request.form.get('g-recaptcha', ''))) == 1:
                return re_error('/error/13')
            else:
                captcha_post('', 0)

            if do_edit_slow_check() == 1:
                return re_error('/error/24')
            
            send = flask.request.form.get('send', '')
            agree = flask.request.form.get('copyright_agreement', '')
            
            if do_edit_send_check(send) == 1:
                return re_error('/error/37')
            
            if do_edit_text_bottom_check_box_check(agree) == 1:
                return re_error('/error/29')

            if do_edit_filter(data[0][0]) == 1:
                return re_error('/error/21')
            
            curs.execute(db_change("select data from other where name = 'document_content_max_length'"))
            db_data = curs.fetchall()
            if db_data and db_data[0][0] != '':
                if int(number_check(db_data[0][0])) < len(data[0][0]):
                    return re_error('/error/44')

            curs.execute(db_change("select data from data where title = ?"), [name])
            data_old = curs.fetchall()
            if data_old:
                leng = leng_check(len(data_old[0][0]), len(data[0][0]))
                curs.execute(db_change("update data set data = ? where title = ?"), [data[0][0], name])
            else:
                leng = '+' + str(len(data[0][0]))
                curs.execute(db_change("insert into data (title, data) values (?, ?)"), [name, data[0][0]])

            history_plus(
                name,
                data[0][0],
                get_time(),
                ip_check(),
                flask.request.form.get('send', ''),
                leng,
                t_check = 'r' + str(num),
                mode = 'revert'
            )

            render_set(
                doc_name = name,
                doc_data = data[0][0],
                data_type = 'backlink'
            )

            conn.commit()

            return redirect('/w/' + url_pas(name))
        else:
            if data:
                preview = '<hr class="main_hr"><pre>' + html.escape(data[0][0]) + '</pre>'
            else:
                preview = ''
            
            return easy_minify(flask.render_template(skin_check(),
                imp = [name, wiki_set(), wiki_custom(), wiki_css(['(r' + str(num) + ') (' + load_lang('revert') + ')', 0])],
                data = '''
                    <form method="post">
                        <input placeholder="''' + load_lang('why') + '''" name="send" type="text">
                        <hr class="main_hr">
                        ''' + captcha_get() + ip_warning() + get_edit_text_bottom_check_box() + get_edit_text_bottom() + '''
                        <button type="submit">''' + load_lang('revert') + '''</button>
                    </form>
                ''' + preview,
                menu = [['history/' + url_pas(name), load_lang('history')], ['recent_changes', load_lang('recent_change')]]
            ))