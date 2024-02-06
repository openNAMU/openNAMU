from .tool.func import *

from .api_topic import api_topic, api_topic_thread_pre_render

from .edit import edit_editor

def topic(topic_num = 0, do_type = '', doc_name = 'Test'):
    with get_db_connect() as conn:
        curs = conn.cursor()
        topic_num = str(topic_num)

        topic_acl = acl_check('', 'topic', topic_num)
        topic_view_acl = acl_check('', 'topic_view', topic_num)
        if topic_view_acl == 1:
            return re_error('/ban')

        ip = ip_check()

        if flask.request.method == 'POST' and do_type == '':
            if do_edit_slow_check('thread') == 1:
                return re_error('/error/42')

            name = flask.request.form.get('topic', 'Test')
            sub = flask.request.form.get('title', 'Test')
            
            if do_title_length_check(name) == 1:
                return re_error('/error/38')
            
            if do_title_length_check(sub, 'topic') == 1:
                return re_error('/error/38')
            
            if topic_num == '0':
                curs.execute(db_change("select code from topic order by code + 0 desc limit 1"))
                t_data = curs.fetchall()
                topic_num = str(int(t_data[0][0]) + 1) if t_data else '1'
            
            if flask.request.form.get('content', 'Test') == '':
                return redirect('/thread/' + topic_num)

            if captcha_post(flask.request.form.get('g-recaptcha-response', flask.request.form.get('g-recaptcha', ''))) == 1:
                return re_error('/error/13')
            else:
                captcha_post('', 0)

            today = get_time()

            if topic_acl == 1:
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
                    add_alarm(match, ip, '<a href="/thread/' + topic_num + '#' + num + '">' + html.escape(name) + ' - ' + html.escape(sub) + '#' + num + '</a>')
            
            curs.execute(db_change("select ip from topic where code = ? and id = '1'"), [topic_num])
            ip_data = curs.fetchall()
            if ip_data and ip_or_user(ip_data[0][0]) == 0:
                add_alarm(ip_data[0][0], ip, '<a href="/thread/' + topic_num + '#' + num + '">' + html.escape(name) + ' - ' + html.escape(sub) + '#' + num + '</a>')

            data = flask.request.form.get('content', 'Test').replace('\r', '')
            data = api_topic_thread_pre_render(curs, data, num, ip, topic_num, name, sub)

            do_add_thread(
                topic_num,
                data,
                '',
                num
            )
            do_reload_recent_thread(
                topic_num, 
                today, 
                name, 
                sub
            )

            conn.commit()

            return redirect('/thread/' + topic_num + '#' + num)
        else:
            thread_data = ''

            if topic_num == '0':
                name = load_lang('make_new_topic')
                sub = load_lang('make_new_topic')

                name_value = doc_name
                sub_value = ''
            else:
                curs.execute(db_change("select title, sub from rd where code = ?"), [topic_num])
                name = curs.fetchall()
                if name:
                    sub = name[0][1]
                    name = name[0][0]

                    name_value = name
                    sub_value = sub
                else:
                    return redirect('/')

            acl_display = 'display: none;' if topic_acl == 1 else ''
            name_display = 'display: none;' if topic_num != '0' else ''

            shortcut = '<div class="opennamu_thread_shortcut" id="thread_shortcut">'
            curs.execute(db_change("select id from topic where code = ? order by id + 0 asc"), [topic_num])
            db_data = curs.fetchall()
            for for_a in db_data:
                shortcut += '<a href="#' + for_a[0] + '">#' + for_a[0] + '</a> '
            
            shortcut += '</div>'

            top_topic = ''
            json_data = json.loads(api_topic(int(topic_num), 'top', '', 'render').data)
            top_topic += json_data['data'] if 'data' in json_data else ''
            
            main_topic = ''
            json_data = json.loads(api_topic(int(topic_num), 'normal', '', 'render').data)
            main_topic += json_data['data'] if 'data' in json_data else ''

            return easy_minify(flask.render_template(skin_check(),
                imp = [name, wiki_set(), wiki_custom(), wiki_css(['(' + load_lang('discussion') + ')', 0])],
                data = '''
                    <style id="opennamu_remove_blind">
                        .opennamu_comment_blind_js {
                            display: none;
                        }
                    </style>
                    <input type="checkbox" onclick="opennamu_do_remove_blind_thread();" checked> ''' + load_lang('remove_blind_thread') + '''

                    ''' + shortcut + '''
                    <h2 id="topic_top_title">''' + html.escape(sub) + '''</h2>
                    
                    <div id="top_topic">''' + top_topic + '''</div>
                    <div id="main_topic">''' + main_topic + '''</div>
                    <div id="plus_topic"></div>
                    
                    <a href="/thread/''' + topic_num + '/tool">(' + load_lang('topic_tool') + ''')</a>
                    <hr class="main_hr">
                    
                    <form style="''' + acl_display + '''" method="post">
                        <div style="''' + name_display + '''">
                            <input placeholder="''' + load_lang('document_name') + '''" name="topic" value="''' + html.escape(name_value) + '''">
                            <hr class="main_hr">
                            <input placeholder="''' + load_lang('discussion_name') + '''" name="title" value="''' + html.escape(sub_value) + '''">
                            <hr class="main_hr">
                        </div>
                        
                        ''' + edit_editor(curs, ip, thread_data, 'thread') + '''
                    </form>
                ''',
                menu = [['topic/' + url_pas(name), load_lang('list')]]
            ))