from .tool.func import *

from .go_api_topic import api_topic_thread_pre_render

from .edit import edit_editor

async def topic(topic_num = 0, do_type = '', doc_name = 'Test'):
    with get_db_connect() as conn:
        curs = conn.cursor()
        topic_num = str(topic_num)

        if topic_num == '0':
            name = get_lang(conn, 'make_new_topic')
            sub = get_lang(conn, 'make_new_topic')

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
                return redirect(conn, '/')
                
        topic_acl = await acl_check(name_value, 'topic', topic_num)
        topic_view_acl = await acl_check('', 'topic_view', topic_num)
        if topic_view_acl == 1:
            return await re_error(conn, 0)
        elif topic_num == '0':
            if await acl_check('', 'discuss_make_new_thread', topic_num) == 1:
                return await re_error(conn, 0)

        ip = ip_check()

        if flask.request.method == 'POST' and do_type == '':
            if await do_edit_slow_check(conn, 'thread') == 1:
                return await re_error(conn, 42)

            name = flask.request.form.get('topic', 'Test')
            sub = flask.request.form.get('title', 'Test')
            data = flask.request.form.get('content', 'Test').replace('\r', '')
            
            if do_title_length_check(conn, name) == 1:
                return await re_error(conn, 38)
            
            if do_title_length_check(conn, sub, 'topic') == 1:
                return await re_error(conn, 38)
            
            if await do_edit_filter(conn, sub) == 1:
                return await re_error(conn, 21)
            
            if await do_edit_filter(conn, data) == 1:
                return await re_error(conn, 21)
            
            if topic_num == '0':
                curs.execute(db_change("select code from topic order by code + 0 desc limit 1"))
                t_data = curs.fetchall()
                topic_num = str(int(t_data[0][0]) + 1) if t_data else '1'
            
            if flask.request.form.get('content', 'Test') == '':
                return redirect(conn, '/thread/' + topic_num)

            if await captcha_post(conn, flask.request.form.get('g-recaptcha-response', flask.request.form.get('g-recaptcha', ''))) == 1:
                return await re_error(conn, 13)

            today = get_time()

            if topic_acl == 1:
                return await re_error(conn, 0)

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
                    await add_alarm(match, ip, '<a href="/thread/' + topic_num + '#' + num + '">' + html.escape(name) + ' - ' + html.escape(sub) + '#' + num + '</a>')
            
            curs.execute(db_change("select ip from topic where code = ? and id = '1'"), [topic_num])
            ip_data = curs.fetchall()
            if ip_data and ip_or_user(ip_data[0][0]) == 0:
                await add_alarm(ip_data[0][0], ip, '<a href="/thread/' + topic_num + '#' + num + '">' + html.escape(name) + ' - ' + html.escape(sub) + '#' + num + '</a>')

            data = await api_topic_thread_pre_render(conn, data, num, ip, topic_num, name, sub)

            do_add_thread(conn, 
                topic_num,
                data,
                '',
                num
            )
            do_reload_recent_thread(conn, 
                topic_num, 
                today, 
                name, 
                sub
            )

            return redirect(conn, '/thread/' + topic_num + '#' + num)
        else:
            acl_display = 'display: none;' if topic_acl == 1 else ''
            name_display = 'display: none;' if topic_num != '0' else ''

            shortcut = '<div class="opennamu_thread_shortcut" id="thread_shortcut">'
            curs.execute(db_change("select id from topic where code = ? order by id + 0 asc"), [topic_num])
            db_data = curs.fetchall()
            for for_a in db_data:
                shortcut += '<a href="#' + for_a[0] + '">#' + for_a[0] + '</a> '
            
            shortcut += '</div>'

            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [name, await wiki_set(), await wiki_custom(conn), wiki_css(['(' + get_lang(conn, 'discussion') + ')', 0])],
                data = '''
                    <script defer src="/views/main_css/js/route/topic.js''' + cache_v() + '''"></script>
                    <style id="opennamu_list_hidden_style">.opennamu_list_hidden { display: none; }</style>
                    <label><input type="checkbox" onclick="opennamu_list_hidden_remove();" checked> ''' + get_lang(conn, 'remove_hidden') + '''</label>
                    <hr class="main_hr">

                    ''' + shortcut + '''
                    <h2 id="topic_top_title">''' + html.escape(sub) + '''</h2>
                    
                    <div id="opennamu_top_thread"></div>
                    <div id="opennamu_main_thread"></div>
                    <div id="opennamu_reload_thread"></div>
                    <script>
                        window.addEventListener("DOMContentLoaded", function() { 
                            opennamu_get_thread("''' + topic_num + '''", "top");
                            opennamu_get_thread("''' + topic_num + '''");
                        });
                    </script>

                    <a href="javascript:opennamu_thread_blind();">(''' + get_lang(conn, 'hide') + ''' | ''' + get_lang(conn, 'hide_release') + ''')</a> 
                    <a href="javascript:opennamu_thread_delete();">(''' + get_lang(conn, 'delete') + ''')</a>
                    <a href="/thread/''' + topic_num + '/tool">(' + get_lang(conn, 'topic_tool') + ''')</a>
                    <hr class="main_hr">
                    
                    <form style="''' + acl_display + '''" method="post">
                        <div style="''' + name_display + '''">
                            <input placeholder="''' + get_lang(conn, 'document_name') + '''" name="topic" value="''' + html.escape(name_value) + '''">
                            <hr class="main_hr">
                            <input placeholder="''' + get_lang(conn, 'discussion_name') + '''" name="title" value="''' + html.escape(sub_value) + '''">
                            <hr class="main_hr">
                        </div>
                        
                        ''' + await edit_editor(conn, ip, '', 'thread') + '''
                    </form>
                ''',
                menu = [['topic/' + url_pas(name), get_lang(conn, 'list')]]
            ))