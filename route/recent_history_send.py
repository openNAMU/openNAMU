from .tool.func import *

def recent_history_send(name = 'Test', rev = 1):
    with get_db_connect() as conn:
        curs = conn.cursor()

        num = str(rev)

        if admin_check() != 1:
            return re_error('/error/3')

        if flask.request.method == 'POST':
            admin_check(None, 'send edit ' + name + ' r' + num)

            curs.execute(db_change("select send from history where title = ? and id = ?"), [name, num])
            if curs.fetchall():
                curs.execute(db_change("update history set send = ? where title = ? and id = ?"), [
                    flask.request.form.get('send', ''),
                    name, 
                    num
                ])

            conn.commit()

            return redirect('/history/' + url_pas(name))
        else:
            curs.execute(db_change("select send from history where title = ? and id = ?"), [name, num])
            send = curs.fetchall()
            if send:
                send = send[0][0]

                return easy_minify(flask.render_template(skin_check(),
                    imp = [name, wiki_set(), wiki_custom(), wiki_css(['(' + load_lang('send_edit') + ') (r' + num + ')', 0])],
                    data = '''
                        <form method="post">
                            <span>''' + load_lang('delete_warning') + '''</span>
                            <hr class="main_hr">
                            <input value="''' + html.escape(send) + '''" name="send">
                            <hr class="main_hr">
                            <button type="submit">''' + load_lang('edit') + '''</button>
                        </form>
                    ''',
                    menu = [['history/' + url_pas(name), load_lang('return')]]
                ))
            else:
                return redirect('/history/' + url_pas(name))