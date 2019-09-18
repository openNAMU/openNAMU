from .tool.func import *

def recent_changes_2(conn, name, tool):
    curs = conn.cursor()

    if flask.request.method == 'POST':
        return redirect(
            '/diff/' + url_pas(name) +
            '?first=' + flask.request.form.get('b', '1') +
            '&second=' + flask.request.form.get('a', '1')
        )
    else:
        one_admin = admin_check(1)
        six_admin = admin_check(6)
        
        ban = ''
        select = ''

        div = '''
            <hr class=\"main_hr\">
            <table id="main_table_set">
                <tbody>
                    <tr>
        '''
        
        if name:
            num = int(number_check(flask.request.args.get('num', '1')))
            if num * 50 > 0:
                sql_num = num * 50 - 50
            else:
                sql_num = 0      

            if tool == 'history':
                div += '''
                    <td id="main_table_width">''' + load_lang('version') + '''</td>
                    <td id="main_table_width">''' + load_lang('editor') + '''</td>
                    <td id="main_table_width">''' + load_lang('time') + '''</td></tr>
                '''
                
                # 기본적인 move만 구현
                tool_select = flask.request.args.get('tool', None)
                if tool_select:
                    if tool_select == 'move':
                        curs.execute('''
                            select id, title, date, ip, send, leng from history
                            where send like ? or send like ?
                            order by id + 0 desc
                            limit ?, '50'
                        ''', ['%(<a>' + name +'</a>%', '%<a>' + name + '</a> move)', str(sql_num)])
                    else:
                        curs.execute('''
                            select id, title, date, ip, send, leng from history
                            where title = ?
                            order by id + 0 desc
                            limit ?, '50'
                        ''', [name, str(sql_num)])
                else:
                    curs.execute('''
                        select id, title, date, ip, send, leng from history
                        where title = ?
                        order by id + 0 desc
                        limit ?, '50'
                    ''', [name, str(sql_num)])
            else:
                div +=  '''
                        <td id="main_table_width">''' + load_lang('document_name') + '''</td>
                        <td id="main_table_width">''' + load_lang('editor') + '''</td>
                        <td id="main_table_width">''' + load_lang('time') + '''</td>
                    </tr>
                '''

                div = '<a href="/topic_record/' + url_pas(name) + '">(' + load_lang('discussion') + ')</a><hr class=\"main_hr\">' + div
                
                curs.execute("select id, title, date, ip, send, leng from history where ip = ? order by date desc limit ?, '50'", [name, str(sql_num)])
        else:
            num = int(number_check(flask.request.args.get('num', '1')))
            if num * 50 > 0:
                sql_num = num * 50 - 50
            else:
                sql_num = 0            
            
            div +=  '''
                    <td id="main_table_width">''' + load_lang('document_name') + '''</td>
                    <td id="main_table_width">''' + load_lang('editor') + '''</td>
                    <td id="main_table_width">''' + load_lang('time') + '''</td>
                </tr>
            '''
            if flask.request.args.get('set', 'normal') == 'user':
                curs.execute('''
                    select id, title, date, ip, send, leng from history 
                    where title like 'user:%' 
                    order by date desc 
                    limit ?, 50
                ''', [str(sql_num)])
            else:
                div = '<a href="?set=user">(' + load_lang('user_document') + ')</a>' + div

                curs.execute('''
                    select id, title, date, ip, send, leng from history 
                    where not title like 'user:%' 
                    order by date desc 
                    limit ?, 50
                ''', [str(sql_num)])

        data_list = curs.fetchall()
        for data in data_list:    
            select += '<option value="' + data[0] + '">' + data[0] + '</option>'     
            send = '<br>'
            
            if data[4]:
                if not re.search("^(?: *)$", data[4]):
                    send = data[4]
            
            if re.search("\+", data[5]):
                leng = '<span style="color:green;">(' + data[5] + ')</span>'
            elif re.search("\-", data[5]):
                leng = '<span style="color:red;">(' + data[5] + ')</span>'
            else:
                leng = '<span style="color:gray;">(' + data[5] + ')</span>'
                
            ip = ip_pas(data[3])
            if tool == 'history':
                m_tool = '<a href="/history_tool/' + url_pas(data[1]) + '?num=' + data[0] + '&type=history">(' + load_lang('tool') + ')</a>'
            else:
                m_tool = '<a href="/history_tool/' + url_pas(data[1]) + '?num=' + data[0] + '">(' + load_lang('tool') + ')</a>'
            
            style = ['', '']
            date = data[2]

            curs.execute('''
                select title from history
                where title = ? and id = ? and hide = 'O'
            ''', [data[1], data[0]])
            hide = curs.fetchall()
            
            if six_admin == 1:
                if hide:                    
                    style[0] = 'id="toron_color_grey"'
                    style[1] = 'id="toron_color_grey"'
                    
                    if send == '<br>':
                        send = '(' + load_lang('hide') + ')'
                    else:
                        send += ' (' + load_lang('hide') + ')'
            elif not hide:
                pass
            else:
                ip = ''
                ban = ''
                date = ''

                send = '(' + load_lang('hide') + ')'

                style[0] = 'style="display: none;"'
                style[1] = 'id="toron_color_grey"'

            if tool == 'history':
                title = '<a href="/w/' + url_pas(name) + '?num=' + data[0] + '">r' + data[0] + '</a> '
            else:
                title = '<a href="/w/' + url_pas(data[1]) + '">' + html.escape(data[1]) + '</a> '
                title += '<a href="/history/' + url_pas(data[1]) + '">(r' + data[0] + ')</a> '

            div +=  '''
                <tr ''' + style[0] + '''>
                    <td>''' + title + m_tool + ' ' + leng + '''</td>
                    <td>''' + ip + ban + '''</td>
                    <td>''' + date + '''</td>
                </tr>
                <tr ''' + style[1] + '''>
                    <td colspan="3">''' + send_parser(send) + '''</td>
                </tr>
            '''

        div +=  '''
                </tbody>
            </table>
        '''
        sub = ''

        if name:
            if tool == 'history':
                if not tool_select:
                    div = '''
                        <a href="?tool=move">(''' + load_lang('move') + ''')</a>
                        <hr class=\"main_hr\">
                    ''' + div
                    
                div = '''
                    <form method="post">
                        <select name="a">''' + select + '''</select> <select name="b">''' + select + '''</select>
                        <button type="submit">''' + load_lang('compare') + '''</button>
                    </form>
                    <hr class=\"main_hr\">
                ''' + div
                title = name
                
                sub += ' (' + load_lang('history') + ')'
                
                menu = [['w/' + url_pas(name), load_lang('document')], ['raw/' + url_pas(name), load_lang('raw')]]
                
                div += next_fix('/history/' + url_pas(name) + '?num=', num, data_list)
            else:
                curs.execute("select end from ban where block = ?", [name])
                if curs.fetchall():
                    sub += ' (' + load_lang('blocked') + ')'

                title = load_lang('edit_record')
                
                menu = [['other', load_lang('other')], ['user', load_lang('user')], ['count/' + url_pas(name), load_lang('count')]]
                
                div += next_fix('/record/' + url_pas(name) + '?num=', num, data_list)
        else:
            menu = 0
            title = load_lang('recent_change')
                
            div += next_fix('/recent_changes?num=', num, data_list)

            if flask.request.args.get('set', 'normal') == 'user':
                sub = ' (' + load_lang('user') + ')'
                menu = [['recent_changes', load_lang('return')]]
        
        if sub == '':
            sub = 0
                
        return easy_minify(flask.render_template(skin_check(), 
            imp = [title, wiki_set(), custom(), other2([sub, 0])],
            data = div,
            menu = menu
        ))