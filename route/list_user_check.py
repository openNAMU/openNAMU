from .tool.func import *

async def list_user_check(name = 'test', plus_name = None, arg_num = 1, do_type = 'normal'):
    with get_db_connect() as conn:
        curs = conn.cursor()

        plus_id = plus_name

        check_type = do_type if do_type in ['simple', 'normal'] else 'normal'
        check_type = '' if check_type == 'normal' else check_type

        num = arg_num
        sql_num = (num * 50 - 50) if num * 50 > 0 else 0

        if await acl_check(tool = 'all_admin_auth', ip = name) != 1 or (plus_id and await acl_check(tool = 'all_admin_auth', ip = plus_id) != 1):
            if await acl_check('', 'owner_auth', '', '') == 1:
                return await re_error(conn, 4)

        div = ''

        if await acl_check(tool = 'check_auth', memo = (check_type + ' ' if check_type != '' else '') + 'check (' + name + ')') == 1:
            return await re_error(conn, 3)

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
                        "where " + ('ip' if ip_or_user(name) == 1 else 'name') + " = ? or " + ('ip' if ip_or_user(plus_id) == 1 else 'name') + " = ?"
                    ""), [name, plus_id])
                    all_ip_count = len(curs.fetchall())

                    curs.execute(db_change("" + \
                        "select distinct ip from ua_d " + \
                        "where " + ('ip' if ip_or_user(name) == 1 else 'name') + " = ?" + \
                    ""), [name])
                    a_ip_count = len(curs.fetchall())

                    curs.execute(db_change("" + \
                        "select distinct ip from ua_d " + \
                        "where " + ('ip' if ip_or_user(plus_id) == 1 else 'name') + " = ?"
                    ""), [plus_id])
                    b_ip_count = len(curs.fetchall())

                    if a_ip_count + b_ip_count != all_ip_count:
                        div += get_lang(conn, 'same_ip_exist') + '<hr class="main_hr">'    
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
                        '<a href="/manager/14/' + url_pas(name) + '">(' + get_lang(conn, 'compare') + ')</a> ' + \
                        '<a href="/list/user/check/' + url_pas(name) + '/simple">(' + get_lang(conn, 'simple_check') + ')</a>' + \
                        '<hr class="main_hr">' + \
                    '' + div
                else:
                    div = '' + \
                        '<a href="/list/user/check/' + url_pas(name) + '">(' + name + ')</a> ' + \
                        '<a href="/list/user/check/' + url_pas(plus_id) + '">(' + plus_id + ')</a>' + \
                        '<hr class="main_hr">' + \
                    '' + div

                div += '''
                    <table id="main_table_set">
                        <tbody>
                            <tr id="main_table_top_tr">
                                <td id="main_table_width">''' + get_lang(conn, 'name') + '''</td>
                                <td id="main_table_width">''' + get_lang(conn, 'ip') + '''</td>
                                <td id="main_table_width">''' + get_lang(conn, 'time') + '''</td>
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
                                <a href="/list/user/check/''' + url_pas(data[0]) + '''">''' + data[0] + '''</a>
                                <a href="/list/user/check/delete/''' + url_pas(data[0]) + '/' + url_pas(data[1]) + '/' + url_pas(data[3]) + '/' + ('0' if ip_or_user(name) == 0 else '1') + '''">
                                    (''' + get_lang(conn, 'delete') + ''')
                                </a>
                            </td>
                            <td><a href="/list/user/check/''' + url_pas(data[1]) + '''">''' + data[1] + '''</a></td>
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

            if plus_id:
                div += get_next_page_bottom(conn, 
                    '/list/user/check/' + url_pas(name) + '/normal/{}/' + url_pas(plus_id), 
                    num, 
                    record
                )
            else:
                div += get_next_page_bottom(conn, 
                    '/list/user/check/' + url_pas(name) + '/normal/{}', 
                    num, 
                    record
                )

            if plus_id:
                name += ', ' + plus_id

            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [name, await wiki_set(), await wiki_custom(conn), wiki_css(['(' + get_lang(conn, 'check') + ')', 0])],
                data = div,
                menu = [['manager', get_lang(conn, 'return')]]
            ))
        else:
            curs.execute(db_change("" + \
                "select distinct " + ('name' if ip_or_user(name) == 1 else 'ip') + " from ua_d " + \
                "where " + ('ip' if ip_or_user(name) == 1 else 'name') + " = ? "
                "order by today desc limit ?, 50" + \
            ""), [name, sql_num])
            record = curs.fetchall()

            div = ''
            for for_a in record:
                div += '<li><a href="/list/user/check/' + url_pas(for_a[0]) + '/simple">' + for_a[0] + '</a></li>'

            if div != '':
                div = '<ul>' + div + '</ul>'
                div += get_next_page_bottom(conn, 
                    '/list/user/check/' + url_pas(name) + '/' + check_type + '/{}', 
                    num, 
                    record
                )

            div = '' + \
                '<a href="/list/user/check/' + url_pas(name) + '/normal">(' + get_lang(conn, 'check') + ')</a>' + \
            '' + div

            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [name, await wiki_set(), await wiki_custom(conn), wiki_css(['(' + get_lang(conn, 'simple_check') + ')', 0])],
                data = div,
                menu = [['check/' + url_pas(name), get_lang(conn, 'return')]]
            ))