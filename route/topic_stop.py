from .tool.func import *
import pymysql

def topic_stop_2(conn, name, sub, tool):
    curs = conn.cursor()

    if admin_check(3, 'topic ' + tool + ' (' + name + ' - ' + sub + ')') != 1:
        return re_error('/error/3')

    ip = ip_check()
    time = get_time()

    if flask.request.method == 'POST' or tool == 'agree' or tool == 'stop':
        curs.execute("select id from topic where title = %s and sub = %s order by id + 0 desc limit 1", [name, sub])
        topic_check = curs.fetchall()
        if topic_check:
            if tool == 'agree':
                curs.execute("select title from rd where title = %s and sub = %s and agree = 'O'", [name, sub])
                if curs.fetchall():
                    curs.execute("insert into topic (id, title, sub, data, date, ip, block, top) values (%s, %s, %s, 'ok', %s, %s, '', '1')", [str(int(topic_check[0][0]) + 1), name, sub, time, ip])
                    curs.execute("update rd set agree = '' where title = %s and sub = %s", [name, sub])
                else:
                    curs.execute("insert into topic (id, title, sub, data, date, ip, block, top) values (%s, %s, %s, 'break', %s, %s, '', '1')", [str(int(topic_check[0][0]) + 1), name, sub, time, ip])
                    curs.execute("update rd set agree = 'O' where title = %s and sub = %s", [name, sub])
            else:
                # 이 쪽 잘못 짜여 있음 내일 수정 필요
                if tool == 'close':
                    why = flask.request.form.get('why', '')
                    set_list = [
                        'O', 
                        'S', 
                        'close' + (('[br][br]why : ' + why) if why else ''), 
                        'open' + (('[br][br]why : ' + why) if why else '')
                    ]
                elif tool == 'stop':
                    set_list = [
                        '', 
                        'O', 
                        'stop', 
                        'restart'
                    ]
                else:
                    return redirect('/topic/' + url_pas(name) + '/sub/' + url_pas(sub))

                curs.execute("select title from rd where title = %s and sub = %s and stop = %s", [name, sub, set_list[1]])
                if curs.fetchall():
                    curs.execute("insert into topic (id, title, sub, data, date, ip, block, top) values (%s, %s, %s, %s, %s, %s, '', '1')", [
                        str(int(topic_check[0][0]) + 1),
                        name, 
                        sub, 
                        set_list[3], 
                        time, 
                        ip
                    ])
                    curs.execute("update rd set stop = '' where title = %s and sub = %s", [name, sub])
                else:
                    curs.execute("insert into topic (id, title, sub, data, date, ip, block, top) values (%s, %s, %s, %s, %s, %s, '', '1')", [
                        str(int(topic_check[0][0]) + 1), 
                        name, 
                        sub, 
                        set_list[2], 
                        time, 
                        ip
                    ])
                    curs.execute("update rd set stop = %s where title = %s and sub = %s", [set_list[0], name, sub])
            
            rd_plus(name, sub, time)
            
            conn.commit()

        return redirect('/topic/' + url_pas(name) + '/sub/' + url_pas(sub))    
    else:
        return easy_minify(flask.render_template(skin_check(),
            imp = [name, wiki_set(), custom(), other2([' (' + tool + ')', 0])],
            data = '''
                <span>''' + load_lang('markup_enabled') + '''</span>
                <hr class=\"main_hr\">
                <form method="post">
                    <input placeholder="''' + load_lang('why') + '''" name="why" type="text">
                    <hr class=\"main_hr\">
                    <button type="submit">''' + load_lang('save') + '''</button>
                </form>
            ''',
            menu = [['topic/' + url_pas(name) + '/sub/' + url_pas(sub), load_lang('return')]]
        ))
