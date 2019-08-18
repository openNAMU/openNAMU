from .tool.func import *
import pymysql

def topic_close_list_2(conn, name, tool):
    curs = conn.cursor()
    
    div = ''
    
    if flask.request.method == 'POST':
        t_num = ''
        
        while 1:
            curs.execute("select title from topic where title = %s and sub = %s limit 1", [name, flask.request.form.get('topic', None) + t_num])
            if curs.fetchall():
                if t_num == '':
                    t_num = ' 2'
                else:
                    t_num = ' ' + str(int(t_num.replace(' ', '')) + 1)
            else:
                break

        return redirect('/topic/' + url_pas(name) + '/sub/' + url_pas(flask.request.form.get('topic', 'test') + t_num))
    else:
        plus = ''
        menu = [['topic/' + url_pas(name), load_lang('return')]]
        
        if tool == 'close':
            curs.execute("select sub from rd where title = %s and stop = 'O' order by sub asc", [name])
            
            sub = load_lang('closed_discussion')
        elif tool == 'agree':
            curs.execute("select sub from rd where title = %s and agree = 'O' order by sub asc", [name])
            
            sub = load_lang('agreed_discussion')
        else:
            curs.execute("select sub from rd where title = %s order by date desc", [name])
            
            sub = load_lang('discussion_list')
            
            menu = [['w/' + url_pas(name), load_lang('document')]]
            
            plus =  '''
                <a href="/topic/''' + url_pas(name) + '''/close">(''' + load_lang('closed_discussion') + ''')</a> <a href="/topic/''' + url_pas(name) + '''/agree">(''' + load_lang('agreed_discussion') + ''')</a>
                <hr class=\"main_hr\">
                <input placeholder="''' + load_lang('discussion_name') + '''" name="topic" type="text">
                <hr class=\"main_hr\">
                <button type="submit">''' + load_lang('go') + '''</button>
            '''

        t_num = 0
        for data in curs.fetchall():
            t_num += 1
            
            curs.execute("select data, date, ip, block from topic where title = %s and sub = %s and id = '1'", [name, data[0]])
            if curs.fetchall():                
                it_p = 0
                
                if sub == load_lang('discussion_list'):
                    curs.execute("select title from rd where title = %s and sub = %s and stop = 'O' order by sub asc", [name, data[0]])
                    if curs.fetchall():
                        it_p = 1
                
                if it_p != 1:
                    curs.execute("select id from topic where title = %s and sub = %s order by date desc limit 1", [name, data[0]])
                    t_data = curs.fetchall()
                
                    div += '''
                        <h2><a href="/topic/''' + url_pas(name) + '''/sub/''' + url_pas(data[0]) + '''">''' + str(t_num) + '''. ''' + data[0] + '''</a></h2>
                        <div id="topic_pre_''' + str(t_num) + '''"></div>
                        <div id="topic_back_pre_''' + str(t_num) + '''"></div>
                        <script>
                            topic_list_load("''' + name + '''", "''' + data[0] + '''", "1", "topic_pre_''' + str(t_num) + '''");
                            if("''' + str(t_num) + '''" !== "1") {
                                topic_list_load("''' + name + '''", "''' + data[0] + '''", "''' + t_data[0][0] + '''", "topic_back_pre_''' + str(t_num) + '''");
                            }
                        </script>
                    '''

        if div == '':
            plus = re.sub('^<br>', '', plus)
        
        return easy_minify(flask.render_template(skin_check(), 
            imp = [name, wiki_set(), custom(), other2([' (' + sub + ')', 0])],
            data = '<form method="post">' + div + plus + '</form>',
            menu = menu
        ))
