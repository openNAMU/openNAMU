from .tool.func import *

def applications_2(conn):
    # 만들다만 느낌이니 수정 필요
    curs = conn.cursor()

    div = ''

    if admin_check() != 1:
        return re_error('/ban')

    curs.execute(db_change('select data from other where name = "requires_approval"'))
    requires_approval = curs.fetchall()
    if requires_approval and requires_approval[0][0] != 'on':
        div += '<p>' + load_lang('approval_requirement_disabled') + '</p>'

    if flask.request.method == 'GET':
        curs.execute(db_change(
            'select data from user_set where name = "application"'
        ))
        db_data = curs.fetchall()
        if db_data:
            div += '' + \
                '<p>' + load_lang('all_register_num') + ' : ' + str(len(db_data)) + '</p>' + \
                '<hr class="main_hr">' + \
            ''

            for application in db_data:
                application = json.loads(application)
                
                question = application['question']
                if not question:
                    question = ''
                    
                answer = application['answer']
                if not answer:
                    answer = ''
                
                email = application['email']
                if not email:
                    email = ''
                
                div += '''
                    <form method=\"post\">
                        <table>
                            <tbody>
                                <tr>
                                    <td>''' + load_lang('id') + '''</td>
                                    <td>''' + application['id'] + '''</td>
                                </tr>
                                <tr>
                                    <td>''' + load_lang('application_time') + '''</td>
                                    <td>''' + application['date'] + '''</td>
                                </tr>
                                <tr>
                                    <td>''' + load_lang('approval_question') + '''</td>
                                    <td>''' + html.escape(question) + '''</td>
                                </tr>
                                <tr>
                                    <td>''' + load_lang('answer') + '''</td>
                                    <td>''' + html.escape(answer) + '''</td>
                                </tr>
                                <tr>
                                    <td>''' + load_lang('email') + '''</td>
                                    <td>''' + html.escape(email) + '''</td>
                                </tr>
                                <tr>
                                    <td colspan="2" style="text-align: center;">
                                        <button type="submit" 
                                                name="approve" 
                                                value="''' + application['id'] + '''">
                                            ''' + load_lang('approve') + '''
                                        </button>
                                        <button type="submit" 
                                                name="decline" 
                                                value="''' + application['id'] + '''">
                                            ''' + load_lang('decline') + '''
                                        </button>
                                    </td>
                                </tr>
                            </tbody>
                        </table>
                    </form>
                    <br>
                '''
        else:
            div += load_lang('no_applications_now')

        return easy_minify(flask.render_template(skin_check(),
            imp = [load_lang('application_list'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
            data = div,
            menu = [['other', load_lang('return')]]
        ))
    else:
        if flask.request.form.get('approve', '') != '':
            curs.execute(db_change(
                'select data from user_set where id = ? and name = "application"'
            ), [
                flask.request.form.get('approve', '')
            ])
            application = curs.fetchall()
            if not application:
                return re_error('/error/26')
            else:
                application = json.loads(application[0][0])
            
            curs.execute(db_change(
                "insert into user_set (id, name, data) values (?, 'pw', ?)"
            ), [
                application['id'],
                application['pw']
            ])
            curs.execute(db_change(
                "insert into user_set (id, name, data) values (?, 'acl', 'user')"
            ), [
                application['id']
            ])
            curs.execute(db_change(
                "insert into user_set (id, name, data) values (?, 'date', ?)"
            ), [
                application['id'],
                application['date']
            ])
            curs.execute(db_change(
                "insert into user_set (id, name, data) values (?, 'encode', ?)"
            ), [
                application['id'],
                application['encode']
            ])
            curs.execute(db_change(
                "insert into user_set (name, id, data) values ('approval_question', ?, ?)"
            ), [
                application['id'], 
                application['question']
            ])
            curs.execute(db_change(
                "insert into user_set (name, id, data) " + \
                "values ('approval_question_answer', ?, ?)"
            ), [
                application['id'], 
                application['answer']
            ])
            
            ua_plus(
                application['id'], 
                application['ip'], 
                application['ua'], 
                application['date']
            )
            
            if application['email'] != '':
                curs.execute(db_change(
                    "insert into user_set (name, id, data) values ('email', ?, ?)"
                ), [
                    application['id'], 
                    application['email']
                ])
            
            curs.execute(db_change(
                'delete from user_set where id = ? and name = "application"'
            ), [
                application['id']
            ])
            conn.commit()
        elif flask.request.form.get('decline', '') != '':
             curs.execute(db_change(
                'delete from user_set where id = ? and name = "application"'
            ), [
                flask.request.form.get('decline', '')
            ])
            conn.commit()

        return redirect('/applications')