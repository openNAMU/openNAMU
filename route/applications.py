from .tool.func import *

def applications_2(conn):
    curs = conn.cursor()

    div = ''
    admin = admin_check()

    if admin != 1:
        return re_error('/ban')

    curs.execute(db_change('select data from other where name = "requires_approval"'))
    requires_approval = curs.fetchall()
    if requires_approval and requires_approval[0][0] != 'on':
        div += '<p>' + load_lang('approval_requirement_disabled') + '</p>'

    if flask.request.method == 'GET':
        curs.execute(db_change('select id, date, question, answer, token from user_application'))
        db_data = curs.fetchall()
        
        if db_data:
            div += '''
                <form method=\"post\">
                    <table id=\"main_table_set\">
                        <tbody>
                            <tr>
                                <td>''' + load_lang('id') + '''</td>
                                <td>''' + load_lang('application_time') + '''</td>
                                <td>''' + load_lang('approval_question') + '''</td>
                                <td>''' + load_lang('answer') + '''</td>
                                <td>''' + load_lang('approve_or_decline') + '''</td>
                            </tr>
            '''
            for application in db_data:
                question = application[2]
                answer = application[3]
                if not question or question == '':
                    question = '(질문 없음)'
                if not answer or answer == '':
                    answer = '(대답 없음)'
                div += '''
                    <tr>
                        <td>''' + application[0] + '''</td>
                        <td>''' + application[1] + '''</td>
                        <td>''' + question + '''</td>
                        <td>''' + answer + '''</td>
                        <td>
                            <button type=\"submit\" name=\"approve\" value=\"''' + application[4] + '''\">''' + load_lang('approve') + '''</button>
                            <button type=\"submit\" name=\"decline\" value=\"''' + application[4] + '''\">''' + load_lang('decline') + '''</button>
                        </td>
                    </tr>
                '''
            div += '</tbody></table></form>'
        else:
            div += load_lang('no_applications_now')
    else:
        if flask.request.form.get('approve', '') != '':
            curs.execute(db_change('select id, pw, date, encode, question, answer, ip, ua from user_application where token = ?'), [flask.request.form.get('approve', '')])
            application = curs.fetchall()
            if not application:
                return re_error('/error/26')
            
            application = application[0]

            curs.execute(db_change("select id from user where id = ?"), [application[0]])
            if curs.fetchall():
                return re_error('/error/6')
            
            curs.execute(db_change("insert into user (id, pw, acl, date, encode) values (?, ?, 'user', ?, ?)"), [application[0], application[1], application[2], application[3]])
            curs.execute(db_change("insert into user_set (name, id, data) values ('approval_question', ?, ?)"), [application[0], application[4]])
            curs.execute(db_change("insert into user_set (name, id, data) values ('approval_question_answer', ?, ?)"), [application[0], application[5]])
            curs.execute(db_change("insert into ua_d (name, ip, ua, today, sub) values (?, ?, ?, ?, '')"), [application[0], application[6], application[7], application[2]])
            curs.execute(db_change('delete from user_application where token = ?'), [flask.request.form.get('approve', '')])
 

            conn.commit()
        elif flask.request.form.get('decline', '') != '':
            curs.execute(db_change('delete from user_application where token = ?'), [flask.request.form.get('decline', '')])
            conn.commit()
        return redirect('/applications')

    return easy_minify(flask.render_template(skin_check(),
        imp = [load_lang('application_list'), wiki_set(), custom(), other2([0, 0])],
        data = div,
        menu = [['other', load_lang('return')]]
    ))