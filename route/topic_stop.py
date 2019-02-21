from .tool.func import *

def topic_stop_2(conn, name, sub, tool):
    curs = conn.cursor()

    if tool == 'close':
        set_list = [
            'O', 
            'S', 
            load_lang('close', 1), 
            load_lang('open', 1)
        ]
    elif tool == 'stop':
        set_list = [
            '', 
            'O', 
            load_lang('stop', 1), 
            load_lang('restart', 1)
        ]
    elif tool == 'agree':
        pass
    else:
        return redirect('/topic/' + url_pas(name) + '/sub/' + url_pas(sub))

    if admin_check(3, 'topic ' + tool + ' (' + name + ' - ' + sub + ')') != 1:
        return re_error('/error/3')

    ip = ip_check()
    time = get_time()
    
    curs.execute("select id from topic where title = ? and sub = ? order by id + 0 desc limit 1", [name, sub])
    topic_check = curs.fetchall()
    if topic_check:
        if tool == 'agree':
            curs.execute("select title from rd where title = ? and sub = ? and agree = 'O'", [name, sub])
            if curs.fetchall():
                curs.execute("insert into topic (id, title, sub, data, date, ip, block, top) values (?, ?, ?, '" + load_lang('agreement', 1) + " X', ?, ?, '', '1')", [str(int(topic_check[0][0]) + 1), name, sub, time, ip])
                curs.execute("update rd set agree = '' where title = ? and sub = ?", [name, sub])
            else:
                curs.execute("insert into topic (id, title, sub, data, date, ip, block, top) values (?, ?, ?, '" + load_lang('agreement', 1) + " O', ?, ?, '', '1')", [str(int(topic_check[0][0]) + 1), name, sub, time, ip])
                curs.execute("update rd set agree = 'O' where title = ? and sub = ?", [name, sub])
        else:
            curs.execute("select title from rd where title = ? and sub = ? and stop = ?", [name, sub, set_list[0]])
            if curs.fetchall():
                curs.execute("insert into topic (id, title, sub, data, date, ip, block, top) values (?, ?, ?, ?, ?, ?, '', '1')", [str(int(topic_check[0][0]) + 1), name, sub, set_list[3], time, ip])
                curs.execute("update rd set stop = '' where title = ? and sub = ?", [name, sub])
            else:
                curs.execute("insert into topic (id, title, sub, data, date, ip, block, top) values (?, ?, ?, ?, ?, ?, '', '1')", [str(int(topic_check[0][0]) + 1), name, sub, set_list[2], time, ip])
                curs.execute("update rd set stop = ? where title = ? and sub = ?", [set_list[0], name, sub])
        
        rd_plus(name, sub, time)
        
        conn.commit()
        
    return redirect('/topic/' + url_pas(name) + '/sub/' + url_pas(sub))    
