from .tool.func import *

def topic_close_list_2(conn, name):
    curs = conn.cursor()

    div = ''
    tool = flask.request.args.get('tool', '')

    plus = ''
    menu = [['topic/' + url_pas(name), load_lang('return')]]

    if tool == 'close':
        curs.execute(db_change("select code, sub from rd where title = ? and stop = 'O' order by sub asc"), [name])

        sub = load_lang('closed_discussion')
    elif tool == 'agree':
        curs.execute(db_change("select code, sub from rd where title = ? and agree = 'O' order by sub asc"), [name])

        sub = load_lang('agreed_discussion')
    else:
        sub = load_lang('discussion_list')
        menu = [['w/' + url_pas(name), load_lang('document')]]

        if acl_check(name, 'topic', None) == 1:
            display = 'display: none;'
        else:
            display = ''

        curs.execute(db_change("select code from topic order by code + 0 desc limit 1"))
        t_data = curs.fetchall()
        if t_data:
            topic_num = str(int(t_data[0][0]) + 1)
        else:
            topic_num = '1'

        plus = '''
            <a href="/topic/''' + url_pas(name) + '?tool=close">(' + load_lang('closed_discussion') + ')</a> <a href="/topic/' + url_pas(name) + '?tool=agree">(' + load_lang('agreed_discussion') + ''')</a>
            <hr class="main_hr">
            <form style="''' + display + '" method="post" action="/thread/' + topic_num + '''">
                <input placeholder="''' + load_lang('discussion_name') + '''" name="title">
                <hr class="main_hr">
                <textarea rows="10" id="content" placeholder="''' + load_lang('content') + '''" name="content"></textarea>
                <hr class="main_hr">
                ''' + captcha_get() + (ip_warring() if display == '' else '') + '''
                <input style="display: none;" name="topic" value="''' + name + '''">
                <button type="submit">''' + load_lang('send') + '''</button>
                <button id="preview" type="button" onclick="load_preview(\'''' + url_pas(name) + '\')">' + load_lang('preview') + '''</button>
            </form>
            <hr class="main_hr">
            <div id="see_preview"></div>
        '''

        curs.execute(db_change("select code, sub from rd where title = ? and stop != 'O' order by date desc"), [name])

    for data in curs.fetchall():
        curs.execute(db_change("select id from topic where code = ? order by id + 0 desc limit 1"), [data[0]])
        t_data = curs.fetchall()

        div += '''
            <h2><a href="/thread/''' + data[0] + '">' + data[0] + '. ' + data[1] + '''</a></h2>
            <div id="topic_pre_''' + data[0] + '''"></div>
            <div id="topic_back_pre_''' + data[0] + '''"></div>
            <script>
                new_topic_load(''' + data[0] + ', 0, 1, "?num=1", "topic_pre_' + data[0] + '''");
                if(''' + t_data[0][0] + ''' !== 1) {
                    new_topic_load(''' + data[0] + ', 0, 1, "?num=' + t_data[0][0] + '", "topic_back_pre_' + data[0] + '''");
                }
            </script>
        '''

    if div == '':
        plus = re.sub(r'^<br>', '', plus)

    return easy_minify(flask.render_template(skin_check(),
        imp = [name, wiki_set(), custom(), other2(['(' + sub + ')', 0])],
        data = div + plus,
        menu = menu
    ))
