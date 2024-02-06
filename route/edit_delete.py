from .tool.func import *

def edit_delete(name):
    with get_db_connect() as conn:
        curs = conn.cursor()

        ip = ip_check()
        if acl_check(name, 'document_delete') == 1:
            return re_error('/ban')

        curs.execute(db_change("select title from data where title = ?"), [name])
        if not curs.fetchall():
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

            curs.execute(db_change("select data from data where title = ?"), [name])
            data = curs.fetchall()
            if data:
                today = get_time()
                leng = '-' + str(len(data[0][0]))

                history_plus(
                    name,
                    '',
                    today,
                    ip,
                    send,
                    leng,
                    t_check = 'delete',
                    mode = 'delete'
                )

                curs.execute(db_change("select title, link from back where title = ? and not type = 'cat' and not type = 'no'"), [name])
                for data in curs.fetchall():
                    curs.execute(db_change("insert into back (title, link, type, data) values (?, ?, 'no', '')"), [data[0], data[1]])

                curs.execute(db_change("delete from back where link = ?"), [name])
                curs.execute(db_change("delete from data where title = ?"), [name])

                conn.commit()

            return redirect('/w/' + url_pas(name))
        else:            
            return easy_minify(flask.render_template(skin_check(),
                imp = [name, wiki_set(), wiki_custom(), wiki_css(['(' + load_lang('delete') + ')', 0])],
                data = '''
                    <form method="post">
                        <input placeholder="''' + load_lang('why') + '''" name="send">
                        <hr class="main_hr">
                        ''' + captcha_get() + ip_warning() + get_edit_text_bottom_check_box() + get_edit_text_bottom() + '''
                        <button type="submit">''' + load_lang('delete') + '''</button>
                    </form>
                ''',
                menu = [['w/' + url_pas(name), load_lang('return')]]
            ))