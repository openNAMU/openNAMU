from .tool.func import *

def login_check_key_2(tool):
    

    if flask.request.method == 'POST':
        if tool == 'check_pass_key':
            if 'c_id' in flask.session and flask.session['c_key'] == flask.request.form.get('key', None):
                hashed = pw_encode(flask.session['c_key'])

                sqlQuery("update user set pw = ? where id = ?", [hashed, flask.session['c_id']])
                sqlQuery("commit")

                d_id = flask.session['c_id']
                pw = flask.session['c_key']

                flask.session.pop('c_id', None)
                flask.session.pop('c_key', None)

                sqlQuery('select data from other where name = "reset_user_text"')
                sql_d = sqlQuery("fetchall")
                if sql_d and sql_d[0][0] != '':
                    b_text = sql_d[0][0] + '<hr class=\"main_hr\">'
                else:
                    b_text = ''

                return easy_minify(flask.render_template(skin_check(),    
                    imp = [load_lang('reset_user_ok'), wiki_set(), custom(), other2([0, 0])],
                    data = b_text + load_lang('id') + ' : ' + d_id + '<br>' + load_lang('password') + ' : ' + pw,
                    menu = [['user', load_lang('return')]]
                ))
            else:
                return redirect('/pass_find')
        else:
            ip = ip_check()
            
            if 'c_id' in flask.session and flask.session['c_key'] == flask.request.form.get('key', None):
                sqlQuery('select data from other where name = "encode"')
                db_data = sqlQuery("fetchall")
                
                if tool == 'check_key':
                    sqlQuery("select id from user limit 1")
                    if not sqlQuery("fetchall"):
                        sqlQuery("insert into user (id, pw, acl, date, encode) values (?, ?, 'owner', ?, ?)", [
                            flask.session['c_id'], 
                            flask.session['c_pw'], 
                            get_time(), 
                            db_data[0][0]
                        ])
    
                        first = 1
                    else:
                        sqlQuery("insert into user (id, pw, acl, date, encode) values (?, ?, 'user', ?, ?)", [
                            flask.session['c_id'], 
                            flask.session['c_pw'], 
                            get_time(), 
                            db_data[0][0]
                        ])
    
                        first = 0
    
                    agent = flask.request.headers.get('User-Agent')
    
                    sqlQuery("insert into user_set (name, id, data) values ('email', ?, ?)", [
                        flask.session['c_id'], 
                        flask.session['c_email']
                    ])
                    sqlQuery("insert into ua_d (name, ip, ua, today, sub) values (?, ?, ?, ?, '')", [
                        flask.session['c_id'], 
                        ip, 
                        agent, 
                        get_time()
                    ])
    
                    flask.session['state'] = 1
                    flask.session['id'] = flask.session['c_id']
                    flask.session['head'] = ''
                            
                    sqlQuery("commit")
                else:
                    sqlQuery('delete from user_set where name = "email" and id = ?', [ip])
                    sqlQuery('insert into user_set (name, id, data) values ("email", ?, ?)', [ip, flask.session['c_email']])
                    
                    first = 0
                          
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

                return redirect('/user')
    else:
        sqlQuery('select data from other where name = "check_key_text"')
        sql_d = sqlQuery("fetchall")
        if sql_d and sql_d[0][0] != '':
            b_text = sql_d[0][0] + '<hr class=\"main_hr\">'
        else:
            b_text = ''

        return easy_minify(flask.render_template(skin_check(),    
            imp = [load_lang('check_key'), wiki_set(), custom(), other2([0, 0])],
            data =  '''
                <form method="post">
                    ''' + b_text + '''
                    <input placeholder="''' + load_lang('key') + '''" name="key" type="text">
                    <hr class=\"main_hr\">
                    <button type="submit">''' + load_lang('save') + '''</button>
                </form>
            ''',
            menu = [['user', load_lang('return')]]
        ))