from .tool.func import *

from .tool.func import *

def login_register_submit_2(conn):
    curs = conn.cursor()
    
    if not 'submit_id' in flask.session:
        return redirect('/register')
    
    curs.execute(db_change('select data from other where name = "approval_question"'))
    sql_data = curs.fetchall()
    if not sql_data:
        return redirect('/register')
    
    data_que = sql_data[0][0]
    
    if flask.request.method == 'POST':
        curs.execute(db_change('select data from other where name = "encode"'))
        data_encode = curs.fetchall()
        data_encode = data_encode[0][0]
        
        user_ip = ip_check()
        user_agent = flask.request.headers.get('User-Agent', '')
        
        user_app_data = {}
        user_app_data['id'] = flask.session['submit_id']
        user_app_data['pw'] = pw_encode(flask.session['submit_pw'])
        user_app_data['encode'] = data_encode
        user_app_data['question'] = data_que
        user_app_data['answer'] = flask.request.form.get('answer', '')
        
        if 'submit_email' in flask.session:
            user_app_data['email'] = flask.session['submit_email']
        else:
            user_app_data['email'] = ''
            
        curs.execute(db_change(
            "insert into user_set (id, name, data) values (?, ?, ?)"
        ), [
            flask.session['submit_id'],
            'application',
            json.dumps(user_app_data)
        ])
        conn.commit()
        
        return redirect('/')
    else:
        return easy_minify(flask.render_template(skin_check(),
            imp = [load_lang('approval_question'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
            data = '''
                <form method="post">
                    ''' + load_lang('approval_question') + ' : ' + data_que + '''
                    <hr class="main_hr">
                    <input placeholder="''' + load_lang('approval_question') + '''" name="answer">
                    <hr class="main_hr">
                    <button type="submit">''' + load_lang('save') + '''</button>
                </form>
            ''',
            menu = [['user', load_lang('return')]]
        ))
