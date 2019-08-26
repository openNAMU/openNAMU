from .tool.func import *
import pymysql

def topic_2(conn, name, sub):
    curs = conn.cursor()
    
    ban = topic_check(name, sub)
    admin = admin_check(3)

    curs.execute("select id from topic where title = %s and sub = %s limit 1", [name, sub])
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
        
        curs.execute("select id from topic where title = %s and sub = %s order by id + 0 desc limit 1", [name, sub])
        old_num = curs.fetchall()
        if old_num:
            num = int(old_num[0][0]) + 1
        else:
            num = 1

        match = re.search('^user:([^/]+)', name)
        if match:
            y_check = 0
            if ip_or_user(match.groups()[0]) == 1:
                curs.execute("select ip from history where ip = %s limit 1", [match.groups()[0]])
                u_data = curs.fetchall()
                if u_data:
                    y_check = 1
                else:
                    curs.execute("select ip from topic where ip = %s limit 1", [match.groups()[0]])
                    u_data = curs.fetchall()
                    if u_data:
                        y_check = 1
            else:
                curs.execute("select id from user where id = %s", [match.groups()[0]])
                u_data = curs.fetchall()
                if u_data:
                    y_check = 1

            if y_check == 1:
                curs.execute('insert into alarm (name, data, date) values (%s, %s, %s)', [
                    match.groups()[0], 
                    ip + ' - <a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '">' + load_lang('user_discussion', 1) + '</a>', 
                    today
                ])
        
        cate_re = re.compile('\[\[((?:분류|category):(?:(?:(?!\]\]).)*))\]\]', re.I)
        data = cate_re.sub('[br]', flask.request.form.get('content', 'Test'))
        
        for rd_data in re.findall("(?:#([0-9]+))", data):
            curs.execute("select ip from topic where title = %s and sub = %s and id = %s", [name, sub, rd_data])
            ip_data = curs.fetchall()
            if ip_data and ip_or_user(ip_data[0][0]) == 0:
                curs.execute('insert into alarm (name, data, date) values (%s, %s, %s)', [ip_data[0][0], ip + ' - <a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '#' + num + '">' + load_lang('discussion', 1) + '</a>', today])
            
        data = re.sub("(?P<in>#(?:[0-9]+))", '[[\g<in>]]', data)

        data = savemark(data)

        rd_plus(name, sub, today)

        curs.execute("insert into topic (id, title, sub, data, date, ip, block, top) values (%s, %s, %s, %s, %s, %s, '', '')", [num, name, sub, data, today, ip])
        conn.commit()
        
        return redirect('/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '#reload')
    else:
        data = ''
    
        curs.execute("select stop from rd where title = %s and sub = %s and stop != ''", [name, sub])
        close_data = curs.fetchall()
        
        if close_data and admin != 1:
            display = 'display: none;'
        else:
            display = ''

        data += '''
            <div id="top_topic"></div>
            <div id="main_topic"></div>
            <div id="plus_topic"></div>
            <script>topic_top_load("''' + name + '''", "''' + sub + '''");</script>
            <a id="reload" href="javascript:void(0);" onclick="topic_reload();">(''' + load_lang('reload') + ''')</a> <a href="/topic/''' + url_pas(name) + '''/sub/''' + url_pas(sub) + '''/tool">(''' + load_lang('topic_tool') + ''')</a>
            <hr class=\"main_hr\">
            <form style="''' + display + '''" method="post">
                <textarea style="height: 100px;" name="content"></textarea>
                <hr class=\"main_hr\">
                ''' + captcha_get() + (ip_warring() if display == '' else '') + '''
                <button type="submit">''' + load_lang('send') + '''</button>
            </form>
        '''

        return easy_minify(flask.render_template(skin_check(), 
            imp = [name, wiki_set(), custom(), other2([' (' + load_lang('discussion') + ')', 0])],
            data = '''
                <h2 id="topic_top_title">''' + sub + '''</h2>
                ''' + data,
            menu = [['topic/' + url_pas(name), load_lang('list')]]
        ))
