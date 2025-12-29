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

            return await render_template(
                await get_lang('check_key'),
                '''
                    <form method="post">
                        ''' + b_text + '''
                        <input class="__ON_INPUT__" placeholder="''' + await get_lang('key') + '''" name="key" type="text">
                        <hr class="main_hr">
                        <button class="__ON_BUTTON__" type="submit">''' + await get_lang('save') + '''</button>
                    </form>
                ''',
                0,
                [['user', await get_lang('return')]]
            )
