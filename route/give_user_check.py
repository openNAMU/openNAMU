from .tool.func import *

def give_user_check_2(conn, name):
    curs = conn.cursor()

    curs.execute(db_change("select acl from user where id = ? or id = ?"), [name, flask.request.args.get('plus', '-')])
    user = curs.fetchall()
    if user and user[0][0] != 'user':
        if admin_check() != 1:
            return re_error('/error/4')

    if admin_check(4, 'check (' + name + ')') != 1:
        return re_error('/error/3')

    num = int(number_check(flask.request.args.get('num', '1')))
    if num * 50 > 0:
        sql_num = num * 50 - 50
    else:
        sql_num = 0

    div = ''
    plus_id = flask.request.args.get('plus', None)

    if ip_or_user(name) == 0:
        curs.execute(db_change("select data from user_set where name = \"approval_question\" and id = ?"), [name])
        approval_question = curs.fetchall()
        if approval_question and approval_question[0][0]:
            curs.execute(db_change("select data from user_set where name = \"approval_question_answer\" and id = ?"), [name])
            approval_question_answer = curs.fetchall()
            if approval_question_answer and approval_question_answer[0]:
                div = '''
                    <table id="main_table_set">
                        <tbody>
                            <tr>
                                <td>Q</td>
                                <td>''' + approval_question[0][0] + '''</td>
                                <td>A</td>
                                <td>''' + approval_question_answer[0][0] + '''</td>
                            </tr>
                        </tbody>
                    </table>
                    <hr class=\"main_hr\">
                '''

    if plus_id:
        plus = "or " + ('ip' if ip_or_user(plus_id) == 1 else 'name') + " = ? "
        set_list = [name, plus_id, sql_num]
    else:
        plus = ''
        set_list = [name, sql_num]

    curs.execute(db_change("" + \
        "select name, ip, ua, today from ua_d " + \
        "where " + ('ip' if ip_or_user(name) == 1 else 'name') + " = ? " + \
        plus + \
        "order by today desc limit ?, 50" + \
    ""), set_list)

    record = curs.fetchall()
    if record:
        if not plus_id:
            div = '' + \
                '<a href="/manager/14?plus=' + url_pas(name) + '">(' + load_lang('compare') + ')</a>' + \
                '<hr class=\"main_hr\">' + \
            '' + div
        else:
            div = '' + \
                '<a href="/check/' + url_pas(name) + '">(' + name + ')</a> ' + \
                '<a href="/check/' + url_pas(plus_id) + '">(' + plus_id + ')</a>' + \
                '<hr class=\"main_hr\">' + \
            '' + div

        div += '''
            <table id="main_table_set">
                <tbody>
                    <tr>
                        <td id="main_table_width">''' + load_lang('name') + '''</td>
                        <td id="main_table_width">''' + load_lang('ip') + '''</td>
                        <td id="main_table_width">''' + load_lang('time') + '''</td>
                    </tr>
        '''

        set_n = 0
        for data in record:
            if data[2]:
                if len(data[2]) > 300:
                    ua = '' + \
                        '<a href="javascript:void();" onclick="document.getElementById(\'check_' + str(set_n) + '\').style.display=\'block\';">(300+)</a>' + \
                        '<div id="check_' + str(set_n) + '" style="display:none;">' + html.escape(data[2]) + '</div>' + \
                    ''
                    set_n += 1
                else:
                    ua = html.escape(data[2])
            else:
                ua = '<br>'

            div += '''
                <tr>
                    <td>''' + ip_pas(data[0]) + '''</td>
                    <td>''' + ip_pas(data[1]) + '''</td>
                    <td>''' + data[3] + '''</td>
                </tr>
                <tr>
                    <td colspan="3">''' + ua + '''</td>
                </tr>
            '''

        div += '''
                </tbody>
            </table>
        '''
    else:
        return re_error('/error/2')

    div += next_fix(
        '/check/' + url_pas(name) + ('?plus=' + plus_id if plus_id else '') + '&num=', 
        num, 
        record
    )

    return easy_minify(flask.render_template(skin_check(),
        imp = [load_lang('check'), wiki_set(), custom(), other2([0, 0])],
        data = div,
        menu = [['manager', load_lang('return')]]
    ))
