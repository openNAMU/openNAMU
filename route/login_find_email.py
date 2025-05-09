from .tool.func import *

# 개편 필요
async def login_find_email(tool):
    with get_db_connect() as conn:
        curs = conn.cursor()
        
        if flask.request.method == 'POST':
            re_set_list = ['c_id', 'c_pw', 'c_ans', 'c_que', 'c_key', 'c_type']
        
            if tool == 'email_change':
                flask.session['c_key'] = load_random_key(32)
                flask.session['c_id'] = ip_check()
                flask.session['c_type'] = 'email_change'
            elif tool == 'pass_find':
                user_id = flask.request.form.get('id', '')
                user_email = flask.request.form.get('email', '')
        
                flask.session['c_key'] = load_random_key(32)
                flask.session['c_id'] = user_id
                flask.session['c_type'] = 'pass_find'
            else:
                if not 'c_type' in flask.session:
                    return redirect(conn, '/register')
        
            if tool != 'pass_find':
                user_email = flask.request.form.get('email', '')
                email_data = re.search(r'@([^@]+)$', user_email)
                if email_data:
                    curs.execute(db_change("select html from html_filter where html = ? and kind = 'email'"), [email_data.group(1)])
                    if not curs.fetchall():
                        for i in re_set_list:
                            flask.session.pop(i, None)
                        
                        return redirect(conn, '/filter/email_filter')
                else:
                    for i in re_set_list:
                        flask.session.pop(i, None)
                    
                    return await re_error(conn, 36)
        
            curs.execute(db_change('select data from other where name = "email_title"'))
            sql_d = curs.fetchall()
            t_text = html.escape(sql_d[0][0]) if sql_d and sql_d[0][0] != '' else ((await wiki_set())[0] + ' key')
        
            curs.execute(db_change('select data from other where name = "email_text"'))
            sql_d = curs.fetchall()
            i_text = (html.escape(sql_d[0][0]) + '\n\nKey : ' + flask.session['c_key']) if sql_d and sql_d[0][0] != '' else ('Key : ' + flask.session['c_key'])
            
            if tool == 'pass_find':
                curs.execute(db_change("select id from user_set where id = ? and name = 'email' and data = ?"), [user_id, user_email])
                if not curs.fetchall():
                    return await re_error(conn, 12)
                    
                if await send_email(conn, user_email, t_text, i_text) == 0:
                    return await re_error(conn, 18)
        
                return redirect(conn, '/pass_find/email')
            else:
                curs.execute(db_change('select id from user_set where name = "email" and data = ?'), [user_email])
                if curs.fetchall():
                    for i in re_set_list:
                        flask.session.pop(i, None)
        
                    return await re_error(conn, 35)
                
                if await send_email(conn, user_email, t_text, i_text) == 0:
                    for i in re_set_list:
                        flask.session.pop(i, None)
        
                    return await re_error(conn, 18)
        
                flask.session['c_email'] = user_email
        
                return redirect(conn, '/pass_find/email')
        else:
            if tool == 'pass_find':
                curs.execute(db_change('select data from other where name = "password_search_text"'))
                sql_d = curs.fetchall()
                b_text = (sql_d[0][0] + '<hr class="main_hr">') if sql_d and sql_d[0][0] != '' else ''
        
                return easy_minify(conn, flask.render_template(skin_check(conn),
                    imp = [get_lang(conn, 'password_search'), await wiki_set(), await wiki_custom(conn), wiki_css(['(' + get_lang(conn, 'email') + ')', 0])],
                    data = b_text + '''
                        <form method="post">
                            <input placeholder="''' + get_lang(conn, 'id') + '''" name="id" type="text">
                            <hr class="main_hr">
                            <input placeholder="''' + get_lang(conn, 'email') + '''" name="email" type="text">
                            <hr class="main_hr">
                            <button type="submit">''' + get_lang(conn, 'save') + '''</button>
                        </form>
                    ''',
                    menu = [['user', get_lang(conn, 'return')]]
                ))
            else:
                if tool == 'need_email' and not 'c_type' in flask.session:
                    return redirect(conn, '/register')
        
                curs.execute(db_change('select data from other where name = "email_insert_text"'))
                sql_d = curs.fetchall()
                b_text = (sql_d[0][0] + '<hr class="main_hr">') if sql_d and sql_d[0][0] != '' else ''
        
                return easy_minify(conn, flask.render_template(skin_check(conn),
                    imp = [get_lang(conn, 'email'), await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
                    data = '''
                        <a href="/filter/email_filter">(''' + get_lang(conn, 'email_filter_list') + ''')</a>
                        <hr class="main_hr">
                        ''' + b_text + '''
                        <form method="post">
                            <input placeholder="''' + get_lang(conn, 'email') + '''" name="email" type="text">
                            <hr class="main_hr">
                            <button type="submit">''' + get_lang(conn, 'save') + '''</button>
                        </form>
                    ''',
                    menu = [['user', get_lang(conn, 'return')]]
                ))