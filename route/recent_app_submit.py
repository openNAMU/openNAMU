from .tool.func import *

async def recent_app_submit():
    with get_db_connect() as conn:
        curs = conn.cursor()

        div = ''

        curs.execute(db_change('select data from other where name = "requires_approval"'))
        requires_approval = curs.fetchall()
        if requires_approval and requires_approval[0][0] != 'on':
            div += await get_lang('approval_requirement_disabled')

        if flask.request.method == 'GET':
            curs.execute(db_change('select data from user_set where name = "application"'))
            db_data = curs.fetchall()
            if db_data:
                div += '' + \
                    await get_lang('all_register_num') + ' : ' + str(len(db_data)) + \
                    '<hr class="main_hr">' + \
                ''

                div += '''
                    <table id="main_table_set">
                        <tr id="main_table_top_tr">
                            <td id="main_table_width_half">''' + await get_lang('id') + '''</td>
                            <td id="main_table_width_half">''' + await get_lang('email') + '''</td>
                        </tr>
                        <tr id="main_table_top_tr">
                            <td>''' + await get_lang('approval_question') + '''</td>
                            <td>''' + await get_lang('answer') + '''</td>
                        </tr>                        
                '''

                for application in db_data:
                    application = json_loads(application[0])

                    if 'question' in application:
                        question = html.escape(application['question'])
                        question = question if question != '' else '<br>'
                    else:
                        question = '<br>'

                    if 'answer' in application:
                        answer = html.escape(application['answer'])
                        answer = answer if answer != '' else '<br>'
                    else:
                        answer = '<br>'


                    if 'email' in application:
                        email = html.escape(application['email'])
                        email = email if email != '' else '<br>'
                    else:
                        email = '<br>'

                    div += '''
                        <form method="post">
                            <tr>
                                <td>''' + application['id'] + '''</td>
                                <td>''' + email + '''</td>
                            </tr>
                            <tr>
                                <td>''' + question + '''</td>
                                <td>''' + answer + '''</td>
                            </tr>
                            <tr>
                                <td colspan="3">
                                    <button type="submit" 
                                            id="opennamu_save_button"
                                            name="approve" 
                                            value="''' + application['id'] + '''">
                                        ''' + await get_lang('approve') + '''
                                    </button>
                                    <button type="submit" 
                                            name="decline" 
                                            value="''' + application['id'] + '''">
                                        ''' + await get_lang('decline') + '''
                                    </button>
                                </td>
                            </tr>
                        </form>
                    '''

                div += '</table>'
            else:
                div += await get_lang('no_applications_now')

            return await render_template(
                await get_lang('application_list'),
                div,
                0,
                [['other', await get_lang('return')]]
            )
        else:
            if await acl_check(tool = 'ban_auth', memo = 'app submit') == 1:
                return await re_error(conn, 0)

            if flask.request.form.get('approve', '') != '':
                curs.execute(db_change('select data from user_set where id = ? and name = "application"'), [flask.request.form.get('approve', '')])
                application = curs.fetchall()
                if not application:
                    return await re_error(conn, 26)
                else:
                    application = json_loads(application[0][0])

                add_user(conn, application['id'], application['pw'], application['email'], application['encode'])

                curs.execute(db_change("insert into user_set (name, id, data) values ('approval_question', ?, ?)"), [application['id'], application['question']])
                curs.execute(db_change("insert into user_set (name, id, data) values ('approval_question_answer', ?, ?)"), [application['id'], application['answer']])

                curs.execute(db_change('delete from user_set where id = ? and name = "application"'), [application['id']])
            elif flask.request.form.get('decline', '') != '':
                curs.execute(db_change('delete from user_set where id = ? and name = "application"'), [flask.request.form.get('decline', '')])

            return redirect(conn, '/app_submit')
