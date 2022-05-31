from .tool.func import *

def topic(topic_num = 0):
    with get_db_connect() as conn:
        curs = conn.cursor()
        topic_num = str(topic_num)

        if flask.request.method == 'POST':
            name = flask.request.form.get('topic', 'Test')
            sub = flask.request.form.get('title', 'Test')
            
            if topic_num == '0':
                curs.execute(db_change("select code from topic order by code + 0 desc limit 1"))
                t_data = curs.fetchall()
                topic_num = str(int(t_data[0][0]) + 1) if t_data else '1'
        else:
            if topic_num == '0':
                name = load_lang('make_new_topic')
                sub = load_lang('make_new_topic')
            else:
                curs.execute(db_change("select title, sub from rd where code = ?"), [topic_num])
                name = curs.fetchall()
                if name:
                    sub = name[0][1]
                    name = name[0][0]
                else:
                    return redirect('/')

        ban = acl_check(name, 'topic', topic_num)

        if flask.request.method == 'POST':
            if flask.request.form.get('content', 'Test') == '':
                return redirect('/thread/' + topic_num)

            if captcha_post(flask.request.form.get('g-recaptcha-response', flask.request.form.get('g-recaptcha', ''))) == 1:
                return re_error('/error/13')
            else:
                captcha_post('', 0)

            ip = ip_check()
            today = get_time()

            if ban == 1:
                return re_error('/ban')

            curs.execute(db_change("select id from topic where code = ? order by id + 0 desc limit 1"), [topic_num])
            old_num = curs.fetchall()
            num = str((int(old_num[0][0]) + 1) if old_num else 1)

            match = re.search(r'^user:([^/]+)', name)
            if match:
                match = match.group(1)
                y_check = 0
                if ip_or_user(match) == 1:
                    curs.execute(db_change("select ip from history where ip = ? limit 1"), [match])
                    u_data = curs.fetchall()
                    if u_data:
                        y_check = 1
                    else:
                        curs.execute(db_change("select ip from topic where ip = ? limit 1"), [match])
                        u_data = curs.fetchall()
                        if u_data:
                            y_check = 1
                else:
                    curs.execute(db_change("select id from user_set where id = ?"), [match])
                    u_data = curs.fetchall()
                    if u_data:
                        y_check = 1

                if y_check == 1:
                    add_alarm(match, ip + ' | <a href="/thread/' + topic_num + '#' + num + '">' + name + ' | ' + sub + ' | #' + num + '</a>')

            cate_re = re.compile(r'\[\[((?:분류|category):(?:(?:(?!\]\]).)*))\]\]', re.I)
            data = cate_re.sub('[br]', flask.request.form.get('content', 'Test').replace('\r', ''))

            for rd_data in re.findall(r"(?: |\n|^)(#(?:[0-9]+))(?: |\n|$)", data):
                curs.execute(db_change("select ip from topic where code = ? and id = ?"), [topic_num, rd_data])
                ip_data = curs.fetchall()
                if ip_data and ip_or_user(ip_data[0][0]) == 0:
                    add_alarm(ip_data[0][0], ip + ' | <a href="/thread/' + topic_num + '#' + num + '">' + name + ' | ' + sub + ' | #' + num + '</a>')

            for rd_data in re.findall(r"(?: |\n|^)@((?:[^ ]+))(?: |\n|$)", data):
                curs.execute(db_change("select ip from history where ip = ? limit 1"), [rd_data])
                ip_data = curs.fetchall()
                if not ip_data:
                    curs.execute(db_change("select ip from topic where ip = ? limit 1"), [rd_data])
                    ip_data = curs.fetchall()

                if ip_data and ip_or_user(ip_data[0][0]) == 0:
                    add_alarm(ip_data[0][0], ip + ' | <a href="/thread/' + topic_num + '#' + num + '">' + name + ' | ' + sub + ' | #' + num + '</a>')

            data = re.sub(r"( |\n|^)(#(?:[0-9]+))( |\n|$)", '\g<1><topic_a>\g<2></topic_a>\g<3>', data)
            data = re.sub(r"( |\n|^)(@(?:[^ ]+))( |\n|$)", '\g<1><topic_call>\g<2></topic_call>\g<3>', data)

            rd_plus(topic_num, today, name, sub)
            curs.execute(db_change("insert into topic (id, data, date, ip, code) values (?, ?, ?, ?, ?)"), [
                num,
                data,
                today,
                ip,
                topic_num
            ])
            conn.commit()

            return redirect('/thread/' + topic_num + '#' + num)
        else:
            display = 'display: none;' if ban == 1 else ''
            data_input_topic_name = ''
            if topic_num == '0':
                data_input_topic_name = '' + \
                    '<input placeholder="' + load_lang('discussion_name') + '" name="title">' + \
                    '<hr class="main_hr">' + \
                    '<input placeholder="' + load_lang('document_name') + '" name="topic">' + \
                    '<hr class="main_hr">' + \
                ''
                
            curs.execute(db_change('select data from other where name = "topic_text"'))
            sql_d = curs.fetchall()
            topic_text = html.escape(sql_d[0][0]) if sql_d and sql_d[0][0] != '' else load_lang('content')

            return easy_minify(flask.render_template(skin_check(),
                imp = [name, wiki_set(), wiki_custom(), wiki_css(['(' + load_lang('discussion') + ')', 0])],
                data = '''
                    <h2 id="topic_top_title">''' + html.escape(sub) + '''</h2>
                    <div id="top_topic"></div>
                    <div id="main_topic"></div>
                    <div id="plus_topic"></div>
                    <script>opennamu_do_thread_make("''' + topic_num + '''");</script>
                    <a href="/thread/''' + topic_num + '/tool">(' + load_lang('topic_tool') + ''')</a>
                    <hr class="main_hr">
                    <form style="''' + display + '''" method="post">
                        ''' + data_input_topic_name + '''
                        <textarea id="textarea_edit_view" class="opennamu_comment_textarea" placeholder="''' + topic_text + '''" name="content"></textarea>
                        <hr class="main_hr">
                        ''' + captcha_get() + (ip_warning() if display == '' else '') + '''
                        <input style="display: none;" name="topic" value="''' + name + '''">
                        <input style="display: none;" name="title" value="''' + sub + '''">
                        <button id="save" type="submit">''' + load_lang('send') + '''</button>
                        <button id="preview" type="button" onclick="load_preview(\'\')">''' + load_lang('preview') + '''</button>
                    </form>
                    <hr class="main_hr">
                    <div id="see_preview"></div>
                ''',
                menu = [['topic/' + url_pas(name), load_lang('list')]]
            ))