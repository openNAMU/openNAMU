from .tool.func import *

def view_raw_2(conn, name, sub_title, num):
    curs = conn.cursor()

    if acl_check(name, 'render') == 1:
        return re_error('/ban')
    
    v_name = name
    sub = ' (' + load_lang('raw') + ')'
    
    if not num:
        num = flask.request.args.get('num', None)
        if num:
            num = int(number_check(num))
    
    if not sub_title and num:
        curs.execute("select title from history where title = ? and id = ? and hide = 'O'", [name, str(num)])
        if curs.fetchall() and admin_check(6) != 1:
            return re_error('/error/3')
        
        curs.execute("select data from history where title = ? and id = ?", [name, str(num)])
        
        sub += ' (r' + str(num) + ')'

        menu = [['history/' + url_pas(name), load_lang('history')]]
    elif sub_title:
        if admin_check(6) != 1:
            curs.execute("select data from topic where id = ? and title = ? and sub = ? and block = ''", [str(num), name, sub_title])
        else:
            curs.execute("select data from topic where id = ? and title = ? and sub = ?", [str(num), name, sub_title])
        
        v_name = load_lang('discussion_raw')
        sub = ' (#' + str(num) + ')'

        menu = [['topic/' + url_pas(name) + '/sub/' + url_pas(sub_title) + '#' + str(num), load_lang('discussion')], ['topic/' + url_pas(name) + '/sub/' + url_pas(sub_title) + '/admin/' + str(num), load_lang('return')]]
    else:
        curs.execute("select data from data where title = ?", [name])
        
        menu = [['w/' + url_pas(name), load_lang('return')]]

    data = curs.fetchall()
    if data:
        p_data = html.escape(data[0][0])
        p_data = '<textarea readonly rows="25">' + p_data + '</textarea>'
        
        return easy_minify(flask.render_template(skin_check(), 
            imp = [v_name, wiki_set(), custom(), other2([sub, 0])],
            data = p_data,
            menu = menu
        ))
    else:
        return re_error('/error/3')