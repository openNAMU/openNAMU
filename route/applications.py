from .tool.func import *

def applications_2(conn):
    curs = conn.cursor()

    div = ''

    if admin_check() != 1:
        return re_error('/ban')

    curs.execute(db_change('select data from other where name = "requires_approval"'))
    requires_approval = curs.fetchall()
    if requires_approval and requires_approval[0][0] != 'on':
        div += load_lang('approval_requirement_disabled')

    if flask.request.method == 'GET':
        curs.execute(db_change(
            'select data from user_set where name = "application"'
        ))
        db_data = curs.fetchall()
        if db_data:
            div += '' + \
                load_lang('all_register_num') + ' : ' + str(len(db_data)) + \
                '<hr class="main_hr">' + \
            ''
            
            div += '''
                <table id="main_table_set">
                    <tr id="main_table_top_tr">
                        <td id="main_table_width">''' + load_lang('id') + '''</td>
                        <td id="main_table_width">''' + load_lang('email') + '''</td>
                        <td id="main_table_width">''' + load_lang('application_time') + '''</td>
                    </tr>
                    <tr id="main_table_top_tr">
                        <td>''' + load_lang('approval_question') + '''</td>
                        <td colspan="2">''' + load_lang('answer') + '''</td>
                    </tr>                        
            '''

            for application in db_data:
                application = json.loads(application[0])
                
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
                            <td>''' + application['date'] + '''</td>
                        </tr>
                        <tr>
                            <td>''' + question + '''</td>
                            <td colspan="2">''' + answer + '''</td>
                        </tr>
                        <tr>
                            <td colspan="3">
                                <button type="submit" 
                                        id="save"
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
                    </form>
                '''
                
            div += '</table>'
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
