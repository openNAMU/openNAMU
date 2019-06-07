from .tool.func import *

def edit_move_2(conn, name):
    curs = conn.cursor()

    if acl_check(name) == 1:
        return re_error('/ban')

    if flask.request.method == 'POST':
        if captcha_post(flask.request.form.get('g-recaptcha-response', '')) == 1:
            return re_error('/error/13')
        else:
            captcha_post('', 0)

        curs.execute("select title from history where title = ?", [flask.request.form.get('title', None)])
        if curs.fetchall():
            if admin_check(None, 'merge documents') == 1:
                curs.execute("select data from data where title = ?", [flask.request.form.get('title', None)])
                data = curs.fetchall()
                if data:            
                    curs.execute("delete from data where title = ?", [flask.request.form.get('title', None)])
                    curs.execute("delete from back where link = ?", [flask.request.form.get('title', None)])
                
                curs.execute("select data from data where title = ?", [name])
                data = curs.fetchall()
                if data:            
                    curs.execute("update data set title = ? where title = ?", [flask.request.form.get('title', None), name])
                    curs.execute("update back set link = ? where link = ?", [flask.request.form.get('title', None), name])
                    
                    data_in = data[0][0]
                else:
                    data_in = ''

                history_plus(
                    name, 
                    data_in, 
                    get_time(), 
                    ip_check(), 
                    flask.request.form.get('send', ''), 
                    '0',
                    'marge <a>' + name + '</a> - <a>' + flask.request.form.get('title', 'test') + '</a> move'
                )

                curs.execute("update back set type = 'no' where title = ? and not type = 'cat' and not type = 'no'", [name])
                curs.execute("delete from back where title = ? and not type = 'cat' and type = 'no'", [flask.request.form.get('title', None)])

                curs.execute("select id from history where title = ? order by id + 0 desc limit 1", [flask.request.form.get('title', None)])
                data = curs.fetchall()
                
                num = data[0][0]

                curs.execute("select id from history where title = ? order by id + 0 asc", [name])
                data = curs.fetchall()
                for move in data:
                    curs.execute("update history set title = ?, id = ? where title = ? and id = ?", [flask.request.form.get('title', None), str(int(num) + int(move[0])), name, move[0]])

                conn.commit()

                return redirect('/w/' + url_pas(flask.request.form.get('title', None)))
            else:
                return re_error('/error/19')
        else:
            curs.execute("select data from data where title = ?", [name])
            data = curs.fetchall()
            if data:            
                curs.execute("update data set title = ? where title = ?", [flask.request.form.get('title', None), name])
                curs.execute("update back set link = ? where link = ?", [flask.request.form.get('title', None), name])
                
                data_in = data[0][0]
            else:
                data_in = ''
                
            history_plus(
                name, 
                data_in, 
                get_time(), 
                ip_check(), 
                flask.request.form.get('send', ''), 
                '0',
                '<a>' + name + '</a> - <a>' + flask.request.form.get('title', 'test') + '</a> move'
            )
            
            curs.execute("update back set type = 'no' where title = ? and not type = 'cat' and not type = 'no'", [name])
            curs.execute("delete from back where title = ? and not type = 'cat' and type = 'no'", [flask.request.form.get('title', None)])

            curs.execute("update history set title = ? where title = ?", [flask.request.form.get('title', None), name])
            conn.commit()

            return redirect('/w/' + url_pas(flask.request.form.get('title', None)))
    else:            
        return easy_minify(flask.render_template(skin_check(), 
            imp = [name, wiki_set(), custom(), other2([' (' + load_lang('move') + ')', 0])],
            data =  '''
                    <form method="post">
                        ''' + ip_warring() + '''
                        <input placeholder="''' + load_lang('document_name') + '" value="' + name + '''" name="title" type="text">
                        <hr class=\"main_hr\">
                        <input placeholder="''' + load_lang('why') + '''" name="send" type="text">
                        <hr class=\"main_hr\">
                        ''' + captcha_get() + '''
                        <button type="submit">''' + load_lang('move') + '''</button>
                    </form>
                    ''',
            menu = [['w/' + url_pas(name), load_lang('return')]]
        ))