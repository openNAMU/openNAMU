from .tool.func import *

def block_log_2(conn, name, tool):
    curs = conn.cursor()

    num = int(number_check(flask.request.args.get('num', '1')))
    if num * 50 > 0:
        sql_num = num * 50 - 50
    else:
        sql_num = 0
    
    div = '''
        <table id="main_table_set">
            <tbody>
                <tr>
                    <td id="main_table_width">''' + load_lang('blocked') + '''</td>
                    <td id="main_table_width">''' + load_lang('admin') + '''</td>
                    <td id="main_table_width">''' + load_lang('period') + '''</td>
                </tr>
    '''
    
    data_list = ''
    
    if not name:
        div = '''
            <a href="/manager/11">(''' + load_lang('blocked') + ''')</a> <a href="/manager/12">(''' + load_lang('admin') + ''')</a>
            <hr class=\"main_hr\">
        ''' + div
        
        sub = 0
        menu = 0
        
        curs.execute("select why, block, blocker, end, today from rb order by today desc limit ?, '50'", [str(sql_num)])
    else:
        menu = [['block_log', load_lang('normal')]]
        
        if tool == 'block_user':
            sub = ' (' + load_lang('blocked') + ')'
            
            curs.execute("select why, block, blocker, end, today from rb where block = ? order by today desc limit ?, '50'", [name, str(sql_num)])
        else:
            sub = ' (' + load_lang('admin') + ')'
            
            curs.execute("select why, block, blocker, end, today from rb where blocker = ? order by today desc limit ?, '50'", [name, str(sql_num)])

    if data_list == '':
        data_list = curs.fetchall()

    for data in data_list:
        why = html.escape(data[0])
        if why == '':
            why = '<br>'
        
        band = re.search("^([0-9]{1,3}\.[0-9]{1,3})$", data[1])
        if band:
            ip = data[1] + ' (' + load_lang('range') + ')'
        else:
            ip = ip_pas(data[1])

        if data[3] != '':
            end = data[3]
        else:
            end = load_lang('limitless') + ''
            
        div += '''
            <tr>
                <td>''' + ip + '''</td>
                <td>''' + ip_pas(data[2]) + '''</td>
                <td>
                    start : ''' + data[4] + '''
                    <br>
                    end : ''' + end + '''
                </td>
            </tr>
            <tr>
                <td colspan="3">''' + why + '''</td>
            </tr>
        '''

    div += '</tbody></table>'
    
    if not name:
        div += next_fix('/block_log?num=', num, data_list)
    else:
        div += next_fix('/' + url_pas(tool) + '/' + url_pas(name) + '?num=', num, data_list)
                
    return easy_minify(flask.render_template(skin_check(), 
        imp = [load_lang('recent_ban'), wiki_set(), custom(), other2([sub, 0])],
        data = div,
        menu = menu
    ))