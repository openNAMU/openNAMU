from .tool.func import *

def topic_tool_2(conn, name, sub):
    curs = conn.cursor()
    
    curs.execute("select id from topic where title = ? and sub = ? limit 1", [name, sub])
    topic_exist = curs.fetchall()
    if not topic_exist:
        return re_error('/topic/' + url_pas(name) + '/sub/' + url_pas(sub))

    all_data = ''

    if admin_check(3) == 1:
        all_data = '<h2>' + load_lang('topic_state') + '</h2><ul><li><a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/tool/close">'

        curs.execute("select title from rd where title = ? and sub = ? and stop = 'O'", [name, sub])
        if curs.fetchall():
            all_data += load_lang('topic_open')
        else:
            all_data += load_lang('topic_close')
        
        all_data += '</a></li><li><a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/tool/stop">'

        curs.execute("select title from rd where title = ? and sub = ? and stop = 'S'", [name, sub])
        if curs.fetchall():
            all_data += load_lang('topic_restart')
        else:
            all_data += load_lang('topic_stop')
            
        all_data += '</a></li><li><a href="/topic/' + url_pas(name) + '/sub/' + url_pas(sub) + '/tool/agree">'
        
        curs.execute("select title from rd where title = ? and sub = ? and agree = 'O'", [name, sub])
        if curs.fetchall():
            all_data += load_lang('topic_destruction')
        else:
            all_data += load_lang('topic_agreement')
        
        all_data += '</a></li></ul>'

    all_data += '<h2>' + load_lang('tool') + '</h2><ul><li><a id="reload" href="javascript:void(0);" onclick="req_alarm();">' + load_lang('use_push_alarm') + '</a></li></ul>'

    return easy_minify(flask.render_template(skin_check(), 
        imp = [name, wiki_set(), custom(), other2([' (' + load_lang('topic_tool') + ')', 0])],
        data = all_data,
        menu = [['topic/' + url_pas(name) + '/sub/' + url_pas(sub), load_lang('return')]]
    ))