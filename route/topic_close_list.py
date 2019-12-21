from .tool.func import *

def topic_close_list_2(conn, name):
    curs = conn.cursor()

    div = ''
    tool = flask.request.args.get('tool', '')

    plus = ''
    menu = [['topic/' + url_pas(name), load_lang('return')]]

    if tool == 'close':
        curs.execute(db_change("select sub from rd where title = ? and stop = 'O' order by sub asc"), [name])

        sub = load_lang('closed_discussion')
    elif tool == 'agree':
        curs.execute(db_change("select sub from rd where title = ? and agree = 'O' order by sub asc"), [name])

        sub = load_lang('agreed_discussion')
    else:
        sub = load_lang('discussion_list')
        menu = [['w/' + url_pas(name), load_lang('document')]]

        if acl_check(name, 'topic', sub) == 1:
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
            <hr class=\"main_hr\">
            <form style="''' + display + '" method="post" action="/thread/' + topic_num + '''">
                <input placeholder="''' + load_lang('discussion_name') + '''" name="title">
                <hr class=\"main_hr\">
                <textarea rows="10" id="content" placeholder="''' + load_lang('content') + '''" name="content"></textarea>
                <hr class=\"main_hr\">
                ''' + captcha_get() + (ip_warring() if display == '' else '') + '''
                <input style="display: none;" name="topic" value="''' + name + '''">
                <button type="submit">''' + load_lang('send') + '''</button>
                <button id="preview" type="button" onclick="load_preview(\'''' + url_pas(name) + '\')">' + load_lang('preview') + '''</button>
            </form>
            <hr class=\"main_hr\">
            <div id="see_preview"></div>
        '''

        curs.execute(db_change("select sub from rd where title = ? order by date desc"), [name])

    t_num = 0
    for data in curs.fetchall():
        t_num += 1

        curs.execute(db_change("select code from topic where title = ? and sub = ? and id = '1'"), [name, data[0]])
        first_topic = curs.fetchall()
        if first_topic:
            it_p = 0

            if tool == '':
                curs.execute(db_change("select title from rd where title = ? and sub = ? and stop = 'O' order by sub asc"), [name, data[0]])
                if curs.fetchall():
                    it_p = 1

            if it_p != 1:
                curs.execute(db_change("select id from topic where title = ? and sub = ? order by date desc limit 1"), [name, data[0]])
                t_data = curs.fetchall()

                div += '''
                    <h2><a href="/thread/''' + first_topic[0][0] + '">' + str(t_num) + '. ' + data[0] + '''</a></h2>
                    <div id="topic_pre_''' + str(t_num) + '''"></div>
                    <div id="topic_back_pre_''' + str(t_num) + '''"></div>
                    <script>
                        topic_list_load(''' + first_topic[0][0] + ', 1, "topic_pre_' + str(t_num) + '''");
                        if(''' + str(t_data[0][0]) + ''' !== 1) {
                            topic_list_load(''' + first_topic[0][0] + ', ' + t_data[0][0] + ', "topic_back_pre_' + str(t_num) + '''");
                        }
                    </script>
                '''

    if div == '':
        plus = re.sub('^<br>', '', plus)

    return easy_minify(flask.render_template(skin_check(),
        imp = [name, wiki_set(), custom(), other2([' (' + sub + ')', 0])],
        data = div + plus,
        menu = menu
    ))
