from .tool.func import *

def topic_2(conn, name, sub):
    curs = conn.cursor()
    
    ban = topic_check(name, sub)
    admin = admin_check(3)

    curs.execute("select id from topic where title = ? and sub = ? limit 1", [name, sub])
    topic_exist = curs.fetchall()
    if not topic_exist and len(sub) > 256:
        return re_error('/error/11')
    
    if flask.request.method == 'POST':
        if captcha_post(flask.request.form.get('g-recaptcha-response', '')) == 1:
            return re_error('/error/13')
        else:
            captcha_post('', 0)

        ip = ip_check()
        today = get_time()
        
        if ban == 1:
            return re_error('/ban')
        
        curs.execute("select id from topic where title = ? and sub = ? order by id + 0 desc limit 1", [name, sub])
        old_num = curs.fetchall()
        if old_num:
            num = int(old_num[0][0]) + 1
        else:
            num = 1

        match = re.search('^user:([^/]+)', name)
        if match:
            curs.execute('insert into alarm (name, data, date) values (?, ?, ?)', [match.groups()[0], ip + ' - <a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '">' + load_lang('user_discussion', 1) + '</a>', today])
        
        cate_re = re.compile('\[\[((?:분류|category):(?:(?:(?!\]\]).)*))\]\]', re.I)
        data = cate_re.sub('[br]', flask.request.form.get('content', 'Test'))
        for rd_data in re.findall("(?:#([0-9]+))", data):
            curs.execute("select ip from topic where title = ? and sub = ? and id = ?", [name, sub, rd_data])
            ip_data = curs.fetchall()
            if ip_data and ip_or_user(ip_data[0][0]) == 0:
                curs.execute('insert into alarm (name, data, date) values (?, ?, ?)', [ip_data[0][0], ip + ' - <a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '#' + str(num) + '">' + load_lang('discussion', 1) + '</a>', today])
            
        data = re.sub("(?P<in>#(?:[0-9]+))", '[[\g<in>]]', data)

        data = savemark(data)

        rd_plus(name, sub, today)

        curs.execute("insert into topic (id, title, sub, data, date, ip, block, top) values (?, ?, ?, ?, ?, ?, '', '')", [str(num), name, sub, data, today, ip])
        conn.commit()
        
        return redirect('/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '#reload')
    else:
        curs.execute("select title from rd where title = ? and sub = ? and stop = 'O'", [name, sub])
        close_data = curs.fetchall()
        
        curs.execute("select title from rd where title = ? and sub = ? and stop = 'S'", [name, sub])
        stop_data = curs.fetchall()
        
        display = ''
        all_data = ''
        data = ''
        number = 1
        
        if admin == 1 and topic_exist:
            if close_data:
                all_data += '<a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/tool/close">(' + load_lang('open') + ')</a> '
            else:
                all_data += '<a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/tool/close">(' + load_lang('close') + ')</a> '
            
            if stop_data:
                all_data += '<a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/tool/stop">(' + load_lang('restart') + ')</a> '
            else:
                all_data += '<a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/tool/stop">(' + load_lang('stop') + ')</a> '

            curs.execute("select title from rd where title = ? and sub = ? and agree = 'O'", [name, sub])
            if curs.fetchall():
                all_data += '<a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/tool/agree">(' + load_lang('destruction') + ')</a>'
            else:
                all_data += '<a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/tool/agree">(' + load_lang('agreement') + ')</a>'
            
            all_data += '<hr class=\"main_hr\">'
        
        if (close_data or stop_data) and admin != 1:
            display = 'display: none;'
        
        curs.execute("select data, id, date, ip, block, top from topic where title = ? and sub = ? order by id + 0 asc", [name, sub])
        topic = curs.fetchall()
        
        curs.execute("select data, id, date, ip from topic where title = ? and sub = ? and top = 'O' order by id + 0 asc", [name, sub])
        for topic_data in curs.fetchall():                   
            who_plus = ''
            
            curs.execute("select who from re_admin where what = ? order by time desc limit 1", ['notice (' + name + ' - ' + sub + '#' + topic_data[1] + ')'])
            topic_data_top = curs.fetchall()
            if topic_data_top:
                who_plus += ' <span style="margin-right: 5px;">@' + topic_data_top[0][0] + ' </span>'
                                
            all_data += '''
                <table id="toron">
                    <tbody>
                        <tr>
                            <td id="toron_color_red">
                                <a href="#''' + topic_data[1] + '''">
                                    #''' + topic_data[1] + '''
                                </a> ''' + ip_pas(topic_data[3]) + who_plus + ''' <span style="float: right;">''' + topic_data[2] + '''</span>
                            </td>
                        </tr>
                        <tr>
                            <td>''' + render_set(data = topic_data[0]) + '''</td>
                        </tr>
                    </tbody>
                </table>
                <br>
            '''    

        for topic_data in topic:
            user_write = topic_data[0]

            if number == 1:
                start = topic_data[3]

            if topic_data[4] == 'O':
                blind_data = 'id="toron_color_grey"'
                
                if admin != 1:
                    curs.execute("select who from re_admin where what = ? order by time desc limit 1", ['blind (' + name + ' - ' + sub + '#' + str(number) + ')'])
                    who_blind = curs.fetchall()
                    if who_blind:
                        user_write = '[[user:' + who_blind[0][0] + ']] ' + load_lang('hide')
                    else:
                        user_write = load_lang('hide')
            else:
                blind_data = ''

            user_write = render_set(data = user_write)
            ip = ip_pas(topic_data[3])
            
            curs.execute('select acl from user where id = ?', [topic_data[3]])
            user_acl = curs.fetchall()
            if user_acl and user_acl[0][0] != 'user':
                ip += ' <a href="javascript:void(0);" title="' + load_lang('admin') + '">★</a>'

            if admin == 1 or blind_data == '':
                ip += ' <a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/admin/' + str(number) + '">(' + load_lang('discussion_tool') + ')</a>'

            curs.execute("select end from ban where block = ?", [topic_data[3]])
            if curs.fetchall():
                ip += ' <a href="javascript:void(0);" title="' + load_lang('blocked') + '">†</a>'
                    
            if topic_data[5] == '1':
                color = '_blue'
            elif topic_data[3] == start:
                color = '_green'
            else:
                color = ''
                
            if user_write == '':
                user_write = '<br>'
                         
            all_data += '''
                <table id="toron">
                    <tbody>
                        <tr>
                            <td id="toron_color''' + color + '''">
                                <a href="javascript:void(0);" id="''' + str(number) + '">#' + str(number) + '</a> ' + ip + '''</span>
                            </td>
                        </tr>
                        <tr ''' + blind_data + '''>
                            <td>''' + user_write + '''</td>
                        </tr>
                    </tbody>
                </table>
                <br>
            '''
            number += 1

        if ban != 1 or admin == 1:
            data += '''
                <div id="plus"></div>
                <script>topic_load("''' + name + '''", "''' + sub + '''");</script>
                <a id="reload" href="javascript:void(0);" onclick="location.href.endsWith(\'#reload\')? location.reload(true):location.href=\'#reload\'">(''' + load_lang('reload') + ''')</a>
                <form style="''' + display + '''" method="post">
                <br>
                <textarea style="height: 100px;" name="content"></textarea>
                <hr class=\"main_hr\">
            ''' + captcha_get()
            
            if display == '':
                data += ip_warring()

            data += '''
                    <button type="submit">''' + load_lang('send') + '''</button>
                </form>
            '''

        return easy_minify(flask.render_template(skin_check(), 
            imp = [name, wiki_set(), custom(), other2([' (' + load_lang('discussion') + ')', 0])],
            data = '<h2 id="topic_top_title">' + sub + '</h2>' + all_data + data,
            menu = [['topic/' + url_pas(name), load_lang('list')]]
        ))
