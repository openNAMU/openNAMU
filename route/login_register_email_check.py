from .tool.func import *

async def login_register_email_check():
    with get_db_connect() as conn:
        curs = conn.cursor()

        if not 'reg_email' in flask.session:
            return redirect(conn, '/register')

        if  flask.request.method == 'POST':
            input_key = flask.request.form.get('key', '')

            if flask.session['reg_key'] != input_key:
                return redirect(conn, '/register')

            curs.execute(db_change('select data from other where name = "requires_approval"'))
            sql_data = curs.fetchall()
            if sql_data and sql_data[0][0] != '':
                flask.session['submit_id'] = flask.session['reg_id']
                flask.session['submit_pw'] = flask.session['reg_pw']
                flask.session['submit_email'] = flask.session['reg_email']

                return redirect(conn, '/register/submit')

            add_user(conn, 
                flask.session['reg_id'],
                flask.session['reg_pw'],
                flask.session['reg_email']
            )

            return redirect(conn, '/login')
        else:
            curs.execute(db_change('select data from other where name = "check_key_text"'))
            sql_d = curs.fetchall()
            b_text = (sql_d[0][0] + '<hr class="main_hr">') if sql_d and sql_d[0][0] != '' else ''

            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [get_lang(conn, 'check_key'), await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
                data = '''
                    <form method="post">
                        ''' + b_text + '''
                        <input placeholder="''' + get_lang(conn, 'key') + '''" name="key" type="text">
                        <hr class="main_hr">
                        <button type="submit">''' + get_lang(conn, 'save') + '''</button>
                    </form>
                ''',
                menu = [['user', get_lang(conn, 'return')]]
            ))