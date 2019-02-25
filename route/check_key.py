from .tool.func import *

def check_key_2(conn, tool):
    curs = conn.cursor()

    if flask.request.method == 'POST':
        if tool == 'check_key':
            if 'c_id' in flask.session and flask.session['c_key'] == flask.request.form.get('key', None):
                curs.execute('select data from other where name = "encode"')
                db_data = curs.fetchall()

                curs.execute("select id from user limit 1")
                if not curs.fetchall():
                    curs.execute("insert into user (id, pw, acl, date, encode) values (?, ?, 'owner', ?)", [flask.session['c_id'], flask.session['c_pw'], get_time(), db_data[0][0]])

                    first = 1
                else:
                    curs.execute("insert into user (id, pw, acl, date, encode) values (?, ?, 'user', ?)", [flask.session['c_id'], flask.session['c_pw'], get_time(), db_data[0][0]])

                    first = 0

                ip = ip_check()
                agent = flask.request.headers.get('User-Agent')

                curs.execute("insert into user_set (name, id, data) values ('email', ?, ?)", [flask.session['c_id'], flask.session['c_email']])
                curs.execute("insert into ua_d (name, ip, ua, today, sub) values (?, ?, ?, ?, '')", [flask.session['c_id'], ip, agent, get_time()])

                flask.session['state'] = 1
                flask.session['id'] = flask.session['c_id']
                flask.session['head'] = ''
                        
                conn.commit()
                
                flask.session.pop('c_id', None)
                flask.session.pop('c_pw', None)
                flask.session.pop('c_key', None)
                flask.session.pop('c_email', None)

                if first == 0:
                    return redirect('/change')
                else:
                    return redirect('/setting')
            else:
                flask.session.pop('c_id', None)
                flask.session.pop('c_pw', None)
                flask.session.pop('c_key', None)
                flask.session.pop('c_email', None)

                return redirect('/register')
        else:
            if 'c_id' in flask.session and flask.session['c_key'] == flask.request.form.get('key', None):
                hashed = pw_encode(flask.session['c_key'])
                curs.execute("update user set pw = ? where id = ?", [hashed, flask.session['c_id']])

                d_id = flask.session['c_id']
                pw = flask.session['c_key']

                flask.session.pop('c_id', None)
                flask.session.pop('c_key', None)

                return easy_minify(flask.render_template(skin_check(),    
                    imp = ['check', wiki_set(), custom(), other2([0, 0])],
                    data =  '''
                            ''' + load_lang('id') + ' : ' + d_id + '''
                            <br>
                            ''' + load_lang('password') + ' : ' + pw + '''
                            ''',
                    menu = [['user', load_lang('return')]]
                ))
            else:
                return redirect('/pass_find')
    else:
        return easy_minify(flask.render_template(skin_check(),    
            imp = ['check', wiki_set(), custom(), other2([0, 0])],
            data =  '''
                    <form method="post">
                        <input placeholder="''' + load_lang('key') + '''" name="key" type="text">
                        <hr class=\"main_hr\">
                        <button type="submit">''' + load_lang('save') + '''</button>
                    </form>
                    ''',
            menu = [['user', load_lang('return')]]
        ))