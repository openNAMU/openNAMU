from .tool.func import *
import pymysql

def give_user_check_2(conn, name):
    curs = conn.cursor()

    curs.execute("select acl from user where id = %s or id = %s", [name, flask.request.args.get('plus', '-')])
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
    
    if flask.request.args.get('plus', None):
        end_check = 1
    
        if ip_or_user(name) == 1:
            if ip_or_user(flask.request.args.get('plus', None)) == 1:
                curs.execute("select name, ip, ua, today from ua_d where ip = %s or ip = %s order by today desc limit %s, 50", [name, flask.request.args.get('plus', None), sql_num])
            else:
                curs.execute("select name, ip, ua, today from ua_d where ip = %s or name = %s order by today desc limit %s, 50", [name, flask.request.args.get('plus', None), sql_num])
        else:
            if ip_or_user(flask.request.args.get('plus', None)) == 1:
                curs.execute("select name, ip, ua, today from ua_d where name = %s or ip = %s order by today desc limit %s, 50", [name, flask.request.args.get('plus', None), sql_num])
            else:
                curs.execute("select name, ip, ua, today from ua_d where name = %s or name = %s order by today desc limit %s, 50", [name, flask.request.args.get('plus', None), sql_num])
    else:
        end_check = 0
        
        if ip_or_user(name) == 1:
            curs.execute("select name, ip, ua, today from ua_d where ip = %s order by today desc limit %s, 50", [name, sql_num])
        else:
            curs.execute("select name, ip, ua, today from ua_d where name = %s order by today desc limit %s, 50", [name, sql_num])
    
    record = curs.fetchall()
    if record:
        if not flask.request.args.get('plus', None):
            div = '<a href="/manager/14?plus=' + url_pas(name) + '">(' + load_lang('compare') + ')</a><hr class=\"main_hr\">'
        else:
            div = '<a href="/check/' + url_pas(name) + '">(' + name + ')</a> <a href="/check/' + url_pas(flask.request.args.get('plus', None)) + '">(' + flask.request.args.get('plus', None) + ')</a><hr class=\"main_hr\">'

        div +=  '''
                <table id="main_table_set">
                    <tbody>
                        <tr>
                            <td id="main_table_width">''' + load_lang('name') + '''</td>
                            <td id="main_table_width">ip</td>
                            <td id="main_table_width">''' + load_lang('time') + '''</td>
                        </tr>
                '''
        
        for data in record:
            if data[2]:
                ua = data[2]
            else:
                ua = '<br>'

            div +=  '''
                    <tr>
                        <td>''' + ip_pas(data[0]) + '''</td>
                        <td>''' + ip_pas(data[1]) + '''</td>
                        <td>''' + data[3] + '''</td>
                    </tr>
                    <tr>
                        <td colspan="3">''' + ua + '''</td>
                    </tr>
                    '''
        
        div +=  '''
                    </tbody>
                </table>
                '''
    else:
        return re_error('/error/2')
        
    if end_check == 1:
        div += next_fix('/check/' + url_pas(name) + '?plus=' + flask.request.args.get('plus', None) + '&num=', num, record)
    else:
        div += next_fix('/check/' + url_pas(name) + '?num=', num, record)
            
    return easy_minify(flask.render_template(skin_check(),    
        imp = [load_lang('check'), wiki_set(), custom(), other2([0, 0])],
        data = div,
        menu = [['manager', load_lang('return')]]
    ))