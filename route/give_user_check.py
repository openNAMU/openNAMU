from .tool.func import *

def give_user_check_2(conn, name):
    curs = conn.cursor()

    plus_id = flask.request.args.get('plus', None)

    if admin_check('all', None, name) == 1 or (plus_id and admin_check('all', None, plus_id) == 1):
        if admin_check() != 1:
            return re_error('/error/4')

    num = int(number_check(flask.request.args.get('num', '1')))
    sql_num = (num * 50 - 50) if num * 50 > 0 else 0
    
    div = ''
    check_type = flask.request.args.get('type', '')

    if admin_check(4, (check_type + ' ' if check_type != '' else '') + 'check (' + name + ')') != 1:
        return re_error('/error/3')

    if check_type == '':
        if ip_or_user(name) == 0:
            curs.execute(db_change("select data from user_set where name = \"approval_question\" and id = ?"), [name])
            approval_question = curs.fetchall()
            if approval_question and approval_question[0][0]:
                curs.execute(db_change("select data from user_set where name = \"approval_question_answer\" and id = ?"), [name])
                approval_question_answer = curs.fetchall()
                if approval_question_answer and approval_question_answer[0]:
                    div += '''
                        <table id="main_table_set">
                            <tbody>
                                <tr id="main_table_top_tr">
                                    <td>Q</td>
                                    <td>''' + approval_question[0][0] + '''</td>
                                    <td>A</td>
                                    <td>''' + approval_question_answer[0][0] + '''</td>
                                </tr>
                            </tbody>
                        </table>
                        <hr class="main_hr">
                    '''

        if plus_id:
            plus = "or " + ('ip' if ip_or_user(plus_id) == 1 else 'name') + " = ? "
            set_list = [name, plus_id, sql_num]
            
            if num == 1:
                curs.execute(db_change("" + \
                    "select distinct ip from ua_d " + \
                    "where " + ('ip' if ip_or_user(name) == 1 else 'name') + " = ? or " + ('ip' if ip_or_user(plus_id) == 1 else 'name') + " = ? "
                ""), [name, plus_id])
                all_ip_count = len(curs.fetchall())
                
                curs.execute(db_change("" + \
                    "select distinct ip from ua_d " + \
                    "where " + ('ip' if ip_or_user(name) == 1 else 'name') + " = ?" + \
                ""), [name])
                a_ip_count = len(curs.fetchall())
                
                curs.execute(db_change("" + \
                    "select distinct ip from ua_d " + \
                    "where " + ('ip' if ip_or_user(plus_id) == 1 else 'name') + " = ? "
                ""), [plus_id])
                b_ip_count = len(curs.fetchall())
                
                if a_ip_count + b_ip_count != all_ip_count:
                    div += load_lang('same_ip_exist') + '<hr class="main_hr">'    
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
                    '<a href="/manager/14?plus=' + url_pas(name) + '">(' + load_lang('compare') + ')</a> ' + \
                    '<a href="/check/' + url_pas(name) + '?type=simple">(' + load_lang('simple_check') + ')</a>' + \
                    '<hr class="main_hr">' + \
                '' + div
            else:
                div = '' + \
                    '<a href="/check/' + url_pas(name) + '">(' + name + ')</a> ' + \
                    '<a href="/check/' + url_pas(plus_id) + '">(' + plus_id + ')</a>' + \
                    '<hr class="main_hr">' + \
                '' + div

            div += '''
                <table id="main_table_set">
                    <tbody>
                        <tr id="main_table_top_tr">
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
                        <td>
                            <a href="/check/''' + url_pas(data[0]) + '''">''' + data[0] + '''</a>
                            <a  href="/check_delete''' + \
                                '''?name=''' + url_pas(data[0]) + \
                                '''&ip=''' + url_pas(data[1]) + \
                                '''&time=''' + url_pas(data[3].replace(' ', '').replace(':', '').replace('-', '')) + \
                                '''&return_type=''' + ('0' if ip_or_user(name) == 0 else '1') + '''">
                                (''' + load_lang('delete') + ''')
                            </a>
                        </td>
                        <td><a href="/check/''' + url_pas(data[1]) + '''">''' + data[1] + '''</a></td>
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
            
        div += next_fix(
            '/check/' + url_pas(name) + ('?plus=' + plus_id + '&num=' if plus_id else '?num='), 
            num, 
            record
        )

        return easy_minify(flask.render_template(skin_check(),
            imp = [load_lang('check'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
            data = div,
            menu = [['manager', load_lang('return')]]
        ))
    else:
        curs.execute(db_change("" + \
            "select distinct " + ('name' if ip_or_user(name) == 1 else 'ip') + " from ua_d " + \
            "where " + ('ip' if ip_or_user(name) == 1 else 'name') + " = ? "
            "order by today desc limit ?, 50" + \
        ""), [name, sql_num])
        record = curs.fetchall()

        div = ''
        for i in record:
            div += '<li><a href="/check/' + url_pas(i[0]) + '?type=simple">' + i[0] + '</a></li>'

        if div != '':
            div = '<ul class="inside_ul">' + div + '</ul>'
            div += next_fix(
                '/check/' + url_pas(name) + '?type=' + check_type + '&num=', 
                num, 
                record
            )

        return easy_minify(flask.render_template(skin_check(),
            imp = [name, wiki_set(), wiki_custom(), wiki_css(['(' + load_lang('simple_check') + ')', 0])],
            data = div,
            menu = [['check/' + url_pas(name), load_lang('return')]]
        ))