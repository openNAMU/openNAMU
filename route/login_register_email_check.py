from .tool.func import *

def login_register_email_check_2():
    with get_db_connect() as conn:
        curs = conn.cursor()

        if not 'reg_email' in flask.session:
            return redirect('/register')

        if  flask.request.method == 'POST':
            input_key = flask.request.form.get('key', '')

            if flask.session['reg_key'] != input_key:
                return redirect('/register')

            curs.execute(db_change('select data from other where name = "requires_approval"'))
            sql_data = curs.fetchall()
            if sql_data and sql_data[0][0] != '':
                flask.session['submit_id'] = flask.session['reg_id']
                flask.session['submit_pw'] = flask.session['reg_pw']
                flask.session['submit_email'] = flask.session['reg_email']

                return redirect('/register/submit')

            add_user(
                flask.session['reg_id'],
                flask.session['reg_pw'],
                flask.session['reg_email']
            )

            return redirect('/login')
        else:
            curs.execute(db_change('select data from other where name = "check_key_text"'))
            sql_d = curs.fetchall()
            b_text = (sql_d[0][0] + '<hr class="main_hr">') if sql_d and sql_d[0][0] != '' else ''

            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('check_key'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
                data = '''
                    <form method="post">
                        ''' + b_text + '''
                        <input placeholder="''' + load_lang('key') + '''" name="key" type="text">
                        <hr class="main_hr">
                        <button type="submit">''' + load_lang('save') + '''</button>
                    </form>
                ''',
                menu = [['user', load_lang('return')]]
            ))