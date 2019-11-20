from .tool.func import *

def list_user_topic_2(conn, name):
    curs = conn.cursor()

    num = int(number_check(flask.request.args.get('num', '1')))
    if num * 50 > 0:
        sql_num = num * 50 - 50
    else:
        sql_num = 0
    
    one_admin = admin_check(1)

    div =   '''
            <table id="main_table_set">
                <tbody>
                    <tr>
                        <td id="main_table_width">''' + load_lang('discussion_name') + '''</td>
                        <td id="main_table_width">''' + load_lang('writer') + '''</td>
                        <td id="main_table_width">''' + load_lang('time') + '''</td>
                    </tr>
            '''
    
    curs.execute(db_change("select title, id, sub, ip, date from topic where ip = ? order by date desc limit ?, 50"), [name, sql_num])
    data_list = curs.fetchall()
    for data in data_list:
        title = html.escape(data[0])
        sub = html.escape(data[2])
        
        if one_admin == 1:
            curs.execute(db_change("select * from ban where block = ?"), [data[3]])
            if curs.fetchall():
                ban = ' <a href="/ban/' + url_pas(data[3]) + '">(' + load_lang('release') + ')</a>'
            else:
                ban = ' <a href="/ban/' + url_pas(data[3]) + '">(' + load_lang('ban') + ')</a>'
        else:
            ban = ''
            
        div += '<tr><td><a href="/topic/' + url_pas(data[0]) + '/sub/' + url_pas(data[2]) + '#' + data[1] + '">' + title + '#' + data[1] + '</a> (' + sub + ')</td>'
        div += '<td>' + ip_pas(data[3]) + ban + '</td><td>' + data[4] + '</td></tr>'

    div += '</tbody></table>'
    div += next_fix('/topic_record/' + url_pas(name) + '?num=', num, data_list)      
    
    curs.execute(db_change("select end from ban where block = ?"), [name])
    if curs.fetchall():
        sub = ' (' + load_lang('blocked') + ')'
    else:
        sub = 0 
    
    return easy_minify(flask.render_template(skin_check(), 
        imp = [load_lang('discussion_record'), wiki_set(), custom(), other2([sub, 0])],
        data = div,
        menu = [['other', load_lang('other')], ['user', load_lang('user')], ['count/' + url_pas(name), load_lang('count')], ['record/' + url_pas(name), load_lang('record')]]
    ))