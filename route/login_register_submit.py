from .tool.func import *

def login_register_submit_2(conn):
    approval_question = ''

    curs.execute(db_change('select data from other where name = "requires_approval"'))
    requires_approval = curs.fetchall()
    requires_approval = requires_approval and requires_approval[0][0] == 'on'
    requires_approval = None if admin == 1 else requires_approval
    if requires_approval:
        curs.execute(db_change('select data from other where name = "approval_question"'))
        data = curs.fetchall()
        if data and data[0][0] != '':
            approval_question = '''
                <hr class="main_hr">
                <span>''' + load_lang('approval_question') + ' : ' + data[0][0] + '''<span>
                <hr class="main_hr">
                <input placeholder="''' + load_lang('approval_question') + '''" name="approval_question_answer" type="text">
                <hr class="main_hr">
            '''

    ans_q = flask.request.form.get('approval_question_answer', '')

    curs.execute(db_change('select data from other where name = "requires_approval"'))
    requires_approval = curs.fetchall()
    requires_approval = requires_approval and requires_approval[0][0] == 'on'
    requires_approval = None if admin == 1 else requires_approval
    if requires_approval:
        curs.execute(db_change('select data from other where name = "approval_question"'))
        approval_question = curs.fetchall()
        approval_question = approval_question[0][0] if approval_question and approval_question[0][0] else ''
    else:
        approval_question = ''

            if requires_approval:
                flask.session['c_ans'] = flask.request.form.get('approval_question_answer', '')
                flask.session['c_que'] = approval_question